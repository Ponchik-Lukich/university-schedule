package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"time"
)

func main() {
	url := "https://time.is/"
	var prevContent string

	for {
		content, err := fetchContent(url)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if prevContent != content {
			fmt.Println("Website content has changed!")
		}

		if prevContent == content {
			fmt.Println("Nothing has changed!")
		}

		prevContent = content
		time.Sleep(time.Second * 5)
	}
}

func fetchContent(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		return "", err
	}

	// Extract the relevant information from the HTML using goquery selectors
	content := doc.Text()

	return content, nil
}
