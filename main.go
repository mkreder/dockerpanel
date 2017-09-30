package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"log"
	"os"
	"path/filepath"

	"github.com/mkreder/dockerpanel/controller"
	"github.com/mkreder/dockerpanel/templates"

)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	templates.WriteHomeTemplate(w)
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	templates.WriteLoginTemplate(w)
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}



func main() {



	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", HomeHandler)
	r.Get("/web", controller.WebHandler)
	r.Get("/login", LoginHandler)

	r.Post("/addweb", controller.AddWeb)
	r.Get("/removeWeb",controller.RemoveWeb)

	port := 9090
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	workDir, _ := os.Getwd()
	distDir := filepath.Join(workDir, "dist")
	jsDir := filepath.Join(workDir, "js")
	vendorDir := filepath.Join(workDir, "vendor")

	FileServer(r, "/dist", http.Dir(distDir))
	FileServer(r, "/js", http.Dir(jsDir))
	FileServer(r, "/vendor", http.Dir(vendorDir))

	log.Printf("Running on %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
