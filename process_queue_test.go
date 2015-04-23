package main

import "testing"
import "time"
import "github.com/calder/babel-lib-go"
import "github.com/calder/go-timeout"

// Enqueue three messages for processing:
//     1. A raw Blob(...)
//     2. A Envelope(badDest, Blob(...))
//     3. A Envelope(dest, Blob(...))
// and attach a fake recipient for <dest>.
//
// The recipient should receive the third, but not the first or second.
func TestProcessQueueForwardsMessage (T *testing.T) {
    log.Info("starting test")

    // Message data
    dest    := babel.Hash1OfData([]byte{1,2,3})
    badDest := babel.Hash1OfData([]byte{2,4,6})
    blob    := NewBlob([]byte{4,5,6})
    badBlob := NewBlob([]byte{8,10,12})

    // Attach a fake recipient
    router := NewRouter()
    recipient := router.openSendQueue(dest.CBR())

    // Enqueue messages
    router.enqueue(badBlob.CBR())
    router.enqueue(EnvelopeFromValues(badDest, badBlob).CBR())
    router.enqueue(EnvelopeFromValues(dest, blob).CBR())

    select {
    case received := <-recipient:
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
