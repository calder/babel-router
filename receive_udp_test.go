package main

import "bytes"
import "encoding/hex"
import "errors"
import "net"
import "testing"
import "time"
import "github.com/calder/go-timeout"

func sendUdp (addrStr string, msg []byte) error {
    // Create UDP connection
    addr, e := net.ResolveUDPAddr("udp", addrStr)
    if e != nil { return e }
    conn, e := net.DialUDP("udp", nil, addr)
    if e != nil { return e }

    // Send message
    n, e := conn.Write(msg)
    if e != nil { return e }
    if n < len(msg) { return errors.New("incomplete send") }
    return nil
}

func TestListenUdp (T *testing.T) {
    log.Info("starting test")

    router := newTestRouter()
    go router.receiveUdp(8124)

    msg := []byte{1,2,3,4,5}
    sendUdp("localhost:8124", msg)

    select {
    case received := <- router.queue:
        if !bytes.Equal(received, msg) {
            T.Log("Error:   ", "received message != sent message")
            T.Log("Sent:    ", hex.EncodeToString(msg))
            T.Log("Received:", hex.EncodeToString(received))
            T.FailNow()
        }
    case <-timeout.Timeout(100 * time.Millisecond):
        T.Log("Error:", "timed out waiting for message")
        T.FailNow()
    }
}
