package utils

var (
	LinkwardenToken string
	InstanceURL     string
)

func SetVars(linkwardenToken string, instanceURL string) {
	LinkwardenToken = linkwardenToken
	InstanceURL = instanceURL
}

func GetLinkwardenToken() string {
	return LinkwardenToken
}

func GetInstanceURL() string {
	return InstanceURL
}
