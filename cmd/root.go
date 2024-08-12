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
    rootCmd.Flags().String("project_name", "", "Project name")
    rootCmd.Flags().String("source_path", "", "Source path string")
    rootCmd.Flags().String("tika_server_url", "", "Apache Tika server URL")
    rootCmd.Flags().String("emb_api_base", "", "Embedding API base URL")
    rootCmd.Flags().String("emb_model_id", "", "Embedding model ID")
    rootCmd.Flags().String("chunk_size", "", "Text split chunk size")
    rootCmd.Flags().String("chunk_overlap", "", "Text split chunk overlap")

    // Database related flags
    rootCmd.Flags().String("type", "redis", "Database type (redis, postgres)")
    rootCmd.Flags().String("url", "", "Database URL")
    rootCmd.Flags().String("index", "", "Database Index")
    rootCmd.Flags().String("username", "", "Database username (required for PostgreSQL)")
    rootCmd.Flags().String("password", "", "Database password (required for PostgreSQL)")
    rootCmd.Flags().String("name", "", "Database name (only required for PostgreSQL)")
    rootCmd.Flags().String("sslmode", "disable", "Database SSL mode (only for PostgreSQL)")

    // Bind flags to Viper
    viper.BindPFlag("mode", rootCmd.Flags().Lookup("mode"))
    viper.BindPFlag("project_name", rootCmd.Flags().Lookup("project_name"))
    viper.BindPFlag("api_str", rootCmd.Flags().Lookup("api_str"))
    viper.BindPFlag("source_path", rootCmd.Flags().Lookup("source_path"))
    viper.BindPFlag("tika_server_url", rootCmd.Flags().Lookup("tika_server_url"))
    viper.BindPFlag("emb_api_base", rootCmd.Flags().Lookup("emb_api_base"))
    viper.BindPFlag("emb_model_id", rootCmd.Flags().Lookup("emb_model_id"))
    viper.BindPFlag("chunk_size", rootCmd.Flags().Lookup("chunk_size"))
    viper.BindPFlag("chunk_overlap", rootCmd.Flags().Lookup("chunk_overlap"))

    // Bind database flags to Viper
    viper.BindPFlag("database.type", rootCmd.Flags().Lookup("type"))
    viper.BindPFlag("database.url", rootCmd.Flags().Lookup("url"))
    viper.BindPFlag("database.index", rootCmd.Flags().Lookup("index"))
    viper.BindPFlag("database.username", rootCmd.Flags().Lookup("username"))
    viper.BindPFlag("database.password", rootCmd.Flags().Lookup("password"))
    viper.BindPFlag("database.name", rootCmd.Flags().Lookup("name"))
    viper.BindPFlag("database.sslmode", rootCmd.Flags().Lookup("sslmode"))

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

    viper.AutomaticEnv() // read in environment variables that match

    if err := viper.ReadInConfig(); err == nil {
        logger.Log.Infof("Using config file: %s", viper.ConfigFileUsed())
    }

    _, err := config.InitConfig()
    if err != nil {
        logger.Log.Fatalf("Error initializing config: %v", err)
        os.Exit(1)
    }
}
