package data

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

//Define sexual domain list that loading from file, just for demo
var SexualDomain []string

func Init() {
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	// Open the file
	file, err := os.Open("./data/domain_blacklist")
	if err != nil {
		log.Printf("Error open data file: %v", err)
		return
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		SexualDomain = append(SexualDomain, line)
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		log.Printf("Error loading data file: %v", err)
		return
	}
	log.Printf("Success loading sexual domain data\n")
}
