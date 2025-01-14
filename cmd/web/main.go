package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

type config struct {
	port int
}

var cfg config

func main() {
	// don't use ports 0 ~ 1023 as it used by OS
	flag.IntVar(&cfg.port, "port", 4000, "HTTP network address")
	// you need to call this *before* you use the addr variable
	// otherwise it will always be the default value ":4000"
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/")})
	mux.Handle(("GET /static"), http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /snippet/view/{id}", snippetView)
	mux.HandleFunc("GET /snippet/create", snippetCreate)
	mux.HandleFunc("POST /snippet/create", snippetCreatePost)

	logger.Info("starting server", "port", cfg.port)

	err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.port), mux)
	logger.Error(err.Error())
	os.Exit(1)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		log.Printf("File system open fail: %v", err)
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		log.Printf("File state checking fail: %v", err)
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		log.Printf("Index file: %s", index)
		if _, err := nfs.fs.Open(index); err != nil {
			log.Printf("Index file open fail: %v", err)
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}
