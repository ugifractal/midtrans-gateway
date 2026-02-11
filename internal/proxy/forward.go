package proxy

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func Forward(r *http.Request, target string) error {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	r.Body = io.NopCloser(bytes.NewReader(body))

	req, err := http.NewRequest(
		r.Method,
		target,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	// copy headers
	for k, v := range r.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}

	req.ContentLength = int64(len(body))
	req.Host = ""

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 400 {
		return errors.New("downstream returned non-2xx")
	}

	return nil
}
