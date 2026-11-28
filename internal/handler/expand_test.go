package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeExpandShortener struct {
	code string
	url  string
}

func (s *fakeExpandShortener) Shorten(_ string) (string, error) {
	return "", nil
}

func (s *fakeExpandShortener) GetURL(code string) (string, bool) {
	if s.code == code {
		return s.url, true
	}
	return "", false
}

func TestHandler_Expand(t *testing.T) {
	type fields struct {
		url  string
		code string
	}

	type want struct {
		statusCode int
		headers    map[string]string
	}

	tests := []struct {
		name   string
		fields fields
		code   string
		want   want
	}{
		{
			name: "positive expand",
			fields: fields{
				url:  "https://practicum.yandex.ru/",
				code: "Vhbd1sr",
			},
			code: "Vhbd1sr",
			want: want{
				statusCode: http.StatusTemporaryRedirect,
				headers: map[string]string{
					"Location": "https://practicum.yandex.ru/",
				},
			},
		},
		{
			name: "not found expand",
			fields: fields{
				url:  "https://practicum.yandex.ru/",
				code: "Vhbd1sr",
			},
			code: "Vhbd2sr",
			want: want{
				statusCode: http.StatusNotFound,
				headers: map[string]string{
					"Location": "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&fakeExpandShortener{url: tt.fields.url, code: tt.fields.code}, host)
			request := httptest.NewRequest(http.MethodGet, "/"+tt.code, nil)
			request.SetPathValue("id", tt.code)
			w := httptest.NewRecorder()
			handler.Expand(w, request)
			result := w.Result()

			defer func() {
				err := result.Body.Close()
				assert.NoError(t, err)
			}()

			for k, v := range tt.want.headers {
				assert.Equal(t, v, result.Header.Get(k))
			}

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
		})
	}
}
