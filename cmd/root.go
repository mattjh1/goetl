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
    rootCmd.Flags().String("database_type", "redis", "Database type (redis, postgres)")
    rootCmd.Flags().String("database_url", "", "Database URL")
    rootCmd.Flags().String("database_index", "", "Database Index")
    rootCmd.Flags().String("database_username", "", "Database username (required for PostgreSQL)")
    rootCmd.Flags().String("database_password", "", "Database password (required for PostgreSQL)")
    rootCmd.Flags().String("database_name", "", "Database name (only required for PostgreSQL)")
    rootCmd.Flags().String("database_sslmode", "disable", "Database SSL mode (only for PostgreSQL)")


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
    viper.BindPFlag("database.database_type", rootCmd.Flags().Lookup("database_type"))
    viper.BindPFlag("database.database_url", rootCmd.Flags().Lookup("database_url"))
    viper.BindPFlag("database.database_index", rootCmd.Flags().Lookup("database_index"))
    viper.BindPFlag("database.database_username", rootCmd.Flags().Lookup("database_username"))
    viper.BindPFlag("database.database_password", rootCmd.Flags().Lookup("database_password"))
    viper.BindPFlag("database.database_name", rootCmd.Flags().Lookup("database_name"))
    viper.BindPFlag("database.database_sslmode", rootCmd.Flags().Lookup("database_sslmode"))

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
