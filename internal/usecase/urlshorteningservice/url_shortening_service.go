package url_shortening_service

import (
	"app/internal/entity"
	"app/internal/usecase/random"
)

func NewService(url string, randomKey random.RandomGenerator) entity.ShortURL {
	return entity.ShortURL{
		OriginalURL: url,
		ID:          randomKey.NewRandomString(),
	}
}
