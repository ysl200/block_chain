package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func main() {
	// 定义请求URL
	baseurl := "http://localhost:8080/add?id=Node_"

	for i := 1; i <= 5; i++ {
		url := baseurl + strconv.Itoa(i)
		// 发送GET请求
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error making GET request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		// 读取响应体
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response body: %v\n", err)
			return
		}

		// 打印响应状态和内容
		fmt.Printf("Response status: %s\n", resp.Status)
		fmt.Printf("Response body: %s\n", string(body))
	}

}
