package tools

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func GetData(url string, cookie *http.Cookie, host string, referer string, headers map[string]string, contentType string) (string, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}else {
		req.Header.Add("Content-Type", "application/json")
	}
	if referer != "" {
		req.Header.Add("referer", referer)
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	if host != "" {
		req.Host = host
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func PostData(url string, data string, cookie *http.Cookie, host, referer string, headers map[string]string, contentType string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	if cookie != nil {
		req.AddCookie(cookie)
	}
	if contentType != "" {
		req.Header.Add("Content-Type", contentType)
	}else {
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	}
	if referer != "" {
		req.Header.Add("referer", referer)
	}
	if host != "" {
		req.Host = host
	}
	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}
	response, err := client.Do(req)
	defer response.Body.Close()
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func CommonRequest(method string, url string, data string, cookie *http.Cookie, host, referer string, headers map[string]string, contentType string) (string, error){
	if method == "POST" {
		return PostData(url, data, cookie, host, referer, headers, contentType)
	} else {
		return GetData(url, cookie, host, referer, headers, contentType)
	}
}