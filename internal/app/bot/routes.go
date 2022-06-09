package vkbot

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Kvertinum01/mangabot/internal/app/remanga"
	"github.com/jung-kurt/gofpdf"
	"github.com/valyala/fasthttp"
	"github.com/zemlyak-l/vkgottle/api"
	"github.com/zemlyak-l/vkgottle/bot"
	"github.com/zemlyak-l/vkgottle/object"
)

type routes struct {
	rapi     *remanga.RemangaAPI
	bot      *bot.Bot
	api      *api.Api
	client   *fasthttp.Client
	uploader *docUploader
}

func (r *routes) search(message object.NewMessage) {
	if message.CmdArgs == nil {
		r.api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Передайте название тайтла",
		})
		return
	}

	mangaName := strings.Join(message.CmdArgs, " ")
	answer := &remanga.SearchAnswer{}
	if err := r.rapi.Search(mangaName, searchCount, answer); err != nil {
		log.Fatal(err)
	}

	resAnswer := "Выберите нужный тайтл и используйте команду /тайтл <id> <том> <глава>:\n"
	for _, value := range answer.Content {
		resAnswer += fmt.Sprintf(
			"• %s (ID %v)\n", value.RusName, value.Dir,
		)
	}

	r.api.MessagesSend(&object.Message{
		PeerID: message.PeerID,
		Text:   resAnswer,
	})
}

func (r *routes) title(message object.NewMessage) {
	if len(message.CmdArgs) < 3 {
		r.api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Используйте команду правильно:\n/тайтл <id> <том> <глава>",
		})
		return
	}
	titleDir := message.CmdArgs[0]
	title := &remanga.TitleInfo{}
	if err := r.rapi.TitleByDir(titleDir, title); err != nil {
		log.Fatal(err)
	}
	if title.Content == nil {
		r.api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Тайтл не найден",
		})
		return
	}
	tomInt, err := strconv.Atoi(message.CmdArgs[1])
	if err != nil {
		r.api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Том должен быть числом",
		})
		return
	}
	r.api.MessagesSend(&object.Message{
		PeerID: message.PeerID,
		Text:   "Поиск главы...",
	})
	chapterStr := message.CmdArgs[2]
	branch_id := title.Content.Branches[0].ID
	for nowPage := 1; ; nowPage++ {
		branch := &remanga.BranchInfo{}
		if err := r.rapi.BranchById(branch_id, nowPage, branch); err != nil {
			log.Fatal(err)
		}
		if branch.Content == nil {
			break
		}
		for _, currChapter := range branch.Content {
			if currChapter.Tome != tomInt || currChapter.Chapter != chapterStr {
				continue
			}

			r.api.MessagesSend(&object.Message{
				PeerID: message.PeerID,
				Text:   "Глава найдена! Создание PDF файла.",
			})

			go r.createAndSend(currChapter, message.PeerID)

			return

		}
	}
	r.api.MessagesSend(&object.Message{
		PeerID: message.PeerID,
		Text:   "Глава не найдена",
	})
}

func (r *routes) createAndSend(currChapter *remanga.BranchContent, peerID int) {
	chapter := &remanga.ChapterInfo{}
	if err := r.rapi.ChapterById(currChapter.ID, chapter); err != nil {
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

		r.client.Do(req, resp)
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
	attachment, err := r.uploader.docUpload(peerID, pdfBytes)
	if err != nil {
		log.Fatal(err)
	}

	r.api.MessagesSend(&object.Message{
		PeerID:     peerID,
		Attachment: attachment,
	})

}
