package misc

type Config struct {
	BufferSize              int     `envconfig:"BUFFER_SIZE" required:"true"`
	Discount                float64 `envconfig:"DISCOUNT" required:"true"`
	DesiredExplorationSpeed float64 `envconfig:"DESIRED_EXPLORATION_SPEED" required:"true"`
	LogLevel                string  `envconfig:"LOG_LEVEL" required:"true"`
	LevelSize               int     `envconfig:"LEVEL_SIZE" required:"true"`
	BucketSize              int     `envconfig:"BUCKET_SIZE" required:"true"`
	SpaceDescFile           string  `envconfig:"SPACE_DESC_FILE" required:"true"`
	CacheTTL                int     `envconfig:"CACHE_TTL" required:"true"`
}
