package models

type Showcase struct {
	ID string `json:"showcaseId"`
	StoreID    string `json:"storeId"`
	Type       string `json:"type"`
	Title      string `json:"title"`
	SLA        int    `json:"sla"`
}
type ShowcaseListResponse struct {
	Showcases []Showcase 
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
	ProductsCount int `json:"productsCount"`
}

type CategoriesResponse []Category