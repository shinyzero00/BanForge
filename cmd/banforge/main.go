package main

import (
	"fmt"
	"os"

	"github.com/d3m0k1d/BanForge/internal/config"
	_ "github.com/d3m0k1d/BanForge/internal/judge"
	"github.com/d3m0k1d/BanForge/internal/logger"
	"github.com/d3m0k1d/BanForge/internal/parser"
	_ "github.com/d3m0k1d/BanForge/internal/storage"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "banforge",
	Short: "IPS log-based written on Golang",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize BanForge",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing BanForge...")
		err := os.Mkdir("/var/log/banforge", 0750)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = os.Mkdir("/etc/banforge", 0750)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = config.CreateConf()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		err = config.FindFirewall()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run BanForge daemon process",
	Run: func(cmd *cobra.Command, args []string) {
		log := logger.New(false)
		log.Info("Starting BanForge daemon")
		//db, err := storage.NewDB()
		//if err != nil {
		//log.Error("Failed to create database", "error", err)
		//}
		cfg, err := config.LoadConfig()
		if err != nil {
			log.Error("Failed to load config", "error", err)
		}
		for service := range cfg.Service {
			if cfg.Service[service].Enabled && cfg.Service[service].Name != "nginx" {
				pars, err := parser.NewScanner(cfg.Service[service].LogPath)
				if err != nil {
					log.Error("Failed to create parser", "error", err)
				}
				go pars.Start()
			}
		}
	},
}

func Init() {

}

func Execute() {
	rootCmd.AddCommand(initCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
