package vkbot

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/Kvertinum01/mangabot/internal/app/remanga"
	"github.com/jung-kurt/gofpdf"
	"github.com/valyala/fasthttp"
	"github.com/zemlyak-l/vkgottle/bot"
	"github.com/zemlyak-l/vkgottle/object"
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

	bot.OnPrivateMessage("/тайтл", func(message object.NewMessage) {
		if len(message.CmdArgs) < 3 {
			bot.Api.MessagesSend(&object.Message{
				PeerID: message.PeerID,
				Text:   "Используйте команду правильно:\n/тайтл <id> <том> <глава>",
			})
			return
		}
		titleDir := message.CmdArgs[0]
		title := &remanga.TitleInfo{}
		if err := rapi.TitleByDir(titleDir, title); err != nil {
			log.Fatal(err)
		}
		if title.Content == nil {
			bot.Api.MessagesSend(&object.Message{
				PeerID: message.PeerID,
				Text:   "Тайтл не найден",
			})
			return
		}
		tomInt, err := strconv.Atoi(message.CmdArgs[1])
		if err != nil {
			bot.Api.MessagesSend(&object.Message{
				PeerID: message.PeerID,
				Text:   "Том должен быть числом",
			})
			return
		}
		bot.Api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Поиск главы...",
		})
		chapterStr := message.CmdArgs[2]
		branch_id := title.Content.Branches[0].ID
		for nowPage := 1; ; nowPage++ {
			branch := &remanga.BranchInfo{}
			if err := rapi.BranchById(branch_id, nowPage, branch); err != nil {
				log.Fatal(err)
			}
			if branch.Content == nil {
				break
			}
			for _, currChapter := range branch.Content {
				if currChapter.Tome != tomInt || currChapter.Chapter != chapterStr {
					continue
				}

				bot.Api.MessagesSend(&object.Message{
					PeerID: message.PeerID,
					Text:   "Глава найдена! Создание PDF файла.",
				})

				go func(currChapter *remanga.BranchContent, peerID int) {
					chapter := &remanga.ChapterInfo{}
					if err := rapi.ChapterById(currChapter.ID, chapter); err != nil {
						log.Fatal(err)
					}
					links, allSizes, currH := linksByPages(chapter.Content.Pages)

					tp := gofpdf.ImageOptions{ImageType: "jpeg"}

					var heightSum float64

					currW := float64(chapter.Content.Pages[1][1].Width)
					pdf := gofpdf.NewCustom(&gofpdf.InitType{
						OrientationStr: "P",
						UnitStr:        "mm",
						SizeStr:        "A4",
						Size: gofpdf.SizeType{
							Wd: currW,
							Ht: currH,
						},
						FontDirStr: "",
					})

					for ind, imgLink := range links {
						req := fasthttp.AcquireRequest()
						resp := fasthttp.AcquireResponse()

						req.Header.SetMethod("GET")
						req.SetRequestURI(imgLink)

						fasthttpClient.Do(req, resp)
						body := resp.Body()
						reader := bytes.NewReader(body)

						sizes := allSizes[ind]
						ratio := sizes[0] / currW
						resH := sizes[1] / ratio

						if heightSum == 0 {
							heightSum = sizes[1]
						}

						var resY float64

						if resH+heightSum < currH {
							resY = heightSum
							heightSum += resH
						} else {
							heightSum = resH
							pdf.AddPage()
						}

						pdf.RegisterImageOptionsReader(imgLink, tp, reader)
						pdf.Image(imgLink, 0, resY, currW, resH, false, tp.ImageType, 0, "")

						fasthttp.ReleaseRequest(req)
						fasthttp.ReleaseResponse(resp)
					}
					pdfBuffer := &bytes.Buffer{}
					if err := pdf.Output(pdfBuffer); err != nil {
						log.Fatal(err)
					}
					pdfBytes := pdfBuffer.Bytes()
					attachment, err := uploader.docUpload(peerID, pdfBytes)
					if err != nil {
						log.Fatal(err)
					}

					bot.Api.MessagesSend(&object.Message{
						PeerID:     message.PeerID,
						Attachment: attachment,
					})

				}(currChapter, message.PeerID)

				return

			}
		}
		bot.Api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Глава не найдена",
		})
	})

	bot.RunSync()
	return nil
}
