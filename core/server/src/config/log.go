package config

import (
	"goland/server/src/controller"
	"io"
	"log"
	"os"
)

func NewGolandLoggers(out io.Writer) *controller.GolandLoggers {
	return &controller.GolandLoggers{
		Info:   log.New(out, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		Warn:   log.New(out, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		Error:  log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		System: log.New(out, "SYSTEM: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

var StaticLoggers = NewGolandLoggers(os.Stdout)
