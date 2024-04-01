package htmlparser

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	links := []Link{}

	doc, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("failed to parse html: %v", err)
	}

	nodes := linkNodes(doc)
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}

	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	ret := []*html.Node{}
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, linkNodes(c)...)
	}
	return ret
}

func buildLink(n *html.Node) Link {
	var link Link
	for _, a := range n.Attr {
		if a.Key == "href" {
			link.Href = a.Val
			link.Text = text(n)
			break
		}
	}
	return link
}

func text(n *html.Node) string {
	var ret string
	if n.Type == html.TextNode {
		return n.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c)
	}
	return strings.Join(strings.Fields(ret), " ")
}
