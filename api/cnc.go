package api

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

	return nil
}
