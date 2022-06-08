package remanga

type TitleInfo struct {
	Msg     string `json:"msg"`
	Content *struct {
		ID  int `json:"id"`
		Img struct {
			High string `json:"high"`
			Mid  string `json:"mid"`
			Low  string `json:"low"`
		} `json:"img"`
		EnName      string `json:"en_name"`
		RusName     string `json:"rus_name"`
		AnotherName string `json:"another_name"`
		Dir         string `json:"dir"`
		Description string `json:"description"`
		IssueYear   int    `json:"issue_year"`
		AvgRating   string `json:"avg_rating"`
		AdminRating string `json:"admin_rating"`
		CountRating int    `json:"count_rating"`
		AgeLimit    int    `json:"age_limit"`
		Status      struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"status"`
		CountBookmarks int `json:"count_bookmarks"`
		TotalVotes     int `json:"total_votes"`
		TotalViews     int `json:"total_views"`
		Type           struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"type"`
		Genres []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"genres"`
		Categories []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"categories"`
		Publishers []struct {
			ID      int    `json:"id"`
			Name    string `json:"name"`
			Img     string `json:"img"`
			Dir     string `json:"dir"`
			Tagline string `json:"tagline"`
			Type    string `json:"type"`
		} `json:"publishers"`
		BookmarkType interface{} `json:"bookmark_type"`
		Branches     []struct {
			ID         int    `json:"id"`
			Img        string `json:"img"`
			Publishers []struct {
				ID      int    `json:"id"`
				Name    string `json:"name"`
				Img     string `json:"img"`
				Dir     string `json:"dir"`
				Tagline string `json:"tagline"`
				Type    string `json:"type"`
			} `json:"publishers"`
			Subscribed    bool `json:"subscribed"`
			TotalVotes    int  `json:"total_votes"`
			CountChapters int  `json:"count_chapters"`
		} `json:"branches"`
		CountChapters int `json:"count_chapters"`
		FirstChapter  struct {
			ID      int    `json:"id"`
			Tome    int    `json:"tome"`
			Chapter string `json:"chapter"`
		} `json:"first_chapter"`
		ContinueReading interface{} `json:"continue_reading"`
		IsLicensed      bool        `json:"is_licensed"`
		NewlateID       interface{} `json:"newlate_id"`
		NewlateTitle    interface{} `json:"newlate_title"`
		Related         interface{} `json:"related"`
		Uploaded        int         `json:"uploaded"`
		CanPostComments bool        `json:"can_post_comments"`
		Adaptation      interface{} `json:"adaptation"`
	} `json:"content"`
	Props *struct {
		BookmarkTypes []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"bookmark_types"`
		AgeLimit []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"age_limit"`
		CanUploadChapters bool        `json:"can_upload_chapters"`
		CanUpdate         bool        `json:"can_update"`
		CanPinComment     bool        `json:"can_pin_comment"`
		PromoOffer        interface{} `json:"promo_offer"`
		AdminLink         interface{} `json:"admin_link"`
		PanelLink         interface{} `json:"panel_link"`
	} `json:"props"`
}

func (api *RemangaAPI) TitleByDir(dir string, target *TitleInfo) error {
	resUrl := api.Url + "titles/" + dir + "/"
	return api.Get(resUrl, target)
}
