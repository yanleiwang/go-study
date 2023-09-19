package http_

import (
	"fmt"
	"io"
	"net/http"
	"testing"
)

func TestTips(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		readBodyOnce(writer, request)
	})

	http.ListenAndServe("localhost:8080", nil)

}

func readBodyOnce(writer http.ResponseWriter, request *http.Request) {
	body, err := io.ReadAll(request.Body)
	if err != nil {
		fmt.Fprintf(writer, "read body failed: %v", err)
		return
	}

	fmt.Fprintf(writer, "read the data: %s\n", string(body))

	// 再次读取， 啥也读不到， 但是也不会报错
	body, err = io.ReadAll(request.Body)
	if err != nil {
		// 不会进来这里
		fmt.Fprintf(writer, "read the data one more time got error: %v", err)
		return
	}
	fmt.Fprintf(writer, "read the data one more time: [%s] and read data length %d \n", string(body), len(body))

}
