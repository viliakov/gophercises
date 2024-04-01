package htmlparser_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/viliakov/gophercises/htmlparser"
)

func TestExtractLinks(t *testing.T) {

	tests := []struct {
		name     string
		htmlFile string
		want     []htmlparser.Link
		wantErr  bool
	}{
		{
			name:     "ex1",
			htmlFile: "test_data/ex1.html",
			want: []htmlparser.Link{
				{
					Href: "/other-page",
					Text: "A link to another page",
				},
			},
			wantErr: false,
		},
		{
			name:     "ex2",
			htmlFile: "test_data/ex2.html",
			want: []htmlparser.Link{
				{
					Href: "https://www.twitter.com/joncalhoun",
					Text: "Check me out on twitter",
				},
				{
					Href: "https://github.com/gophercises",
					Text: "Gophercises is on Github!",
				},
			},
			wantErr: false,
		},
		{
			name:     "ex3",
			htmlFile: "test_data/ex3.html",
			want: []htmlparser.Link{
				{
					Href: "#",
					Text: "Login",
				},
				{
					Href: "/lost",
					Text: "Lost? Need help?",
				},
				{
					Href: "https://twitter.com/marcusolsson",
					Text: "@marcusolsson",
				},
			},
			wantErr: false,
		},
		{
			name:     "ex4",
			htmlFile: "test_data/ex4.html",
			want: []htmlparser.Link{
				{
					Href: "/dog-cat",
					Text: "dog cat",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		file, err := os.Open(tt.htmlFile)
		if err != nil {
			t.Fatalf("failed to open a test file %s: %v", tt.htmlFile, err)
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := htmlparser.Parse(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
