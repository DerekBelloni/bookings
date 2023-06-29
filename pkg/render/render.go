package render

import (
	"bytes"
	"fmt"
	"github.com/derekbelloni/bookings/pkg/config"
	"github.com/derekbelloni/bookings/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	var templateCache map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// get requested template from cache
	myTemplate, ok := templateCache[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buffer := new(bytes.Buffer)

	td = AddDefaultData(td)

	err := myTemplate.Execute(buffer, td)
	if err != nil {
		log.Println(err)
	}
	// render the template
	_, err = buffer.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing template to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// go to file system, get everything that ends in *.page.tmpl
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with all files *.page.tmpl
	for _, page := range pages {
		// get the file name
		name := filepath.Base(page)
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}

		if len(matches) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = templateSet
	}
	return myCache, nil
}

// Original Way for Handling Templates in Tutorial
//var templateCache = make(map[string]*template.Template)

//func RenderTemplate(w http.ResponseWriter, t string) {
//	var tmpl *template.Template
//	var err error
//
//	// check to see if we already have the template in our cache
//	_, inMap := templateCache[t]
//	if !inMap {
//		// need to create template
//		log.Println("creating template and adding to cache")
//		err = createTemplateCache(t)
//		if err != nil {
//			log.Println(err)
//		}
//	} else {
//		// we have the template in the cache
//		log.Println("using cached template")
//	}
//
//	tmpl = templateCache[t]
//	err = tmpl.Execute(w, nil)
//	if err != nil {
//		log.Println(err)
//	}
//}
//
//func createTemplateCache(t string) error {
//	templates := []string{
//		fmt.Sprintf("./templates/%s", t),
//		"./templates/base.layout.tmpl",
//	}
//
//	// parse the template
//	tmpl, err := template.ParseFiles(templates...)
//	if err != nil {
//		return err
//	}
//	// add template to cache (map)
//	templateCache[t] = tmpl
//
//	return nil
//}
