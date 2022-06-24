package vkbot

import (
	"github.com/Kvertinum01/mangabot/internal/app/remanga"
)

type PageLinks struct {
	strPages  []string
	sizes     [][2]float64
	currSizes [2]float64
}

func linksByPages(pages [][]remanga.ChapterPage) *PageLinks {
	linksObj := &PageLinks{
		strPages:  make([]string, 0),
		sizes:     make([][2]float64, 0),
		currSizes: [2]float64{100000, 0},
	}
	for _, page := range pages {
		for _, part := range page {
			floatSizes := [2]float64{
				float64(part.Width), float64(part.Height),
			}
			if floatSizes[1] > linksObj.currSizes[1] {
				linksObj.currSizes[1] = floatSizes[1]
			}
			if floatSizes[0] < linksObj.currSizes[0] {
				linksObj.currSizes[0] = floatSizes[0]
			}
			linksObj.strPages = append(linksObj.strPages, part.Link)
			linksObj.sizes = append(linksObj.sizes, floatSizes)
		}
	}
	return linksObj
}
