package main

import "testing"
import "time"
import "github.com/calder/babel-lib-go"
import "github.com/calder/go-timeout"

func TestProcessQueueForwardsMessage (T *testing.T) {
    log.Info("starting test")

    router := NewRouter()
    // TODO: attach distributer to the router
    router.queue <- babel.Wrap(babel.TYPE, BLOB, []byte{1,2,3,4,5})

    select {
    // TODO: get message from the test distributer
    case <-timeout.Timeout(100 * time.Millisecond):
    }
}
