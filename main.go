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
	"github.com/mkreder/dockerpanel/login"
	"github.com/mkreder/dockerpanel/templates"
	"github.com/mkreder/dockerpanel/worker"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	userName := login.GetUNombreUsuario(r)
	if userName != "" {
		templates.WriteHomeTemplate(w, tools.GetRunningContainers())
	} else {
		templates.WriteLoginTemplate(w, "")
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("El servidor de archivos no permite parametros URL.")
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
	r.Get("/removeWeb", controller.RemoveWeb)

	r.Get("/dns", controller.ZonaHandler)
	r.Post("/dns", controller.AddRegion)
	r.Get("/removeZona", controller.RemoveZona)
	r.Get("/editRegistros", controller.RegistroHandler)
	r.Post("/registros", controller.AddRegistro)
	r.Get("/registros", controller.AddRegistro)
	r.Get("/removeRegistro", controller.RemoveRegistro)

	r.Get("/ftp", controller.FtpHandler)
	r.Post("/ftp", controller.AddUsuarioFtp)
	r.Get("/removeUsuarioFtp", controller.RemoveUsuarioFtp)
	r.Post("/ftpconfig", controller.FtpConfigHandler)
	r.Get("/ftpconfig", controller.FtpConfigHandler)

	r.Get("/bd", controller.BDHandler)
	r.Post("/bd", controller.AddBD)
	r.Get("/removeBd", controller.RemoveBD)
	r.Post("/ubd", controller.AddUsuarioBD)
	r.Get("/ubd", controller.AddUsuarioBD)
	r.Get("/removeUsuarioBD", controller.RemoveUsuarioBD)
	r.Post("/addubd", controller.AssociateBD)
	r.Get("/addubd", controller.AssociateBD)
	r.Post("/removeubd", controller.DisassociateBD)
	r.Get("/removeubd", controller.DisassociateBD)
	r.Post("/addbdip", controller.AddIP)
	r.Get("/addbdip", controller.AddIP)
	r.Post("/removebdip", controller.RemoveIP)
	r.Get("/removebdip", controller.RemoveIP)

	r.Get("/mail", controller.MailHandler)
	r.Post("/mail", controller.AddDominio)
	r.Get("/removeDominio", controller.RemoveDominio)
	r.Get("/editListas", controller.ListaHandler)
	r.Post("/addLista", controller.AddLista)
	r.Get("/addLista", controller.AddLista)
	r.Get("/removeLista", controller.RemoveLista)
	r.Get("/editCuentas", controller.CuentaHandler)
	r.Post("/addCuenta", controller.AddCuenta)
	r.Get("/addCuenta", controller.AddCuenta)
	r.Get("/removeCuenta", controller.RemoveCuenta)

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

	log.Printf("Aplicación iniciada en el puerto %s\n", srv.Addr)
	log.Fatal(srv.ListenAndServe())

}
