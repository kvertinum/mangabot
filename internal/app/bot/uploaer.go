package vkbot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"

	"github.com/valyala/fasthttp"
	"github.com/zemlyak-l/vkgottle/api"
)

type docsSaveRes struct {
	Response *struct {
		Document struct {
			OwnerID int `json:"owner_id"`
			ID      int `json:"id"`
		} `json:"doc"`
	} `json:"response"`
}

type docsSaveVal struct {
	File  string `schema:"file"`
	Title string `schema:"title"`
}

type getUploadServerRes struct {
	Response *struct {
		UploadURL string `json:"upload_url"`
	} `json:"response"`
}

type getUploadServerVal struct {
	Type   string `schema:"type"`
	PeerID int    `schema:"peer_id"`
}

type docUploader struct {
	client *fasthttp.Client
	api    *api.Api
}

func (u *docUploader) docUpload(peerID int, docName string, file []byte) (string, error) {
	uploadServerRes := &getUploadServerRes{}
	if err := u.api.Request("docs.getMessagesUploadServer", getUploadServerVal{
		Type: "doc", PeerID: peerID,
	}, uploadServerRes); err != nil {
		return "", err
	}

	uploadedFile := &struct {
		File string `json:"file"`
	}{}

	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	bodyBufer := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBufer)
	fileWriter, err := bodyWriter.CreateFormFile("file", "chapter.pdf")
	fileWriter.Write(file)
	if err != nil {
		return "", err
	}
	bodyWriter.Close()
	contentType := bodyWriter.FormDataContentType()

	req.Header.SetMethod("POST")
	req.Header.SetContentType(contentType)
	req.SetRequestURI(uploadServerRes.Response.UploadURL)
	req.SetBody(bodyBufer.Bytes())

	u.client.Do(req, resp)
	body := resp.Body()

	if err := json.Unmarshal(body, uploadedFile); err != nil {
		return "", err
	}

	finalDocument := &docsSaveRes{}
	if err := u.api.Request("docs.save", docsSaveVal{
		File: uploadedFile.File, Title: docName,
	}, finalDocument); err != nil {
		return "", err
	}
	doc := finalDocument.Response.Document
	answer := fmt.Sprintf("doc%v_%v", doc.OwnerID, doc.ID)

	return answer, nil
}
