package main

type Router struct {
    queue chan []byte
}

func NewRouter () *Router {
    router := &Router{make(chan []byte)}
    go router.processQueue()
    return router
}

func newTestRouter () *Router {
    return &Router{make(chan []byte)}
}
