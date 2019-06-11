package httptool

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

/*
GET httptool.GET
*/
func GET(url string, params map[string]string) []byte {
	if params != nil {
		var i = 0
		for k, v := range params {
			if i == 0 {
				url = url + "?" + k + "=" + v
			} else {
				url = url + "&" + k + "=" + v
			}
			i++
		}
	}
	fmt.Println(url)
	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Error("request get err,", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("get read body err,", err)
		return nil
	}
	return body
}

/*
POST httptool.POST
*/
func POST(url string, params map[string]string, data interface{}) []byte {
	if params != nil {
		var i = 0
		for k, v := range params {
			if i == 0 {
				url = url + "?" + k + "=" + v
			} else {
				url = url + "&" + k + "=" + v
			}
			i++
		}
	}
	client := &http.Client{}
	msg, _ := json.Marshal(data)
	request, _ := http.NewRequest("POST", url, bytes.NewReader(msg))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Connection", "Keep-Alive")
	response, err := client.Do(request)
	if err != nil {
		log.Error("request post err,", err)
		return nil
	}
	body, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		log.Error("post read body err,", err)
		return nil
	}

	return body
}
