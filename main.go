package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	url := "https://time.is/"
	var prevContent [16]byte

	for {
		content, err := calculateHash(url)
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
