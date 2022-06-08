package vkbot

import (
	"github.com/Kvertinum01/mangabot/internal/app/remanga"
	"github.com/signintech/gopdf"
)

func linksByPages(pages [][]remanga.ChapterPage) ([]string, [][2]float64, float64) {
	currH := 0
	strPages := make([]string, 0)
	sizes := make([][2]float64, 0)
	for _, page := range pages {
		for _, part := range page {
			if part.Height > currH {
				currH = part.Height
			}
			strPages = append(strPages, part.Link)
			sizes = append(sizes, [2]float64{
				float64(part.Width), float64(part.Height),
			})
		}
	}
	return strPages, sizes, float64(currH)
}

func pdfByLinks([]string) gopdf.GoPdf {
	return gopdf.GoPdf{}
}
