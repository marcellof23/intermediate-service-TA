package cmd

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
			boot.Clients = cfg.Clients

			sigchan := make(chan os.Signal, 1)
			signal.Notify(sigchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

			dep, err := boot.InitDependencies(cfg)
			if err != nil {
				log.Fatal(err)
			}
			router := api.InitRoutes(dep, sigchan)

			srv := &http.Server{
				Addr:    cfg.Server.Addr,
				Handler: router,
			}
			go srv.ListenAndServe()

			<-sigchan

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := srv.Shutdown(ctx); err != nil {
				log.Fatal("Server Shutdown:", err)
			}
			log.Println("Server exiting")
		},
	}

	rootCmd.AddCommand(apiCmd)
}
