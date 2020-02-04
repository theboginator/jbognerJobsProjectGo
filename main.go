package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
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

func getJobs(url string) *http.Response { //Get jobs using a provided URL, then return them as *http.response
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Request went horribly wrong: ", err)
		return nil
	}
	client := http.Client{}
	reply, err := client.Do(req)
	if err != nil {
		log.Fatal("Reply went horribly wrong: ", err)
		return nil
	}
	return reply
}

func writePostings(response []posting, outputFile io.Writer) { //Write provided postings array to provided text file
	for i := 0; i < len(response); i++ {
		fmt.Println("Title: ", response[i].Title)
	}
	for i := 0; i < len(response); i++ { //Print each posting and its data count to the text file
		fmt.Fprintln(outputFile, "Title: ", response[i].Title)
		//outputFile.Write(response[index], " : ", response[element])
	}
	fmt.Println("\nAttempted to write results to 'postings.txt'.") //Declare an attempt was made to write the file
}

func main() {
	var response []posting  //Create array of postings
	var data *http.Response //Create variable to hold JSON reply from Github
	var ctr = 1
	outputFile, err := os.Create("postings.txt") //Create a 'postings.txt' file to write our data to
	if err != nil {                              //handle file creation error
		log.Fatal("There was a problem creating the file. ", err)
	}
	defer outputFile.Close() //Create the text file for our answers
	var res = 1
	for res == 1 {
		urlstring := "https://jobs.github.com/positions.json?description=&location=&page=" + strconv.Itoa(ctr)
		ctr++
		url := fmt.Sprintf(urlstring) //setup url
		data = getJobs(url)
		body, err := ioutil.ReadAll(data.Body)
		if err != nil {
			log.Fatal("Reading went horribly wrong: ", err)
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Unmarshal function went horribly wrong: ", err)
		}
		writePostings(response, outputFile)
		cmp := []byte{91, 93} //This is the array that appears when no further data is incoming.
		res = bytes.Compare(body, cmp)
	}
	fmt.Println("At the end.")
}
