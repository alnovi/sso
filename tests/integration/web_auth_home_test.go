package integration

import (
	"net/http"
	"net/http/httptest"
)

func (s *TestSuite) TestWebAuthHome() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := s.App.Server.NewContext(req, rec)

	err := s.SendToServer(s.App.Provider.WebAuth().Home, c)

	s.Assert().NoError(err)
	s.Assert().Equal(http.StatusFound, rec.Code)
}
