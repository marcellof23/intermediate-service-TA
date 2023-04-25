package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/intermediate-service-ta/api"
	"github.com/intermediate-service-ta/boot"
)

func init() {
	var apiCmd = &cobra.Command{
		Use:   "api",
		Short: "Run an API server",
		Run: func(cmd *cobra.Command, args []string) {
			configfile := file
			if len(args) != 0 {
				configfile = args[0]
			}

			os.Mkdir("backup", os.ModePerm)

			cfg, err := boot.LoadConfig(configfile)
			if err != nil {
				log.Fatal(err)
			}

			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			dep, err := boot.InitDependencies(cfg)
			if err != nil {
				log.Fatal(err)
			}
			router := api.InitRoutes(dep, sigchan)

			go router.Run(cfg.Server.Addr)
			<-sigchan

			fmt.Println("Shutting down...")
		},
	}

	rootCmd.AddCommand(apiCmd)
}
