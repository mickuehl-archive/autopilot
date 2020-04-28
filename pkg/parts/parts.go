package parts

import (
	log "github.com/majordomusio/log15"
)

var (
	logger log.Logger
)

func init() {
	logger = log.New("module", "parts")
}
