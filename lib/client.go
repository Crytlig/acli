package lib

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadApiKey() (string, error) {
	// Check for OPENAI_API_KEY environment variable
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		// Otherwise load via .env
		err := godotenv.Load(".env")
		if err != nil {
			return "", fmt.Errorf("error loading .env file")
		}
		// Check again to see if environment variable has been loaded
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if apiKey == "" {
		return "", fmt.Errorf("API Key is not set. Please set OPENAI_API_KEY environment variable or .env file")
	}

	return apiKey, nil
}

const openai_backend = "https://api.openai.com"

type Client struct {
	apiKey  string
	baseUrl string
}

type NewClientOpts struct {
	ApiKey string
	Url    string
}

func NewClient(opts *NewClientOpts) (*Client, error) {

	if opts == nil || opts.ApiKey == "" {
		return nil, fmt.Errorf("invalid options for creating NewClient")
	}

	if opts.Url == "" {
		opts.Url = openai_backend
	}

	client := &Client{
		apiKey:  opts.ApiKey,
		baseUrl: opts.Url,
	}

	return client, nil
}
