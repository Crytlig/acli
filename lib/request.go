package lib

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/crytlig/acli/bubble"

	"github.com/urfave/cli/v2"
	"golang.design/x/clipboard"
)

func HandleRequest(c *cli.Context, query string, debugMode bool) error {
	// Handle recursive calls
	quitting := bubble.Quitting

	if quitting {
		return nil
	}

	apiKey, err := LoadApiKey()
	if err != nil {
		log.Fatalf("Failed to load API key: %v", err)
	}

	client, err := NewClient(&NewClientOpts{ApiKey: apiKey})
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	model := GPT3Dot5Turbo0613.Name
	resp, err := client.GenerateOutput(context.Background(), model, query)
	if err != nil {
		log.Fatalf("Failed to generate output: %v", err)
	}

	var cliCmd map[string]string
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.FunctionCall.Arguments), &cliCmd); err != nil {
		return fmt.Errorf("failed to unmarshal CLI command: %v", err)
	}

	aiCmd := cliCmd["short_and_concise"]
	bubble.AcceptInput(fmt.Sprintf("Suggested: %s", aiCmd))

	var cmd *exec.Cmd

	switch bubble.UserChoice {
	case "Accept":
		cliName := strings.Split(aiCmd, " ")[0]
		cliPath, err := exec.LookPath(cliName)
		if err != nil {
			log.Fatalf("Unable to locate requested CLI: %v", err)
		}

		cmdArgs := prepareCommandArgs(aiCmd)
		// Prints the actual command. TUI hides this, so it's useful to have saved
		fmt.Println(aiCmd)
		if debugMode {
			log.Println("Location of binary:", cliPath)
		}

		cmd = exec.Command(cliPath, cmdArgs...)
		if cmd != nil && debugMode {
			log.Printf("Full execed command: %s", cmd.String())
		}

		executeCommand(cmd, debugMode)
	case "Retry":
		if err := HandleRequest(c, c.Args().First(), debugMode); err != nil {
			log.Fatalf("Retry failed: %v", err)
		}

	case "Copy to clipboard":
		err := clipboard.Init()
		if err != nil {
			panic(err)
		}
		clipboard.Write(clipboard.FmtText, []byte(aiCmd))
	case "Rephrase":
		bubble.UserTextInputs.PrevQuery = query
		bubble.RephraseInput()
		// We have to check again to ensure we're not forcing another recursive call since multiple TUIs are opened
		// if the user asks for rephrasing twice or thrice or more
		if bubble.Quitting {
			return nil
		}

		if err := HandleRequest(c, bubble.UserTextInputs.UserText, debugMode); err != nil {
			log.Fatalf("Rephrase failed: %v", err)
		}
		return nil

	// User presses q to quit
	case "":
		os.Exit(0)
	// Add default. Yet to be encountered
	default:
		log.Fatalf("Unsupported user choice: %s", bubble.UserChoice)
	}
	return nil
}

func prepareCommandArgs(cmd string) []string {
	noSingleQuotes := strings.ReplaceAll(cmd, "'", "")

	// final is the list of command arguments after splitting the command
	// example: az get application id of app registration myapp123 -> final = ["get", "application", "id", "of", "app", "registration", "myapp123"]
	// note: there could be spaces in the command arguments
	final := strings.Split(strings.SplitN(noSingleQuotes, " ", 2)[1], " ")
	return final
}

func executeCommand(cmd *exec.Cmd, debugMode bool) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	fmt.Println(stdout.String())
	if err := cmd.Run(); err != nil {
		if debugMode {
			fmt.Println("Unable to fire command")
			fmt.Println("Error:", err)
			fmt.Println("Stderr:", stderr.String())
		}
		log.Fatal(err)
	}
	fmt.Println(stdout.String())
}
