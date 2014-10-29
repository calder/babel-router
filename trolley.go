package main

import "log"
import "strings"
import "time"
import "github.com/calder/babel-encoding-go"
import "github.com/calder/fiddle"
import "gopkg.in/mgo.v2"
import "gopkg.in/mgo.v2/bson"

/*******************
***   Database   ***
*******************/

type Message struct {
    Date    *time.Time
    To      []byte
    Content []byte
}

func storeMessage (msg *Message) {
    db.C("messages").Insert(msg)
}

func getMessages (to *fiddle.Bits) *mgo.Iter {
    return db.C("messages").Find(&bson.M{"to":to.RawBytes()}).Iter()
}

var db *mgo.Database

func init () {
    session, e := mgo.Dial("localhost")
    if e != nil { log.Fatal(e) }
    session.SetMode(mgo.Monotonic, true)
    db = session.DB(conf.Db["name"])

    db.C("messages").EnsureIndex(mgo.Index{
        Key:    []string{"id"},
        Unique: true,
        Sparse: true,
    })
    db.C("messages").EnsureIndex(mgo.Index{Key:[]string{"date"}})
    db.C("messages").EnsureIndex(mgo.Index{Key:[]string{"to"}})
}

/****************
***   Pipes   ***
****************/

type Pipe func(babel.Any)

func newPipe (id string, pipe string) {
    // Parse pipe string
    d := strings.Index(pipe, "://")
    if d == -1 { log.Fatal("Invalid pipe: ", pipe) }
    typ := pipe[:d]
    arg := pipe[d+3:]

    // Create the pipe
    switch typ {
    case "udp":
        if pipes[id] == nil { pipes[id] = make(map[string]Pipe) }
        pipes[id][pipe] = func (pkt babel.Any) {
            log.Println("Sent:", pkt)
            babel.SendUdp(arg, pkt)
        }
    default:
        log.Fatal("Unkown pipe type: ", typ)
    }

    // Send backlogged messages
    msgs := getMessages(fiddle.FromRawHex(id))
    var msg Message
    for msgs.Next(&msg) {
        fun := pipes[id][pipe]
        pkt, e := babel.DecodeBytes(msg.Content)
        if e != nil { log.Println("Warning: ", e); continue }
        fun(pkt)
    }
}

func pipeMessage (to *babel.Id1, dat *fiddle.Bits) {
    id := to.Dat.RawHex()
    for pipe := range pipes[id] {
        fun := pipes[id][pipe]
        fun(dat)
    }
}

var pipes map[string]map[string]Pipe

func init () {
    pipes = make(map[string]map[string]Pipe)
    for id := range conf.Pipes {
        for pipe := range conf.Pipes[id] {
            newPipe(id, conf.Pipes[id][pipe])
        }
    }
}

/********************
***   Receivers   ***
********************/

func init () {
    for i := range conf.Receivers {
        r := conf.Receivers[i]
        d := strings.Index(r, "://")
        if d == -1 { log.Fatal("Invalid receiver: ", r) }
        typ := r[:d]
        arg := r[d+3:]
        switch typ {
        case "udp":
            _, e := babel.ReceiveUdp(arg, 1048576, handle, babel.ErrorLogger)
            if e != nil { log.Fatal(e) }
        default:
            log.Fatal("Unkown receiver type: ", typ)
        }
    }
}

/*************************
***   Packet Handler   ***
*************************/

type Context struct {
    Date    *time.Time
    To      *babel.Id1
    Content *fiddle.Bits
}

func handle (pkt babel.Any) {
    log.Println("Received:", pkt)
    peel(pkt, &Context{})
}

func peel (pkt babel.Any, c *Context) {
    switch pkt.(type) {
    case *babel.UdpSub:
        peelUdpSub(pkt.(*babel.UdpSub), c)
    case *babel.Message:
        peelMessage(pkt.(*babel.Message), c)
    default:
        log.Println("    Discarded: unkown packet type")
    }
}

func peelUdpSub (sub *babel.UdpSub, c *Context) {
    id := sub.Id1.Dat.RawHex()
    addr := sub.Addr.Dat
    newPipe(id, "udp://"+addr)
}

func peelMessage (msg *babel.Message, c *Context) {
    c.To = msg.To
    c.Content = msg.Dat
    send(c)
}

func send (c *Context) {
    msg := &Message{
        Date:    c.Date,
        To:      c.To.Dat.RawBytes(),
        Content: c.Content.Bytes(),
    }
    storeMessage(msg)
    pipeMessage(c.To, c.Content)
}

/**************************
***   Test UDP Sender   ***
**************************/

func init () {
    // for {
    //     e := babel.Send(&babel.Id1{babel.NIL}, &babel.Unicode{"Ohai world!"})
    //     if e != nil { log.Println("Warning: ", e) }
    //     time.Sleep(time.Second)
    // }
}

/***************
***   Main   ***
***************/

func init () {
    loadConfig()
}

func main () {
    select{}
}