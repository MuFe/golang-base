package wx

const (
	redirectOauth       = "/connect/oauth2/authorize"
	webAppRedirectOauth = "/connect/qrconnect"
)

//GetRedirectURL 获取跳转的url地址
func GetRedirectURL(redirectURI, scope, state, appId string) (string, error) {
	api := baseAuthURL + redirectOauth
	queries := requestQueries{
		"redirect_uri":  redirectURI,
		"response_type": "code",
		"appid":         appId,
		"scope":         scope,
		"state":         state,
	}

	url, err := encodeURL(api, queries)
	if err != nil {
		return "", err
	}
	return url + "#wechat_redirect", nil
}

//GetWebAppRedirectURL 获取网页应用跳转的url地址
func GetWebAppRedirectURL(redirectURI, scope, state, appId string) (string, error) {
	api := baseAuthURL + webAppRedirectOauth
	queries := requestQueries{
		"redirect_uri":  redirectURI,
		"response_type": "code",
		"appid":         appId,
		"scope":         scope,
		"state":         state,
	}

	url, err := encodeURL(api, queries)
	if err != nil {
		return "", err
	}
	return url + "#wechat_redirect", nil
}
