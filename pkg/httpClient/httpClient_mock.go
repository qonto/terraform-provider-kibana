package httpClient

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpClientMock struct {
	ShouldFail bool
	Resp       []FakeResponse
}

type FakeResponse struct {
	Status  int
	Headers http.Header
	Payload []byte
}

func (c *HttpClientMock) Get(url string, headers map[string]string) ([]byte, int, error) {
	if c.ShouldFail {
		return []byte(""), 400, errors.New("Failed to get resource")
	}
	resp := c.PopPayload()
	return resp.Payload, resp.Status, nil
}
func (c *HttpClientMock) GetReturnHeaders(url string, headers map[string]string) ([]byte, http.Header, int, error) {
	if c.ShouldFail {
		return []byte(""), http.Header{}, 400, errors.New("Failed to get resource")
	}
	resp := c.PopPayload()
	return resp.Payload, resp.Headers, resp.Status, nil
}

func (c *HttpClientMock) GetReturnReader(url string, headers map[string]string) (io.ReadCloser, int, error) {
	if c.ShouldFail {
		return nil, 400, errors.New("Failed to get resource")
	}
	resp := c.PopPayload()
	return ioutil.NopCloser(bytes.NewReader(resp.Payload)), resp.Status, nil
}

func (c *HttpClientMock) Delete(url string, headers map[string]string) ([]byte, int, error) {
	if c.ShouldFail {
		return []byte(""), 400, errors.New("Failed to delete resource")
	}
	resp := c.PopPayload()
	return resp.Payload, resp.Status, nil
}

func (c *HttpClientMock) Put(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error) {
	if c.ShouldFail {
		return []byte(""), 400, errors.New("Failed to put resource")
	}
	resp := c.PopPayload()
	return resp.Payload, resp.Status, nil
}

func (c *HttpClientMock) Post(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error) {
	if c.ShouldFail {
		return []byte(""), 400, errors.New("Failed to get resource")
	}
	resp := c.PopPayload()
	return resp.Payload, resp.Status, nil
}

func (c *HttpClientMock) Patch(url string, headers map[string]string, jsonBody []byte) ([]byte, int, error) {
	if c.ShouldFail {
		return []byte(""), 400, errors.New("Failed to get resource")
	}
	resp := c.PopPayload()
	return resp.Payload, resp.Status, nil
}

func (c *HttpClientMock) PopPayload() (ret FakeResponse) {
	if len(c.Resp) > 0 {
		ret = c.Resp[0]
		c.Resp = c.Resp[1:]
		return ret
	}
	return FakeResponse{Status: 200}
}

func (c *HttpClientMock) PushPayload(resp FakeResponse) {
	c.Resp = append(c.Resp, resp)
}
