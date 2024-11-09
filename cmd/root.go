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
		// Set persistent flags for configurations shared across all subcommands
		rootCmd.PersistentFlags().String("mode", "", "Mode of operation (e.g., development, production)")
		rootCmd.PersistentFlags().StringP("project_name", "p", "", "Project name")
		rootCmd.PersistentFlags().StringP("source_path", "s", "", "Source path string")
		rootCmd.PersistentFlags().StringP("glob_pattern", "g", "", "File pattern (default is *.*)")
		rootCmd.PersistentFlags().String("tika_server_url", "", "Apache Tika server URL")
		rootCmd.PersistentFlags().String("emb_api_base", "", "Embedding API base URL")
		rootCmd.PersistentFlags().String("emb_model_id", "", "Embedding model ID")
		rootCmd.PersistentFlags().String("chunk_size", "", "Text split chunk size")
		rootCmd.PersistentFlags().String("chunk_overlap", "", "Text split chunk overlap")

		// Database-related persistent flags
		rootCmd.PersistentFlags().StringP("type", "t", "redis", "Database type (redis, postgres)")
		rootCmd.PersistentFlags().StringP("url", "u", "", "Database URL")
		rootCmd.PersistentFlags().StringP("index", "i", "", "Database Index")
		rootCmd.PersistentFlags().StringP("username", "U", "", "Database username (required for PostgreSQL)")
		rootCmd.PersistentFlags().StringP("password", "P", "", "Database password (required for PostgreSQL)")
		rootCmd.PersistentFlags().String("name", "", "Database name (only required for PostgreSQL)")
		rootCmd.PersistentFlags().String("sslmode", "disable", "Database SSL mode (only for PostgreSQL)")

		// Bind persistent flags to Viper keys
		viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode"))
		viper.BindPFlag("project_name", rootCmd.PersistentFlags().Lookup("project_name"))
		viper.BindPFlag("source_path", rootCmd.PersistentFlags().Lookup("source_path"))
		viper.BindPFlag("glob_pattern", rootCmd.PersistentFlags().Lookup("glob_pattern"))
		viper.BindPFlag("tika_server_url", rootCmd.PersistentFlags().Lookup("tika_server_url"))
		viper.BindPFlag("emb_api_base", rootCmd.PersistentFlags().Lookup("emb_api_base"))
		viper.BindPFlag("emb_model_id", rootCmd.PersistentFlags().Lookup("emb_model_id"))
		viper.BindPFlag("chunk_size", rootCmd.PersistentFlags().Lookup("chunk_size"))
		viper.BindPFlag("chunk_overlap", rootCmd.PersistentFlags().Lookup("chunk_overlap"))

		// Bind database-related persistent flags to Viper keys
		viper.BindPFlag("database.type", rootCmd.PersistentFlags().Lookup("type"))
		viper.BindPFlag("database.url", rootCmd.PersistentFlags().Lookup("url"))
		viper.BindPFlag("database.index", rootCmd.PersistentFlags().Lookup("index"))
		viper.BindPFlag("database.username", rootCmd.PersistentFlags().Lookup("username"))
		viper.BindPFlag("database.password", rootCmd.PersistentFlags().Lookup("password"))
		viper.BindPFlag("database.name", rootCmd.PersistentFlags().Lookup("name"))
		viper.BindPFlag("database.sslmode", rootCmd.PersistentFlags().Lookup("sslmode"))

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
