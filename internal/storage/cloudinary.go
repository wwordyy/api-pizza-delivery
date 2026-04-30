package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	MaxImageBytes        = 10 << 20 // 10 MiB
	cloudinaryFolder     = "pizza-delivery/products"
)

var allowedImageTypes = map[string]struct{}{
	"image/jpeg": {},
	"image/png":  {},
	"image/webp": {},
	"image/gif":  {},
}

var (
	ErrNotConfigured = errors.New("загрузка изображений не настроена (Cloudinary)")
	ErrNotAnImage    = errors.New("недопустимый тип файла (нужно изображение)")
	ErrTooLarge      = errors.New("файл превышает допустимый размер")
)

type Cloudinary struct {
	cld *cloudinary.Cloudinary
}

func New(cloudName, apiKey, apiSecret string) (*Cloudinary, error) {
	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, ErrNotConfigured
	}
	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("cloudinary init: %w", err)
	}
	return &Cloudinary{cld: cld}, nil
}

func (c *Cloudinary) UploadImage(ctx context.Context, r io.Reader) (string, error) {
	if c == nil || c.cld == nil {
		return "", ErrNotConfigured
	}

	data, err := io.ReadAll(io.LimitReader(r, MaxImageBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > MaxImageBytes {
		return "", ErrTooLarge
	}
	if len(data) == 0 {
		return "", ErrNotAnImage
	}

	sniffN := len(data)
	if sniffN > 512 {
		sniffN = 512
	}
	ct := http.DetectContentType(data[:sniffN])
	if _, ok := allowedImageTypes[ct]; !ok {
		return "", ErrNotAnImage
	}

	t := true
	res, err := c.cld.Upload.Upload(ctx, bytes.NewReader(data), uploader.UploadParams{
		Folder:         cloudinaryFolder,
		ResourceType:   "image",
		UseFilename:    &t,
		UniqueFilename: &t,
	})
	if err != nil {
		return "", err
	}
	if res.SecureURL == "" {
		return "", errors.New("пустой URL в ответе Cloudinary")
	}
	return res.SecureURL, nil
}
