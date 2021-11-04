package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gorilla/mux"

	"strings"

	"xteve/src"
)

// Name
const Name = "Tube"

// Version
const Version = "2.2.0.0200"

// DBVersion
const DBVersion = "2.1.0"

// APIVersion
const APIVersion = "1.1.0"

var homeDirectory = fmt.Sprintf("%s%s.%s%s", src.GetUserHomeDirectory(), string(os.PathSeparator), strings.ToLower(Name), string(os.PathSeparator))
var samplePath = fmt.Sprintf("%spath%sto%sxteve%s", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator))
var sampleRestore = fmt.Sprintf("%spath%sto%sfile%s", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator))

var configFolder = flag.String("config", "", ": Config Folder        ["+samplePath+"] (default: "+homeDirectory+")")
var port = flag.String("port", "", ": Server port          [34400] (default: 34400)")
var restore = flag.String("restore", "", ": Restore from backup  ["+sampleRestore+"xteve_backup.zip]")
var debug = flag.Int("debug", 0, ": Debug level          [0 - 3] (default: 0)")
var info = flag.Bool("info", false, ": Show system info")
var h = flag.Bool("h", false, ": Show help")

// Aktiviert den Entwicklungsmodus. Für den Webserver werden dann die lokalen Dateien verwendet.
var dev = flag.Bool("dev", false, ": Activates the developer mode, the source code must be available. The local files for the web interface are used.")

// var listen = flag.String("listen", ":8080", "Listen on address")

//go:embed frontend/build
var embeded embed.FS

func main() {

	// Build-Nummer von der Versionsnummer trennen
	var build = strings.Split(Version, ".")

	mode := flag.String("mode", "proxy", "Mode to serve REACT site: proxy, dir, pkger, embed ")
	// Proxy mode proxies all other connections to the npm server
	proxy := flag.String("proxy", "http://localhost:3000/", "Address to proxy requests to")
	// Dir follows general filesystem pathing rules
	dir := flag.String("dir", "./frontend/build/", "Directory where the static built app resides")
	// Embed uses the new 1.16 embed functions to offer what pkger does
	embed := flag.String("embed", "frontend/build", "Directory where the static built embeded app resides (1.16+)")

	var system = &src.System
	system.APIVersion = APIVersion
	system.Build = build[len(build)-1:][0]
	system.DBVersion = DBVersion
	system.Name = Name
	system.Version = strings.Join(build[0:len(build)-1], ".")

	// Panic !!!
	defer func() {

		if r := recover(); r != nil {

			fmt.Println()
			fmt.Println("* * * * * FATAL ERROR * * * * *")
			fmt.Println("OS:  ", runtime.GOOS)
			fmt.Println("Arch:", runtime.GOARCH)
			fmt.Println("Err: ", r)
			fmt.Println()

			pc := make([]uintptr, 20)
			runtime.Callers(2, pc)

			for i := range pc {

				if runtime.FuncForPC(pc[i]) != nil {

					f := runtime.FuncForPC(pc[i])
					file, line := f.FileLine(pc[i])

					if string(file)[0:1] != "?" {
						fmt.Printf("%s:%d %s\n", filepath.Base(file), line, f.Name())
					}

				}

			}

			fmt.Println()
			fmt.Println("* * * * * * * * * * * * * * * *")

		}

	}()

	flag.Parse()

	if *h {
		flag.Usage()
		return
	}

	system.Dev = *dev

	// Systeminformationen anzeigen
	if *info {

		system.Flag.Info = true

		err := src.Init()
		if err != nil {
			src.ShowError(err, 0)
			os.Exit(0)
		}

		src.ShowSystemInfo()
		return

	}

	// Webserver Port
	if len(*port) > 0 {
		system.Flag.Port = *port
	}

	// Debug Level
	system.Flag.Debug = *debug
	if system.Flag.Debug > 3 {
		flag.Usage()
		return
	}

	// Speicherort für die Konfigurationsdateien
	if len(*configFolder) > 0 {
		system.Folder.Config = *configFolder
	}

	// Backup wiederherstellen
	if len(*restore) > 0 {

		system.Flag.Restore = *restore

		err := src.Init()
		if err != nil {
			src.ShowError(err, 0)
			os.Exit(0)
		}

		err = src.XteveRestoreFromCLI(*restore)
		if err != nil {
			src.ShowError(err, 0)
		}

		os.Exit(0)
	}

	err := src.Init()
	if err != nil {
		src.ShowError(err, 0)
		os.Exit(0)
	}

	err = src.StartSystem(false)
	if err != nil {
		src.ShowError(err, 0)
		os.Exit(0)
	}

	err = src.InitMaintenance()
	if err != nil {
		src.ShowError(err, 0)
		os.Exit(0)
	}

	// Basic ServeMux and API that just sends the time
	r := mux.NewRouter()

	// mux.HandleFunc("/api", basicAPI)
	r.HandleFunc("/", src.Index)
	r.HandleFunc("/stream/", src.Stream)
	r.HandleFunc("/xmltv/", src.Tube)
	r.HandleFunc("/m3u/", src.Tube)
	r.HandleFunc("/web/", src.Web)
	r.HandleFunc("/download/", src.Download)
	r.HandleFunc("/api/", src.API)
	r.HandleFunc("/images/", src.Images)
	r.HandleFunc("/data_images/", src.DataImages)

	//Broken Out (NEW)
	r.HandleFunc("/api/status", src.GetStatus).Methods("GET")
	r.HandleFunc("/api/backup", src.Backup).Methods("GET")
	r.HandleFunc("/api/config", src.Config).Methods("POST")
	r.HandleFunc("/api/filter", src.PlaylistFilter).Methods("POST")

	// The React serve magic
	switch *mode {
	case "proxy":
		// Proxy mode is most useful for development
		// Preserves live-reload
		u, err := url.Parse(*proxy)
		if err != nil {
			log.Fatalf("Cannot parse proxy address: %s", err)
		}
		r.Handle("/", httputil.NewSingleHostReverseProxy(u))
	case "dir":
		// Dir mode is useful if you build your react app but don't want to embed it in the binary, such as Docker deploys
		r.Handle("/", http.FileServer(EmbedDir{http.Dir(*dir)}))
	case "embed":
		// Embed uses the new 1.16+ Embed functionality
		filesystem := fs.FS(embeded)
		static, err := fs.Sub(filesystem, *embed)
		if err != nil {
			log.Fatal("Cannot open filesystem", err)
		}
		r.Handle("/", http.FileServer(EmbedDir{http.FS(static)}))
	default:
		// Any other mode would assume you have a reverse proxy, like nginx, that filters traffic
		log.Println("No react mode; this only works if you have a frontend reverse proxy")
	}
	s := &http.Server{
		Addr:         *port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Println(s.ListenAndServe())

}

// EmbedDir provides a convenience method to default requests back to /index.html, allowing react-router to work correctly
type EmbedDir struct {
	http.FileSystem
}

// Open implementation of http.FileSystem that falls back to serving /index.html, allowing react-router to operate
func (d EmbedDir) Open(name string) (http.File, error) {
	if f, err := d.FileSystem.Open(name); err == nil {
		return f, err
	} else {
		return d.FileSystem.Open("/index.html")
	}
}
