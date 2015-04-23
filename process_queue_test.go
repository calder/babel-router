package main

import "testing"
import "time"
import "github.com/calder/babel-lib-go"
import "github.com/calder/go-timeout"

func TestProcessQueueForwardsMessage (T *testing.T) {
    log.Info("starting test")

    router := NewRouter()
    dest := babel.Hash1OfData([]byte{1,2,3,4,5})
    blob := NewBlob([]byte{6,7,8,9,10})

    // Attach a fake recipient
    recipient := router.openSendQueue(dest.CBR())

    // Enqueue a message to the fake recipient
    router.enqueue(EnvelopeFromValues(dest, blob).CBR())

    select {
    case received := <-recipient:
        // typ, data, _, e := babel.Unwrap(babel.TYPE, received)
        // if e != nil {
        //     T.Log("Error:", "error decoding message")
        //     T.Log("Receieved:", received)
        //     T.Log("Original: ", blob.CBR())
        //     T.FailNow()
        // }
        // if typ != BLOB {
        //     T.Log("Error:", "received message type != BLOB")
        //     T.Log("Type: ", THex(typ))
        //     T.Log("BLOB: ", THex(BLOB))
        //     T.FailNow()
        // }
        blob2 := DecodeBlob(received)
        if !blob2.Equal(blob) {
            T.Log("Error:    ", "received blob != sent blob")
            T.Log("Receieved:", Hex(blob2.CBR()))
            T.Log("Original: ", Hex(blob.CBR()))
            T.FailNow()
        }
    case <-timeout.Timeout(100 * time.Millisecond):
        T.Log("Error:", "timed out waiting for message.")
        T.FailNow()
    }
}
