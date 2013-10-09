package main

import "encoding/json"
import "io/ioutil"
import "log"
import "os"
import "strings"
import "time"
import "github.com/calder/babel"
import "github.com/calder/fiddle"
import "github.com/calder/vita"
import "labix.org/v2/mgo"

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

/******************
***   Decoder   ***
******************/

var dec *babel.Decoder

func init () {
    dec = babel.NewDecoder()
}

/*************************
***   Packet Handler   ***
*************************/

func handle (bin babel.Bin) {
    switch bin.(type) {
    case *babel.MsgBin:
        handleMsg(bin.(*babel.MsgBin))
    default:
        log.Println("    Discarded: unkown packet type.")
    }
}

func handleMsg (msg *babel.MsgBin) {
    switch msg.To.(type) {
    case *babel.IdBin:
        storeMessage(msg.To.(*babel.IdBin).Dat, msg.Dat.Encode())
    default:
        log.Println("    Discarded: unkown recipient type.")
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
            _, e := vita.ReceiveUdp(arg, 1048576, dec, handle, vita.ErrorLogger)
            if e != nil { log.Fatal(e) }
        default:
            log.Fatal("Unkown receiver type: ", typ)
        }
    }
}

/**************************
***   Test UDP Sender   ***
**************************/

func init () {
    for {
        e := vita.Send(&babel.IdBin{babel.NIL}, &babel.UnicodeBin{"Ohai world!"})
        if e != nil { log.Println("Warning: ", e) }
        time.Sleep(time.Second)
    }
}

/***************
***   Main   ***
***************/

func main () {
    select{}
}