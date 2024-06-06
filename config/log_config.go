package config

import (
	"io"
	"log"
	"os"
)

//nolint:gochecknoglobals
var Debug = false

//nolint:gochecknoinits
func init() {
	if !Debug {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}
}
