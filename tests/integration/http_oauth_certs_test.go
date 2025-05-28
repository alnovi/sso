package integration

import (
	"net/http"
	"net/http/httptest"

	"github.com/alnovi/sso/internal/transport/http/controller/oauth"
)

func (s *TestSuite) TestHttpOauthCerts() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctrl := oauth.NewCertsController(s.app.Provider.Certs())

	c := s.app.HttpServer.NewContext(req, rec)

	err := s.sendToServer(ctrl.Certs, c)
	s.Assert().NoError(err, MsgNotAssertError)
	s.Assert().Equal(http.StatusOK, rec.Code, MsgNotAssertCode)
	s.Assert().Contains(rec.Body.String(), "RSA", MsgNotAssertBody)
}
