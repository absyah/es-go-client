package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"es-go-client/domain"
)

type ESClient struct {
	baseURL string
}

func NewClient(host string) *ESClient {
	return &ESClient{host}
}

func (c *ESClient) Ping() error {
	response, err := http.Get(c.baseURL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed to ping Elasticsearch: %v", err)
	}

	log.Println("debug ping response: ", string(responseBody))

	return nil
}

func (c *ESClient) InsertIndex(e *domain.Book) error {
	body, _ := json.Marshal(e)

	id := strconv.Itoa(e.Id)
	req, err := http.NewRequest("PUT", c.baseURL+"/book/_doc/"+id, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to init insert index request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to insert index: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed read insert index response: %v", err)
	}

	log.Println("debug insert index response: ", string(responseBody))

	return nil
}

func (c *ESClient) UpdateIndex(e *domain.Book) error {
	body, _ := json.Marshal(map[string]*domain.Book{
		"doc": e,
	})

	id := strconv.Itoa(e.Id)
	req, err := http.NewRequest("POST", c.baseURL+"/book/_update/"+id, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to init update index request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update index: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed read update index response: %v", err)
	}

	log.Println("debug update index response: ", string(responseBody))

	return nil
}

func (c *ESClient) DeleteIndex(id int) error {

	req, err := http.NewRequest("DELETE", c.baseURL+"/book/_doc/"+strconv.Itoa(id), nil)
	if err != nil {
		return fmt.Errorf("failed to make a delete index request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete index: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("failed read delete index response: %v", err)
	}

	log.Println("debug delete index response: ", string(responseBody))

	return nil
}

func (c *ESClient) Search(keyword string) ([]*domain.Book, error) {
	query := fmt.Sprintf(`
	{
		"query": {
			"bool": {
				"should": [
					{"term": { "title": "%s" }},
					{"term": { "author": "%s" }}
				]
			}
		}
	}
	`, keyword, keyword)

	req, err := http.NewRequest("GET", c.baseURL+"/book/_search", strings.NewReader(query))
	if err != nil {
		return nil, fmt.Errorf("failed to init search request: %v", err)
	}

	httpClient := http.Client{}
	req.Header.Add("Content-type", "application/json")
	response, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %v", err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read search response: %v", err)
	}

	var searchHits domain.SearchHits
	if err := json.Unmarshal(responseBody, &searchHits); err != nil {
		return nil, fmt.Errorf("failed to unmarshal search response: %v", err)
	}

	books := []*domain.Book{}
	for _, hit := range searchHits.Hits.Hits {
		books = append(books, hit.Source)
	}

	return books, nil
}
