package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/karldane/mcp-framework/framework"
	"github.com/karldane/slack-mcp/slack"
)

func main() {
	writeEnabled := flag.Bool("write-enabled", false, "Enable write tools")
	flag.Parse()

	server := slack.NewServer()

	if framework.HandleScanFlag(server.Server) {
		return
	}

	botToken := os.Getenv("SLACK_BOT_TOKEN")
	if botToken == "" {
		fmt.Fprintln(os.Stderr, "Note: SLACK_BOT_TOKEN not set - server starting in disconnected mode")
		fmt.Fprintln(os.Stderr, "Tools are registered but require SLACK_BOT_TOKEN to execute")
	}

	userToken := os.Getenv("SLACK_USER_TOKEN")
	if userToken == "" {
		fmt.Fprintln(os.Stderr, "Note: SLACK_USER_TOKEN not set - search tools will not be registered")
	}

	server.SetWriteEnabled(*writeEnabled)
	server.Initialize()
	server.Start()
}
