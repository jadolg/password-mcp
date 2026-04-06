package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
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

func handleGeneratePassword(_ context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	length, err := req.RequireInt("length")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}
	if length <= 0 {
		return mcp.NewToolResultError("password length must be greater than 0"), nil
	}

	charset, err := req.RequireString("charset")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	var selectedCharset string
	switch charset {
	case "letters":
		selectedCharset = letters
	case "lettersAndNumbers":
		selectedCharset = lettersAndNumbers
	case "lettersNumbersSymbols":
		selectedCharset = lettersNumbersSymbols
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid charset: %s", charset)), nil
	}

	caseOpt, err := req.RequireString("case")
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	password, err := generatePassword(length, selectedCharset)
	if err != nil {
		return nil, err
	}

	switch caseOpt {
	case "uppercase":
		password = strings.ToUpper(password)
	case "lowercase":
		password = strings.ToLower(password)
	case "mixed":
		// do nothing
	default:
		return mcp.NewToolResultError(fmt.Sprintf("invalid case: %s", caseOpt)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Generated password: %s", password)), nil
}

func main() {
	useHTTP := flag.Bool("http", false, "Use HTTP transport instead of stdio")
	addr := flag.String("addr", ":8080", "Address to listen on when using HTTP transport")
	flag.Parse()

	s := server.NewMCPServer("Password Generator Tool", "1.0.0",
		server.WithToolCapabilities(true),
	)

	s.AddTool(
		mcp.NewTool("generate_password",
			mcp.WithDescription("Generate a password with a given length, charset, and case"),
			mcp.WithReadOnlyHintAnnotation(true),
			mcp.WithDestructiveHintAnnotation(false),
			mcp.WithOpenWorldHintAnnotation(false),
			mcp.WithNumber("length",
				mcp.Required(),
				mcp.Description("The desired length of the password"),
			),
			mcp.WithString("charset",
				mcp.Required(),
				mcp.Description("The character set to use"),
				mcp.Enum("letters", "lettersAndNumbers", "lettersNumbersSymbols"),
			),
			mcp.WithString("case",
				mcp.Required(),
				mcp.Description("The case to use for the password"),
				mcp.Enum("uppercase", "lowercase", "mixed"),
			),
		),
		handleGeneratePassword,
	)

	if *useHTTP {
		httpServer := server.NewStreamableHTTPServer(s)
		log.Printf("Listening on %s/mcp", *addr)
		if err := httpServer.Start(*addr); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
