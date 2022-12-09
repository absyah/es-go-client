package domain

type Book struct {
	Id     int    `json:"id"`
	Title  string `json:"title,omitempty"`
	Author string `json:"author,omitempty"`
}

type SearchHits struct {
	Hits struct {
		Hits []*struct {
			Source *Book `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
