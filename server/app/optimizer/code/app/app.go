package app

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"runtime"
	"tail.server/app/optimizer/code/misc"
	"tail.server/app/optimizer/code/space"
	"time"

	"github.com/kelseyhightower/envconfig"
)

var (
	Spaces   map[string]*space.Space
	Cache    *misc.Cache[string, space.ExploreData]
	CacheTTL time.Duration
)

func initGlobalVars(cfg misc.Config) {
	Cache = misc.New[string, space.ExploreData]()
	CacheTTL = time.Duration(cfg.CacheTTL) * time.Second

	var err error
	Spaces, err = space.LoadSpaces(cfg)
	if err != nil {
		log.Error().Msgf("%s", err.Error())
		panic(fmt.Sprintf("error on spaces.LoadSpaces() %s", err.Error()))
	}
}

func init() {
	cfg := misc.Config{}
	if err := envconfig.Process("", &cfg); err != nil {
		panic(fmt.Sprintf("error during config initialization %s", err.Error()))
	}
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic("error during log initialization")
	}
	zerolog.SetGlobalLevel(logLevel)
	zerolog.TimeFieldFormat = time.StampNano

	log.Debug().Msgf(
		"Server will be using %d of %d available threads for goroutine calls",
		runtime.GOMAXPROCS(0),
		runtime.NumCPU())

	log.Debug().Msgf("Config %+v", cfg)

	initGlobalVars(cfg)

	// Background periodic task to learn winning curve
	for _, s := range Spaces {
		s.BackgroundTask()
	}
	log.Debug().Msgf("Initialization is done")
}

func App() (errReturn error) {
	http.HandleFunc("/optimize", optimizeHandler)
	http.HandleFunc("/feedback", feedbackHandler)
	http.HandleFunc("/space", spaceHandler)
	log.Debug().Msg("Server listening on :8000")
	return http.ListenAndServe(":8000", nil)
}
