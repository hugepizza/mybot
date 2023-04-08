package gmap

import "testing"

func TestRoute_ShareUrl(t *testing.T) {
	route := Route{
		Name: "test",
		PassingPoints: []*Point{
			{
				Name: "BGH Campo Sioco Terminal",
			}, {
				Name: "LEO Hotel",
			}, {
				Name: "Roxas National High School",
			},
			{
				Name: "Santo Rosario Barangay Hall",
			},
		},
	}
	link := route.ShareUrl()
	if link != "https://www.google.com/maps/dir/?api=1&waypoints=1.000000,1.000000&waypoints=2.000000,2.000000&travelmode=driver" {
		t.Errorf("invalid link: %s", link)
	}
}

func TestPoint_ShareUrl(t *testing.T) {
	point := Point{
		Name: "BGH Campo Sioco Terminal",
	}
	link := point.ShareUrl()
	if link != "https://www.google.com/maps/search/BGH+Campo+Sioco+Terminal" {
		t.Errorf("invalid link: %s", link)
	}
}
