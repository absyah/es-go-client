package delivery

import (
	"encoding/json"
	"es-go-client/domain"
	"es-go-client/repository"
	"net/http"
	"strconv"
)

type Server struct {
	ESClient *repository.ESClient
}

func (s *Server) InsertIndexHandler(w http.ResponseWriter, r *http.Request) {
	var book *domain.Book
	json.NewDecoder(r.Body).Decode(&book)
	err := s.ESClient.InsertIndex(book)
	if err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, book)
}

func (s *Server) UpdateIndexHandler(w http.ResponseWriter, r *http.Request) {
	var book *domain.Book
	json.NewDecoder(r.Body).Decode(&book)
	if err := s.ESClient.UpdateIndex(book); err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, book)
}

func (s *Server) DeleteIndexHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	if err := s.ESClient.DeleteIndex(id); err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, domain.Book{Id: id})
}

func (s *Server) SearchHandler(w http.ResponseWriter, r *http.Request) {
	keyword := r.FormValue("keyword")
	books, err := s.ESClient.Search(keyword)
	if err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, books)
}

func (s *Server) PingHandler(w http.ResponseWriter, r *http.Request) {
	if err := s.ESClient.Ping(); err != nil {
		writeResponseInternalError(w, err)
		return
	}
	writeResponseOK(w, map[string]string{
		"msg": "pong",
	})
}

func writeResponseOK(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeResponse(w, response)
}

func writeResponseInternalError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	writeResponse(w, map[string]interface{}{
		"error": err,
	})
}

func writeResponse(w http.ResponseWriter, response interface{}) {
	json.NewEncoder(w).Encode(response)
}
