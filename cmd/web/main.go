package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/derekbelloni/bookings/pkg/config"
	"github.com/derekbelloni/bookings/pkg/handlers"
	"github.com/derekbelloni/bookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	//change this to true when in production
	app.InProduction = false

	// initialize a new session
	session = scs.New()
	// set the duration of a session
	session.Lifetime = 24 * time.Hour
	// set parameters for cookies
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // in production this should be set to true

	app.Session = session

	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("can not create template cache")
	}
	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	fmt.Printf("Starting application on port %v \n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
