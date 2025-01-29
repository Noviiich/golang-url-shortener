package handlers

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Noviiich/golang-url-shortener/domain"
	"github.com/Noviiich/golang-url-shortener/types"
)

type UrlShortenerHandler struct {
	link *domain.Link
}

type RequestBody struct {
	Long string `json:"long"`
}

func NewUrlShortenerHandler(link *domain.Link) *UrlShortenerHandler {
	return &UrlShortenerHandler{link: link}
}

func (h *UrlShortenerHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var requestBody RequestBody
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
	}

	link := types.Link{
		OriginalURL: requestBody.Long,
		ShortId:     generateShortUrl(requestBody.Long),
		CreateAt:    time.Now(),
		Clicks:      0,
	}

	err := h.link.Create(context.Background(), link)
	if err != nil {
		http.Error(w, "Failed to save link", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(link)

}

func (h *UrlShortenerHandler) GetOriginalUrl(w http.ResponseWriter, r *http.Request) {
	shortLink := r.URL.Path[1:]
	longLink, err := h.link.Get(context.Background(), shortLink)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"long": longLink.Long,
	}
	json.NewEncoder(w).Encode(response)

}

func generateShortUrl(longUrl string) string {
	hasher := md5.New()
	hasher.Write([]byte(longUrl))
	shortUrl := hex.EncodeToString(hasher.Sum(nil))
	return shortUrl[:8]
}
