package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type HTTPParams struct {
	Url     string
	Method  string
	Body    map[string]any
	Headers map[string]string
}

type ATHTTP struct{}

func NewATHTTP() *ATHTTP {
	return &ATHTTP{}
}

func (a *ATHTTP) Comm(
	params HTTPParams,
) ([]byte, error) {

	// var resp *http.Response
	// var err error
	// var body []byte

	// switch params.Method {
	// case "GET":
	// 	resp, err = http.Get(params.Url)
	// case "POST":
	// 	resp, err = http.Post(params.Url, bytes.NewBuffer(body))
	// default:
	// 	return nil, errors.New("invalid method")
	// }

	// if err != nil {
	// 	return nil, errors.New("failed to send request")
	// }

	// defer resp.Body.Close()

	// body, err = io.ReadAll(resp.Body)
	// if err != nil {
	// 	return nil, errors.New("failed to read response body")
	// }

	// return body, nil
	var body io.Reader

	if params.Body != nil {
		jsonData, _ := json.Marshal(params.Body)
		body = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(params.Method, params.Url, body)
	if err != nil {
		return nil, err
	}

	for key, value := range params.Headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return resBody, nil
}
