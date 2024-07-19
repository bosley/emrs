package api

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

func HttpCNC(binding string, uiKey string, info *HttpsInfo) CNCApi {
	return newHttpController(
		Options{
			Binding:     binding,
			AssetId:     "",
			AccessToken: uiKey,
		},
		info,
	)
}

func (c *httpController) Shutdown() error {

	opts := c.opts

	// TODO: This big-ol blob has to have an existing solution thats not as clunky
	//          basically... we bing to addr:port, but http requires the schema (http/s) to
	//          be prefixed. a regular address without the prefix will get caught by ParseRequestURI
	//          but it thinkgs localhost is fine, though http requres the prefix on localhost as well
	// This issue is one that will be all around every http api so it needs to be ironed

	_, err := url.ParseRequestURI(opts.Binding)
	if err != nil || strings.HasPrefix(opts.Binding, "localhost") {
		if c.https != nil {
			opts.Binding = fmt.Sprintf("https://%s", opts.Binding)
		} else {
			opts.Binding = fmt.Sprintf("http://%s", opts.Binding)
		}
		_, e := url.ParseRequestURI(opts.Binding)
		if e != nil {
			slog.Error("failed to make the given binding work for CNC", "error", e.Error())
			return e
		}
	}

	request, err := buildHttpPostRequest("/cnc/shutdown", "", []byte{}, opts)
	if err != nil {
		return err
	}

	client := newHttpClient(c.https)

	result, err := client.Do(request)
	if err != nil {
		return err
	}

	if result.StatusCode != http.StatusOK {
		return ErrUnexpectedStatusCode
	}
	return nil
}
