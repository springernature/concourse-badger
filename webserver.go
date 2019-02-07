package main

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type Concourse struct {
	PgUsername, PgPassword, PgHost string
}

func (c Concourse) PipelineStatus(team, pipeline string) (status string, err error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/concourse?sslmode=disable", c.PgUsername, c.PgPassword, c.PgHost)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return
	}

	cql := fmt.Sprintf(`SELECT b.status
FROM builds b
INNER JOIN pipelines p ON p.id = b.pipeline_id
INNER JOIN teams t ON t.id = b.team_id
WHERE t.name = '%s' AND p.name = '%s' AND b.completed = TRUE
ORDER BY end_time DESC
LIMIT 1`, team, pipeline)

	err = db.QueryRow(cql).Scan(&status)
	if err != nil {
		return
	}
	err = db.Close()
	return
}

func main() {
	port := os.Getenv("PORT")

	concourse := Concourse{
		PgUsername: os.Getenv("POSTGRES_USERNAME"),
		PgPassword: os.Getenv("POSTGRES_PASSWORD"),
		PgHost:     os.Getenv("POSTGRES_HOST"),
	}

	handler := func(w http.ResponseWriter, r *http.Request) {

		pathSegments := strings.Split(r.URL.Path, "/")
		if len(pathSegments) != 3 {
			w.WriteHeader(404)
			return
		}

		status, err := concourse.PipelineStatus(pathSegments[1], pathSegments[2])
		if err != nil {
			status = "unknown"
		}

		w.Header().Set("Content-type", "image/svg+xml")
		w.Header().Set("Cache-Control", "max-age=60")
		fmt.Fprint(w, badge(status))
	}

	http.HandleFunc("/", handler)
	fmt.Println("Running on http://localhost:" + port)
	http.ListenAndServe(":"+port, nil)
}

func badge(status string) string {

	var color string
	switch status {
	case "succeeded":
		color = `#44cc11`
		status = "passed"
	case "failed":
		color = `#e05d44`
	case "aborted":
		color = `#8f4b2d`
	case "errored":
		color = `#fe7d37`
	default:
		color = `#9f9f9f`
		status = "unknown"
	}

	const svg = `<?xml version="1.0" encoding="UTF-8"?>
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="88" height="20">
  <linearGradient id="b" x2="0" y2="100%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
  <clipPath id="a">
    <rect width="88" height="20" rx="3" fill="#fff"/>
  </clipPath>
  <g clip-path="url(#a)">
    <path fill="#555" d="M0 0h37v20H0z"/>
    <path fill="{{ .FillColor }}" d="M37 0h51v20H37z"/>
    <path fill="url(#b)" d="M0 0h88v20H0z"/>
  </g>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
    <text x="195" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="270">build</text>
    <text x="195" y="140" transform="scale(.1)" textLength="270">build</text>
    <text x="615" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="410">{{ .Status }}</text>
    <text x="615" y="140" transform="scale(.1)" textLength="410">{{ .Status }}</text>
  </g>
</svg>`

	tmpl, err := template.New("Badge").Parse(svg)
	if err != nil {
		panic(err)
	}

	buffer := &bytes.Buffer{}

	err = tmpl.Execute(buffer, struct {
		FillColor, Status string
	}{
		color,
		status,
	})

	if err != nil {
		panic(err)
	}

	return buffer.String()

}
