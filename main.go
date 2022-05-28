package main

import (
	"flag"
	"github.com/gjrdiesel/cache-proxy/config-factory"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var Cache map[string][]byte

func main() {
	Program := filepath.Base(os.Args[0])
	Cache = make(map[string][]byte)
	c := config_factory.New(Program)
	config := c.Settings()

	redo := flag.Bool("ask", false, "rerun configuration setup")
	flag.Parse()
	if *redo != false {
		config = c.RedoConfiguration()
	}

	log.Printf("Serving %s (siphoning %s) on HTTP port: %s\n", Program, *config.SiphonUrl, *config.Port)
	log.Println("Pass a `--ask` to rerun configuration")
	log.Println("Listening at http://127.0.0.1:" + *config.Port)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		url := *config.SiphonUrl + request.URL.Path
		u := request.URL.Host + request.RequestURI

		if Cache[url] != nil {
			writer.Write(Cache[url])
			log.Println("HIT " + request.Method + ": " + u + " -> " + url)
			return
		}

		log.Println("MISS " + request.Method + ": " + u + " -> " + url)
		req, _ := http.NewRequest(request.Method, url, nil)
		res, _ := http.DefaultClient.Do(req)
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		writer.WriteHeader(res.StatusCode)
		writer.Write(body)

		Cache[url] = body
	})

	log.Fatal(http.ListenAndServe(":"+*config.Port, nil))
}
