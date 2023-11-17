package main

import (
	"log"
	"sync"

	"github.com/Cedi-Search/Cedi-Search-Engine/crawler"
	"github.com/Cedi-Search/Cedi-Search-Engine/data"
	"github.com/Cedi-Search/Cedi-Search-Engine/database"
	"github.com/Cedi-Search/Cedi-Search-Engine/deus"
	"github.com/Cedi-Search/Cedi-Search-Engine/ishtari"
	"github.com/Cedi-Search/Cedi-Search-Engine/jiji"
	"github.com/Cedi-Search/Cedi-Search-Engine/jumia"
	"github.com/joho/godotenv"
)

func main() {

	log.Println("[+] Startup")

	wg := sync.WaitGroup{}

	godotenv.Load()

	database := database.NewDatabase()

	database.Init()

	sniffers := []data.Sniffer{
		jumia.NewSniffer(database),
		jiji.NewSniffer(database),
		deus.NewSniffer(database),
		ishtari.NewSniffer(database),
	}

	indexers := []data.Indexer{
		jumia.NewIndexer(database),
		jiji.NewIndexer(database),
		deus.NewIndexer(database),
		ishtari.NewIndexer(database),
	}

	wg.Add(len(sniffers))
	for _, sniffer := range sniffers {
		go sniffer.Sniff(&wg)
	}

	wg.Add(len(indexers))
	for _, indexer := range indexers {
		go indexer.Index(&wg)
	}

	crawler := crawler.NewCrawler(database)
	crawler.Crawl()

	wg.Wait()
}
