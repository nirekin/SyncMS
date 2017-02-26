package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func downloadIncident(id string) []byte {
	defer traceTime("downloadIncident")()
	TraceActivity.Printf("Downloading incident : %s\n", id)
	fmt.Printf("Downloading incident : %s\n", id)
	var content []byte
	content = downloadUrlContent("http://www.gunviolencearchive.org/incident/" + id)

	f, err := os.Create(fmt.Sprintf("%s/incidents_%s.txt", IncidentsDir, id))
	if f != nil {
		defer f.Close()
	}
	check(err)
	_, err = f.Write(content)
	check(err)
	return content
}

func downloadPage(year, page int) []byte {
	defer traceTime("downloadPage")()
	TraceActivity.Printf("Downloading pages for : %d_%d\n", year, page)
	var content []byte
	content = downloadUrlContent(buildYearBasedUrl(year, page))

	if year == 0 {
		year = time.Now().Year()
	}

	f, err := os.Create(fmt.Sprintf("%s/ms_y_%d_p_%d.txt", PagesDir, year, page))
	if f != nil {
		defer f.Close()
	}
	check(err)
	_, err = f.Write(content)
	check(err)
	return content
}

func downloadUrlContent(url string) []byte {
	TraceActivity.Printf("downloadUrlContent url : --%s--\n", url)
	resp, err := http.Get(url)

	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		TraceError.Printf("DownloadUrlContent error reading the url...\n")
	}

	check(err)
	c, err := ioutil.ReadAll(resp.Body)
	check(err)
	return c
}

func buildYearBasedUrl(year int, page int) string {
	url := root_url
	if year > 0 {
		url = fmt.Sprintf("%s?year=%d", url, year)
		if page > 0 {
			url = fmt.Sprintf("%s&page=%d", url, page)
		}
	} else {
		if page > 0 {
			url = fmt.Sprintf("%s?page=%d", url, page)
		}
	}
	TraceActivity.Printf("Looking for base year: %d, url : %s\n", year, url)
	return url
}
