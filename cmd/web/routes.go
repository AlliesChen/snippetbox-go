package main

import (
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./ui/static/"), app.logger})
	mux.Handle(("GET /static"), http.NotFoundHandler())
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.home))
	mux.Handle("GET /snippet/view/{id}", dynamic.ThenFunc(app.snippetView))
	mux.Handle("GET /snippet/create", dynamic.ThenFunc(app.snippetCreate))
	mux.Handle("POST /snippet/create", dynamic.ThenFunc(app.snippetCreatePost))

	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, commonHeaders)
	return standardMiddleware.Then(mux)
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
