package api

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

var ErrInvalidEmrsURL = errors.New("malformed emrs url")

type EmrsAddress struct {
	Asset  string
	Route  []string
	Server string
}

func (a *EmrsAddress) From(url string) error {

	urlServ := strings.Split(url, `@`)

	if len(urlServ) <= 0 {
		slog.Error("empty emrs url given")
		return ErrInvalidEmrsURL
	}

	if len(urlServ) > 2 {
		slog.Error("emrs url contains extraneous `@`", "url", url)
		return ErrInvalidEmrsURL
	}

	newServerAddress := ""

	if len(urlServ) == 2 {
		newServerAddress = urlServ[1]
	}

	assetAction := strings.Split(urlServ[0], `:`)

	if len(assetAction) != 2 {
		slog.Error("invalid ASSET:ACTION in emrs url", "given", urlServ[0])
		return ErrInvalidEmrsURL
	}

	newRoute, err := DecomposeRoute(assetAction[1])
	if err != nil {
		slog.Error("failed to decompose route", "route", assetAction[1], "error", err.Error())
		return ErrInvalidEmrsURL
	}

	a.Asset = assetAction[0]
	a.Route = newRoute
	a.Server = newServerAddress

	return nil
}

func (a *EmrsAddress) ToUrl() (string, error) {

	if !ValidateCommonChunk(a.Asset) {
		slog.Error("invalid asset name", "name", a.Asset)
		return "", ErrInvalidEmrsURL
	}

	r, err := ComposeRoute(a.Route)
	if err != nil {
		slog.Error("invalid route", "error", err.Error())
		return "", ErrInvalidEmrsURL
	}

	result := fmt.Sprintf("%s:%s", a.Asset, r)

	if len(strings.TrimSpace(a.Server)) != 0 {
		result = fmt.Sprintf("%s@%s", result, a.Server)
	}
	return result, nil
}
