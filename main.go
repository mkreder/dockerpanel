package main

import (
	"fmt"
	"github.com/mkreder/dockerpanel/tools"
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
	"github.com/mkreder/dockerpanel/login"
	"github.com/mkreder/dockerpanel/worker"

)


func HomeHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUserName(r)
	if userName != "" {
		templates.WriteHomeTemplate(w,tools.GetRunningContainers())
	} else {
		templates.WriteLoginTemplate(w,"")
	}
}


func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("File server does not allow URL parameters.")
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
	r.Post("/web", controller.AddWeb)
	r.Get("/removeWeb",controller.RemoveWeb)

	r.Get("/dns", controller.ZoneHandler)
	r.Post("/dns", controller.AddZone)
	r.Get("/removeZone",controller.RemoveZone)
	r.Get("/editRecords",controller.RecordHandler)
	r.Post("/records",controller.AddRecord)
	r.Get("/records",controller.AddRecord)
	r.Get("/removeRecord",controller.RemoveRecord)

	r.Get("/ftp", controller.FtpHandler)
	r.Post("/ftp", controller.AddFTPUser)
	r.Get("/removeFTPUser",controller.RemoveFTPUser)
	r.Post("/ftpconfig", controller.FtpConfigHandler)
	r.Get("/ftpconfig", controller.FtpConfigHandler)

	r.Get("/db", controller.DatabaseHandler)
	r.Post("/db",controller.AddDatabase)
	r.Get("/removeDatabase",controller.RemoveDatabase)
	r.Post("/userdb",controller.AddDatabaseUser)
	r.Get("/userdb",controller.AddDatabaseUser)
	r.Get("/removeDatabaseUser",controller.RemoveDatabaseUser)
	r.Post("/adduserdb",controller.AssociateDatabase)
	r.Get("/adduserdb",controller.AssociateDatabase)
	r.Post("/removeuserdb",controller.DisassociateDatabase)
	r.Get("/removeuserdb",controller.DisassociateDatabase)
	r.Post("/adddbip",controller.AddDatabaseIP)
	r.Get("/adddbip",controller.AddDatabaseIP)
	r.Post("/removedbip",controller.RemoveDatabaseIP)
	r.Get("/removedbip",controller.RemoveDatabaseIP)

	r.Get("/mail",controller.MailHandler)
	r.Post("/mail",controller.AddDomain)
	r.Get("/removeDomain",controller.RemoveDomain)
	r.Get("/editLists",controller.ListHandler)
	r.Post("/addList",controller.AddList)
	r.Get("/addList",controller.AddList)
	r.Get("/removeList",controller.RemoveList)
	r.Get("/editAccounts",controller.AccountHandler)
	r.Post("/addAccount",controller.AddAccount)
	r.Get("/addAccount",controller.AddAccount)
	r.Get("/removeAccount",controller.RemoveAccount)


	r.Get("/login", login.LoginHandler)
	r.Post("/login", login.LoginHandler)

	r.Get("/profile", controller.ProfileHandler)
	r.Post("/profile", controller.ProfileHandler)

	r.Get("/logout", login.LogoutHandler)

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

	go worker.RunDNSWorker()
	go worker.RunDBWorker()
	go worker.RunMailWorker()
	go worker.RunWebWorker()
	go worker.RunFTPWorker()

	log.Printf("Application started on port %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
