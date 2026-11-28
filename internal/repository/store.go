package repository

import (
	"shortener/internal/model"
	"sync"
)

type Store struct {
	mu          sync.RWMutex
	linksByCode map[string]model.Link
	codesByURL  map[string]string
}

func NewStore() *Store {
	return &Store{
		linksByCode: make(map[string]model.Link),
		codesByURL:  make(map[string]string),
	}
}

func (st *Store) Add(link model.Link) (string, bool) {
	st.mu.Lock()
	defer st.mu.Unlock()
	savedCode, hasCode := st.codesByURL[link.Href]
	if hasCode {
		return savedCode, true
	}
	_, hasLink := st.linksByCode[link.Code]
	if hasLink {
		return "", false
	}
	st.linksByCode[link.Code] = link
	st.codesByURL[link.Href] = link.Code
	return link.Code, true
}

func (st *Store) Get(code string) (string, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	link, ok := st.linksByCode[code]
	return link.Href, ok
}
