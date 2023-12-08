package wx

import (
	"encoding/json"
	"github.com/mufe/golang-base/camp/wx/util"
	"github.com/mufe/golang-base/camp/xlog"
)

const (
	apiLogin              = "/sns/jscode2session"
	apiGetAccessToken     = "/cgi-bin/token"
	apiGetTicket          = "/cgi-bin/ticket/getticket"
	apiGetUserAccessToken = "/sns/oauth2/access_token"
	apiGetUserInfo        = "/sns/userinfo"
)

// LoginResponse 返回给用户的数据
type LoginResponse struct {
	CommonError
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	// 用户在开放平台的唯一标识符
	// 只在满足一定条件的情况下返回
	UnionID string `json:"unionid"`
}

// Login 登录凭证校验。通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程。
//
// appID 小程序 appID
// secret 小程序的 app secret
// code 小程序登录时获取的 code
func Login(appID, secret, code string) (*LoginResponse, error) {
	api := baseURL + apiLogin

	return login(appID, secret, code, api)
}

func login(appID, secret, code, api string) (*LoginResponse, error) {
	queries := requestQueries{
		"appid":      appID,
		"secret":     secret,
		"js_code":    code,
		"grant_type": "authorization_code",
	}

	url, err := encodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	res := new(LoginResponse)
	if err := getJSON(url, res); err != nil {
		return nil, err
	}

	if res.ErrCode != 0 {
		xlog.ErrorP(url)
		return nil, xlog.Error(res.ErrMSG)
	}
	return res, nil
}

// TokenResponse 获取 access_token 成功返回数据
type TokenResponse struct {
	CommonError
	AccessToken string `json:"access_token"` // 获取到的凭证
	ExpiresIn   uint   `json:"expires_in"`   // 凭证有效时间，单位：秒。目前是7200秒之内的值。
}

// TicketResponse 获取 ticket 成功返回数据
type TicketResponse struct {
	CommonError
	Ticket    string `json:"ticket"`     // 获取到的凭证
	ExpiresIn uint   `json:"expires_in"` // 凭证有效时间，单位：秒。目前是7200秒之内的值。
}

// GetAccessToken 获取小程序全局唯一后台接口调用凭据（access_token）。
// 调调用绝大多数后台接口时都需使用 access_token，开发者需要进行妥善保存，注意缓存。
func GetAccessToken(appID, secret string) (*TokenResponse, error) {
	api := baseURL + apiGetAccessToken
	return getAccessToken(appID, secret, api)
}

func getAccessToken(appID, secret, api string) (*TokenResponse, error) {

	queries := requestQueries{
		"appid":      appID,
		"secret":     secret,
		"grant_type": "client_credential",
	}

	url, err := encodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	res := new(TokenResponse)
	if err := getJSON(url, res); err != nil {
		return nil, err
	}

	return res, nil
}

// 通过AccessToken获取Ticket
func GetTicket(accessToken string) (*TicketResponse, error) {
	api := baseURL + apiGetTicket
	return getTicket(accessToken, api)
}

func getTicket(accessToken, api string) (*TicketResponse, error) {

	queries := requestQueries{
		"access_token": accessToken,
		"type":         "jsapi",
	}

	url, err := encodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	res := new(TicketResponse)
	if err := getJSON(url, res); err != nil {
		return nil, err
	}
	xlog.ErrorP(url)
	return res, nil
}

// ResAccessToken 获取用户授权access_token的返回结果
type ResAccessToken struct {
	CommonError

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`

	// UnionID 只有在用户将公众号绑定到微信开放平台帐号后，才会出现该字段。
	// 公众号文档 https://mp.weixin.qq.com/wiki?t=resource/res_main&id=mp1421140842
	UnionID string `json:"unionid"`
}

// GetUserAccessToken 通过网页授权的code 换取access_token
func GetUserAccessToken(appID, secret, code string) (*ResAccessToken, error) {
	api := baseURL + apiGetUserAccessToken
	return getUserAccessToken(appID, secret, code, api)
}

func getUserAccessToken(appId, secret, code, api string) (*ResAccessToken, error) {

	queries := requestQueries{
		"appid":      appId,
		"secret":     secret,
		"code":       code,
		"grant_type": "authorization_code",
	}

	url, err := encodeURL(api, queries)
	if err != nil {
		xlog.ErrorP(err)
		return nil, err
	}

	res := new(ResAccessToken)
	if err := getJSON(url, res); err != nil {

		return nil, err
	}
	xlog.ErrorP(res)
	return res, nil
}

//UserInfo 用户授权获取到用户信息
type UserInfo struct {
	CommonError
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int64    `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}

// GetUserInfo 通过access_token换UserInfo
func GetUserInfo(appId, secret, code string) (*UserInfo, error) {
	result, err := GetUserAccessToken(appId, secret, code)
	if err != nil {
		return nil, err
	}
	api := baseURL + apiGetUserInfo
	return getUserInfo(result.AccessToken, result.OpenID, api)
}

func getUserInfo(accessToken, openId, api string) (*UserInfo, error) {

	queries := requestQueries{
		"access_token": accessToken,
		"openid":       openId,
		"lang":         "zh_CN",
	}

	url, err := encodeURL(api, queries)
	if err != nil {
		return nil, err
	}

	res := new(UserInfo)
	if err := getJSON(url, res); err != nil {
		return nil, err
	}

	return res, nil
}

// PhoneNumber 解密后的用户手机号码信息
type PhoneNumber struct {
	PhoneNumber     string `json:"phoneNumber"`
	PurePhoneNumber string `json:"purePhoneNumber"`
	CountryCode     string `json:"countryCode"`
}

// DecryptPhoneNumber 解密手机号码
//
// @ssk 通过 Login 向微信服务端请求得到的 session_key
// @data 小程序通过 api 得到的加密数据(encryptedData)
// @iv 小程序通过 api 得到的初始向量(iv)
func DecryptPhoneNumber(ssk, data, iv string) (phone PhoneNumber, err error) {
	bts, err := util.CBCDecrypt(ssk, data, iv)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &phone)
	return
}
