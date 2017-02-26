package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	root_url         = "http://www.gunviolencearchive.org/reports/mass-shooting"
	root_url_short   = "http://www.gunviolencearchive.org"
	PagesDir         = "./pages"
	IncidentsDir     = "./incidents"
	JsonDir          = "./json"
	JsonIncidentsDir = "./incidents_json"
	lowestYear       = 2014
)

var (
	bCleanAll      bool
	bCleanIncident bool
	bAllYears      bool
)

type Content struct {
	Incidents []*Incident
}

type rootPage struct {
	year    int
	content []byte
}

type IncidentList struct {
	IncidentList []Incident `json:"List"`
}

type Incident struct {
	Date        time.Time `json:"Date"`
	ProcessDate time.Time `json:"ProcessDate"`
	State       string    `json:"State"`
	StateCode   string    `json:"StateCode"`
	City        string    `json:"City"`
	Direction   string    `json:"Direction"`
	Killed      int       `json:"Killed"`
	Injured     int       `json:"Injured"`
	Total       int       `json:"Total"`
	IncidentId  string    `json:"IncidentId"`
	Source      string    `json:"Source"`
	Latitude    string    `json:"Latitude"`
	Longitude   string    `json:"Longitude"`
}

type DonwLoadedIncident struct {
	incident Incident
	content  []byte
}

func Init() {
	initLogger()
	initStates()
}

func main() {
	flag.BoolVar(&bCleanAll, "cleanAll", false, "clean all the content previously loaded")
	flag.BoolVar(&bCleanIncident, "cleanIncidents", false, "clean all the incident content previously loaded")
	flag.BoolVar(&bAllYears, "allYears", false, "load all years since 2013")
	flag.Parse()

	Init()
	checkDir(bCleanIncident, IncidentsDir)
	checkDir(bCleanAll, PagesDir, JsonIncidentsDir, JsonDir)
	runIt()
}

func checkDir(b bool, paths ...string) {
	defer traceTime("checkDir")()
	for _, path := range paths {
		if b {
			os.RemoveAll(path)
			TraceActivity.Printf("The directory : %s has been cleaned\n", path)
		}

		if _, err := os.Stat(path); os.IsNotExist(err) {
			TraceActivity.Printf("The directory : %s has been created\n", path)
			os.Mkdir(path, 0777)
		}
	}
}

func runIt() {
	defer traceTime("runIt")()
	years := getYears()
	rootPages := []rootPage{}
	for _, year := range years {
		TraceActivity.Printf("Years to process %d\n", year)
		currentYear := time.Now().Year()
		if year == currentYear {
			year = 0
		}
		c := rootPage{year, downloadPage(year, 0)}
		rootPages = append(rootPages, c)
	}
	getYearPages(rootPages)
}

func getYears() []int {
	var result []int
	currentYear := time.Now().Year()
	if bAllYears {
		nbYears := currentYear - lowestYear + 1
		result = []int{}
		for i := 0; i <= nbYears-1; i++ {
			result = append(result, currentYear)
			currentYear--
		}
	} else {
		result = []int{currentYear}
	}
	return result
}

func getYearPages(rootPages []rootPage) {
	allPages := []rootPage{}
	for _, rp := range rootPages {
		allPages = append(allPages, rp)
		TraceActivity.Printf("downloadPages, key: %s \n", rp.year)

		c, err := html.Parse(bytes.NewReader(rp.content))
		check(err)
		lastPage := traverseLastPage(c)
		TraceActivity.Printf("get last page %d \n", lastPage)

		keyParts := bytes.Split([]byte(strconv.Itoa(rp.year)), []byte("_"))
		year, err := strconv.Atoi(string(keyParts[0]))
		check(err)
		for i := lastPage; i > 0; i-- {
			TraceActivity.Printf("looking for page  %d_%d \n", rp.year, i)
			if year == time.Now().Year() {
				year = 0
			}
			page := rootPage{year, downloadPage(rp.year, i)}
			allPages = append(allPages, page)
		}
	}
	incidentList := parsePagesYears(allPages)
	getIncidents(incidentList)
}

func parsePagesYears(allPages []rootPage) (incidentList []Incident) {
	defer traceTime("parsePagesYears")()
	TraceActivity.Printf("parsePagesYears\n")
	incidentList = make([]Incident, 0)

	done := make(chan bool, len(allPages))
	incidents := make(chan Incident)

	for _, pageContent := range allPages {
		//fmt.Printf("parsing page year \n")
		go traverseIncidentLinks(pageContent.content, done, incidents)
	}

	j := 0
	for {
		select {
		case incident := <-incidents:
			incidentList = append(incidentList, incident)
			break
		case <-done:
			j++
			if j == len(allPages) {
				return
			}
		}
	}
	return
}

func getIncidents(incidents []Incident) {
	defer traceTime("getIncidents")()
	downloadedIncidents := make(map[string]*DonwLoadedIncident)
	for _, incident := range incidents {
		fileName := fmt.Sprintf("%s/incidents_%s.txt", IncidentsDir, incident.IncidentId)
		if _, err := os.Stat(fileName); os.IsNotExist(err) {
			if strings.TrimSpace(incident.IncidentId) != "" {
				downloadedIncident := &DonwLoadedIncident{incident, downloadIncident(incident.IncidentId)}
				downloadedIncidents[incident.IncidentId] = downloadedIncident
			}
		} else {
			b, err := ioutil.ReadFile(fileName)
			check(err)
			downloadedIncident := &DonwLoadedIncident{incident, b}
			downloadedIncidents[incident.IncidentId] = downloadedIncident
		}
	}
	completeIncidentFiles(downloadedIncidents)
	buildGlobalJsonFile(downloadedIncidents)
}

func completeIncidentFiles(downloadedIncidents map[string]*DonwLoadedIncident) {
	for _, downloadIncident := range downloadedIncidents {
		parseLocation(downloadIncident)
		incidentJson, _ := json.Marshal(downloadIncident.incident)

		fileName := fmt.Sprintf("./%s/incident_%s.json", JsonIncidentsDir, downloadIncident.incident.IncidentId)
		fmt.Printf("json file name %s\n", fileName)

		f, err := os.Create(fileName)
		check(err)

		_, err = f.Write([]byte(string(incidentJson)))
		check(err)
		f.Close()
	}
}

func buildGlobalJsonFile(downloadedIncidents map[string]*DonwLoadedIncident) {
	//fmt.Printf("buildGlobalJsonFile... \n")
	//fmt.Printf("total file... %d\n", len(downloadedIncidents))
	defer traceTime("buildGlobalJsonFile")()

	l := &IncidentList{}
	for _, downloadIncident := range downloadedIncidents {
		l.IncidentList = append(l.IncidentList, downloadIncident.incident)
	}
	//fmt.Printf("list length... %d\n", len(l.IncidentList))

	incidentJson, _ := json.Marshal(l)

	fileName := fmt.Sprintf("./%s/%s", "json", "incidents_all.json")
	//fmt.Printf("json file name %s\n", fileName)

	f, err := os.Create(fileName)
	check(err)

	_, err = f.Write([]byte(string(incidentJson)))
	check(err)
	f.Close()
}
