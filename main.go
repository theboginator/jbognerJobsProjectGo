package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gen2brain/dlgs"
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

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Text    string   `xml:",chardata"`
	A10     string   `xml:"a10,attr"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Text        string `xml:",chardata"`
		Os          string `xml:"os,attr"`
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Image       struct {
			Text  string `xml:",chardata"`
			URL   string `xml:"url"`
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"image"`
		TotalResults string `xml:"totalResults"`
		Item         []struct {
			Text string `xml:",chardata"`
			Guid struct {
				Text        string `xml:",chardata"`
				IsPermaLink string `xml:"isPermaLink,attr"`
			} `xml:"guid"`
			Link   string `xml:"link"`
			Author struct {
				Text string `xml:",chardata"`
				Name string `xml:"name"`
			} `xml:"author"`
			Category    []string `xml:"category"`
			Title       string   `xml:"title"`
			Description string   `xml:"description"`
			PubDate     string   `xml:"pubDate"`
			Updated     string   `xml:"updated"`
			Location    struct {
				Text  string `xml:",chardata"`
				Xmlns string `xml:"xmlns,attr"`
			} `xml:"location"`
		} `xml:"item"`
	} `xml:"channel"`
}

func setup_database() *sql.DB { //Create the database
	database, _ := sql.Open("sqlite3", "./jobsdata.db")
	statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS jobsdata (id INTEGER PRIMARY KEY, fulltime TEXT, url TEXT, created TEXT, company TEXT, website TEXT, location TEXT, lat REAL, long REAL, description TEXT)")
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

func getContent(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
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

func insert_stackposting(database *sql.DB, job RSS) { //Insert a job into the database
	for i := range job.Channel.Item { //Print each posting and its data count to the text file
		statement, _ := database.Prepare("INSERT INTO jobsdata (company, url, created, location, description) VALUES (?, ?, ?, ?, ?)")
		//TODO: Sanitize inputs before insertion
		_, err := statement.Exec(string(job.Channel.Item[i].Author.Name), string(job.Channel.Item[i].Link), string(job.Channel.Item[i].PubDate), string(job.Channel.Item[i].Location.Text), string(job.Channel.Item[i].Description))
		if err != nil {
			fmt.Errorf("Encountered error: %v", err)
		}
	}
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
	fmt.Println("Time for stackoverflow jobs.")
	stackdata, err := getContent("https://stackoverflow.com/jobs/feed")
	if err != nil {
		log.Printf("Failed to get XML: %v", err)
	}
	fmt.Println("Received XML:")
	var posting RSS
	err = xml.Unmarshal(stackdata, &posting)
	if err != nil {
		log.Printf("Big OOF while unmarshal is happening: %v", err)
	}
	for index := range posting.Channel.Item {
		fmt.Println("Title: ", posting.Channel.Item[index].Title)
	}
	//jobsdb := setup_database() //Set up the jobs database
	insert_stackposting(jobsdb, posting) //Write the responses for this page to the database
	//TODO: Geocode DB entries
	//TODO: Get list of all cities/states that exists db into array
	_, _, err = dlgs.List("List", "Select item from list:", []string{"Bug", "New Feature", "Improvement"}) //This WILL get the list of available cities/countries for job narrowing
	if err != nil {
		panic(err)
	}
	//TODO: Take items returned from selection and then pull from database
	//TODO: Take selected items and generate a .csv
	//TODO: Call python code that will make plot from .csv
}
