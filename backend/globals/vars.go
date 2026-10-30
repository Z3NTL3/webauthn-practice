package globals

import (
	"sync"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/spf13/viper"
)

var WEBAUTH *webauthn.WebAuthn

// Concurrency safe in memory store.
// Remember, this is a simple demo, therefore no database is used.
// Data will not be persisted between restarts.
var STORAGE = new(sync.Map)

func init() {
	config := &webauthn.Config{
		RPDisplayName: "Localhost",
		RPID:          "localhost",
		RPOrigins:     []string{viper.GetString("FRONTEND_PORT")},
	}

	webauth, err := webauthn.New(config)
	if err != nil {
		panic(err)
	}

	WEBAUTH = webauth
}
