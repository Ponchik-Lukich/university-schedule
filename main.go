package main

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

var websites = []string{
	"https://time.is/",
}

var prevContent = []string{
	"",
}

var prevXpath = []string{
	"",
}

var wg sync.WaitGroup

func main() {

	for i, url := range websites {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			result := false
			data, _ := getWebsiteData(url)
			doc, _ := goquery.NewDocumentFromReader(strings.NewReader(data))
			body := doc.Find("*").First()
			removeJunk(body)
			if prevXpath[i] != "" {
				body = getXpathData(body, prevXpath[i])
			}
			data = body.Text()
			hash, _ := getWebsiteHash(data)
			if hash != prevContent[i] {
				prevContent[i] = hash
				result = true
			}

			//content, err := calculateHash(url)
			//fmt.Println(content)
			//if err != nil {
			//	println(err)
			//	return
			//}
			//
			fmt.Println(hash)
			if result {
				fmt.Printf("Website %s changed!\n", url)
			} else {
				fmt.Printf("Nothing on %s\n", url)
			}
			//
			//prevContent[i] = content
		}(i, url)
	}
	wg.Wait()
}

func calculateHash(url string) ([16]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return [16]byte{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return [16]byte{}, err
	}
	fmt.Println(resp.Body)

	hash := md5.Sum(body)

	return hash, nil
}

func getWebsiteHash(data string) (string, error) {
	hasher := sha1.New()
	_, err := hasher.Write([]byte(data))
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func getWebsiteData(url string) (string, error) {
	method := "GET"

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic("Error closing body. Here's why: " + err.Error())
		}
	}(res.Body)
	if res.Status != "200 OK" {
		return "", errors.New("status code is not 200")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func removeJunk(data *goquery.Selection) {
	// get attributes
	//println(data.Attr("class"))
	data = data.RemoveAttr("class")
	//println(data.Attr("class"))
	data = data.RemoveAttr("id")
	data.Contents().Each(func(i int, s *goquery.Selection) {
		if s.Is("script") {
			s.Remove()
		} else if s.Is("style") {
			s.Remove()
		} else {
			removeJunk(s)
		}
	})
}
func getXpathData(body *goquery.Selection, xpath string) *goquery.Selection {
	// get only body
	// remove classes and ids, keep only text recursively
	// remove /html/body from xpath
	xpath = strings.ReplaceAll(xpath, "/html/body", "")

	return body.Find(xpath)
}
