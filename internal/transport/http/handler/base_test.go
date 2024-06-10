package handler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestBaseHandler_RequestIsJSON(t *testing.T) {
	testCases := []struct {
		name    string
		headers map[string]string
		isJson  bool
	}{
		{
			name:   "Standard request",
			isJson: false,
		},
		{
			name:    "Content-Type = application/json",
			headers: map[string]string{"Content-Type": "application/json"},
			isJson:  true,
		},
		{
			name:    "X-Requested-With = XMLHttpRequest",
			headers: map[string]string{"X-Requested-With": "XMLHttpRequest"},
			isJson:  true,
		},
	}

	b := &BaseHandler{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			for k, v := range tc.headers {
				req.Header.Add(k, v)
			}
			assert.Equal(t, tc.isJson, b.RequestIsJSON(req))
		})
	}
}

func TestBaseHandler_MustClientId(t *testing.T) {
	exp := "12345"

	b := &BaseHandler{}
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	c.Set("client_id", exp)

	assert.NotPanics(t, func() {
		assert.Equal(t, exp, b.MustClientId(c))
	})
}

func TestBaseHandler_MustClientId_Panic(t *testing.T) {
	b := &BaseHandler{}
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	assert.Panics(t, func() {
		b.MustClientId(c)
	})
}

func TestBaseHandler_MustUserId(t *testing.T) {
	exp := "12345"

	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())
	c.Set("user_id", exp)

	b := &BaseHandler{}

	assert.Equal(t, exp, b.MustUserId(c))
}

func TestBaseHandler_MustUserId_Panic(t *testing.T) {
	b := &BaseHandler{}
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	assert.Panics(t, func() {
		b.MustUserId(c)
	})
}

func TestBaseHandler_StatusText(t *testing.T) {
	testCases := []struct {
		name  string
		codes [2]int
		empty bool
	}{
		{
			name:  "Invalid code",
			codes: [2]int{0, 0},
			empty: true,
		},
		{
			name:  "Codes 100x",
			codes: [2]int{100, 103},
			empty: false,
		},
		{
			name:  "Codes 200x",
			codes: [2]int{200, 208},
			empty: false,
		},
		{
			name:  "Codes 300x",
			codes: [2]int{300, 305},
			empty: false,
		},
		{
			name:  "Codes 400x",
			codes: [2]int{400, 418},
			empty: false,
		},
		{
			name:  "Codes 420x",
			codes: [2]int{421, 426},
			empty: false,
		},
		{
			name:  "Codes 500x",
			codes: [2]int{500, 508},
			empty: false,
		},
	}

	b := &BaseHandler{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for code := tc.codes[0]; code <= tc.codes[1]; code++ {
				if tc.empty {
					assert.Empty(t, b.StatusText(code), fmt.Sprintf("not emty text for code: %d", code))
				} else {
					assert.NotEmpty(t, b.StatusText(code), fmt.Sprintf("emty text for code: %d", code))
				}
			}
		})
	}
}
