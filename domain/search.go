package domain

type Publisher struct {
	Title string `json:"title"`
}

type Author struct {
	Name string `json:"name"`
}

type Book struct {
	ImageName  string    `json:"image_name"`
	Publishers Publisher `json:"publishers"`
	ID         string    `json:"id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Slug       string    `json:"slug"`
	Authors    []Author  `json:"authors"`
}

type SearchResult struct {
	Books []Book `json:"books"`
}
