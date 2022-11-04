package httpaddtransport

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

// Structure which realises http.RoundTripper interface
type Transport struct {
	Headers http.Header // Default headers for each sent http request
	Log     bool        // Enable logging

	tr     http.RoundTripper
	logger *log.Logger
}

// Initalize new Transport and return it
func New(
	tripper http.RoundTripper,
	logger *log.Logger,
) (*Transport, error) {

	if logger == nil {
		logger = log.Default()
	}

	return &Transport{
		tr:      tripper,
		logger:  logger,
		Headers: http.Header{},
	}, nil
}

// Default function for RoundTrip
func (v *Transport) RoundTrip(req *http.Request) (*http.Response, error) {

	for k, v := range v.Headers {
		if len(v) != 0 {
			req.Header.Set(k, v[0])
		} else {
			req.Header.Del(k)
		}
	}

	if v.Log {
		req, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("dump request for log error: %s", err))
		}
		v.logger.Println("Do request: ", string(req))
	}

	return v.Tripper().RoundTrip(req)
}

// Returns current top http RoundTripper
func (v *Transport) Tripper() http.RoundTripper {

	if v.tr != nil {
		return v.tr
	}

	if http.DefaultTransport != nil {
		return http.DefaultTransport
	}

	newTransport := http.Transport{}
	v.tr = &newTransport

	return v.tr
}
