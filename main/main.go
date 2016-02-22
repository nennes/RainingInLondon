package main

import (
	"github.com/nennes/RainingInLondon/config"
	"github.com/nennes/RainingInLondon/utils"
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
	pollPeriod  = flag.Duration("poll", 5*time.Minute, "Polling period")
	forecastURL = "http://datapoint.metoffice.gov.uk/public/data/txt/wxfcs/regionalforecast/json/514"
)

func main() {
	flag.Parse()
	tmpl, err := loadTemplate("main")
	utils.ErrorPanic(err)

	// Get the port from the environment variable, or use 8080 if none is found
	port := (os.Getenv("PORT"))
	if len(port) == 0 {
		port = "8080"
	}

	// Handle all requests to the web root ("/")
	// with a Server instance returned by NewServer
	http.Handle("/", NewServer(tmpl, *pollPeriod))

	// Listen on the specified port on any interface
	// This will block until the program is terminated
	utils.ErrorPanic(http.ListenAndServe(":" + port , nil))
}

func loadTemplate(name string) (*template.Template, error) {
	// Load the template specified by the name parameter
	fileBytes, readErr := ioutil.ReadFile(fmt.Sprintf("templates/%s.html", name))
	utils.ErrorPanic(readErr)

	// Create a template object and parse it
	return template.New("tmpl").Parse(string(fileBytes))
}

type Server struct {
	tmpl             *template.Template
	period           time.Duration
	poll_time	 time.Time
	forecastLongTerm *config.ForecastLongTerm
}

// NewServer returns an initialized server.
func NewServer(tmpl *template.Template, period time.Duration) *Server {
	s := &Server{tmpl: tmpl, period: period}
	// Launch a goroutine that keeps polling the target URL
	go s.poll()
	return s
}

// ServeHTTP implements the HTTP user interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Create the structure that will be passed to the template
	data := struct {
		ForecastLongTerm *config.ForecastLongTerm
		PollTime	 string
	}{
		s.forecastLongTerm,
		s.poll_time.Format(time.Kitchen),
	}

	// Execute the template by applying it to the data structure
	err := s.tmpl.Execute(w, data)
	utils.ErrorPanic(err)
}

func (s *Server) poll() {
	for {
		// Get the JSON using an environment variable as the key parameter
		s.forecastLongTerm = fetchJson(fmt.Sprintf("%s?key=%s", forecastURL, os.Getenv("DATAPOINT_KEY")))
		// Update the last polled time
		s.poll_time = time.Now()
		// Sleep for the specified amount of time (passed in or 5 minute default)
		time.Sleep(s.period)
	}
}

func fetchJson(url string) *config.ForecastLongTerm {
	// Do a GET on the passed in URL
	res, err := http.Get(url)
	utils.ErrorPanic(err)
	defer res.Body.Close()

	// Create a decoder that reader from res.Body
	decoder := json.NewDecoder(res.Body)

	// Declare the variable that will hold the JSON data
	// Its type is a struct that reflects the JSON format
	var data config.ForecastLongTerm

	// Read the next JSON-encoded value from its input
	// and store it in the data struct
	err = decoder.Decode(&data)

	return &data
}
