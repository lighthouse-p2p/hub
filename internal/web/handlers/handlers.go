package handlers

import (
	"github.com/lighthouse-p2p/hub/internal/config"
)

// Handlers holds all the handlers used for fiber HTTP
type Handlers struct {
	Cfg *config.Config
}
