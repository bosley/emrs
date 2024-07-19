package api

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
)

const (
	HttpApiVersion = "v0.0.0"
)

var ErrUnexpectedStatusCode = errors.New("received unexpected status code")

type HttpsInfo struct {
	Cert string
	Key  string
}

type httpController struct {
	opts  Options
	https *HttpsInfo
}

func newHttpController(o Options, info *HttpsInfo) *httpController {
	return &httpController{
		opts:  o,
		https: info,
	}
}

func buildHttpPostRequest(endpoint string, route string, data []byte, opt Options) (*http.Request, error) {
	slog.Debug("build post request", "binding", opt.Binding, "asset", opt.AssetId)

	dest, err := url.JoinPath(opt.Binding, endpoint)
	if err != nil {
		return nil, err
	}

	r, err := http.NewRequest("POST", dest, bytes.NewBuffer(data))
	r.Header.Add("Content-Type", "octet-stream")
	r.Header.Add("EMRS-API-Version", HttpApiVersion)
	r.Header.Add("origin", opt.AssetId)
	r.Header.Add("token", opt.AccessToken)
	r.Header.Add("route", route)

	if err != nil {
		return nil, err
	}
	return r, nil
}

func newHttpClient(info *HttpsInfo) *http.Client {

	rootCAs, _ := x509.SystemCertPool()
	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}

	if info != nil {
		f, err := os.Open(info.Cert)
		if err != nil {
			slog.Error("failed to read in cert file for submission request", "error", err.Error())
			os.Exit(1)
		}
		defer f.Close()

		certs, err := io.ReadAll(f)

		if err != nil {
			slog.Error("Failed to append cert to RootCAs",
				"cert", info.Cert,
				"error", err.Error())
			os.Exit(1)
		}
		if ok := rootCAs.AppendCertsFromPEM(certs); !ok {
			slog.Info("No certs appended, using system certs only")
		}
	}

	config := &tls.Config{
		InsecureSkipVerify: info == nil,
		RootCAs:            rootCAs,
	}
	tr := &http.Transport{TLSClientConfig: config}
	return &http.Client{Transport: tr}
}
