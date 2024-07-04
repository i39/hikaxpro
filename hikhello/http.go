package main

import (
	"fmt"
	log "github.com/go-pkgz/lgr"
	"html/template"
	"net/http"
	"path/filepath"
)

func httpPoller() error {
	defer wg.Done()
	// Parse the templates
	mainTplPath := filepath.Join("templates", "main.html")
	partialTplPath := filepath.Join("templates", "partial.html")

	mainTmpl, err := template.ParseFiles(mainTplPath)
	if err != nil {
		return fmt.Errorf("error parsing main template: %v", err)
	}

	partialTmpl, err := template.ParseFiles(partialTplPath)
	if err != nil {
		return fmt.Errorf("error parsing partial template: %v", err)
	}

	// HTTP handler to serve the main template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := mainTmpl.Execute(w, nil)
		if err != nil {
			log.Printf("[ERROR] error execute main template")
		}
	})

	// HTTP handler to serve the zone list as partial HTML
	http.HandleFunc("/zones", func(w http.ResponseWriter, r *http.Request) {
		err := partialTmpl.Execute(w, deviceInfoList)
		if err != nil {
			log.Printf("[ERROR] error execute partial template")
		}
	})
	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Printf("[ERROR] your browser does not support server-sent events (SSE).")
			return
		}

		for {
			select {
			case <-dataChangedToHTTP:
				log.Printf("[DEBUG] polling to http")
				_, err := fmt.Fprintf(w, "data: update\n\n")
				if err != nil {
					log.Printf("[ERROR] Error writing to the response: %v", err)
					return
				}
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}
	})

	addr := listenAddress(opts.HttpListen)
	log.Printf("[DEBUG] listen address %s", addr)
	// Start the HTTP server
	log.Printf("[INFO] Server started at : %v", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		return fmt.Errorf("error starting the server %v", err)
	}
	return nil
}
