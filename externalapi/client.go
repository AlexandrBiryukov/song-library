// externalapi/client.go
package externalapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"song-library/config"
	"song-library/logger"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

func FetchSongDetail(cfg *config.Config, group, song string) (*SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", cfg.APIBaseURL, group, song)
	logger.Log.Debugf("Fetching song details from external API: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch song details: %s", resp.Status)
	}

	var songDetail SongDetail
	err = json.NewDecoder(resp.Body).Decode(&songDetail)
	if err != nil {
		return nil, err
	}

	return &songDetail, nil
}
