package main

import "encoding/json"
import "io/ioutil"
import "log"
import "os"

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

func loadConfig () {
    // Check arguments
    if len(os.Args) != 2 { log.Fatal("Usage: trolley CONFIGFILE") }

    // Read config file
    f, e := ioutil.ReadFile(os.Args[1])
    if e != nil { log.Fatal("Error reading config file: ", e) }

    // Parse config file
    e = json.Unmarshal(f, &conf)
    if e != nil { log.Fatal("Configuration error: ", e) }
}