package remanga

import (
	"encoding/json"
	"errors"
	"strconv"
)

var (
	unkErr = errors.New("unk err")
)

type ChapterPage struct {
	ID            int    `json:"id"`
	Link          string `json:"link"`
	Page          int    `json:"page"`
	Height        int    `json:"height"`
	Width         int    `json:"width"`
	CountComments int    `json:"count_comments"`
}

type ChapterContent struct {
	ID         int         `json:"id"`
	Tome       int         `json:"tome"`
	Chapter    string      `json:"chapter"`
	Name       string      `json:"name"`
	Score      int         `json:"score"`
	ViewID     interface{} `json:"view_id"`
	UploadDate string      `json:"upload_date"`
	IsPaid     bool        `json:"is_paid"`
	TitleID    int         `json:"title_id"`
	VolumeID   interface{} `json:"volume_id"`
	BranchID   int         `json:"branch_id"`
	Price      interface{} `json:"price"`
	PubDate    interface{} `json:"pub_date"`
	Publishers []struct {
		ID             int         `json:"id"`
		Name           string      `json:"name"`
		Dir            string      `json:"dir"`
		ShowDonate     bool        `json:"show_donate"`
		DonatePageText interface{} `json:"donate_page_text"`
		Img            struct {
			High string `json:"high"`
			Mid  string `json:"mid"`
			Low  string `json:"low"`
		} `json:"img"`
		PaidSubscription struct {
			Name        string      `json:"name"`
			Description string      `json:"description"`
			Price       string      `json:"price"`
			Publisher   int         `json:"publisher"`
			User        interface{} `json:"user"`
		} `json:"paid_subscription"`
	} `json:"publishers"`
	Index int             `json:"index"`
	Pages [][]ChapterPage `json:"pages"`
}

type PreChapterContent ChapterContent

type PreChapterInfo struct {
	ID    int             `json:"id"`
	Pages json.RawMessage `json:"pages"`
}

func (c *ChapterContent) UnmarshalJSON(data []byte) error {
	pc := &PreChapterContent{}
	if err := json.Unmarshal(data, pc); err != nil {
		prePage := &struct {
			Pages []ChapterPage `json:"pages"`
		}{}
		if err := json.Unmarshal(data, prePage); err != nil {
			return err
		}
		allPages := make([][]ChapterPage, 1)
		pc.Pages = append(allPages, prePage.Pages)
	}
	*c = ChapterContent(*pc)

	return nil
}

type ChapterInfo struct {
	Msg     string          `json:"msg"`
	Content *ChapterContent `json:"content"`
	Props   *struct{}       `json:"props"`
}

func (api *RemangaAPI) ChapterById(chapterID int, target *ChapterInfo) error {
	chapterStr := strconv.Itoa(chapterID)
	resUrl := api.Url + "titles/chapters/" + chapterStr + "/"

	return api.Get(resUrl, target)
}
