package utils

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func PatchSign(key string, data url.Values) url.Values {
	time := time.Now().Unix()
	md5Key := fmt.Sprintf("%x", md5.Sum([]byte(key)))
	data.Set("request_time", fmt.Sprint(time))
	data.Set("request_token", fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprint(time)+md5Key))))

	return data
}

func Post(url string, data url.Values) string {
	// 超时时间：5秒
	client := &http.Client{Timeout: 20 * time.Second}
	response, err := client.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()

	result, _ := ioutil.ReadAll(response.Body)
	return string(result)
}
