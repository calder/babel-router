package main

import "bytes"
import "encoding/hex"
import "errors"
import "net"
import "testing"
import "time"

func Timeout (timeout time.Duration) chan bool {
    channel := make(chan bool)
    go func () {
        time.Sleep(timeout)
        channel <- true
    }()
    return channel
}

func sendUdp (addrStr string, msg []byte) error {
    // Create UDP connection
    addr, e := net.ResolveUDPAddr("udp", addrStr)
    if e != nil { return e }
    conn, e := net.DialUDP("udp", nil, addr)
    if e != nil { return e }

    // Send the packet
    n, e := conn.Write(msg)
    if e != nil { return e }
    if n < len(msg) { return errors.New("incomplete send") }
    return nil
}

func TestListenUdp (T *testing.T) {
    router := NewRouter()
    go router.ListenUdp(":8124")

    msg := []byte{1,2,3,4,5}
    sendUdp("localhost:8124", msg)

    select {
    case received := <- router.Queue:
        if !bytes.Equal(received, msg) {
            T.Log("Error:   ", "received message != sent message")
            T.Log("Sent:    ", hex.EncodeToString(msg))
            T.Log("Received:", hex.EncodeToString(received))
            T.FailNow()
        }
    case <-Timeout(100 * time.Millisecond):
        T.Log("Error:", "timed out waiting for message")
        T.FailNow()
    }
}