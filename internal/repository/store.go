package repository

import (
	"shortener/internal/model"
	"sync"
)

type Store struct {
	mu          sync.RWMutex
	linksByCode map[string]model.Link
	codesByUrl  map[string]string
}

func NewStore() *Store {
	return &Store{
		linksByCode: make(map[string]model.Link),
		codesByUrl:  make(map[string]string),
	}
}

func (st *Store) Add(new model.Link) (addedCode string, added bool) {
	st.mu.Lock()
	defer st.mu.Unlock()
	savedCode, hasCode := st.codesByUrl[new.Href]
	if hasCode {
		return savedCode, false
	}
	_, hasLink := st.linksByCode[new.Code]
	if hasLink {
		return "", false
	}
	st.linksByCode[new.Code] = new
	st.codesByUrl[new.Href] = new.Code
	return new.Code, true
}

func (st *Store) Get(code string) (string, bool) {
	st.mu.RLock()
	defer st.mu.RUnlock()
	link, ok := st.linksByCode[code]
	return link.Href, ok
}
