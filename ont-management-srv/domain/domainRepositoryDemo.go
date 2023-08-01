package domain

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/vietbui1502/mqtt/ont-management-srv/dto"
)

type DomainRepositoryDemo struct {
	sexualDomain []string
}

func (d DomainRepositoryDemo) FindDomainCategory(domain string) (int, error) {
	category := dto.Unknown

	// For debugging ontInfo
	log.Printf("Domain query:%s\n", domain)

	for _, item := range d.sexualDomain {
		if item == domain {
			category = dto.Pornography
		}
	}

	return category, nil
}

func NewDomainRepositoryDemo() DomainRepositoryDemo {

	var sexualDomain []string
	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	// Open the file
	file, err := os.Open("./data_demo/domain_blacklist")
	if err != nil {
		log.Printf("Error open data file: %v", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sexualDomain = append(sexualDomain, line)
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		log.Printf("Error loading data file: %v", err)
	}
	log.Printf("Success loading sexual domain data\n")

	return DomainRepositoryDemo{sexualDomain}

}
