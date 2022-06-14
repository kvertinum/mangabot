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
	"github.com/zemlyak-l/vkcringe/api"
	"github.com/zemlyak-l/vkcringe/bot"
	"github.com/zemlyak-l/vkcringe/object"
)

const (
	helpLink = "https://vk.cc/cemj4R"
)

type routes struct {
	rapi     *remanga.RemangaAPI
	bot      *bot.Bot
	api      *api.Api
	client   *fasthttp.Client
	uploader *docUploader
}

func (r *routes) help(message object.NewMessage) {
	r.api.MessagesSend(&object.Message{
		PeerID: message.PeerID,
		Text:   "Помощь по командам:\n" + helpLink,
	})
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

	resAnswer := "Выберите нужный тайтл и используйте команду /глава <id> <том> <глава>:\n"
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

func (r *routes) chapter(message object.NewMessage) {
	if len(message.CmdArgs) < 3 {
		r.api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Используйте команду правильно:\n/глава <id> <том> <глава>",
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
	branchID := title.Content.Branches[0].ID
	for nowPage := 1; ; nowPage++ {
		branch := &remanga.BranchInfo{}
		if err := r.rapi.BranchById(branchID, nowPage, 300, branch); err != nil {
			log.Fatal(err)
		}
		if branch.Content == nil {
			break
		}
		for _, currChapter := range branch.Content {
			if currChapter.Tome != tomInt || currChapter.Chapter != chapterStr {
				continue
			}

			if currChapter.IsPaid {
				r.api.MessagesSend(&object.Message{
					PeerID: message.PeerID,
					Text:   "Невозможно скачать платную главу!",
				})
				return
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
	links, allSizes, pdfSizes := linksByPages(chapter.Content.Pages)

	tp := gofpdf.ImageOptions{ImageType: "jpeg"}

	var heightSum float64

	currW := pdfSizes[0]
	currH := pdfSizes[1]
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
	docName := fmt.Sprintf("chapter %v-%s", currChapter.Tome, currChapter.Chapter)
	attachment, err := r.uploader.docUpload(peerID, docName, pdfBytes)
	if err != nil {
		log.Fatal(err)
	}

	textAnswer := "Если вы видиле картинку, порезанную посередине - скажите спасибо Remanga"
	r.api.MessagesSend(&object.Message{
		PeerID:     peerID,
		Text:       textAnswer,
		Attachment: attachment,
	})

}

func (r *routes) chapters(message object.NewMessage) {
	if len(message.CmdArgs) == 0 {
		r.api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Укажите ID тайтла",
		})
	}
	titleName := message.CmdArgs[0]
	title := &remanga.TitleInfo{}
	if err := r.rapi.TitleByDir(titleName, title); err != nil {
		log.Fatal(err)
	}
	if title.Content == nil {
		r.api.MessagesSend(&object.Message{
			PeerID: message.PeerID,
			Text:   "Тайтл не найден",
		})
		return
	}
	branchID := title.Content.Branches[0].ID
	branch := &remanga.BranchInfo{}
	if err := r.rapi.BranchById(branchID, 1, 20, branch); err != nil {
		log.Fatal(err)
	}
	answer := "Последние 20 глав:"
	for _, currBranch := range branch.Content {
		answer += fmt.Sprintf(
			"\nТом %v Глава %s",
			currBranch.Tome, currBranch.Chapter,
		)
		if currBranch.Name != "" {
			answer += " - " + currBranch.Name
		}
		if currBranch.IsPaid {
			answer += " (Платная)"
		}
	}
	r.api.MessagesSend(&object.Message{
		PeerID: message.PeerID,
		Text:   answer,
	})
}
