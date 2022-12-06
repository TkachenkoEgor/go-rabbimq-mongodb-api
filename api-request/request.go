package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// define the URL and method
	url := "http://localhost:8080/login"
	method := "POST"

	payload := strings.NewReader(`{
    "username": "key",
    "password": "key"
}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	// get the token
	token := string(body)

	url2 := "http://localhost:8080/mydate?firstD=1&secondD=99999999&collectionName=testCollection1"
	method2 := "GET"

	req2, err := http.NewRequest(method2, url2, nil)

	if err != nil {
		fmt.Println(err)
		return

	}
	req2.Header.Add("Token", token)

	res2, err := client.Do(req2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res2.Body.Close()

	body2, err := ioutil.ReadAll(res2.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(body2))
}
