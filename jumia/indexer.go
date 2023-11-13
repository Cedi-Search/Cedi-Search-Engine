package jumia

import (
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
	"github.com/owbird/cedisearch/database"
	"github.com/owbird/cedisearch/models"
)

type Indexer interface{}

type IndexerImpl struct {
	db *database.Database
}

func NewIndexer(database *database.Database) *IndexerImpl {
	return &IndexerImpl{
		db: database,
	}
}

func (il *IndexerImpl) Index(wg *sync.WaitGroup) {
	log.Println("[+] Indexing Jumia...")

	pages := il.db.GetCrawledPages("Jumia")

	if len(pages) == 0 {
		log.Println("[+] No pages to index for Jumia!")
		wg.Done()
	}

	for _, page := range pages {
		parsedPage := soup.HTMLParse(page.HTML)

		productName := parsedPage.Find("h1").Text()

		productPriceStirng := parsedPage.Find("span", "class", "-prxs").Text()

		priceParts := strings.Split(productPriceStirng, " ")[1]

		price, err := strconv.ParseFloat(strings.ReplaceAll(priceParts, ",", ""), 64)

		if err != nil {
			log.Fatalln(err, page.URL)
		}

		productRatingText := parsedPage.Find("div", "class", "stars").Text()

		productRatingString := strings.Split(productRatingText, " ")[0]

		rating, err := strconv.ParseFloat(productRatingString, 64)

		if err != nil {
			log.Fatalln(err)
		}

		productDescriptionEl := parsedPage.Find("div", "class", "-mhm")

		productDescription := ""

		if productDescriptionEl.Error == nil {
			productDescription = parsedPage.Find("div", "class", "-mhm").FullText()
		}

		productIDText := parsedPage.Find("li", "class", "-pvxs").FullText()

		productID := strings.Split(productIDText, " ")[1]

		productImagesEl := parsedPage.FindAll("img", "class", "-fw")

		productImages := []string{}

		for _, el := range productImagesEl {
			productImages = append(productImages, el.Attrs()["data-src"])
		}

		productData := models.Product{
			Name:        productName,
			Price:       price,
			Rating:      rating,
			Description: productDescription,
			URL:         page.URL,
			Source:      page.Source,
			ProductID:   productID,
			Images:      productImages,
		}

		il.db.IndexProduct(productData)
		il.db.DeleteFromCrawledPages(page)

	}

	il.Index(wg)

}