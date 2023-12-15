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
const DASHSCOPE_EMBEDDING_URL = "https://dashscope.aliyuncs.com/api/v1/services/embeddings/text-embedding/text-embedding"
const EMBEDDING_QUERY = "query"
const EMBEDDING_DOCUMENT = "document"

type embeddingItem struct {
	TextIndex int32     `json:"text_index"`
	Embedding []float64 `json:"embedding"`
}

type outputObject struct {
	Embeddings []embeddingItem `json:"embeddings"`
}

type returnObject struct {
	Output    outputObject     `json:"output"`
	Usage     map[string]int32 `json:"usage"`
	RequestId string           `json:"request_id"`
}

func main() {
	reqBody := map[string]interface{}{}
	msgs := []string{}
	msgs = append(msgs, "我在厦门上学")
	msgs = append(msgs, "湖北宜昌我就去过一次")
	// msgs = append(msgs, "你来过东北吗？")

	reqBody["model"] = "text-embedding-v2"
	reqBody["input"] = map[string][]string{
		"texts": msgs,
	}
	reqBody["parameters"] = map[string]string{
		"text_type": EMBEDDING_QUERY,
	}

	reqBodyJSON, _ := json.Marshal(reqBody)
	fmt.Printf("reqBodyJSON:%v\n", string(reqBodyJSON))
	payload := bytes.NewReader(reqBodyJSON)
	client := &http.Client{}
	req, err := http.NewRequest("POST", DASHSCOPE_EMBEDDING_URL, payload)

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

	// fmt.Printf("return: %v\n", string(resBody))

	respObj := map[string]interface{}{}
	err = json.Unmarshal(resBody, &respObj)
	if err != nil {
		fmt.Printf("ali_embedding Unmarshal response failed: err %v, resBody:%v\n", err, string(resBody))
		return
	}

	code, exist := respObj["code"]
	if exist {
		fmt.Printf("ali_embedding return error code: code:%v, message:%v, resBody:%v\n", cast.ToString(code), cast.ToString(respObj["message"]), string(resBody))
		return
	}

	returnObj := returnObject{}
	err = json.Unmarshal(resBody, &returnObj)
	if err != nil {
		fmt.Printf("ali_embedding parse returnObj failed: err %v, resBody:%v\n", err, string(resBody))
		return
	}

	if len(returnObj.Output.Embeddings) != len(msgs) {
		fmt.Printf("ali_embedding parse output failed: output len(%v) != input len(%v), resBody:%v\n", len(returnObj.Output.Embeddings), len(msgs), string(resBody))
		return
	}

	fmt.Printf("embeddings len:%v; usage:%v, request_id:%v\n", len(returnObj.Output.Embeddings), returnObj.Usage, returnObj.RequestId)
	for i := 0; i < len(returnObj.Output.Embeddings); i++ {
		fmt.Printf("embeddings[%d] len: %v\n", i, len(returnObj.Output.Embeddings[i].Embedding))
		for j, v := range returnObj.Output.Embeddings[i].Embedding {
			if j < 10 {
				fmt.Printf("[%d][%d] = %v\n", i, j, v)
			}
		}
	}
}
