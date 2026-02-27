// Package main is the entry point for the paperless-ngx-mcp server.
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/chrisallenlane/paperless-ngx-mcp/internal/client"
	"github.com/chrisallenlane/paperless-ngx-mcp/internal/server"
)

func main() {
	paperlessURL := os.Getenv("PAPERLESS_URL")
	if paperlessURL == "" {
		fmt.Fprintln(os.Stderr, "PAPERLESS_URL is required")
		os.Exit(1)
	}

	paperlessToken := os.Getenv("PAPERLESS_TOKEN")
	if paperlessToken == "" {
		fmt.Fprintln(os.Stderr, "PAPERLESS_TOKEN is required")
		os.Exit(1)
	}

	c := client.New(paperlessURL, paperlessToken)
	s := server.New(c)

	log.SetOutput(os.Stderr)

	if err := s.Run(context.Background(), os.Stdin, os.Stdout); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
