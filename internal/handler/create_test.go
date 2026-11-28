package handler

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"net/http/httptest"
	"shortener/internal/service"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const host = "http://localhost:8080"

type fakeCreateShortener struct {
	code string
	err  error
}

func (s *fakeCreateShortener) Shorten(_ string) (string, error) {
	return s.code, s.err
}

func (s *fakeCreateShortener) GetURL(_ string) (string, bool) {
	return "", true
}

type bodyFake struct{}

func (b *bodyFake) Read(p []byte) (n int, err error) {
	return 0, errors.New("fake error")
}

func TestHandler_Create(t *testing.T) {
	type fields struct {
		code string
		err  error
	}

	type want struct {
		contentType string
		statusCode  int
		result      string
	}

	type send struct {
		headers map[string]string
		request string
		body    io.Reader
	}

	tests := []struct {
		name   string
		fields fields
		send   send
		want   want
	}{
		{
			name: "positive shorten",
			fields: fields{
				code: "Vhbd1sr",
				err:  nil,
			},
			send: send{
				request: "/",
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
				body: strings.NewReader("https://practicum.yandex.ru"),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusCreated,
				result:      "http://localhost:8080/Vhbd1sr",
			},
		},
		{
			name:   "negative content type",
			fields: fields{},
			send: send{
				request: "/",
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: strings.NewReader("https://practicum.yandex.ru"),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				result:      "Only text/plain requests are allowed!\n",
			},
		},
		{
			name: "negative invalid Url",
			fields: fields{
				err: service.ErrInvalidURL,
			},
			send: send{
				request: "/",
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
				body: strings.NewReader("practicumyandex-ru"),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				result:      "invalid url\n",
			},
		},
		{
			name: "negative internal error",
			fields: fields{
				err: errors.New("internal error"),
			},
			send: send{
				request: "/",
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
				body: strings.NewReader("https://practicum.yandex.ru"),
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusInternalServerError,
				result:      "internal error\n",
			},
		},
		{
			name:   "negative read body error",
			fields: fields{},
			send: send{
				request: "/",
				headers: map[string]string{
					"Content-Type": "text/plain",
				},
				body: &bodyFake{},
			},
			want: want{
				contentType: "text/plain",
				statusCode:  http.StatusBadRequest,
				result:      "Error read body\n",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(&fakeCreateShortener{code: tt.fields.code, err: tt.fields.err}, host)
			request := httptest.NewRequest(http.MethodPost, tt.send.request, tt.send.body)
			for k, v := range tt.send.headers {
				request.Header.Set(k, v)
			}
			w := httptest.NewRecorder()
			handler.Create(w, request)
			result := w.Result()
			defer func() {
				err := result.Body.Close()
				assert.NoError(t, err)
			}()
			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			contentType, _, err := mime.ParseMediaType(result.Header.Get("Content-Type"))
			require.NoError(t, err)
			assert.Equal(t, tt.want.contentType, contentType)
			bodyBytes, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.result, string(bodyBytes))
		})
	}
}
