package wx

import (
	"encoding/json"
	"errors"
	"github.com/mufe/golang-base/camp/wx/util"
)

const (
	// baseURL 微信请求基础URL
	baseURL     = "https://api.weixin.qq.com"
	baseAuthURL = "https://open.weixin.qq.com"
)

// POST 参数
type requestParams map[string]interface{}

// URL 参数
type requestQueries map[string]string

// Userinfo 解密后的用户信息
type Userinfo struct {
	OpenID   string `json:"openId"`
	Nickname string `json:"nickName"`
	Gender   int    `json:"gender"`
	Province string `json:"province"`
	Language string `json:"language"`
	Country  string `json:"country"`
	City     string `json:"city"`
	Avatar   string `json:"avatarUrl"`
	UnionID  string `json:"unionId"`
}

// DecryptUserInfo 解密用户信息
//
// @rawData 不包括敏感信息的原始数据字符串，用于计算签名。
// @encryptedData 包括敏感数据在内的完整用户信息的加密数据
// @signature 使用 sha1( rawData + session_key ) 得到字符串，用于校验用户信息
// @iv 加密算法的初始向量
// @ssk 微信 session_key
func DecryptUserInfo(rawData, encryptedData, signature, iv, ssk string) (ui Userinfo, err error) {

	if ok := util.Validate(rawData, ssk, signature); !ok {
		err = errors.New("数据校验失败")
		return
	}

	bts, err := util.CBCDecrypt(ssk, encryptedData, iv)
	if err != nil {
		return
	}

	err = json.Unmarshal(bts, &ui)
	return
}
