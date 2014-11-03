package main

import "log"
import "net"

const maxPacketSize = 1024 * 1024

func (r *Router) ListenUdp (addrStr string) {
    addr, e := net.ResolveUDPAddr("udp", addrStr)
    if e != nil { panic(e) }
    conn, e := net.ListenUDP("udp", addr)
    if e != nil { panic(e) }

    buffer := make([]byte, maxPacketSize)
    for {
        n, e := conn.Read(buffer)
        if e != nil { log.Println("UDP read error:", e.Error()); continue }

        msg := make([]byte, n)
        copy(msg, buffer[:n])
        r.Queue <- msg
    }
}