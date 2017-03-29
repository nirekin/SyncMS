package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func traverseLastPage(n *html.Node) int {
	if n.Type == html.ElementNode && n.Data == "li" {
		attrs := n.Attr
		for _, v := range attrs {
			if v.Key == "class" && v.Val == "pager-last last" {
				innerA := n.FirstChild
				attrs = innerA.Attr
				for _, v := range attrs {
					if v.Key == "href" {
						result := bytes.Split([]byte(v.Val), []byte("="))
						if bytes.Contains(result[1], []byte("&")) {
							result = bytes.Split(result[1], []byte("&"))
							i, err := strconv.Atoi(string(result[0]))
							check(err)
							return i
						} else {
							i, err := strconv.Atoi(string(result[1]))
							check(err)
							return i
						}
					}
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		i := traverseLastPage(c)
		if i > -1 {
			return i
		}
	}
	return -1
}

func traverseIncidentLinks(pqgeContent []byte, done chan bool, incidents chan Incident) {
	r := strings.NewReader(string(pqgeContent))
	d := html.NewTokenizer(r)
	var inBody bool
	var inIncident bool
	var inField bool
	var incident Incident
	colCounter := 0
	aCounter := 0
	for {
		tokenType := d.Next()
		if tokenType == html.ErrorToken {
			done <- true
			return
		}
		token := d.Token()
		tokenText := strings.TrimSpace(token.String())

		switch tokenType {
		case html.StartTagToken:
			if token.Data == "tbody" {
				inBody = true
			}
			if token.Data == "tr" {
				inIncident = true
				incident = Incident{}
			}
			if token.Data == "td" {
				inField = true
			}
			if token.Data == "a" {
				if inBody && inIncident && inField && len(tokenText) > 0 {
					switch aCounter {
					case 0:
						attrs := token.Attr
						for _, v := range attrs {
							if v.Key == "href" && strings.HasPrefix(v.Val, "/incident") {
								arr := strings.Split(v.Val, "/")
								incident.IncidentId = arr[len(arr)-1]
								break
							}
						}
						break
					case 1:
						attrs := token.Attr
						for _, v := range attrs {
							if v.Key == "href" {
								incident.Source = v.Val
								break
							}
						}
						break
					}
				}
			}
			break
		case html.TextToken:

			if inBody && inIncident && inField && len(tokenText) > 0 {
				switch colCounter {
				case 0:
					incident.Date = convertDate(tokenText)
					break
				case 1:
					incident.State = tokenText
					break
				case 2:
					incident.City = tokenText
					break
				case 3:
					incident.Direction = tokenText
					break
				case 4:
					i, err := strconv.Atoi(tokenText)
					check(err)
					incident.Killed = i
					break
				case 5:
					i, err := strconv.Atoi(tokenText)
					check(err)
					incident.Injured = i
					setStateCode(&incident)
					incident.Total = incident.Injured + incident.Killed
					break
				}
			}
			break
		case html.EndTagToken:
			if token.Data == "tbody" {
				inBody = false
			}
			if token.Data == "tr" {
				inIncident = false
				colCounter = 0
				aCounter = 0
				incident.ProcessDate = time.Now()
				incidents <- incident
			}
			if token.Data == "td" {
				inField = false
				colCounter++
			}
			if token.Data == "a" {
				aCounter++
			}
		}
	}
	done <- true
}

func convertDate(date string) time.Time {
	layout := "January 2, 2006"
	t, err := time.Parse(layout, date)
	check(err)
	return t
}

func parseLocation(incident *DonwLoadedIncident) {
	r := strings.NewReader(string(incident.content))
	d := html.NewTokenizer(r)
	var inSpan bool

	for {
		tokenType := d.Next()
		if tokenType == html.ErrorToken {
			return
		}
		token := d.Token()
		tokenText := strings.TrimSpace(token.String())

		switch tokenType {
		case html.StartTagToken:
			if token.Data == "span" {
				inSpan = true
			}
			break
		case html.TextToken:
			if inSpan && len(tokenText) > 0 {
				if strings.HasPrefix(tokenText, "Geolocation") {
					fmt.Printf("Geolocation: %s\n", tokenText)
					latLon := readLocation([]byte(tokenText))
					incident.incident.Latitude = latLon[0]
					incident.incident.Longitude = latLon[1]
				}
			}
			break
		case html.EndTagToken:
			if token.Data == "span" {
				inSpan = false
			}
		}
	}
}

func readLocation(s []byte) []string {
	result := bytes.Split(s, []byte(":"))
	result = bytes.Split(result[1], []byte(","))
	if len(result) == 2 {
		return []string{string(bytes.TrimSpace(result[0])), string(bytes.TrimSpace(result[1]))}
	} else {
		return []string{"", ""}
	}
	return []string{"", ""}
}
