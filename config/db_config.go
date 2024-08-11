package config

// TODO: add support for the other structs here.
// Starting off, only redis works

type RedisConfig struct {
    URL string `mapstructure:"database_url"`
    Index string `mapstructure:"index"`
}

type PostgresConfig struct {
    URL      string `mapstructure:"database_url"`
    Username string `mapstructure:"database_username"`
    Password string `mapstructure:"database_password"`
    DBName   string `mapstructure:"database_name"`
    SSLMode  string `mapstructure:"database_sslmode"`
}

type databaseCfg struct {
    Type   string      `mapstructure:"database_type"`
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
