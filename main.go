package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"time"
)

var websites = []string{
	"https://time.is/",
	"https://home.potatohd.ru/departments/2266/exams",
}

func main() {
	// create [16]byte prevContent array for each website
	prevContent := make([][16]byte, len(websites))
	for i := 0; i < len(websites); i++ {
		prevContent = append(prevContent, [16]byte{})
	}

	// check each website in a loop parallel
	for i, url := range websites {
		go func(i int, url string) {
			for {
				content, err := calculateHash(url)
				if err != nil {
					println(err)
					continue
				}

				if prevContent[i] != content {
					fmt.Printf("Website %s changed!\n", url)
				} else {
					fmt.Printf("Nothing on %s\n", url)
				}

				prevContent[i] = content
				time.Sleep(time.Minute)
			}
		}(i, url)
	}
	//time.Sleep(time.Minute)
	select {}
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

	hash := md5.Sum(body)

	return hash, nil
}
