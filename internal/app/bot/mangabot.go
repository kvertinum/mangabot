package vkbot

import (
	"time"

	"github.com/Kvertinum01/mangabot/internal/app/remanga"
	"github.com/valyala/fasthttp"
	"github.com/zemlyak-l/vkgottle/bot"
)

const (
	searchCount = 5
)

func SetupBot(config *Conifg) error {
	rapi := remanga.NewRemangaAPI()
	bot, err := bot.NewBot(config.Token)
	if err != nil {
		return err
	}

	fasthttpClient := &fasthttp.Client{
		ReadTimeout:              5 * time.Second,
		WriteTimeout:             5 * time.Second,
		MaxIdleConnDuration:      time.Minute,
		NoDefaultUserAgentHeader: true,
		MaxResponseBodySize:      4 * 1024 * 1024 * 1024,
	}

	uploader := &docUploader{
		client: fasthttpClient,
		api:    bot.Api,
	}

	privateRoutes := &routes{
		rapi:     rapi,
		api:      bot.Api,
		bot:      bot,
		client:   fasthttpClient,
		uploader: uploader,
	}

	bot.OnPrivateMessage("/поиск", privateRoutes.search)

	bot.OnPrivateMessage("/тайтл", privateRoutes.title)

	bot.RunSync()
	return nil
}
