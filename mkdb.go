package main

import (
	"bufio"
	"os"
	"strings"
	"sync"

	"github.com/blevesearch/bleve"
	log "github.com/sirupsen/logrus"
)

func isValidType(pType string) bool {
	if (pType == TYPE_IPv4) || (pType == TYPE_IPv6) {
		return true
	} else {
		return false
	}
}

func createIndexedAsnDB(pDB, pType, pFile *string) {
	index, err := bleve.Open(*pDB)
	if err != nil {
		mapping := bleve.NewIndexMapping()
		index, err = bleve.New(*pDB, mapping)

		if err != nil {
			log.Fatalf("Failed to create index at %s error: %+v", *pDB, err)
		}
	}

	file, err := os.Open(*pFile)
	if err != nil {
		log.Fatalf("Failed to open CSV file: %s", *pFile)
	}

	defer file.Close()
	defer index.Close()

	if !isValidType(*pType) {
		log.Fatalf("Invalid type: %s", *pType)
	}

	scanner := bufio.NewScanner(file)
	count := 0

	log.Infof("Indexing records to DB: %s", *pDB)

	var wg sync.WaitGroup
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		record := AsnRecord{ID: parts[1],
			Address:      parts[0],
			Organization: parts[2],
			Type:         *pType}

		wg.Add(1)
		go func(r *AsnRecord) {
			defer wg.Done()
			log.Infof("Indexing record with ID: %s", record.ID)
			if err := index.Index(r.ID, r); err != nil {
				log.Fatalf("Failed to index record. Error: %+v", err)
			}
		}(&record)

		count = count + 1
	}

	wg.Wait()
	log.Infof("Indexed %d entries", count)
}
