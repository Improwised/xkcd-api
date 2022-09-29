package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/Improwised/xkcd-api/config"
	"github.com/Improwised/xkcd-api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// GetAPICommandDef runs app
func GetAPICommandDef(cfg config.AppConfig, logger *zap.Logger) cobra.Command {
	apiCommand := cobra.Command{
		Use:   "api",
		Short: "To start api",
		Long:  `To start api`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create fiber app
			app := fiber.New(fiber.Config{})

			// setup routes
			err := routes.Setup(app)
			if err != nil {
				return err
			}

			// Call when SIGINT or SIGTERM received
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, os.Kill)
			go func() {
				_ = <-c
				fmt.Println("Gracefully shutting down...")
				app.Shutdown() /// Stop to accept new connections
			}()

			// Listen on port 3000
			if err := app.Listen(cfg.Port); err != nil {
				log.Panic(err)
			}

			return nil

		},
	}

	return apiCommand
}
