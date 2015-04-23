package main

import "encoding/hex"
import "os"
import "strings"
import "github.com/op/go-logging"

func Hex (bytes []byte) string {
    return strings.ToUpper(hex.EncodeToString(bytes))
}

func SHex (s string) string {
    return Hex([]byte(s))
}

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
