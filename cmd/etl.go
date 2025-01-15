package cmd

import (
	"sync"
	"time"

	"github.com/mattjh1/goetl/config"
	"github.com/mattjh1/goetl/config/logger"
	"github.com/mattjh1/goetl/pkg/extract"
	"github.com/mattjh1/goetl/pkg/load"
	"github.com/mattjh1/goetl/pkg/utils"
	"github.com/mattjh1/goetl/pkg/transform"
	"github.com/spf13/cobra"
	"github.com/schollz/progressbar/v3"
	"github.com/tmc/langchaingo/schema"
)

var etlCmd = &cobra.Command{
	Use:   "etl",
	Short: "Run the ETL process",
	Long:  `Extract, Transform, and Load data.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.InitConfig()
		if err != nil {
			logger.Log.Errorf("Error loading config: %e", err)
			return
		}
		runETL(cfg)
	},
}

func init() {
	rootCmd.AddCommand(etlCmd)
}

func runETL(cfg *config.Config) {
	var wg sync.WaitGroup

	dataCh := make(chan string)
	transformedCh := make(chan schema.Document)
	path := cfg.SourcePath
	globPattern := cfg.GlobPattern
	since := time.Date(1970, 8, 1, 0, 0, 0, 0, time.UTC)

	logger.Log.Info("Starting ETL process...")

	totalFiles, err := utils.CountFilesInPath(path, globPattern)
	if err != nil {
		logger.Log.Errorf("Error counting files: %v", err)
		return
	}

	// Create progress bars
	extractBar := progressbar.Default(int64(totalFiles), "Extracting files")
	transformBar := progressbar.Default(-1, "Transforming data") 

	// Extraction
	wg.Add(1)
	go func() {
		defer wg.Done()
		extract.Extract(dataCh, cfg, path, globPattern, since, extractBar)
		logger.Log.Info("File extraction completed.")
	}()

	// Transformation
	wg.Add(1)
	go func() {
		defer wg.Done()
		transform.Transform(dataCh, transformedCh, cfg, transformBar)
		logger.Log.Info("Data transformation completed.")
	}()

	// Loading
	wg.Add(1)
	go func() {
		defer wg.Done()
			load.Load(transformedCh, cfg)
		logger.Log.Info("Data loading completed.")
	}()

	wg.Wait()
	logger.Log.Info("ETL process completed.")
}
