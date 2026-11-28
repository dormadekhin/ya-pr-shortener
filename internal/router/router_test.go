package router

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRouterHandler struct {
}

func (f *fakeRouterHandler) Create(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Create"))
}

func (f *fakeRouterHandler) Expand(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	w.Write([]byte("Expand " + id))
}

func (f *fakeRouterHandler) Fallback(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Fallback"))
}

func TestNewRouter(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		request string
		want    string
	}{
		{
			name:    "POST root to Create",
			method:  http.MethodPost,
			request: "/",
			want:    "Create",
		},
		{
			name:    "GET root to Fallback",
			method:  http.MethodGet,
			request: "/",
			want:    "Fallback",
		},
		{
			name:    "GET id to Expand",
			method:  http.MethodGet,
			request: "/Vhgt",
			want:    "Expand Vhgt",
		},
		{
			name:    "POST id to Fallback",
			method:  http.MethodPost,
			request: "/Vhgt",
			want:    "Fallback",
		},
		{
			name:    "GET nested path to Fallback",
			method:  http.MethodGet,
			request: "/a/b",
			want:    "Fallback",
		},
		{
			name:    "POST nested path to Fallback",
			method:  http.MethodPost,
			request: "/a/b",
			want:    "Fallback",
		},
		{
			name:    "DELETE nested path to Fallback",
			method:  http.MethodDelete,
			request: "/a/b",
			want:    "Fallback",
		},
		{
			name:    "DELETE root to Fallback",
			method:  http.MethodDelete,
			request: "/",
			want:    "Fallback",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := NewRouter(&fakeRouterHandler{})
			request := httptest.NewRequest(tt.method, tt.request, nil)
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, request)
			result := w.Result()
			defer func() {
				err := result.Body.Close()
				assert.NoError(t, err)
			}()
			bodyBytes, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want, string(bodyBytes))
		})
	}
}
