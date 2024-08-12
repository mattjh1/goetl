package config

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/viper"
)

type Config struct {
    Mode         string        `mapstructure:"mode"`
    SourcePath   string        `mapstructure:"source_path"`
    ProjectName  string        `mapstructure:"project_name"`
    TikaServerURL string       `mapstructure:"tika_server_url"`
    EmbAPIBase   string        `mapstructure:"emb_api_base"`
    EmbModelID   string        `mapstructure:"emb_model_id"`
    ChunkSize    int        	 `mapstructure:"chunk_size"`
    ChunkOverlap int        	 `mapstructure:"chunk_overlap"`
    Database     databaseCfg   `mapstructure:"database"`
}

func InitConfig() (*Config, error) {
    home, err := os.UserHomeDir()
    if err != nil {
        return nil, fmt.Errorf("could not determine home directory: %w", err)
    }

    configDir := filepath.Join(home, ".goetl")
		dataDir := filepath.Join(configDir, "data")
    configFile := filepath.Join(configDir, "goetl.yaml")


    viper.SetConfigFile(configFile)
    viper.SetConfigType("yaml")

    // Set default values
		viper.SetDefault("database.type", "redis")
    viper.SetDefault("database.url", "redis://localhost:6379")
    viper.SetDefault("database.index", "redis_index")
    viper.SetDefault("mode", "development")
    viper.SetDefault("source_path", dataDir)
    viper.SetDefault("project_name", "goetl")
    viper.SetDefault("tika_server_url", "http://localhost:9998")
    viper.SetDefault("emb_api_base", "http://localhost:11434")
    viper.SetDefault("emb_model_id", "nomic-embed-text")
    viper.SetDefault("chunk_size", 512)
    viper.SetDefault("chunk_overlap", 64)

    // Create config directory if it doesn't exist
    if _, err := os.Stat(configDir); os.IsNotExist(err) {
        err := os.Mkdir(configDir, 0755)
        if err != nil {
            return nil, fmt.Errorf("could not create config directory: %w", err)
        }
    }

    // Create data directory if it doesn't exist
    if _, err := os.Stat(dataDir); os.IsNotExist(err) {
        err := os.Mkdir(dataDir, 0755)
        if err != nil {
            return nil, fmt.Errorf("could not create config directory: %w", err)
        }
    }

    // Create config file with defaults if it doesn't exist
    if _, err := os.Stat(configFile); os.IsNotExist(err) {
        if err := viper.SafeWriteConfigAs(configFile); err != nil {
            return nil, fmt.Errorf("could not write default config file: %w", err)
        }
        fmt.Printf("Config file created at %s with default values.\n", configFile)
    }

    // Read in the config file
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("could not read config file: %w", err)
    }

    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, fmt.Errorf("could not unmarshal config: %w", err)
    }

		// Configure the database based on the type
    switch config.Database.Type {
    case "redis":
        redisConfig := RedisConfig{
            URL: viper.GetString("database.url"),
            Index: viper.GetString("database.index"),
        }
        config.Database.Config = redisConfig

    case "postgres":
        postgresConfig := PostgresConfig{
            URL:      viper.GetString("database.url"),
            Username: viper.GetString("database.username"),
            Password: viper.GetString("database.password"),
            DBName:   viper.GetString("database.name"),
            SSLMode:  viper.GetString("database.sslmode"),
        }
        config.Database.Config = postgresConfig

    default:
        return nil, fmt.Errorf("unsupported database type: %s", config.Database.Type)
    }

    return &config, nil
}
