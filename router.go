package main

// import "log"
// import "strings"
// import "time"
// import "github.com/calder/babel-lib-go"

type Router struct {
    Queue chan []byte
}

func NewRouter () *Router {
    return &Router{make(chan []byte)}
}