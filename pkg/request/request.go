package request

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func httpRequest(ctx context.Context, method, url string, timeout int, headers map[string]string, postData []byte) ([]byte,
	error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(postData))
	if err != nil {
		return nil, fmt.Errorf("httpRequest failed, %s", err.Error())
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	if v, ok := headers["Host"]; ok {
		req.Host = v
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("httpResponse failed, %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("httpResponse code %d, %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read http response body failed, %s", err.Error())
	}

	return body, nil
}

func PostWithContext(ctx context.Context, url string, data []byte, timeout int, headers map[string]string) ([]byte, error) {
	return httpRequest(ctx, "POST", url, timeout, headers, data)
}

func GetWithContext(ctx context.Context, url string, timeout int, headers map[string]string) ([]byte, error) {
	return httpRequest(ctx, "GET", url, timeout, headers, nil)
}
