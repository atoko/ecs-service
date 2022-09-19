package controller

import "log"

type GolandLoggers struct {
	Warn   *log.Logger
	Info   *log.Logger
	Error  *log.Logger
	System *log.Logger
}
