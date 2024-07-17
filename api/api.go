package api

type Options struct {
	Binding     string
	AssetId     string // TODO: Consider: we don't need to store info in the server about the asset. we could have a "registration" type that secures nodes from impersonation, but it isn't a hard req for mvp
	AccessToken string // Should be a voucher that corresponds to the type of interface that is expected.
}

type CNCApi interface {
}

type SubmissionApi interface {
}

type StatsApi interface {
}

type controller struct {
}

func CNC(uiKey string) CNCApi {

	// TODO: Using the provided uikey setup the controller to perform command and control functions on /cnc

	return &controller{}
}

func Submissions(opts Options) SubmissionApi {

	return &controller{}
}

func Stats(opts Options) StatsApi {

	return &controller{}
}
