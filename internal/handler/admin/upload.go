package adminHandler

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/storage"
	u "api-pizza-delivery/internal/utils"
	"errors"
	"net/http"
)

type UploadHandler struct {
	cld *storage.Cloudinary
}

func NewUploadHandler(cld *storage.Cloudinary) *UploadHandler {
	return &UploadHandler{cld: cld}
}

// UploadImage принимает multipart form: поле "file" (изображение до 10 MiB).
// Ответ: { "url": "https://res.cloudinary.com/..." } — подставьте url в img_url при создании/редактировании товара.
func (h *UploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	if h.cld == nil {
		u.JsonError(w, "загрузка изображений не настроена (переменные CLOUDINARY_*)", http.StatusServiceUnavailable)
		return
	}

	if err := r.ParseMultipartForm(storage.MaxImageBytes + (1 << 20)); err != nil {
		u.JsonError(w, "неверная multipart-форма", http.StatusBadRequest)
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		u.JsonError(w, "нужно поле file с файлом", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 0 && header.Size > storage.MaxImageBytes {
		u.JsonError(w, storage.ErrTooLarge.Error(), http.StatusRequestEntityTooLarge)
		return
	}

	url, err := h.cld.UploadImage(r.Context(), file)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrTooLarge):
			u.JsonError(w, err.Error(), http.StatusRequestEntityTooLarge)
		case errors.Is(err, storage.ErrNotAnImage):
			u.JsonError(w, err.Error(), http.StatusBadRequest)
		default:
			u.JsonError(w, "ошибка загрузки файла", http.StatusInternalServerError)
		}
		return
	}

	u.JsonResponse(w, dto.UploadImageResponse{URL: url}, http.StatusOK)
}
