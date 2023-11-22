package utils

import (
	"math/rand"
	"time"

	"github.com/anaskhan96/soup"
)

// ShuffleLinks shuffles the order of links.
func ShuffleLinks(links []soup.Root) {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	rand.Shuffle(len(links), func(i, j int) {
		links[i], links[j] = links[j], links[i]
	})
}