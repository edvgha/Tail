package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"tail.server/app/optimizer/code/misc"
	"tail.server/app/optimizer/code/space"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	Spaces   map[string]*space.Space
	Cache    *misc.Cache[string, space.ExploreData]
	CacheTTL time.Duration
	TLog     zerolog.Logger
)

func initLogger(level string) zerolog.Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	logLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		panic("error during log initialization")
	}

	fileLogger, err := os.OpenFile(
		"/tmp/tail_app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic("error on fileLogger")
	}

	var output io.Writer = zerolog.ConsoleWriter{
		Out:           fileLogger,
		NoColor:       false,
		TimeFormat:    time.RFC3339Nano,
		FormatLevel:   func(i interface{}) string { return strings.ToUpper(fmt.Sprintf("[%s]", i)) },
		FormatCaller:  func(i interface{}) string { return filepath.Base(fmt.Sprintf("%s", i)) },
		FormatMessage: func(i interface{}) string { return fmt.Sprintf("| %s |", i) },
		PartsExclude:  []string{zerolog.TimestampFieldName},
	}

	return zerolog.New(output).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()
}

func initGlobalVars(cfg misc.Config) {
	Cache = misc.New[string, space.ExploreData]()
	CacheTTL = time.Duration(cfg.CacheTTL) * time.Second

	var err error
	Spaces, err = space.LoadSpaces(cfg, TLog)
	if err != nil {
		TLog.Error().Msgf("%s", err.Error())
		panic(fmt.Sprintf("error on spaces.LoadSpaces() %s", err.Error()))
	}
}

func init() {
	cfg := misc.Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		panic(fmt.Sprintf("error during config initialization %s", err.Error()))
	}

	TLog = initLogger(cfg.LogLevel)

	TLog.Info().Msgf(
		"Server will be using %d of %d available threads for goroutine calls",
		runtime.GOMAXPROCS(0),
		runtime.NumCPU())

	TLog.Info().Msgf("Config %+v", cfg)

	initGlobalVars(cfg)

	// Background periodic task to learn winning curve
	for _, s := range Spaces {
		s.BackgroundTask()
	}
	TLog.Debug().Msg("Initialization is done")
}

func App() (errReturn error) {
	http.HandleFunc("/optimize", optimizeHandler)
	http.HandleFunc("/feedback", feedbackHandler)
	http.HandleFunc("/space", spaceHandler)
	TLog.Info().Msg("Server listening on :8000")
	return http.ListenAndServe(":8000", nil)
}
