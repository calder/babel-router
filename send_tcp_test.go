package main

import "net"
import "testing"

func TestSendTcp (T *testing.T) {
    log.Info("starting test")

    router := NewRouter()
    go router.sendTcp(8125)

    // Create TCP connection
    _, e := net.Dial("tcp", "localhost:8125")
    if e != nil { panic(e) }

    // TODO: subscribe and stuff
}
