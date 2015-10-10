package main

import (
	"net/http"
	"os"
	"text/template"
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

	http.ListenAndServe(":8000", nil)
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