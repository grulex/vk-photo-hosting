package client

import "fmt"

type VkError struct {
	ErrorCode     int    `json:"error_code"`
	ErrorMsg      string `json:"error_msg"`
	Method        string `json:"method"`
	RequestParams []struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"request_params"`
}

func (e VkError) Error() string {
	if len(e.Method) == 0 {
		return e.ErrorMsg
	}

	return fmt.Sprintf("%v:%v", e.Method, e.ErrorMsg)
}
