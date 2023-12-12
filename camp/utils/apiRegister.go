package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/mufe/golang-base/camp/xlog"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type KongRegister struct{}

type HttpApi struct {
	Pattern string
	Name    string
	Version int32
	Handler func(http.ResponseWriter, *http.Request)
}

func (p KongRegister) StartRegister(api []HttpApi, serviceHost string) bool {
	pathMap := make(map[string]string)
	if serviceHost == "" {
		serviceHost = os.Getenv("SERVICE_80_NAME") + ".service.consul"
	}
	plugins := make([]string, 1)
	plugins[0] = "cors"
	for _, info := range api {
		pathMap[info.Pattern] = ""
		name := info.Name
		if name == "" {
			name = info.Pattern[1 : len(info.Pattern)-0]
			name = strings.ReplaceAll(name, "/", "_")
		}
		updateService(os.Getenv("CONSUL_ADMIN_IP"), name, info.Pattern, serviceHost, plugins)
	}
	getServiceFromTag(os.Getenv("CONSUL_ADMIN_IP"), serviceHost, pathMap)
	return false
}

type KongResult struct {
	Id         string   `json:"id"`
	Path       string   `json:"path"`
	EntityId   string   `json:"entity_id"`
	EntityName string   `json:"entity_name"`
	Tags       []string `json:"tags"`
	Code       int32    `json:"code"`
	Message    string   `json:"message"`
	Name       string   `json:"name"`
}

type KongResults struct {
	Data []KongResult `json:"data"`
	Next string       `json:"next"`
}

func updateService(adminUrl string, name string, pattern string, serviceHost string, plugins []string) error {
	urlStr := adminUrl + "services/" + name
	str, _ := get(urlStr, "GET", url.Values{})
	result := KongResult{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		xlog.ErrorP(str)
		xlog.ErrorP(urlStr)
		return xlog.Error(err.Error() + str)
	} else {
		urlStr := adminUrl + "services/"
		method := "POST"
		postValues := url.Values{}
		postValues.Add("name", name)
		if result.Code == 0 && result.Id != "" {
			method = "PATCH"
			urlStr += result.Id
		}
		postValues.Add("host", serviceHost)
		postValues.Add("path", pattern)
		postValues.Add("tags[]", serviceHost)
		str, _ := get(urlStr, method, postValues)
		result := KongResult{}
		err := json.Unmarshal([]byte(str), &result)
		xlog.Info(str)
		xlog.Info(err)
		if err != nil {
			xlog.ErrorP(err)
			xlog.ErrorP(str)
			return xlog.Error(err.Error() + str)
		} else {
			if result.Code != 0 {
				xlog.ErrorP(result)
				return xlog.Error(str)
			}
		}
		err = createRoute(adminUrl, result.Id, pattern, name)
		if err != nil {
			xlog.ErrorP(err)
			xlog.ErrorP(str)
			return xlog.Error(err.Error() + str)
		}
		err = createPlugin(adminUrl, result.Id, plugins)
		if err != nil {
			xlog.ErrorP(err)
			xlog.ErrorP(str)
			return xlog.Error(err.Error() + str)
		}
		return nil
	}

}

func createRoute(adminUrl string, serviceId string, path string, name string) error {
	err := delRouter(adminUrl, serviceId)
	postValues := url.Values{}
	postValues.Add("paths", path)
	postValues.Add("name", name)
	postValues.Add("methods[]", http.MethodGet)
	postValues.Add("methods[]", http.MethodPost)
	postValues.Add("methods[]", http.MethodPut)
	postValues.Add("methods[]", http.MethodDelete)
	postValues.Add("methods[]", http.MethodOptions)
	postValues.Add("methods[]", http.MethodPatch)
	postValues.Add("methods[]", http.MethodHead)
	str, _ := get(adminUrl+"services/"+serviceId+"/routes", "POST", postValues)
	kResult := KongResult{}
	err = json.Unmarshal([]byte(str), &kResult)

	if err != nil {
		xlog.ErrorP(str)
		xlog.ErrorP(err)
		return err
	} else {
		return nil
	}
}

func createPlugin(adminUrl string, serviceId string, plugins []string) error {
	err := delPlugin(adminUrl, serviceId)
	for _, str := range plugins {
		postValues := url.Values{}
		postValues.Add("name", str)
		str, _ := get(adminUrl+"services/"+serviceId+"/plugins", "POST", postValues)
		result := KongResult{}
		err := json.Unmarshal([]byte(str), &result)
		if err != nil {
			xlog.Info(str)
			return err
		} else {
			return nil
		}
	}
	return err
}

func get(url string, method string, data url.Values) (response string, err error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		xlog.Info(resp)
		return "", err
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			xlog.ErrorP(err)
			return "", xlog.Error("请求token失败")
		}
	}

	response = result.String()
	return response, nil
}

func delService(adminUrl string, serviceId string) error {
	err := delRouter(adminUrl, serviceId)
	if err != nil {
		return err
	}
	err = delPlugin(adminUrl, serviceId)
	if err != nil {
		return err
	}
	str, _ := get(adminUrl+"services/"+serviceId, "DELETE", url.Values{})
	result := KongResults{}
	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		return err
	}
	return nil
}

func delRouter(adminUrl string, serviceId string) error {
	str, _ := get(adminUrl+"services/"+serviceId+"/routes", "GET", url.Values{})
	result := KongResults{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return err
	}
	for _, info := range result.Data {
		str, _ = get(adminUrl+"routes/"+info.Id, "DELETE", url.Values{})
	}
	return nil
}

func delPlugin(adminUrl string, serviceId string) error {
	str, _ := get(adminUrl+"services/"+serviceId+"/plugins", "GET", url.Values{})
	result := KongResults{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return err
	}
	for _, info := range result.Data {
		_, _ = get(adminUrl+"plugins/"+info.Id, "DELETE", url.Values{})
	}
	return nil
}

func getServiceFromTag(adminUrl string, tag string, pathMap map[string]string) error {
	str, _ := get(adminUrl+"tags/"+tag, "GET", url.Values{})
	result := KongResults{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return err
	}

	delServices := make([]string, 0)
	for _, info := range result.Data {
		if info.EntityName == "services" {
			str, _ := get(adminUrl+"services/"+info.EntityId, "GET", url.Values{})
			result := KongResult{}
			err := json.Unmarshal([]byte(str), &result)
			if err != nil {
				return err
			}
			_, ok := pathMap[result.Path]
			if !ok || len(result.Tags) == 0 {
				delServices = append(delServices, result.Id)
			}
		}
	}
	for _, str := range delServices {
		_ = delService(adminUrl, str)
	}
	return nil
}

func DelDefual(adminUrl string) error {
	str, _ := get(adminUrl+"services/", "GET", url.Values{})
	result := KongResults{}
	err := json.Unmarshal([]byte(str), &result)
	if err != nil {
		return err
	}
	delServices := make([]string, 0)
	for _, info := range result.Data {
		fmt.Println(info.Name + "1")
		//if len(info.Tags) == 0 {
		delServices = append(delServices, info.Id)
		//}
	}
	if result.Next != "" {
		str, _ := get(adminUrl+result.Next[1:len(result.Next)], "GET", url.Values{})
		result := KongResults{}
		err := json.Unmarshal([]byte(str), &result)
		if err != nil {
			return err
		}
		for _, info := range result.Data {
			fmt.Println(info.Name + "2")
			//if len(info.Tags) == 0 {
			delServices = append(delServices, info.Id)
			//}
		}
	}
	for _, str := range delServices {
		fmt.Println(str + "3")
		_ = delService(adminUrl, str)
	}
	return nil
}
