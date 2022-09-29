package cli

import (
	"github.com/Improwised/xkcd-api/config"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// Init app initialization
func Init(cfg config.AppConfig, logger *zap.Logger) error {
	apiCmd := GetAPICommandDef(cfg, logger)

	rootCmd := &cobra.Command{Use: "xkcd-api"}
	rootCmd.AddCommand(&apiCmd)
	return rootCmd.Execute()
}
