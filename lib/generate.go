package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (client *Client) PostJSON(endpoint string, payload []byte) (*http.Response, error) {
	url := client.baseUrl + endpoint
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+client.apiKey)

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %v", err)
	}

	return resp, nil
}

type FunctionCallCompletion struct {
	Name      string
	Arguments map[string]string
}

// Use built-in special UnmarshalJSON method https://pkg.go.dev/encoding/json#example-package-CustomMarshalJSON
// FunctionCall implements json.Unmarshaler that gets invoked automatically with a FunctionCall object
func (functionCall *FunctionCallCompletion) UnmarshalJSON(data []byte) error {

	type Alias FunctionCallCompletion

	tempStruct := &struct {
		Arguments json.RawMessage
		*Alias
	}{
		Alias: (*Alias)(functionCall),
	}

	if err := json.Unmarshal(data, &tempStruct); err != nil {
		return err
	}

	err := json.Unmarshal(tempStruct.Arguments, &tempStruct.Arguments)
	if err != nil {
		return err
	}

	return nil
}

// The function structs are simply a way to not call raw json..
// It's based on the raw json from: https://openai.com/blog/function-calling-and-other-api-updates

type CompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role         string `json:"role"`
			Content      any    `json:"content,omitempty"`
			FunctionCall struct {
				Name      string `json:"name"`
				Arguments string `json:"arguments"`
			} `json:"function_call"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Custom function for function calling
type ShortAndConcise struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type Function struct {
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Parameters   FunctionParameter `json:"parameters"`
	FunctionCall FunctionCall      `json:"function_call"`
}

type FunctionParameterProperties struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

type FunctionParameter struct {
	Type       string                                 `json:"type"`
	Properties map[string]FunctionParameterProperties `json:"properties"`
	Required   []string                               `json:"required"`
}
type FunctionCall struct {
	Name string `json:"name"`
}

type FunctionCallPayload struct {
	Model     string     `json:"model"`
	Functions []Function `json:"functions,omitempty"`
	Messages  []Message  `json:"messages"`
}

func (client *Client) GenerateOutput(ctx context.Context, model string, query string) (res CompletionResponse, err error) {
	payload := FunctionCallPayload{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: query,
			},
		},
		Functions: []Function{
			{
				Name:        "acli",
				Description: "Outputs a command according to the request without additional explanation",
				Parameters: FunctionParameter{
					Type: "object",
					Properties: map[string]FunctionParameterProperties{
						"short_and_concise": {
							Type:        "string",
							Description: "A concise command for the users query. Preferably using built-in filtering commands like --json, --query, --filter or a combination.",
						},
					},
					Required: []string{"short_and_concise"},
				},
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return res, fmt.Errorf("error %w", err)
	}

	// Use the PostJSON method to send the request
	resp, err := client.PostJSON("/v1/chat/completions", jsonPayload)
	if err != nil {
		return res, fmt.Errorf("failed to generate output: %v", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return res, fmt.Errorf("failed to read response body: %v", err)
	}

	var response CompletionResponse
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return res, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	return response, err
}
