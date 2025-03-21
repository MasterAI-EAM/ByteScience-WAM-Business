package gpt

import (
	"ByteScience-WAM-Business/conf"
	"ByteScience-WAM-Business/internal/utils"
	"ByteScience-WAM-Business/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

var Client *azopenai.Client

func NewGptClient(options *azopenai.ClientOptions) (*azopenai.Client, error) {
	keyCredential := azcore.NewKeyCredential(conf.GlobalConf.Gpt.ApiKey)
	client, err := azopenai.NewClientWithKeyCredential(conf.GlobalConf.Gpt.Endpoint, keyCredential, options)

	Client = client
	return client, err
}

func GetChatCompletions(messages []azopenai.ChatRequestMessageClassification) (*azopenai.GetChatCompletionsResponse, error) {
	resp, err := Client.GetChatCompletions(context.TODO(), azopenai.ChatCompletionsOptions{
		Messages:       messages,
		DeploymentName: &conf.GlobalConf.Gpt.Model,
	}, nil)

	if err != nil {
		fmt.Println("client.GetChatCompletions err:", err)
		return nil, err
	}

	return &resp, err
}

// ExtractInformationWithGPT4 获取gpt解析excel后的数据  新建excel格式
func ExtractInformationWithGPT4(data []string, prefix string, headerNum int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	headerData := data[:headerNum]
	bodyData := data[headerNum:]
	for _, row := range bodyData {
		combinedData := headerData
		combinedData = append(combinedData, row)
		prompt := fmt.Sprintf("%s%s", prefix, utils.ToJSON(combinedData))

		messages := []azopenai.ChatRequestMessageClassification{
			// 定义对话的系统信息
			&azopenai.ChatRequestSystemMessage{Content: azopenai.NewChatRequestSystemMessageContent("专业的材料实验数据分析助手。")},
			// 用户输入数据
			&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(prompt)},
		}

		resp, err := GetChatCompletions(messages)
		if err != nil {
			logger.Logger.Errorf("[ExtractInformationWithGPT4] GetChatCompletions err: %v", err)
			return nil, err
		}

		// 获取响应内容
		var responseText string
		if len(resp.Choices) > 0 && resp.Choices[0].Message != nil && resp.Choices[0].Message.Content != nil {
			responseText = *resp.Choices[0].Message.Content
		}

		// 清理响应文本
		if len(responseText) > 7 && responseText[:7] == "```json" {
			responseText = responseText[7:]
		}
		if len(responseText) > 3 && responseText[len(responseText)-3:] == "```" {
			responseText = responseText[:len(responseText)-3]
		}

		var responseData map[string]interface{}
		err = json.Unmarshal([]byte(responseText), &responseData)
		if err != nil {
			logger.Logger.Errorf("[ExtractInformationWithGPT4] GetChatCompletions err: %v", err)
			return nil, err
		}

		results = append(results, responseData)
	}

	return results, nil
}
