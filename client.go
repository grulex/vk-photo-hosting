package vkphotohosting

import (
	"context"
	"github.com/grulex/vk-photo-hosting/internal"
	"io"
)

type vkApiClient interface {
	GetUploadServer(ctx context.Context, token string, groupId, albumId uint64) (string, error)
	UploadPhoto(
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
	)
}
