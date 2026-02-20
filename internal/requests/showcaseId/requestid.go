package showcaseid

import (
	"encoding/json"
	"fmt"
	"testovoe/client"
	"testovoe/internal/models"
)

func GetShowcaseID(lat, lon string, token, cookies string) (string, error) {

	url := fmt.Sprintf(
		"https://api-web.samokat.ru/showcases/list?lat=%s&lon=%s",
		lat, lon,
	)

	body, err := client.DoRequest(url, token, cookies)
	if err != nil {
		return "", err
	}

	var data []models.ShowcaseListResponse
	err = json.Unmarshal(body, &data)

	if len(data) == 0 {
		return "", fmt.Errorf("no showcases returned")
	}

	return data[0].Showcases[0].ID, nil
}
