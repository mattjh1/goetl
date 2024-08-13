package config

// TODO: add support for the other structs here.
// Starting off, only redis works

type RedisConfig struct {
    URL string `mapstructure:"url"`
    Index string `mapstructure:"index"`
}

type PostgresConfig struct {
    URL      string `mapstructure:"url"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
    DBName   string `mapstructure:"name"`
    SSLMode  string `mapstructure:"sslmode"`
}

type databaseCfg struct {
    Type   string      `mapstructure:"type"`
    Config interface{} `mapstructure:"-"`
}

func (db *databaseCfg) GetRedisConfig() (*RedisConfig, bool) {
    config, ok := db.Config.(RedisConfig)
    return &config, ok
}

func (db *databaseCfg) GetPostgresConfig() (*PostgresConfig, bool) {
    config, ok := db.Config.(PostgresConfig)
    return &config, ok
}
