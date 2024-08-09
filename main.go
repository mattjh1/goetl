package main

import (
    "github.com/spf13/viper"
    "github.com/mattjh1/goetl/cmd"
    "github.com/mattjh1/goetl/config/logger"
)

func main() {
    // Initialize the logger
    mode := viper.GetString("mode")
    logger.InitLogger(mode)

    // Execute the root command
    cmd.Execute()
}
