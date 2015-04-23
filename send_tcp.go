package main

import "encoding/hex"
import "net"
import "strconv"
import "time"
import "github.com/calder/babel-lib-go"

func handleConnection (conn net.Conn) {
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))

    typ, data, e := babel.ReadMessage(conn)
    if e != nil { log.Debug(e.Error()); return }

    if typ == 0 {
        log.Debug(hex.EncodeToString(data))
    } else {
        log.Debug("unrecognized message type:", typ)
    }
}

func (r *Router) sendTcp (port int) {
    ln, e := net.Listen("tcp", ":" + strconv.Itoa(port))
    if e != nil { panic(e) }

    for {
        conn, e := ln.Accept()
        if e != nil { panic(e) }
        go handleConnection(conn)
    }
}
