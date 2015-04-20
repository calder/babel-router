package main

import "os"
import "github.com/op/go-logging"

var log = logging.MustGetLogger("index")

var logFormat = logging.MustStringFormatter(
    "%{color}%{time:15:04:05.000000} %{shortfunc} ▶ %{message}%{color:reset}")

var logBackend = logging.NewLogBackend(os.Stderr, "", 0)
var logFormattedBackend = logging.NewBackendFormatter(logBackend, logFormat)
var logLeveledBackend = logging.AddModuleLevel(logFormattedBackend)

func init () {
    logLeveledBackend.SetLevel(logging.DEBUG, "")
    logging.SetBackend(logLeveledBackend)
}
