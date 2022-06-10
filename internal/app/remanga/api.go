package remanga

import (
	"encoding/json"
	"net/url"
	"time"

	"github.com/gorilla/schema"
	"github.com/valyala/fasthttp"
)

const (
	baseApiUrl = "https://api.remanga.org/api/"
)

type RemangaAPI struct {
	Url     string
	Encoder *schema.Encoder
	Client  *fasthttp.Client
}

func NewRemangaAPI() *RemangaAPI {
	// Setup RemangaAPI

	return &RemangaAPI{
		Url:     baseApiUrl,
		Encoder: schema.NewEncoder(),
		Client: &fasthttp.Client{
			ReadTimeout:              5 * time.Second,
			WriteTimeout:             5 * time.Second,
			MaxIdleConnDuration:      time.Minute,
			NoDefaultUserAgentHeader: true,
		},
	}
}

func (api *RemangaAPI) Request(methodName string, data interface{}, target interface{}) error {
	// Request to remanga api
	resUrl := api.Url + methodName + "/?"

	urlData := url.Values{}
	if err := api.Encoder.Encode(data, urlData); err != nil {
		return err
	}
	urlEncoded := urlData.Encode()
	resUrl += urlEncoded

	return api.Get(
		resUrl,
		target,
	)
}

func (api *RemangaAPI) Get(url string, target interface{}) error {
	// Create GET request
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod("GET")
	req.SetRequestURI(url)

	api.Client.Do(req, resp)
	body := resp.Body()

	if target == nil || len(body) == 0 {
		return nil
	}

	return json.Unmarshal(body, target)
}
