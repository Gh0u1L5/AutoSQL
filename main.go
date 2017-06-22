package main

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"sort"
	"sync"
	"time"
)

var history = struct {
	*sync.Mutex
	records map[string][]string
}{
	&sync.Mutex{},
	make(map[string][]string),
}

func keys(query url.Values) []string {
	result := make([]string, 0, len(query))
	for key := range query {
		result = append(result, key)
	}
	return result
}

func equal(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for index, value := range a {
		if value != b[index] {
			return false
		}
	}
	return true
}

func scanURL(rawurl string) {
	u, err := url.Parse(rawurl)
	if err != nil || len(u.RawQuery) == 0 {
		return
	}

	history.Lock()
	path := u.Host + u.Path
	params := keys(u.Query())
	sort.Strings(params)
	if equal(params, history.records[path]) {
		history.Unlock()
		return
	} else {
		history.records[path] = params
		history.Unlock()
	}

	var t SQLmapTask
	t.create("")
	t.start(SQLmapTaskInfo{Url: rawurl})
	for {
		if t.status() == "running" {
			time.Sleep(5 * time.Minute)
			continue
		}
		result := t.result().([]interface{})
		if len(result) != 0 {
			log.Printf("\033[0;31mHIT!\033[0m URL: %v, Result: %v\n", rawurl, t.result())
		}
		t.delete()
		break
	}
}

func main() {
	ln, err := net.Listen("tcp", ":1017")
	if err != nil {
		log.Fatal(err)
	}
	http.Serve(ln, http.HandlerFunc(handleRequest))
}
