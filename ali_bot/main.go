package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cast"
)

const DASHSCOPE_API_KEY = "sk-1e105db03d56436c89695fa27ebded5e"
const DASHSCOPE_GENERATION_URL = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
const EMBEDDING_QUERY = "query"
const EMBEDDING_DOCUMENT = "document"

type llmMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type outputObject struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"`
}

type returnObject struct {
	Output    outputObject     `json:"output"`
	Usage     map[string]int32 `json:"usage"`
	RequestId string           `json:"request_id"`
}

func main() {
	reqBody := map[string]interface{}{}

	msgs := []llmMsg{}
	msgs = append(msgs, llmMsg{
		Role:    "user",
		Content: prompt,
	})

	reqBody["model"] = "qwen-plus"
	reqBody["input"] = map[string]interface{}{
		"messages": msgs,
	}
	reqBody["parameters"] = map[string]interface{}{
		"temperature":   0.001,
		"result_format": "text",
	}

	reqBodyJSON, _ := json.Marshal(reqBody)
	fmt.Printf("reqBodyJSON:%v\n", string(reqBodyJSON))
	payload := bytes.NewReader(reqBodyJSON)
	client := &http.Client{}
	req, err := http.NewRequest("POST", DASHSCOPE_GENERATION_URL, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", DASHSCOPE_API_KEY))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("return: %v\n", string(resBody))

	respObj := map[string]interface{}{}
	err = json.Unmarshal(resBody, &respObj)
	if err != nil {
		fmt.Printf("ali_bot Unmarshal response failed: err %v, resBody:%v\n", err, string(resBody))
		return
	}

	code, exist := respObj["code"]
	if exist {
		fmt.Printf("ali_bot return error code: code:%v, message:%v, resBody:%v\n", cast.ToString(code), cast.ToString(respObj["message"]), string(resBody))
		return
	}

	returnObj := returnObject{}
	err = json.Unmarshal(resBody, &returnObj)
	if err != nil {
		fmt.Printf("ali_bot parse returnObj failed: err %v, resBody:%v\n", err, string(resBody))
		return
	}

	if returnObj.Output.FinishReason != "stop" {
		fmt.Printf("ali_bot return failed: FinishReason:%v, body:%v\n", returnObj.Output.FinishReason, string(resBody))
		return
	}

	fmt.Printf("output text:%v; usage:%v, request_id:%v\n", returnObj.Output.Text, returnObj.Usage, returnObj.RequestId)
}
