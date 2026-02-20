package categories

import (
	"encoding/json"
	"fmt"
	"testovoe/client"
	"testovoe/internal/models"
)

func GetCategories(showcaseID, token, cookies string) ([]models.Category, error) {

	url := fmt.Sprintf(
		"https://api-web.samokat.ru/v2/showcases/%s/categories/list",
		showcaseID,
	)

	body, err := client.DoRequest(url, token, cookies)
	if err != nil {
		return nil, err
	}

	var cats models.CategoriesResponse
	if err := json.Unmarshal(body, &cats); err != nil {
		return nil, err
	}

	return cats, nil
}
