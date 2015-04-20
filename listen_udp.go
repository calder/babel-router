package main

import "net"
import "strconv"

const maxPacketSize = 1024 * 1024

// Listen for UDP packets on the given port and add them to the router's queue.
func (r *Router) listenUdp (port int) {
    addr, e := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(port))
    if e != nil { panic(e) }
    conn, e := net.ListenUDP("udp", addr)
    if e != nil { panic(e) }

    buffer := make([]byte, maxPacketSize)
    for {
        n, e := conn.Read(buffer)
        if e != nil { log.Debug("UDP read error:", e.Error()); continue }

        msg := make([]byte, n)
        copy(msg, buffer[:n])

        r.queue <- msg
    }
}
