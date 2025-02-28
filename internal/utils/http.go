package utils

import (
	"ByteScience-WAM-Business/internal/model/dto/ai"
	"ByteScience-WAM-Business/pkg/logger"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var aiUrl = "http://96.9.229.193:52222"
var predictionPath = "/data/inference/prediction"

func SendPredictionRequest(request *ai.ForwardDirectionRequest) (*ai.ForwardDirectionResult, error) {
	// 将请求体序列化为 JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		logger.Logger.Errorf("[SendPredictionRequest] json.Marshal err: %v", err)
		return nil, err
	}

	url := aiUrl + predictionPath
	// 创建一个 POST 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Logger.Errorf("[SendPredictionRequest] failed to create request: %v", err)
		return nil, err
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Logger.Errorf("[SendPredictionRequest] failed to send request: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	// 读取并打印响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Logger.Errorf("[SendPredictionRequest] Error reading response body: %v", err)
		return nil, err
	}

	logger.Logger.Infof("[SendPredictionRequest] response body: %v", string(body))

	// 定义响应结构体
	var response struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}

	// 解析响应内容为 response
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	// 如果 code 不为 0，表示发生了错误
	if response.Code != 0 {
		return nil, NewBusinessError(ExternalRequestError, response.Message)
	}

	// 解析响应内容为 ForwardDirectionResult
	var result ai.ForwardDirectionResult
	if err := json.Unmarshal(body, &result); err != nil {
		logger.Logger.Errorf("[SendPredictionRequest] failed to unmarshal response body: %v", err)
		return nil, err
	}

	return &result, nil
}
