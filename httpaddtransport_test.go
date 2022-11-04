package httpaddtransport_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ciricc/httpaddtransport"
)

func TestAddCustomHeader(t *testing.T) {

	userAgentHeader := "golang/1.19"

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != userAgentHeader {
			t.Errorf("user agent expected %s but got %s", userAgentHeader, r.Header.Get("User-Agent"))
		}
	}))

	httpAddTransport, err := httpaddtransport.New(nil, nil)

	if err != nil {
		t.Error(err)
	}

	httpAddTransport.Headers.Set("User-Agent", userAgentHeader)
	httpAddTransport.Log = true // enabled loggin all sent requests

	httpClient := http.Client{
		Transport: httpAddTransport,
	}

	_, err = httpClient.Get(serv.URL)
	if err != nil {
		t.Error(err)
	}
}

type MyProxyTripper struct {
	tr http.RoundTripper
	http.RoundTripper
	tripped bool
}

func (v *MyProxyTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	v.tripped = true
	return v.tr.RoundTrip(req)
}

func TestCustomRoundTripper(t *testing.T) {
	myProxyTripper := MyProxyTripper{
		tr: http.DefaultTransport,
	}
	httpAddTransport, err := httpaddtransport.New(&myProxyTripper, log.Default())
	if err != nil {
		t.Fatal(err)
	}

	if httpAddTransport.Tripper() != &myProxyTripper {
		t.Fatalf("expected tripper %v but got %v", &myProxyTripper, httpAddTransport.Tripper())
	}

	serv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	}))

	httpClient := http.Client{
		Transport: httpAddTransport,
	}

	_, err = httpClient.Get(serv.URL)
	if err != nil {
		t.Error(err)
	}

	if !myProxyTripper.tripped {
		t.Fatalf("not tripped proxy custom tripper")
	}
}
