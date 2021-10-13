package httpcli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"unsafe"
)

type HttpCli struct {
	httpClient  *http.Client
	httpRequest *http.Request
	Error       error
}

func newHttpClient() *http.Client {
	return &http.Client{Timeout: 5 * time.Second}
}

func Get(url string) *HttpCli {
	request, err := http.NewRequest("GET", url, nil)
	return &HttpCli{
		httpClient:  newHttpClient(),
		httpRequest: request,
		Error:       err,
	}
}

func Post(url string) *HttpCli {
	request, err := http.NewRequest("POST", url, nil)
	return &HttpCli{
		httpClient:  newHttpClient(),
		httpRequest: request,
		Error:       err,
	}
}

func (c *HttpCli) SetTimeOut(timeout time.Duration) *HttpCli {
	c.httpClient.Timeout = timeout
	return c
}

func (c *HttpCli) SetHeader(key, value string) *HttpCli {
	c.httpRequest.Header.Set(key, value)
	return c
}

func (c *HttpCli) BodyWithJson(obj interface{}) *HttpCli {
	if c.Error != nil {
		return c
	}

	buf, err := json.Marshal(obj)
	if err != nil {
		c.Error = err
		return c
	}
	c.httpRequest.Body = ioutil.NopCloser(bytes.NewReader(buf))
	c.httpRequest.ContentLength = int64(len(buf))
	c.httpRequest.Header.Set("Content-Type", "application/json")
	return c
}

func (c *HttpCli) BodyWithBytes(buf []byte) *HttpCli {
	if c.Error != nil {
		return c
	}

	c.httpRequest.Body = ioutil.NopCloser(bytes.NewReader(buf))
	c.httpRequest.ContentLength = int64(len(buf))
	return c
}

func (c *HttpCli) BodyWithForm(form map[string]string) *HttpCli {
	if c.Error != nil {
		return c
	}

	var value url.Values = make(map[string][]string, len(form))
	for k, v := range form {
		value.Add(k, v)
	}
	buf := Str2bytes(value.Encode())

	c.httpRequest.Body = ioutil.NopCloser(bytes.NewReader(buf))
	c.httpRequest.ContentLength = int64(len(buf))
	c.httpRequest.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return c
}

func Str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

// ToJson adds request raw body encoding by JSON.
func (c *HttpCli) ToJson(obj interface{}) error {
	if c.Error != nil {
		return c.Error
	}

	response, err := c.httpClient.Do(c.httpRequest)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("error code: %d ", response.StatusCode))
	}
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, obj)
	if err != nil {
		return err
	}
	return nil
}

func (c *HttpCli) ToBytes() ([]byte, error) {
	if c.Error != nil {
		return nil, c.Error
	}

	response, err := c.httpClient.Do(c.httpRequest)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("error code: %d ", response.StatusCode))
	}
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}
