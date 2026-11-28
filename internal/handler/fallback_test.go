package handler

import (
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Fallback(t *testing.T) {
	h := NewHandler(nil, host)
	request := httptest.NewRequest(http.MethodPost, "/unknown/path", nil)
	w := httptest.NewRecorder()
	h.Fallback(w, request)
	result := w.Result()
	defer func() {
		assert.NoError(t, result.Body.Close())
	}()
	assert.Equal(t, http.StatusBadRequest, result.StatusCode)
	contentType, _, err := mime.ParseMediaType(result.Header.Get("Content-Type"))
	require.NoError(t, err)
	assert.Equal(t, "text/plain", contentType)
	body, err := io.ReadAll(result.Body)
	require.NoError(t, err)
	assert.Equal(t, "Wrong request.\n", string(body))
}
