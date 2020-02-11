package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type posting struct { //Let's define a struct to hold jobs data
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

func setup_database() *sql.DB { //Create the database
	database, _ := sql.Open("sqlite3", "./jobsdata.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS jobsdata (id INTEGER PRIMARY KEY, fulltime TEXT, url TEXT, created TEXT, company TEXT, website TEXT, location TEXT, description TEXT)")
	statement.Exec()
	return database
}

func insert_posting(database *sql.DB, job []posting) { //Insert a job into the database
	for i := 0; i < len(job); i++ { //Print each posting and its data count to the text file
		statement, _ := database.Prepare("INSERT INTO jobsdata (fulltime, url, created, company, website, location, description) VALUES (?, ?, ?, ?, ?, ?, ?)")
		//TODO: Sanitize inputs before insertion
		statement.Exec(job[i].FullTime, job[i].Url, job[i].CreatedAt, job[i].Company, job[i].Company_url, job[i].Location, job[i].Description)
	}
}

func get_jobs(url string) (*http.Response, error) { //Get jobs using a provided URL, then return them as *http.response
	//I wonder if this should return a response and an error to really take advantage of go?
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("Request went horribly wrong: ", err)
		return nil, err
	}
	client := http.Client{}
	reply, err := client.Do(req)
	if err != nil {
		log.Fatal("Reply went horribly wrong: ", err)
		return nil, err
	}
	return reply, nil
}

func write_postings(response []posting, outputFile io.Writer) { //Write provided postings array to provided text file
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
	var ctr = 1             //This will be used to keep track of the page of responses

	outputFile, err := os.Create("postings.txt") //Create a 'postings.txt' file to write our data to
	if err != nil {                              //handle file creation error
		log.Fatal("There was a problem creating the file. ", err)
	}
	defer outputFile.Close() //Create the text file for our answers, and leave it open for the write function

	jobsdb := setup_database() //Set up the jobs database

	var res = 1 //Track whether we're done or not
	for res == 1 {
		urlstring := "https://jobs.github.com/positions.json?description=&location=&page=" + strconv.Itoa(ctr) //Generate the url with the right page number
		ctr++                                                                                                  //Prep for the next page
		url := fmt.Sprintf(urlstring)                                                                          //setup url
		data, err = get_jobs(url)                                                                              //Retrieve data from the API
		if err != nil {
			log.Fatal("Request went horribly wrong: ", err) //Handle a fault
		}
		body, err := ioutil.ReadAll(data.Body) //Get the JSON from data
		if err != nil {                        //Handle a fault
			log.Fatal("Reading went horribly wrong: ", err)
		}
		err = json.Unmarshal(body, &response) //Translate the JSON into our struct
		if err != nil {                       //Hande a fault
			fmt.Println("Unmarshal function went horribly wrong: ", err)
		}
		insert_posting(jobsdb, response)     //Write the responses for this page to the database
		write_postings(response, outputFile) //Write the response data struct to our text file
		cmp := []byte{91, 93}                //This is the array that appears when no further data is incoming
		res = bytes.Compare(body, cmp)       //If we receive the "no further data" response, we're done
	}
	fmt.Println("At the end.")
}
