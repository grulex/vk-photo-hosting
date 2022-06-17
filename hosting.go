package vkphotohosting

import (
	"bytes"
	"context"
	"github.com/grulex/vk-photo-hosting/internal"
	"github.com/grulex/vk-photo-hosting/internal/client"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type vkPhotoHosting struct {
	vkApiClient vkApiClient
	userToken   string
	groupId     uint64
}

// NewHosting create hosting instance, see README.md for describe params
func NewHosting(userToken string, groupId uint64, httpTimeout time.Duration) *vkPhotoHosting {
	httpClient := &http.Client{Timeout: httpTimeout}
	apiClient := client.NewClient(httpClient)
	return &vkPhotoHosting{
		vkApiClient: apiClient,
		userToken:   userToken,
		groupId:     groupId,
	}
}

func (h vkPhotoHosting) UploadByReader(ctx context.Context, albumId uint64, image io.Reader) (id uint64, variants internal.Variants, err error) {
	server, err := h.vkApiClient.GetUploadServer(ctx, h.userToken, h.groupId, albumId)
	if err != nil {
		return
	}

	return h.vkApiClient.UploadPhoto(ctx, h.userToken, server, h.groupId, albumId, image)
}

func (h vkPhotoHosting) UploadByFile(ctx context.Context, albumId uint64, filePath string) (id uint64, variants internal.Variants, err error) {
	file, err := os.Open(filePath)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("photo", filepath.Base(filePath))
	_, err = io.Copy(part, file)
	err = writer.Close()
	if err != nil {
		return
	}

	server, err := h.vkApiClient.GetUploadServer(ctx, h.userToken, h.groupId, albumId)
	if err != nil {
		return
	}

	return h.vkApiClient.UploadPhoto(ctx, h.userToken, server, h.groupId, albumId, body)
}

func (h vkPhotoHosting) UploadByUrl(
	ctx context.Context,
	albumId uint64,
	photoUrl string,
	downloadTimeout time.Duration,
) (
	id uint64,
	variants internal.Variants,
	err error,
) {
	downloadClient := http.Client{Timeout: downloadTimeout}
	req, _ := http.NewRequest("GET", photoUrl, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.2 Safari/605.1.15")
	req = req.WithContext(ctx)
	resp, err := downloadClient.Do(req)
	if err != nil {
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	server, err := h.vkApiClient.GetUploadServer(ctx, h.userToken, h.groupId, albumId)
	if err != nil {
		return
	}

	return h.vkApiClient.UploadPhoto(ctx, h.userToken, server, h.groupId, albumId, resp.Body)
}
