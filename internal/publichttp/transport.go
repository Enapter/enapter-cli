package publichttp

import "net/http"

type AuthTokenTransport struct {
	token string
	next  http.RoundTripper
}

func NewAuthTokenTransport(t http.RoundTripper, token string) http.RoundTripper {
	return &AuthTokenTransport{token: token, next: t}
}

func (t *AuthTokenTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	const header = "X-Enapter-Auth-Token"
	s := cloneRequest(req)
	s.Header.Set(header, t.token)
	return t.next.RoundTrip(s)
}

func cloneRequest(req *http.Request) *http.Request {
	shallow := new(http.Request)
	*shallow = *req
	for k, s := range req.Header {
		shallow.Header[k] = s
	}
	return shallow
}
