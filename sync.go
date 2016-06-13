//export GOROOT=/home/guillaume/go
//export PATH=$PATH:$GOROOT/bin
//export GOPATH=/home/guillaume/gowork

package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	root_url         = "http://www.gunviolencearchive.org/reports/mass-shooting"
	incident_url     = "http://www.gunviolencearchive.org/incident/"
	PagesDir         = "./pages"
	IncidentsDir     = "./incidents"
	JsonIncidentsDir = "./incidents_json"
	incidentJsonFile = ".json"
)

var (
	bbodys = []byte("<tbody>")
	bbodye = []byte("</tbody>")

	bRep1    = []byte("</span></div>")
	bRep2    = []byte("<div>")
	bRep3    = []byte(":")
	bRep4    = []byte("page=")
	bRep5    = []byte("\"")
	bReplEmp = []byte("")
	bSplit3  = []byte(",")
	bSplit4  = []byte("pager-last last")
	bSplit5  = []byte("?")
	bSplit6  = []byte(">")
	bae      = []byte("</a>")
	blie     = []byte("</li>")
	btre     = []byte("</tr>")
	btde     = []byte("</td>")
	bule     = []byte("</ul>")
	btrs     = []byte("<tr>")
	btds     = []byte("<td>")

	classEven = []byte(" class=\"even\"")
	classOdd  = []byte(" class=\"odd\"")
	blis1     = []byte("<li class=\"1 last\">")
	blis2     = []byte("<li class=\"0 first\">")
	blis3     = []byte("<li class=\"0 first last\">")

	blink       = []byte("<ul class=\"links inline links-new-lines\">")
	bincidentas = []byte("<a href=\"/incident/")
	bincident   = []byte("\">View Incident")
	bsource     = []byte("\">View Source")

	bSplit1 = []byte("<span>Geolocation")
	bSplit2 = []byte("<h2>Participants</h2>")

	cptFile            int
	cptFileIncident    int
	cptLine            int
	cptLineMameCounter = 1
	mutex              = &sync.Mutex{}

	stateMap map[string]string
)

type Incident struct {
	Date       string `json:"Date"`
	State      string `json:"State"`
	City       string `json:"City"`
	Direction  string `json:"Direction"`
	Killed     string `json:"Killed"`
	Injured    string `json:"Injured"`
	IncidentId string `json:"IncidentId"`
	Source     string `json:"Source"`
	Latitude   string `json:"Latitude"`
	Longitude  string `json:"Longitude"`
}

func main() {

	stateMap = make(map[string]string)

	stateMap["Alabama"] = "AL"
	stateMap["Alaska"] = "AK"
	stateMap["Arizona"] = "AZ"
	stateMap["Arkansas"] = "AR"
	stateMap["California"] = "CA"
	stateMap["Colorado"] = "CO"
	stateMap["Connecticut"] = "CT"
	stateMap["Delaware"] = "DE"
	stateMap["District of Columbia"] = "DC"
	stateMap["Florida"] = "FL"
	stateMap["Georgia"] = "GA"
	stateMap["Hawaii"] = "HI"
	stateMap["Idaho"] = "ID"
	stateMap["Illinois"] = "IL"
	stateMap["Indiana"] = "IN"
	stateMap["Iowa"] = "IA"
	stateMap["Kansas"] = "KS"
	stateMap["Kentucky"] = "KY"
	stateMap["Louisiana"] = "LA"
	stateMap["Maine"] = "ME"
	stateMap["Maryland"] = "MD"
	stateMap["Massachusetts"] = "MA"
	stateMap["Michigan"] = "MI"
	stateMap["Minnesota"] = "MN"
	stateMap["Mississippi"] = "MS"
	stateMap["Missouri"] = "MO"
	stateMap["Montana"] = "MT"
	stateMap["Nebraska"] = "NE"
	stateMap["Nevada"] = "NV"
	stateMap["New Hampshire"] = "NH"
	stateMap["New Jersey"] = "NJ"
	stateMap["New Mexico"] = "NM"
	stateMap["New York"] = "NY"
	stateMap["North Carolina"] = "NC"
	stateMap["North Dakota"] = "ND"
	stateMap["Ohio"] = "OH"
	stateMap["Oklahoma"] = "OK"
	stateMap["Oregon"] = "OR"
	stateMap["Pennsylvania"] = "PA"
	stateMap["Rhode Island"] = "RI"
	stateMap["South Carolina"] = "SC"
	stateMap["South Dakota"] = "SD"
	stateMap["Tennessee"] = "TN"
	stateMap["Texas"] = "TX"
	stateMap["Utah"] = "UT"
	stateMap["Vermont"] = "VT"
	stateMap["Virginia"] = "VA"
	stateMap["Washington"] = "WA"
	stateMap["West Virginia"] = "WV"
	stateMap["Wisconsin"] = "WI"
	stateMap["Wyoming"] = "WY"

	checkDir(PagesDir, true)
	checkDir(IncidentsDir, false)
	checkDir(JsonIncidentsDir, false)

	initFileName()

	downloadYears()
	parseYears()

	fmt.Printf("year file parsed: %d\n", cptFile)
	fmt.Printf("line parsed: %d\n", cptLine)
	fmt.Printf("incident files parsed: %d\n", cptFileIncident)
}

func initFileName() {
	files, err := ioutil.ReadDir(JsonIncidentsDir)
	check(err)
	for _, f := range files {
		fileName := strings.Split(f.Name(), ".")[0]
		i, _ := strconv.Atoi(fileName)
		if i > cptLineMameCounter {
			cptLineMameCounter = i
		}
	}
	fmt.Printf("cptLineMameCounter initialized to : %d\n", cptLineMameCounter)
}

func parseYears() {
	fmt.Printf("parseYears:\n")
	files, err := ioutil.ReadDir(PagesDir)
	l := len(files)
	fmt.Printf("files to parsed: %d\n", l)
	check(err)

	ch := make(chan bool, l)
	for _, f := range files {
		cptFile++
		content, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", PagesDir, f.Name()))
		check(err)

		result := bytes.Split(content, bbodys)
		if len(result) == 2 {
			result = bytes.Split(result[1], bbodye)
			go parseTBody(ch, result[0])
		} else {
			fmt.Printf("not in: \n")
		}
	}

	for i := 0; i < l; i++ {
		<-ch
		fmt.Printf("done : %d\n", (i + 1))
	}
}

func downloadYears() {
	// Current year
	downloadYearPages(0)
	for year := time.Now().Year() - 1; year > 2012; year-- {
		downloadYearPages(year)
	}
}

func downloadYearPages(year int) {
	page := 0
	err, lastPage := getLastPageIndex(buildYearBaseUrl(year, page))
	check(err)

	for i := 0; i <= lastPage; i++ {
		content := downloadUrlContent(buildYearBaseUrl(year, i))

		f, err := os.Create(fmt.Sprintf("%s/ms_y_%d_p_%d.txt", PagesDir, year, i))
		defer f.Close()
		check(err)

		_, err = f.Write(content)
		check(err)
	}
}

func parseTBody(ch chan bool, tbody []byte) {

	tbody = remove(tbody, bincidentas)
	tbody = remove(tbody, bincident)
	tbody = remove(tbody, bsource)
	tbody = remove(tbody, bbodys)
	tbody = remove(tbody, bbodye)
	tbody = remove(tbody, bae)
	tbody = remove(tbody, blie)
	tbody = remove(tbody, btre)
	tbody = remove(tbody, btde)
	tbody = remove(tbody, bule)
	tbody = remove(tbody, classEven)
	tbody = remove(tbody, classOdd)
	tbody = remove(tbody, blis1)
	tbody = remove(tbody, blis2)
	tbody = remove(tbody, blis3)
	tbody = remove(tbody, blink)

	out := string(tbody[:len(tbody)])
	out = strings.TrimSpace(out)

	if len(out) == 0 {
		fmt.Printf("empty out: \n")
	}

	result := bytes.Split(tbody, btrs)
	if len(result) == 0 {
		fmt.Printf("no split: \n")
		fmt.Printf("out: %s\n", out)
	}
	for i := range result {
		line := result[i]

		if bytes.Contains(line, btds) {
			cptLine++
			tokens := bytes.Split(line, btds)

			var id string
			var source string

			if bytes.Contains(tokens[7], []byte("<a href=\"")) {
				id = string(tokens[7][:len(tokens[7])])

				idPart := strings.Split(id, "<a href=\"")
				id = idPart[0]
				id = strings.Replace(id, "\n", "", -1)

				source = idPart[1]
			} else {
				id = string(tokens[7][:len(tokens[7])])
				source = ""
			}

			id = strings.TrimSpace(id)

			fileName := fmt.Sprintf("%s/%s.txt", IncidentsDir, id)

			if _, err := os.Stat(fileName); os.IsNotExist(err) {
				cptFileIncident++
				urlInc := incident_url + id
				response, err := http.Get(urlInc)
				check(err)

				defer response.Body.Close()
				contents, err := ioutil.ReadAll(response.Body)
				check(err)
				s := string(contents[:len(contents)])
				f, err := os.Create(fileName)
				check(err)
				defer f.Close()
				n2, err := f.Write([]byte(s))
				check(err)

				fmt.Printf("wrote %d bytes for the incident %s\n", n2, id)

				inc := new(Incident)
				inc.Date = strings.TrimSpace(string(tokens[1][:len(tokens[1])]))
				inc.State = stateMap[strings.TrimSpace(string(tokens[2][:len(tokens[2])]))]
				inc.City = strings.TrimSpace(string(tokens[3][:len(tokens[3])]))
				inc.Direction = strings.TrimSpace(string(tokens[4][:len(tokens[4])]))
				inc.Killed = strings.TrimSpace(string(tokens[5][:len(tokens[5])]))
				inc.Injured = strings.TrimSpace(string(tokens[6][:len(tokens[6])]))
				inc.IncidentId = strings.TrimSpace(id)
				inc.Source = strings.TrimSpace(source)

				pos := parseLocation(contents)
				inc.Latitude = pos[0]
				inc.Longitude = pos[1]

				incidentJson, _ := json.Marshal(inc)
				mutex.Lock()
				cptLineMameCounter++
				fileName := fmt.Sprintf("%s/%d.json", JsonIncidentsDir, cptLineMameCounter)
				fmt.Printf("json file name %s\n", fileName)

				mutex.Unlock()

				f, err = os.Create(fileName)
				check(err)

				_, err = f.Write([]byte(string(incidentJson)))
				check(err)
				f.Close()
			}
		}
	}
	ch <- true
}

func parseLocation(s []byte) []string {
	result := bytes.Split(s, bSplit1)
	if len(result) == 2 {
		result := bytes.Split(result[1], bSplit2)
		if len(result) == 2 {
			content := result[0]
			content = remove(content, bRep1)
			content = remove(content, bRep2)
			content = remove(content, bRep3)
			content = bytes.TrimSpace(content)
			result := bytes.Split(content, bSplit3)
			if len(result) == 2 {
				return []string{string(bytes.TrimSpace(result[0])), string(bytes.TrimSpace(result[1]))}
			} else {
				return []string{"", ""}
			}
		}
	}
	return []string{"", ""}
}

func getLastPageIndex(url string) (Error, int) {
	tokens := bytes.Split(downloadUrlContent(url), bSplit4)

	if len(tokens) == 2 {
		tokens = bytes.Split(tokens[1], bSplit5)
		if len(tokens) == 2 {
			tokens = bytes.Split(tokens[1], bSplit6)
			if len(tokens) > 1 {
				s := remove(tokens[0], bRep4)
				s = remove(s, bRep5)
				i, err := strconv.Atoi(string(s))
				if err == nil {
					return nil, i
				} else {
					return errors.New(url + ", cannot convert: " + string(s)), 0
				}
			} else {
				return errors.New(url + ", cannot find >"), 0
			}
		} else {
			return errors.New(url + ", cannot find ?"), 0
		}
	} else {
		return errors.New(url + ", cannot find pager-last last"), 0
	}
	return nil, 0
}

func downloadUrlContent(url string) []byte {
	r, err := http.Get(url)
	defer r.Body.Close()
	check(err)
	c, err := ioutil.ReadAll(r.Body)
	check(err)
	return c
}

func buildYearBaseUrl(year int, page int) string {
	url := root_url
	if year > 0 {
		url = fmt.Sprintf("%ss/%d", url, year)
	}

	if page > 0 {
		url = fmt.Sprintf("%s?page=%d", url, page)
	}
	fmt.Printf("year base year: %d, url : %s\n", year, url)
	return url
}

func remove(origin []byte, sub []byte) []byte {
	return bytes.Replace(origin, sub, bReplEmp, -1)
}

type Error interface {
	Error() string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkDir(path string, clearContent bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
	} else if clearContent {
		files, err := ioutil.ReadDir(path)
		check(err)

		for _, f := range files {
			os.RemoveAll(fmt.Sprintf("%s/%s", path, f.Name))
		}
	}
}
