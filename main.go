package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const OPENAI_API_URL string = "https://api.openai.com/v1/chat/completions"

type Response struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
			Refusal string `json:"refusal"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens        int `json:"prompt_tokens"`
		CompletionTokens    int `json:"completion_tokens"`
		TotalTokens         int `json:"total_tokens"`
		PromptTokensDetails struct {
			CachedTokens int `json:"cached_tokens"`
			AudioTokens  int `json:"audio_tokens"`
		} `json:"prompt_tokens_details"`
	} `json:"usage"`
	CompletionTokenDetails struct {
		ReasoningTokens          int `json:"reasoning_tokens"`
		AudioTokens              int `json:"audio_tokens"`
		AcceptedPredictionTokens int `json:"accepted_prediction_tokens"`
		RejectedPredictionTokens int `json:"rejected_prediction_tokens"`
	} `json:"completion_token_details"`
	ServiceTier       string `json:"service_tier"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func main() {
	if len(os.Args) < 2 {
		usage(1)
	}

	logLevel := getEnvWithDefault("JEEVES_LOG_LEVEL", "info")
	apiKey := getEnv("OPENAI_API_KEY")
	model := getEnvWithDefault("JEEVES_OPENAI_MODEL", "gpt-4o-mini")

	if logLevel == "debug" {
		fmt.Println("Using model:", model)
	}

	prompt := strings.Join(os.Args[1:], " ")
	if logLevel == "debug" {
		fmt.Println("User prompt:", prompt)
	}

	requestBody := map[string]any{
		"model": model,
		"store": true,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	postBody, err := json.Marshal(requestBody)
	if err != nil {
		fmt.Printf("Error marshalling json request body: %v\n", err)
		os.Exit(1)
	}

	req, err := http.NewRequest(http.MethodPost, OPENAI_API_URL, bytes.NewBuffer(postBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("X-OpenAI-Use-Case", "government")
	req.Header.Set("X-OpenAI-Data-Usage-Opt-Out", "true")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request to OpenAI API: %v\n", err)
		os.Exit(1)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		os.Exit(1)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error unmarshalling response body: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n%s\n\n", result.Choices[0].Message.Content)
}

func usage(exitCode int) {
	fmt.Println(`
Usage: jeeves <prompt>

Environment variables:
OPENAI_API_KEY			your OpenAI API key (mandatory)
JEEVES_OPENAI_MODEL		specify the model to use, default is "gpt-4o-mini" (optional)
JEEVES_LOG_LEVEL		sets log level, default is "info" (optional)
			`)
	os.Exit(exitCode)
}

func getEnvWithDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Printf("%s environment variable not found, please export it via the shell and try again\n", key)
		os.Exit(1)
	}
	return value
}
