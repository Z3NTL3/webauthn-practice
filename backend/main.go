package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"charm.land/lipgloss/v2"
	"charm.land/log/v2"
	"github.com/spf13/viper"
	_ "github.com/z3ntl3/webauthn-practice/globals"
	"github.com/z3ntl3/webauthn-practice/handlers"
	"github.com/z3ntl3/webauthn-practice/middleware"
)

func init() {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.Kitchen,
		Prefix:          "WebAuthn Demo",
		Level:           log.DebugLevel,
	})

	styles := log.DefaultStyles()
	styles.Levels[log.ErrorLevel] = lipgloss.NewStyle().
		SetString("ERROR").
		Padding(0, 1, 0, 1).
		Background(lipgloss.Color("255")).
		Foreground(lipgloss.Color("0"))

	// more on middleware/log.go
	log.SetDefault(logger)

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("config", "load", err)
	}

	for _, key := range viper.AllKeys() {
		viper.Set(key, os.Expand(viper.GetString(key), func(variable string) string {
			return viper.GetString(variable)
		}))
	}
}

func main() {
	frontend_url := viper.GetString("FRONTEND_URL")
	backend_url := viper.GetString("BACKEND_URL")

	log.Info("Server Startup", "server", "bootstrap-phase")
	mux := http.NewServeMux()

	mux.HandleFunc("POST /register/create", handlers.RegisterCreateOptions)
	mux.HandleFunc("POST /register/finish", handlers.RegisterFinalize)

	mux.HandleFunc("POST /login/begin", handlers.LoginBegin)
	mux.HandleFunc("POST /login/finalize", handlers.LoginFinalize)

	log.Info("Ready", "URL", backend_url)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("BACKEND_PORT")), middleware.Logger(middleware.CORS(frontend_url, mux))))
}
