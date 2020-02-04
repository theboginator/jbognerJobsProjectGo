package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type posting struct {
	Id          string `json:"id"`
	FullTime    string `json:"type"`
	Url         string `json:"url"`
	CreatedAt   string `json:"created_at"`
	Company     string `json:"company"`
	Company_url string `json:"company_url"`
	Location    string `json:"location"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Logo        string `json:"company_logo"`
}

func main() {
	fmt.Println("Hello world!")
	url := fmt.Sprintf("https://jobs.github.com/positions.json?")
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Request went horribly wrong: ", err)
		return
	}

	client := http.Client{}
	reply, err := client.Do(req)
	if err != nil {
		log.Fatal("Reply went horribly wrong: ", err)
		return
	}

	var response []posting
	/*
		if err := json.NewDecoder(reply.Body).Decode(&response); err != nil {
			log.Println("Decoding went horribly wrong: ", err)
		}
	*/
	body, err := ioutil.ReadAll(reply.Body)

	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("error:", err)
	}

	for i := 0; i < len(response); i++ {
		fmt.Println("Title: ", response[i].Title)
	}

	fmt.Println("At the end.")

}
