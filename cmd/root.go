package cmd

import (
    "os"
    "path/filepath"

    "github.com/spf13/cobra"
    "github.com/spf13/viper"
    "github.com/mattjh1/goetl/config"
    "github.com/mattjh1/goetl/config/logger"
)

var cfgFile string

var rootCmd = &cobra.Command{
    Use:   "goetl",
    Short: "goetl is a simple ETL tool",
    Long:  `goetl is a simple ETL tool to extract, transform, and load data.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        logger.Log.Error(err) 
        os.Exit(1)
    }
}

func init() {
    cobra.OnInitialize(initConfig)

    rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.goetl/goetl.yaml)")
    rootCmd.Flags().String("mode", "", "Mode of operation (e.g., development, production)")
    rootCmd.Flags().String("llm_provider", "", "LLM provider")
    rootCmd.Flags().String("llm_model", "", "LLM model ID")
    rootCmd.Flags().Float64("llm_temperature", 0.0, "LLM temperature")
    rootCmd.Flags().String("project_name", "", "Project name")
    rootCmd.Flags().String("api_str", "", "API string")
    rootCmd.Flags().String("source_path", "", "Source path string")
    rootCmd.Flags().String("openrouter_api_base", "", "OpenRouter API base URL")
    rootCmd.Flags().String("openrouter_api_key", "", "OpenRouter API key")
    rootCmd.Flags().String("ollama_api_base", "", "Ollama API base URL")
    rootCmd.Flags().String("emb_model_id", "", "Embedding model ID")
    rootCmd.Flags().String("neo4j_uri", "", "Neo4j URI")
    rootCmd.Flags().String("neo4j_username", "", "Neo4j username")
    rootCmd.Flags().String("neo4j_password", "", "Neo4j password")
    rootCmd.Flags().StringSlice("backend_cors_origins", nil, "Backend CORS origins")

    // Bind flags to Viper
    viper.BindPFlag("mode", rootCmd.Flags().Lookup("mode"))
    viper.BindPFlag("llm_settings.provider", rootCmd.Flags().Lookup("llm_provider"))
    viper.BindPFlag("llm_settings.model", rootCmd.Flags().Lookup("llm_model"))
    viper.BindPFlag("llm_settings.temperature", rootCmd.Flags().Lookup("llm_temperature"))
    viper.BindPFlag("project_name", rootCmd.Flags().Lookup("project_name"))
    viper.BindPFlag("api_str", rootCmd.Flags().Lookup("api_str"))
    viper.BindPFlag("openrouter_api_base", rootCmd.Flags().Lookup("openrouter_api_base"))
    viper.BindPFlag("openrouter_api_key", rootCmd.Flags().Lookup("openrouter_api_key"))
    viper.BindPFlag("ollama_api_base", rootCmd.Flags().Lookup("ollama_api_base"))
    viper.BindPFlag("emb_model_id", rootCmd.Flags().Lookup("emb_model_id"))
    viper.BindPFlag("neo4j_uri", rootCmd.Flags().Lookup("neo4j_uri"))
    viper.BindPFlag("neo4j_username", rootCmd.Flags().Lookup("neo4j_username"))
    viper.BindPFlag("neo4j_password", rootCmd.Flags().Lookup("neo4j_password"))
    viper.BindPFlag("backend_cors_origins", rootCmd.Flags().Lookup("backend_cors_origins"))
}

func initConfig() {
    if cfgFile != "" {
        viper.SetConfigFile(cfgFile)
    } else {
        home, err := os.UserHomeDir()
        cobra.CheckErr(err)
        viper.AddConfigPath(filepath.Join(home, ".goetl"))
        viper.SetConfigName("goetl")
    }

    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err == nil {
        logger.Log.Infof("Using config file: %s", viper.ConfigFileUsed())
    }

    _, err := config.InitConfig()
    if err != nil {
        logger.Log.Fatalf("Error initializing config: %v", err)
        os.Exit(1)
    }
}

