package cmd

import (
	"azuki774/go-simple-auth-proxy/internal/auth"
	"azuki774/go-simple-auth-proxy/internal/client"
	"azuki774/go-simple-auth-proxy/internal/repository"
	"azuki774/go-simple-auth-proxy/internal/server"
	"context"
	"log/slog"
	"sync"

	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		slog.Info("start called")

		// factory
		srv := server.Server{
			ListenPort:    "8080",
			Client:        &client.Client{ProxyAddr: "http://localhost:8888"},
			Authenticater: &auth.Authenticater{AuthStore: &repository.Store{Mu: &sync.Mutex{}}},
		}

		srv.Start(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
