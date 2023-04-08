package gmap

import (
	"fmt"
	"strings"
)

const TravelModeDriver = "DRIVING"

type Route struct {
	Name          string   `json:"name"`
	PassingPoints []*Point `json:"passingPoints"`
}

type Point struct {
	Name string `json:"name"`
}

func (r Route) ShareUrl() string {
	link := "https://www.google.com/maps/dir"
	for _, p := range r.PassingPoints {
		link += fmt.Sprintf("/%s", strings.ReplaceAll(p.Name, " ", "+"))
	}
	//link += "&travelmode=" + TravelModeDriver
	return link
}

func (p Point) ShareUrl() string {
	return fmt.Sprintf("https://www.google.com/maps/search/%s", strings.ReplaceAll(p.Name, " ", "+"))
}
