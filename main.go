package main

import (
//	"./config"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
	"os"
)

// Command-line flags.
var (
	httpAddr   = flag.String("http", ":8080", "Listen address")
	pollPeriod = flag.Duration("poll", 15*time.Minute, "Polling period")
)

const forecastURL = "http://datapoint.metoffice.gov.uk/public/data/txt/wxfcs/regionalforecast/json/514"

func main() {
	flag.Parse()
	changeURL := fmt.Sprintf("%s?key=%s", forecastURL, os.Getenv("DATAPOINT_KEY"))
	http.Handle("/", NewServer(changeURL, *pollPeriod))
	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}

type Server struct {
	url     string
	period  time.Duration
}

// NewServer returns an initialized server.
func NewServer(url string, period time.Duration) *Server {
	return &Server{url: url, period: period}
}

// ServeHTTP implements the HTTP user interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := struct {
		URL     string
	}{
		s.url,
	}
	err := tmpl.Execute(w, data)
	if err != nil {
		log.Print(err)
	}
}

var tmpl = template.Must(template.New("tmpl").Parse(`
	<!DOCTYPE html>
	<html>
		<body>
			<center>
				<h2>London Weather</h2>
				<h1>
					<a href="{{.URL}}">GET!</a>
				</h1>
			</center>
		</body>
	</html>
`))