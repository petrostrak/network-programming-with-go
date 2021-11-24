package securingcommunicationswithtls

import (
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/net/http2"
)

func TestClientTLS(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil {
			u := "https://" + r.Host + r.RequestURI
			http.Redirect(w, r, u, http.StatusMovedPermanently)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	resp, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; actual status %d", http.StatusOK, resp.StatusCode)
	}

	tp := &http.Transport{
		TLSClientConfig: &tls.Config{
			CurvePreferences: []tls.CurveID{tls.CurveP256},
			MinVersion:       tls.VersionTLS12,
		},
	}

	err = http2.ConfigureTransport(tp)
	if err != nil {
		t.Fatal(err)
	}

	client2 := &http.Client{Transport: tp}

	_, err = client2.Get(ts.URL)
	if err != nil || !strings.Contains(err.Error(), "certificate signed by unknown authority") {
		t.Fatalf("expected unknown authority error; actual: %q", err)
	}

	tp.TLSClientConfig.InsecureSkipVerify = true

	resp, err = client2.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d; actual status %d", http.StatusOK, resp.StatusCode)
	}
}
