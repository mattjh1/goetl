package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/mattjh1/goetl/config"
	"github.com/mattjh1/goetl/pkg/extract"
	"github.com/mattjh1/goetl/pkg/load"
	"github.com/mattjh1/goetl/pkg/transform"
)

var etlCmd = &cobra.Command{
	Use:   "etl",
	Short: "Run the ETL process",
	Long:  `Extract, Transform, and Load data.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.InitConfig()
		if err != nil {
			fmt.Println("Error loading config:", err)
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
	transformedCh := make(chan string)
	// TODO replace with cfg data
	path := cfg.SourcePath
	globPattern := "**/*.*"
	since := time.Date(1970, 8, 1, 0, 0, 0, 0, time.UTC) // Example date

	wg.Add(1)
	go func() {
		defer wg.Done()
		extract.Extract(dataCh, cfg, path, globPattern, since)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		transform.Transform(dataCh, transformedCh, cfg)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		load.Load(transformedCh, cfg)
	}()

	wg.Wait()
	fmt.Println("ETL process completed.")
}
