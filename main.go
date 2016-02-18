package main

import (
	"./config"
	"./utils"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var (
	httpAddr    = flag.String("http", ":8080", "Listen address")
	pollPeriod  = flag.Duration("poll", 5*time.Minute, "Polling period")
	forecastURL = "http://datapoint.metoffice.gov.uk/public/data/txt/wxfcs/regionalforecast/json/514"
)

func main() {
	flag.Parse()
	tmpl, err := loadTemplate("main")
	utils.ErrorPanic(err)

	http.Handle("/", NewServer(tmpl, *pollPeriod))
	utils.ErrorPanic(http.ListenAndServe(*httpAddr, nil))
}

func loadTemplate(name string) (*template.Template, error) {
	fileBytes, readErr := ioutil.ReadFile(fmt.Sprintf("templates/%s.html", name))
	utils.ErrorPanic(readErr)

	t := template.New("tmpl")

	return t.Parse(string(fileBytes))
}

type Server struct {
	tmpl             *template.Template
	period           time.Duration
	forecastLongTerm *config.ForecastLongTerm
}

// NewServer returns an initialized server.
func NewServer(tmpl *template.Template, period time.Duration) *Server {
	s := &Server{tmpl: tmpl, period: period}
	go s.poll()
	return s
}

// ServeHTTP implements the HTTP user interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		ForecastLongTerm *config.ForecastLongTerm
	}{
		s.forecastLongTerm,
	}
	err := s.tmpl.Execute(w, data)
	utils.ErrorPanic(err)
}

func (s *Server) poll() {
	s.forecastLongTerm = fetchJson(fmt.Sprintf("%s?key=%s", forecastURL, os.Getenv("DATAPOINT_KEY")))
	time.Sleep(s.period)
}

func fetchJson(url string) *config.ForecastLongTerm {
	res, err := http.Get(url)
	utils.ErrorPanic(err)
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)

	var data config.ForecastLongTerm
	err = decoder.Decode(&data)

	return &data
}
