package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"

	"github.com/phuwn/go-crawler-example/crawler"
	"github.com/phuwn/go-crawler-example/db"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

type msg struct {
	ID  int
	Err error
}

func worker(client *http.Client, jobs <-chan int, res chan<- msg) {
	for j := range jobs {
		res <- msg{j, crawler.CrawlPancakeSquad(client, j)}
	}
}

func main() {
	db.Start()

	var (
		numPancakeSquad = 10000
		numWorkers      = 100
		client          = &http.Client{Timeout: 5 * time.Second}
		jobs            = make(chan int, numPancakeSquad)
		results         = make(chan msg, numPancakeSquad)
	)

	for w := 0; w < numWorkers; w++ {
		go worker(client, jobs, results)
	}

	for j := 0; j < numPancakeSquad; j++ {
		jobs <- j
	}

	for a := 0; a < numPancakeSquad; a++ {
		msg := <-results
		if msg.Err != nil {
			fmt.Printf("failed to crawl NFT #%d with error %s\n", msg.ID, msg.Err.Error())
			continue
		}
		fmt.Printf("Crawled NFT #%d\n", msg.ID)
	}
}
