package salesforce

type OAuthClientIF interface {
	GetOAuthInfo() (*OAuthInfo, error)
}

type OAuthInfo struct {
	AccessToken string `json:"access_token"`
	InstanceUrl string `json:"instance_url"`
}
