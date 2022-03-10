package httpClient

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpClientAPI interface {
	Get(url string, headers map[string]string) ([]byte, int, error)
	GetReturnHeaders(url string, headers map[string]string) ([]byte, http.Header, int, error)
	GetReturnReader(url string, headers map[string]string) (io.ReadCloser, int, error)
	Post(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error)
	Patch(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error)
	Delete(url string, headers map[string]string) ([]byte, int, error)
	Put(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error)
}

type HttpClient struct {
	api *http.Client
}

func CreateHTTPClient() HttpClientAPI {
	return &HttpClient{api: &http.Client{}}
}

func (c *HttpClient) Get(url string, headers map[string]string) ([]byte, int, error) {
	return c.Request("GET", url, headers)
}

func (c *HttpClient) GetReturnHeaders(url string, headers map[string]string) ([]byte, http.Header, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, 0, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := c.api.Do(req)
	if err != nil {
		return nil, nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	return respBody, resp.Header, resp.StatusCode, err
}

func (c *HttpClient) GetReturnReader(url string, headers map[string]string) (io.ReadCloser, int, error) {
	return c.RequestReader("GET", url, headers)
}

func (c *HttpClient) Delete(url string, headers map[string]string) ([]byte, int, error) {
	return c.Request("DELETE", url, headers)
}

func (c *HttpClient) Put(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error) {
	return c.RequestJson("PUT", url, headers, jsonBody)
}

func (c *HttpClient) Post(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error) {
	return c.RequestJson("POST", url, headers, jsonBody)
}

func (c *HttpClient) Patch(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error) {
	return c.RequestJson("PATCH", url, headers, jsonBody)
}

func (c *HttpClient) Request(method string, url string, headers map[string]string) ([]byte, int, error) {
	respBodyReader, statusCode, err := c.RequestReader(method, url, headers)
	if err != nil {
		return nil, 0, err
	}
	defer respBodyReader.Close()
	respBody, err := ioutil.ReadAll(respBodyReader)
	return respBody, statusCode, err
}

func (c *HttpClient) RequestJson(method, url string, headers map[string]string, jsonBody []byte) ([]byte, int, error) {
	if !isValidUrl(url) {
		return nil, 0, fmt.Errorf("Invalid url  %s", url)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, 0, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := c.api.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	return respBody, resp.StatusCode, err
}

func (c *HttpClient) RequestReader(method string, url string, headers map[string]string) (io.ReadCloser, int, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, 0, err
	}
	for key, value := range headers {
		req.Header.Add(key, value)
	}
	resp, err := c.api.Do(req)
	if err != nil {
		return nil, 0, err
	}
	return resp.Body, resp.StatusCode, err
}

func isValidUrl(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	u, err := url.Parse(toTest)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return false
	}
	return true
}
