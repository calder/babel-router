package main

import "encoding/json"
import "io/ioutil"
import "log"
import "os"
import "strings"
import "time"
import "github.com/calder/babel"
import "github.com/calder/fiddle"
import "labix.org/v2/mgo"
import "labix.org/v2/mgo/bson"

/*****************
***   Config   ***
*****************/

type Config struct {
    Ids         map[string]map[string]string
    Db          map[string]string
    Pipes       map[string][]string
    PipeServers []string
    Receivers   []string
}

var conf *Config

func init () {
    // Check arguments
    if len(os.Args) != 2 { log.Fatal("Usage: trolley CONFIGFILE") }

    // Read config file
    f, e := ioutil.ReadFile(os.Args[1])
    if e != nil { log.Fatal("Error reading config file: ", e) }

    // Parse config file
    e = json.Unmarshal(f, &conf)
    if e != nil { log.Fatal("Configuration error: ", e) }
}

/*******************
***   Database   ***
*******************/

type Message struct {
    Date    time.Time
    To      []byte
    Content []byte
}

func storeMessage (to *fiddle.Bits, dat *fiddle.Bits) {
    msg := &Message{time.Now(), to.RawBytes(), dat.Bytes()}
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
        pipes[id][pipe] = func (pkt babel.Any) { babel.SendUdp(arg, pkt) }
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

func pipeMessage (to *babel.Id, msg babel.Any) {
    id := to.Dat.RawHex()
    for pipe := range pipes[id] {
        fun := pipes[id][pipe]
        fun(msg)
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

func handle (pkt babel.Any) {
    switch pkt.(type) {
    case *babel.UdpSub:
        handleUdpSub(pkt.(*babel.UdpSub))
    case *babel.Message:
        handleMessage(pkt.(*babel.Message))
    default:
        log.Println("    Discarded: unkown packet type")
    }
}

func handleUdpSub (sub *babel.UdpSub) {
    id := sub.Id.Dat.RawHex()
    addr := sub.Addr.Dat
    newPipe(id, "udp://"+addr)
}

func handleMessage (msg *babel.Message) {
    storeMessage(msg.To.Dat, babel.EncodeUnsafe(msg.Dat))
    pipeMessage(msg.To, msg.Dat)
}

/**************************
***   Test UDP Sender   ***
**************************/

func init () {
    // for {
    //     e := babel.Send(&babel.Id{babel.NIL}, &babel.Unicode{"Ohai world!"})
    //     if e != nil { log.Println("Warning: ", e) }
    //     time.Sleep(time.Second)
    // }
}

/***************
***   Main   ***
***************/

func main () {
    select{}
}