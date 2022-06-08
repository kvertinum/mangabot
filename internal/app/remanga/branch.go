package remanga

type BranchContent struct {
	ID         int         `json:"id"`
	Tome       int         `json:"tome"`
	Chapter    string      `json:"chapter"`
	Name       string      `json:"name"`
	Score      int         `json:"score"`
	Rated      interface{} `json:"rated"`
	Viewed     interface{} `json:"viewed"`
	UploadDate string      `json:"upload_date"`
	IsPaid     bool        `json:"is_paid"`
	IsBought   interface{} `json:"is_bought"`
	Price      string      `json:"price"`
	PubDate    string      `json:"pub_date"`
	Publishers []struct {
		Name string `json:"name"`
		Dir  string `json:"dir"`
	} `json:"publishers"`
	Index    int         `json:"index"`
	VolumeID interface{} `json:"volume_id"`
}

type BranchInfo struct {
	Msg     string           `json:"msg"`
	Content []*BranchContent `json:"content"`
	Props   *struct {
		Page     int `json:"page"`
		BranchID int `json:"branch_id"`
	} `json:"props"`
}

type BranchValue struct {
	BranchID int `schema:"branch_id"`
	Page     int `schema:"page"`
	Count    int `schema:"count"`
}

func (api *RemangaAPI) BranchById(branchID int, page int, target *BranchInfo) error {
	return api.Request(
		"titles/chapters", BranchValue{
			BranchID: branchID,
			Page:     page,
			Count:    300,
		}, target,
	)
}
