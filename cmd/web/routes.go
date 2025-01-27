package main

import (
	"log/slog"
	"net/http"
	"path/filepath"
)

func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/"), app.logger})

	mux.Handle(("GET /static"), http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	mux.HandleFunc("GET /{$}", app.home)
	mux.HandleFunc("GET /snippet/view/{id}", app.snippetView)
	mux.HandleFunc("GET /snippet/create", app.snippetCreate)
	mux.HandleFunc("POST /snippet/create", app.snippetCreatePost)

	return mux
}

type neuteredFileSystem struct {
	fs     http.FileSystem
	logger *slog.Logger
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		nfs.logger.Error(err.Error(), "path", path)
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		nfs.logger.Error(err.Error(), "path", path)
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			nfs.logger.Error(err.Error(), "index.html", index)
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}

	return f, nil
}
