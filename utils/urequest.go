package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// http.get
func HttpGet(url string, retObj interface{}) error {
	return HttpGetWithToken(url, "", retObj)
}

func HttpDelete(url string, retObj interface{}) error {
	req, _ := http.NewRequest("DELETE", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}
	err = json.Unmarshal(body, retObj)
	return err
}

// http.get with token
func HttpGetWithToken(url, token string, retObj interface{}) error {
	req, _ := http.NewRequest("GET", url, nil)
	if token != "" {
		req.Header.Add("Authorization", token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}
	err = json.Unmarshal(body, retObj)
	return err
}

// http.post
func HttpPost(url string, params interface{}, retObj interface{}) error {
	buf := new(bytes.Buffer)
	err := json.NewEncoder(buf).Encode(params)
	if err != nil {
		log.Printf("Encode json failed: %+v\n", err)
		return err
	}
	resp, err := http.Post(url, "application/json; charset=utf-8", buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return errors.New(string(body))
	}
	err = json.Unmarshal(body, retObj)
	return err
}

// http post with header
func HttpPostWithHeader(urlStr string, header map[string]string, body interface{}, retObj interface{}) error {
	bodyArray, err := json.Marshal(body)
	fmt.Printf("%s", string(bodyArray))
	if err != nil {
		return err
	}
	respByte, err := sendRequest(urlStr, header, bytes.NewReader(bodyArray))
	if err != nil {
		return err
	}
	err = json.Unmarshal(respByte, &retObj)
	if err != nil {
		return fmt.Errorf("response json Unmarshall failed, http response body: %s", respByte)
	}
	return nil
}

func sendRequest(urlStr string, header map[string]string, body io.Reader) ([]byte, error) {
	request, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, fmt.Errorf("webRequest failed: %v", err)
	}
	for k, v := range header {
		request.Header.Set(k, v)
	}
	client := http.Client{}
	resp, err := client.Do(request)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("doReqeust failed: %v", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("http StatusCode is not in 2XX range, StatusCode: %d, Body: %s", resp.StatusCode, resp.Body)
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("readAll failed: %v", err)
	}
	jsonStr := string(respBytes)
	jsonByte := []byte(jsonStr)
	jsonByteLen := len(jsonByte)
	return jsonByte[:jsonByteLen], nil
}
