package main

import (
	"net/http"
	"os"
	"text/template"
	"bufio"
	"strings"
	"log"
)

func main() {
	templates := populateTemplates()
	http.HandleFunc("/",
		func(w http.ResponseWriter, req *http.Request) {
			requestFile := req.URL.Path[1:]
			template := templates.Lookup(requestFile + ".html")
			if template != nil {
				template.Execute(w, nil)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		})

	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)

	http.ListenAndServe(":8000", nil)
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	log.Println("requesting file", path)
	var contentType string
	if strings.HasSuffix(path, ".css") {
		log.Println("setting content type to css")
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		log.Println("setting content type to png")
		contentType = "image/png"
	} else {
		log.Println("setting content type to plain text")
		contentType = "text/plain"
	}

	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		w.Header().Add("Content-Type", contentType)
		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		log.Println("requested file not found")
		w.WriteHeader(http.StatusNotFound)
	}
}

func populateTemplates() *template.Template {
	result := template.New("templates")

	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)

	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath + "/" + pathInfo.Name())
		}
	}

	result.ParseFiles(*templatePaths...)
	return result
}