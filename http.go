package tools

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)

func GetToken(method, rawUrl, uri, user, password string) (token string, err error) {
	var (
		params url.Values
		resp   *http.Response
		body   []byte
	)
	switch method {
	case http.MethodGet:
		params.Add("account", user)
		params.Add("password", password)
	case http.MethodPost:
		if body, err = json.Marshal(map[string]any{
			"account":  user,
			"password": password,
		}); err != nil {
			err = fmt.Errorf("parse request body failed:%s", err)
		}
	default:
		return
	}
	if resp, err = DoRequest(method, rawUrl, uri, &params, nil, bytes.NewBuffer(body)); err == nil {
		token = resp.Header.Get("Token")
	}
	return
}

func DoGet(rawUrl, uri, token string, params *url.Values, result any) (code int, err error) {
	var (
		resp   *http.Response
		header = http.Header{}
	)
	if token != "" {
		header.Add("Token", token)
	}
	if resp, err = DoRequest(http.MethodGet, rawUrl, uri, params, &header, nil); err == nil {
		code, err = ParseResponse(resp, result)
	}
	return
}

func DoPost(rawUrl, uri string, token string, params *url.Values, body any, result any) (code int, err error) {
	var (
		content []byte
		resp    *http.Response
		header  = http.Header{}
	)
	if content, err = json.Marshal(body); err == nil {
		if token != "" {
			header.Set("Token", token)
		}
		header.Set("Content-Type", "application/json")
		if resp, err = DoRequest(http.MethodPost, rawUrl, uri, params, &header, bytes.NewBuffer(content)); err == nil {
			code, err = ParseResponse(resp, result)
		}
	} else {
		err = fmt.Errorf("parse request body failed:%s", err)
	}
	return
}

func DoRequest(method, rawUrl, uri string, params *url.Values, header *http.Header, body io.Reader) (resp *http.Response, err error) {
	if uri != "" {
		rawUrl = fmt.Sprintf("%s/%s", rawUrl, strings.TrimLeft(uri, "/"))
	}
	if params != nil && len(*params) > 0 {
		if _url, e := url.Parse(rawUrl); e == nil {
			_url.RawQuery = params.Encode()
			rawUrl = _url.String()
		} else {
			err = fmt.Errorf("parse rawUrl failed:%s", e)
			return
		}
	}
	var req *http.Request
	if req, err = http.NewRequest(method, rawUrl, body); err == nil {
		if header != nil && len(*header) > 0 {
			req.Header = *header
		}
		if resp, err = http.DefaultClient.Do(req); err != nil {
			err = fmt.Errorf("send request failed:%s", err)
		}
	} else {
		err = fmt.Errorf("create request failed:%s", err)
	}
	return
}

func ParseResponse(resp *http.Response, result any) (code int, err error) {
	defer resp.Body.Close()
	code = resp.StatusCode
	var body []byte
	if code == http.StatusOK {
		if body, err = io.ReadAll(resp.Body); err == nil {
			if err = json.Unmarshal(body, result); err != nil {
				err = fmt.Errorf("parse response body failed:%s", err)
			}
		} else {
			err = fmt.Errorf("read response body failed:%s", err)
		}
	}
	return
}

func LocalIP() (ip string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	addr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(addr.String(), ":")[0]
	return
}
