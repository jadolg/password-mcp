package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

const (
	letters               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lettersAndNumbers     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	lettersNumbersSymbols = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+"
)

func generatePassword(length int, charset string) (string, error) {
	password := make([]byte, length)
	for i := range password {
		idxBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", fmt.Errorf("crypto/rand error: %w", err)
		}
		password[i] = charset[idxBig.Int64()]
	}
	return string(password), nil
}

type GeneratePasswordArguments struct {
	Length  int    `json:"length" jsonschema:"required,description=The desired length of the password"`
	Charset string `json:"charset" jsonschema:"required,enum=letters,enum=lettersAndNumbers,enum=lettersNumbersSymbols,description=The character set to use (letters, lettersAndNumbers, lettersNumbersSymbols)"`
	Case    string `json:"case" jsonschema:"required,enum=uppercase,enum=lowercase,enum=mixed,description=The case to use for the password (uppercase, lowercase, mixed)"`
}

func main() {
	done := make(chan struct{})

	transport := stdio.NewStdioServerTransport()
	server := mcp.NewServer(transport, mcp.WithName("Password Generator Tool"), mcp.WithVersion("1.0.0"))

	err := server.RegisterTool("generate_password", "Generate a password with a given length, charset, and case", func(arguments GeneratePasswordArguments) (*mcp.ToolResponse, error) {
		var selectedCharset string
		switch arguments.Charset {
		case "letters":
			selectedCharset = letters
		case "lettersAndNumbers":
			selectedCharset = lettersAndNumbers
		case "lettersNumbersSymbols":
			selectedCharset = lettersNumbersSymbols
		default:
			return nil, fmt.Errorf("invalid charset: %s", arguments.Charset)
		}

		if arguments.Length <= 0 {
			return nil, fmt.Errorf("password length must be greater than 0")
		}

		password, err := generatePassword(arguments.Length, selectedCharset)
		if err != nil {
			return nil, err
		}

		switch arguments.Case {
		case "uppercase":
			password = strings.ToUpper(password)
		case "lowercase":
			password = strings.ToLower(password)
		case "mixed":
			// do nothing
		default:
			return nil, fmt.Errorf("invalid case: %s", arguments.Case)
		}

		return mcp.NewToolResponse(mcp.NewTextContent(fmt.Sprintf("Generated password: %s", password))), nil
	})

	if err != nil {
		log.Fatalf("Failed to register tool: %v", err)
	}

	if err := server.Serve(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
	<-done
}
