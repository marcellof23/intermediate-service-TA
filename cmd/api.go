package cmd

import (
	"log"

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

			cfg, err := boot.LoadConfig(configfile)
			if err != nil {
				log.Fatal(err)
			}

			dep, err := boot.InitDependencies(cfg)
			if err != nil {
				log.Fatal(err)
			}
			router := api.InitRoutes(dep)

			err = router.Run(cfg.Server.Addr)
			if err != nil {
				log.Fatal(err)
			}

		},
	}

	rootCmd.AddCommand(apiCmd)
}
