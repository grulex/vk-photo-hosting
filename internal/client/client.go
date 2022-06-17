package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/grulex/vk-photo-hosting/internal"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	pkgUrl "net/url"
	"strconv"
)

const (
	apiBaseUrl = "https://api.vk.com/method/"
	apiVersion = "5.131"

	getUploadServerMethod = "photos.getUploadServer"
	savePhotoMethod       = "photos.save"
)

type errorResponse struct {
	Error *VkError `json:"error"`
}

type vkApiClient struct {
	httpClient *http.Client
}

func NewClient(httpClient *http.Client) *vkApiClient {
	return &vkApiClient{httpClient: httpClient}
}

func (c vkApiClient) GetUploadServer(ctx context.Context, token string, groupId, albumId uint64) (string, error) {
	url := fmt.Sprintf(
		"%s%s?group_id=%s&album_id=%s&v=%s&access_token=%s",
		apiBaseUrl,
		getUploadServerMethod,
		strconv.FormatInt(int64(groupId), 10),
		strconv.FormatInt(int64(albumId), 10),
		apiVersion,
		token,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}
	body, err := c.execReqWithCheckErr(req)
	if err != nil {
		return "", err
	}
	var response struct {
		Response struct {
			UploadUrl string `json:"upload_url"`
		} `json:"response"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Response.UploadUrl, nil
}

func (c vkApiClient) UploadPhoto(
	ctx context.Context,
	token,
	uploadServer string,
	groupId uint64,
	albumId uint64,
	image io.Reader,
) (
	id uint64,
	photoVariants internal.Variants,
	err error,
) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	fw, err := writer.CreateFormFile("photo", "photo.jpg")
	if err != nil {
		return
	}
	if _, err = io.Copy(fw, image); err != nil {
		return
	}
	err = writer.Close()
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(ctx, "POST", uploadServer, &buffer)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	uploadRespBody, err := c.execReqWithCheckErr(req)
	if err != nil {
		return
	}

	var response struct {
		Server     int64  `json:"server"`
		PhotosList string `json:"photos_list"`
		Hash       string `json:"hash"`
	}
	err = json.Unmarshal(uploadRespBody, &response)
	if err != nil {
		return
	}

	params := pkgUrl.Values{}
	params.Add("server", strconv.FormatInt(response.Server, 10))
	params.Add("photos_list", response.PhotosList)
	params.Add("hash", response.Hash)
	params.Add("album_id", strconv.FormatInt(int64(albumId), 10))
	params.Add("group_id", strconv.FormatInt(int64(groupId), 10))
	encodedParams := params.Encode()

	urlSave := fmt.Sprintf(
		"%s%s?v=%s&access_token=%s&%s",
		apiBaseUrl,
		savePhotoMethod,
		apiVersion,
		token,
		encodedParams,
	)
	reqSave, err := http.NewRequestWithContext(ctx, "GET", urlSave, nil)
	if err != nil {
		return
	}
	bodySave, err := c.execReqWithCheckErr(reqSave)
	if err != nil {
		return
	}

	var responseSave struct {
		Response []*struct {
			Id    uint64            `json:"id"`
			Sizes internal.Variants `json:"sizes"`
		} `json:"response"`
	}

	err = json.Unmarshal(bodySave, &responseSave)
	if err != nil {
		return
	}

	return responseSave.Response[0].Id, responseSave.Response[0].Sizes, nil
}

func (c vkApiClient) execReqWithCheckErr(req *http.Request) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	errResp := &errorResponse{}
	err = json.Unmarshal(body, &errResp)
	if err != nil {
		return nil, errors.Wrapf(err, "error when unmarshal response %v", string(body))
	}

	if errResp.Error != nil {
		return nil, errResp.Error
	}

	return body, nil
}
