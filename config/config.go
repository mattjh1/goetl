package config

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/spf13/viper"
)

type LLMSettings struct {
    Provider    string  `mapstructure:"provider"`
    Model       string  `mapstructure:"model"`
    Temperature float64 `mapstructure:"temperature"`
}

type ProviderConfig struct {
    BaseURL string `mapstructure:"base_url"`
    APIKey  string `mapstructure:"api_key"`
}

type Config struct {
    Mode              string                    `mapstructure:"mode"`
    LLMSettings       LLMSettings               `mapstructure:"llm_settings"`
    ProjectName       string                    `mapstructure:"project_name"`
    APIStr            string                    `mapstructure:"api_str"`
    SourcePath        string                    `mapstructure:"source_path"`
		TikaServerURL			string										`mapstructure:"tika_server_url"`
    OpenRouterAPIBase string                    `mapstructure:"openrouter_api_base"`
    OpenRouterAPIKey  string                    `mapstructure:"openrouter_api_key"`
    OllamaAPIBase     string                    `mapstructure:"ollama_api_base"`
    EmbModelID        string                    `mapstructure:"emb_model_id"`
    Neo4jURI          string                    `mapstructure:"neo4j_uri"`
    Neo4jUsername     string                    `mapstructure:"neo4j_username"`
    Neo4jPassword     string                    `mapstructure:"neo4j_password"`
    Providers         map[string]ProviderConfig `mapstructure:"providers"`
    BackendCORSOrigins []string                 `mapstructure:"backend_cors_origins"`
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
    viper.SetDefault("mode", "development")
    viper.SetDefault("source_path", dataDir)
    viper.SetDefault("llm_settings.provider", "ollama")
    viper.SetDefault("llm_settings.model", "phi3:14b")
    viper.SetDefault("llm_settings.temperature", 0.0)
    viper.SetDefault("project_name", "goetl")
    viper.SetDefault("tika_server_url", "goetl")
    viper.SetDefault("openrouter_api_base", "https://openrouter.ai/api/v1")
    viper.SetDefault("openrouter_api_key", "your_openrouter_api_key")
    viper.SetDefault("ollama_api_base", "http://localhost:11434")
    viper.SetDefault("emb_model_id", "nomic-embed-text")
    viper.SetDefault("neo4j_uri", "bolt://localhost:7687")
    viper.SetDefault("neo4j_username", "neo4j")
    viper.SetDefault("neo4j_password", "your_neo4j_password")
    viper.SetDefault("providers", map[string]ProviderConfig{
        "openrouter": {BaseURL: "https://openrouter.ai/api/v1", APIKey: "your_openrouter_api_key"},
        "ollama":     {BaseURL: "http://localhost:11434/v1", APIKey: "ollama"},
        "openai":     {BaseURL: "https://api.openai.com", APIKey: "your_openai_api_key"},
    })
    viper.SetDefault("backend_cors_origins", []string{"*"})

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

    return &config, nil
}
