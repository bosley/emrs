package api

import (
  "net/http"
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
	request, err := buildHttpPostRequest("/cnc/shutdown", "", []byte{}, c.opts)
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
