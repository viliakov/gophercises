/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: server,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().String("port", ":8080", "Port to listen on")
	serverCmd.Flags().String("storyfile", "gopher.json", "File with the story")
}

func server(cmd *cobra.Command, args []string) {
	port := cmd.Flags().Lookup("port").Value.String()
	storyFilename := cmd.Flags().Lookup("storyfile").Value.String()

	storyFile, err := os.Open(storyFilename)
	if err != nil {
		log.Fatalf("can't open file (%s) with the story: %v\n", storyFilename, err)
	}

	// Unmarshal JSON
	story := make(map[string]storyArc)
	dec := json.NewDecoder(storyFile)
	err = dec.Decode(&story)
	if err != nil {
		log.Fatalf("can't parse JSON: %v", err)
	}

	for arcName, arcDetails := range story {
		var uri string
		if arcName == "intro" {
			uri = "/"
		} else {
			uri = "/" + arcName
		}
		log.Printf("Adding %s handler with title %q\n", uri, arcDetails.Title)
		http.Handle(uri, &arcHandler{
			name: arcName,
			arc:  arcDetails,
		})
	}

	log.Printf("Starting server on %s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

type storyArc struct {
	Title   string   `json:"title"`
	Story   []string `json:"story"`
	Options []struct {
		Text string `json:"text"`
		Arc  string `json:"arc"`
	} `json:"options,omitempty"`
}

type arcHandler struct {
	name string
	arc  storyArc
}

func (s *arcHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s\n", r.Method, r.URL.Path)
	renderedPage, err := renderStoryHTML(s.name, s.arc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("can't get rendered page for %s: %v", s.name, err)
		return
	}

	_, err = io.WriteString(w, renderedPage)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("failed to send page: %v", err)
		return
	}
}

func renderStoryHTML(title string, s storyArc) (string, error) {
	const tpl = `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="utf-8">
		<title>Choose Your Own Adventure</title>
	</head>
	<body>
		<section class="page">
			<h1>{{.Title}}</h1>
				{{range .Story}}<p>{{ . }}</p>{{end}}
			<ul>
				{{range .Options}}<li><a href="./{{ .Arc }}">{{ .Text }}</a></li>{{end}}
			</ul>
		</section>
	</body>
</html>`

	tmpl, err := template.New(title).Parse(tpl)
	if err != nil {
		return "", err
	}

	var renderedPage strings.Builder
	err = tmpl.Execute(&renderedPage, s)
	if err != nil {
		return "", err
	}

	return renderedPage.String(), nil
}
