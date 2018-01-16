package loggers

import ()

type LoggerServer interface {
	New(addr string)
	ListenAndServe(shutdownChan chan bool)
	DoneChan() chan bool
}
