package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/blevesearch/bleve"
	log "github.com/sirupsen/logrus"
)

func isValidType(pType string) bool {
	if (pType == TYPE_IPv4) || (pType == TYPE_IPv6) {
		return true
	}

	return false
}

func getRecordIndexID(r *AsnRecord) string {
	h := sha1.New()
	h.Write([]byte(fmt.Sprintf("%s_%s", r.ID, r.Address)))

	return hex.EncodeToString(h.Sum(nil))
}

func removeQuotes(s string) string {
	return strings.ReplaceAll(s, "\"", "")
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

	// var wg sync.WaitGroup
	// channel := make(chan *AsnRecord)

	// // Start the indexers
	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go func(c chan *AsnRecord, workerId int) {
	// 		defer wg.Done()

	// 		log.Infof("Indexer worker running: %d", workerId)
	// 		for r := range c {
	// 			riID := getRecordIndexID(r)
	// 			if err := index.Index(riID, r); err != nil {
	// 				log.Fatalf("Failed to index record. Error: %+v", err)
	// 			}
	// 		}
	// 	}(channel, i)
	// }

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		record := AsnRecord{ID: parts[1],
			Address:      parts[0],
			Organization: removeQuotes(parts[2]),
			Type:         *pType}

		riID := getRecordIndexID(&record)
		index.Index(riID, record)

		// wg.Add(1)
		// go func(r *AsnRecord) {
		// 	defer wg.Done()

		// 	riID := getRecordIndexID(r)
		// 	log.Infof("Indexing record with ASNID: %s ID: %s", r.ID, riID)

		// 	if err := index.Index(riID, r); err != nil {
		// 		log.Fatalf("Failed to index record. Error: %+v", err)
		// 	}
		// }(&record)

		// channel <- &record
		count = count + 1
	}

	// close(channel)
	// wg.Wait()

	log.Infof("Indexed %d entries", count)
}
