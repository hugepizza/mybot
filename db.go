package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"mybot/gmap"
)

func QueryRoutes(db *sql.DB) ([]*Route, error) {
	rows, err := db.Query("SELECT `index`,name,note,passingPoints,verified FROM routes ORDER BY `index` LIMIT 1000 Offset 0")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []*Route
	for rows.Next() {
		var (
			index                     int
			name, note, passingPoints string
			verified                  int
		)
		if err := rows.Scan(&index, &name, &note, &passingPoints, &verified); err != nil {
			return nil, err
		}
		var points []*gmap.Point
		if err := json.Unmarshal([]byte(passingPoints), &points); err != nil {
			return nil, err
		}
		routes = append(routes, &Route{
			Route: gmap.Route{
				Name:          name,
				PassingPoints: points,
			},
			Index:     index,
			Note:      note,
			Validated: verified >= 5,
		})
	}
	return routes, nil
}

func IncrVerified(db *sql.DB, index int) error {
	_, err := db.Exec("UPDATE routes SET verified = verified + 1 WHERE `index` = ?", index)
	return err
}
