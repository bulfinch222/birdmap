package taxonomy

import (
	"birdmap/internal"
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Service) getTaxonomy(key string) ([]internal.Bird, error) {
	url := fmt.Sprintf("https://api.ebird.org/v2/ref/taxonomy/ebird?key=%v&fmt=json", key)
	resp, err := http.Get(url)

	if err != nil {
		return  nil, fmt.Errorf("Не удалось получить данные таксономии с EBird: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Ошибка: при запросе к таксономии EBird вернул %s", resp.Status)
	}

	var tax []internal.Bird
	err = json.NewDecoder(resp.Body).Decode(&tax)
	if err != nil {
		return nil, fmt.Errorf("Не удалось расшифровать JSON с таксонометрией %w", err)
	}
	return tax, nil
} 