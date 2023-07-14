package cloudapi

import (
	"fmt"
	"io"
	"net/http"

	"github.com/shurcooL/graphql"
)

type Client struct {
	client *graphql.Client
}

func NewClientWithURL(httpClient *http.Client, host string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &Client{
		client: graphql.NewClient(host, httpClient),
	}
}

type CredentialsTransport struct {
	tripper http.RoundTripper
	token   string
	version string
}

func NewCredentialsTransport(t http.RoundTripper, token, version string) http.RoundTripper {
	return CredentialsTransport{
		tripper: t,
		token:   token,
		version: version,
	}
}

func (t CredentialsTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	newReq := new(http.Request)
	*newReq = *r

	newReq.Header = make(http.Header, len(r.Header))
	for k, s := range r.Header {
		newReq.Header[k] = s
	}

	newReq.Header.Set("Authorization", "Bearer "+t.token)
	newReq.Header.Set("X-ENAPTER-CLI-VERSION", t.version)

	return t.tripper.RoundTrip(newReq)
}

type CLIMessageWriterTransport struct {
	tripper http.RoundTripper
	writer  io.Writer
}

func NewCLIMessageWriterTransport(t http.RoundTripper, w io.Writer) http.RoundTripper {
	return CLIMessageWriterTransport{
		tripper: t,
		writer:  w,
	}
}

func (t CLIMessageWriterTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	resp, err := t.tripper.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if msg := resp.Header.Get("X-ENAPTER-CLI-MESSAGE"); msg != "" {
		fmt.Fprintln(t.writer, msg)
	}

	return resp, nil
}
