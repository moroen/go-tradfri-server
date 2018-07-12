package main

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"runtime"
	"time"

	"github.com/coreos/go-systemd/daemon"
	coap "github.com/moroen/go-tradfricoap"
)

type Todo struct {
	Name      string    `json:"name"`
	Completed bool      `json:"completed"`
	Due       time.Time `json:"due"`
}

const port = "8085"

func getExecDir() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	return path.Dir(currentFilePath)
}

func main() {

	err := coap.LoadConfig()
	if err != nil {
		panic("\nNo config found!")
	}

	router := NewRouter()
	log.Printf("Starting server - listening on port %s", port)
	daemon.SdNotify(false, "READY=1")
	go func() {
		interval, err := daemon.SdWatchdogEnabled(false)
		if err != nil || interval == 0 {
			return
		}
		for {
			_, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s", port)) // Check if server is running
			if err == nil {
				daemon.SdNotify(false, "WATCHDOG=1")
			}
			time.Sleep(interval)
		}
	}()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))
}

func Logger(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"%s\t%s\t%s\t%s",
			r.Method,
			r.RequestURI,
			name,
			time.Since(start),
		)
	})
}
