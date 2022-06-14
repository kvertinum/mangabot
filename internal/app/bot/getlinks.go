package vkbot

import (
	"github.com/Kvertinum01/mangabot/internal/app/remanga"
)

func linksByPages(pages [][]remanga.ChapterPage) ([]string, [][2]float64, [2]float64) {
	strPages := make([]string, 0)
	sizes := make([][2]float64, 0)
	biggestSizes := [2]float64{100000, 0}
	for _, page := range pages {
		for _, part := range page {
			floatSizes := [2]float64{
				float64(part.Width), float64(part.Height),
			}
			if floatSizes[1] > biggestSizes[1] {
				biggestSizes[1] = floatSizes[1]
			}
			if floatSizes[0] < biggestSizes[0] {
				biggestSizes[0] = floatSizes[0]
			}
			strPages = append(strPages, part.Link)
			sizes = append(sizes, floatSizes)
		}
	}
	return strPages, sizes, biggestSizes
}
