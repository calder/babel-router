package main

import "testing"
import "time"
import "github.com/calder/babel-lib-go"

func TestProcessQueueDecodesMessage (T *testing.T) {
    testDecoderCalled := make(chan bool)
    testType := babel.Tag("CDE9C47D")
    testDecoder := func (bytes []byte) (res babel.Any, err error) {
        testDecoderCalled <- true
        return babel.NewInt32(123), nil
    }
    babel.AddType(testType, testDecoder)

    router := NewRouter()
    router.queue <- testType

    select {
    case <-testDecoderCalled:
    case <-Timeout(100 * time.Millisecond):
        T.Log("Error:", "custom type decoder never called")
        T.FailNow()
    }
}