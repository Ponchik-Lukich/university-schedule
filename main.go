package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"sync"
)

var websites = []string{
	"https://time.is/",
	"https://home.potatohd.ru/departments/2266/exams",
}

var wg sync.WaitGroup

func main() {
	prevContent := make([][16]byte, len(websites))
	for i := 0; i < len(websites); i++ {
		prevContent = append(prevContent, [16]byte{})
	}

	for i, url := range websites {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			content, err := calculateHash(url)
			if err != nil {
				println(err)
				return
			}

			if prevContent[i] != content {
				fmt.Printf("Website %s changed!\n", url)
			} else {
				fmt.Printf("Nothing on %s\n", url)
			}

			prevContent[i] = content
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

	hash := md5.Sum(body)

	return hash, nil
}
