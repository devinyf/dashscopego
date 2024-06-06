package config

import (
	"io"
	"log"
	"os"
)

var Debug = false

func init() {
	if !Debug {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stderr)
	}
}
