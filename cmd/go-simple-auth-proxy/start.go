package cmd

import (
	"azuki774/go-simple-auth-proxy/internal/auth"
	"azuki774/go-simple-auth-proxy/internal/client"
	"azuki774/go-simple-auth-proxy/internal/repository"
	"azuki774/go-simple-auth-proxy/internal/server"
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
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

		if err := configLoad(); err != nil {
			slog.Error("config error", "err", err)
			os.Exit(1)
		}
		slog.Info("config loaded")

		basicAuthLoad()
		slog.Info("basic auth loaded")

		// factory
		srv := server.Server{
			ListenPort:    startConfig.Port,
			Client:        &client.Client{ProxyAddr: startConfig.ProxyAddress},
			Authenticater: &auth.Authenticater{AuthStore: &repository.StoreInMemory{Mu: &sync.Mutex{}, BasicAuthStore: basicAuthMap}},
			ExporterPort:  startConfig.ExporterPort,
		}

		// ready check
		err := srv.CheckReadiness()
		if err != nil {
			os.Exit(1)
		}
		slog.Info("proxy ready ok")

		srv.Start(context.Background())
	},
}

type StartConfig struct {
	Version        int      `toml:"conf-version"`
	Port           string   `toml:"server_port"`
	BasicAuthList  []string `toml:"basicauth"`
	ProxyAddress   string   `toml:"proxy_address"`
	CookieLifeTime int      `toml:"cookie_lifetime"`

	ExporterPort string `toml:"exporter_port"`
}

var startConfig StartConfig
var startConfigDir string
var basicAuthMap map[string]string

func configLoad() (err error) {
	_, err = toml.DecodeFile(startConfigDir, &startConfig)
	if err != nil {
		return err
	}
	return nil
}

func basicAuthLoad() {
	basicAuthMap = make(map[string]string)
	for _, v := range startConfig.BasicAuthList {
		userpass := strings.Split(v, ":")
		basicAuthMap[userpass[0]] = userpass[1]
	}
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
	startCmd.Flags().StringVarP(&startConfigDir, "config", "c", "config.toml", "config directory")
}
