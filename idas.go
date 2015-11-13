package main

import (
	"net/http"
	"log"
	"fmt"
	"sync/atomic"
	"io/ioutil"
	"strconv"
	"time"
	"flag"
	"github.com/satori/go.uuid"
)

const addr = ":1718"

var filename = flag.String("f", "last.txt", "filename to get/set counter")
var memOnly = flag.Bool("m", false, "ignore file, use memory only")

func main() {
	flag.Parse()

	var id uint64

	if !*memOnly {
		id = readcounter()

		ticker := time.NewTicker(time.Second * 1)
		go func() {
			last := id
			for _ = range ticker.C {
				if last != id {
					writecounter(id)
					last = id
				}
			}
		}()
	}

	http.Handle("/id", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ID(w, req, atomic.AddUint64(&id, 1))
	}))

	http.Handle("/guid", http.HandlerFunc(GUID))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func writecounter(current uint64) {
	log.Println("writing counter: ", current);
	err := ioutil.WriteFile(*filename, []byte(strconv.FormatUint(current, 10)), 0644)
	if err != nil {
		log.Fatal("error: ", err)
	}
}

func readcounter() uint64 {
	var id uint64
	d, err := ioutil.ReadFile(*filename)
	if err != nil {
		log.Println("resetting counter")
	} else {
		data := string(d)
		id, err = strconv.ParseUint(data, 10, 64)
		if err != nil {
			log.Printf("error parsing data: %s, error: %s", data, err)
			id = 0
		}
		log.Println("current counter: ", id)
	}

	return id
}

func ID(w http.ResponseWriter, req *http.Request, current uint64) {
	fmt.Fprintln(w, current)
}

func GUID(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "%s", uuid.NewV4())
}
