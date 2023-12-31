package conn

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
)

var CliCommands = []*cli.Command{
	{
		Name:        "list",
		Aliases:     []string{"ls"},
		Description: "List all connected remote server",
		Action: func(ctx *cli.Context) error {
			var servers []CliConnection
			raw, _ := json.Marshal(viper.Get("servers"))
			_ = json.Unmarshal(raw, &servers)

			log.Info().Msgf("There are %d server(s) connected in total.", len(servers))
			for idx, server := range servers {
				log.Info().Msgf("%d) %s: %s", idx+1, server.ID, server.Url)
			}

			return nil
		},
	},
	{
		Name:        "connect",
		Aliases:     []string{"add"},
		Description: "Connect and save configuration of remote server",
		ArgsUsage:   "<id> <server url> <credential>",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() < 3 {
				return fmt.Errorf("must have three arguments: <id> <server url> <credential>")
			}

			c := CliConnection{
				ID:         ctx.Args().Get(0),
				Url:        ctx.Args().Get(1),
				Credential: ctx.Args().Get(2),
			}

			if err := c.CheckConnectivity(); err != nil {
				return fmt.Errorf("couldn't connect server: %s", err.Error())
			} else {
				var servers []CliConnection
				raw, _ := json.Marshal(viper.Get("servers"))
				_ = json.Unmarshal(raw, &servers)
				viper.Set("servers", append(servers, c))

				if err := viper.WriteConfig(); err != nil {
					return err
				} else {
					log.Info().Msg("Successfully connected a new remote server, enter \"rds ls\" to get more info.")
					return nil
				}
			}
		},
	},
	{
		Name:        "disconnect",
		Aliases:     []string{"remove"},
		Description: "Remove a remote server configuration",
		ArgsUsage:   "<id>",
		Action: func(ctx *cli.Context) error {
			if ctx.Args().Len() < 1 {
				return fmt.Errorf("must have more one arguments: <server url>")
			}

			var servers []CliConnection
			raw, _ := json.Marshal(viper.Get("servers"))
			_ = json.Unmarshal(raw, &servers)
			viper.Set("servers", lo.Filter(servers, func(item CliConnection, idx int) bool {
				return item.ID != ctx.Args().Get(0)
			}))

			if err := viper.WriteConfig(); err != nil {
				return err
			} else {
				log.Info().Msg("Successfully disconnected a remote server, enter \"rds ls\" to get more info.")
				return nil
			}
		},
	},
}
