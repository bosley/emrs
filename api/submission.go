package api

import (
	"net/http"
)

func HttpSubmissions(opts Options, info *HttpsInfo) SubmissionApi {
	return newHttpController(opts, info)
}

func (c *httpController) Submit(route string, data []byte) error {

	request, err := buildHttpPostRequest("/submit/event", route, data, c.opts)
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
