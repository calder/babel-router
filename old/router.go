package main

import "log"
import "strings"
import "time"
import "github.com/calder/babel-lib-go"

type Context struct {
    Time    *time.Time
    To      *babel.Id1
    Content *fiddle.Bits
}

func handle (pkt babel.Any) {
    log.Println("Received:", pkt)
    handle(pkt, &Context{})
}

func handle (pkt babel.Any, c *Context) {
    switch pkt.(type) {
    case *babel.UdpSub:
        handleUdpSub(pkt.(*babel.UdpSub), c)
    case *babel.Message:
        handleMessage(pkt.(*babel.Message), c)
    default:
        log.Println("    Discarded: unkown packet type")
    }
}

func handleUdpSub (sub *babel.UdpSub, c *Context) {
    id := sub.Id1.Dat.RawHex()
    addr := sub.Addr.Dat
    newPipe(id, "udp://"+addr)
}