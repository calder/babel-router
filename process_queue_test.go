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
    router.enqueue(EnvelopeFromValues(dest, blob).CBR())

    select {
    // TODO: get message from the test distributer
    case <-timeout.Timeout(100 * time.Millisecond):
    }
}
