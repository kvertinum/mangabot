package remanga

type SearchValue struct {
	Query string `schema:"query"`
	Count int    `schema:"count"`
}

type SearchAnswer struct {
	Msg     string `json:"msg"`
	Content []*struct {
		ID           int         `json:"id"`
		EnName       string      `json:"en_name"`
		RusName      string      `json:"rus_name"`
		Dir          string      `json:"dir"`
		BookmarkType interface{} `json:"bookmark_type"`
		Img          struct {
			High string `json:"high"`
			Mid  string `json:"mid"`
			Low  string `json:"low"`
		} `json:"img"`
		IssueYear     int    `json:"issue_year"`
		AvgRating     string `json:"avg_rating"`
		Type          int    `json:"type"`
		CountChapters int    `json:"count_chapters"`
	} `json:"content"`
	Props *struct {
		TotalItems int `json:"total_items"`
		TotalPages int `json:"total_pages"`
		Page       int `json:"page"`
	} `json:"props"`
}

func (api *RemangaAPI) Search(mangaName string, count int, target *SearchAnswer) error {
	return api.Request("search", SearchValue{
		Query: mangaName,
		Count: count,
	}, target)
}
