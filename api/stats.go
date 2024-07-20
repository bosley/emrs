package api

import (
	"bytes"
	"encoding/json"
	"net/url"
	"time"
)

type StatsResponse struct {
	Uptime time.Duration `json:uptime`
}

func HttpStats(opts Options, info *HttpsInfo) StatsApi {
	return newHttpController(opts, info)
}

// Retrieve the duration of time a remote server has been running
// If the server can not be reached, an error can be expected
func (c *httpController) GetUptime() (time.Duration, error) {

	var t time.Duration

	dest, err := url.JoinPath(c.opts.Binding, "/stat")
	if err != nil {
		return t, err
	}

	client := newHttpClient(c.https)

	response, err := client.Get(dest)
	if err != nil {
		return t, err
	}

	defer response.Body.Close()

	data := new(bytes.Buffer)
	data.ReadFrom(response.Body)

	var result StatsResponse

	if err := json.Unmarshal(data.Bytes(), &result); err != nil {
		return t, err
	}

	t = result.Uptime

	return t, nil
}
