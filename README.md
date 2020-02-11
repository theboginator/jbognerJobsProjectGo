# jbognerJobsProjectGo
Project 1: Reads info from Github Jobs API
To run this project, simply run main.go.
It will create a postings.txt file that lists all the job posting titles retrieved from the Github Jobs API
Sprint 2 Updates:
Running main.go will still print jobs titles to a text file, but now also places the data into an SQLite table.
It is necessary to run "go get github.com/mattn/go-sqlite3" (If running on Windows, this will require TDM-GCC).
Mattn's go-sqlite3 package provides an easy way to interface with sqlite database