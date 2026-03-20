# Slack MCP Server

A Go-based MCP (Model Context Protocol) server providing comprehensive access to Slack workspace management, messaging, user lookup, and search capabilities.

## Overview

This server implements 30 Slack tools organized into logical categories:
- **Channel Management**: List, create, archive, and manage channels
- **Messaging**: Post messages, reply to threads, add reactions
- **Conversations**: Manage direct messages and conversations
- **Search**: Search messages and files across the workspace (requires user token)
- **User Management**: List users, get user info and presence

## Installation

### Prerequisites

- Go 1.21 or later
- Slack Bot Token (xoxb-...)
- Optional: Slack User Token (xoxp-...) for search functionality

### Building from Source

```bash
git clone https://github.com/karldane/slack-mcp.git
cd slack-mcp
make
```

This will automatically download dependencies and build a stripped binary.

#### Build Options

```bash
make              # Download deps and build (default)
make deps         # Download dependencies only
make build        # Build binary only (assumes deps exist)
make build-all    # Build for Linux, macOS, and Windows
make test         # Run tests
make clean        # Remove build artifacts
make install      # Install to GOPATH/bin
make help         # Show all options
```

### Environment Setup

Set the following environment variables:

```bash
# Required for most operations
export SLACK_BOT_TOKEN="xoxb-your-bot-token"

# Optional - required only for search tools
export SLACK_USER_TOKEN="xoxp-your-user-token"
```

### Obtaining Tokens

1. **Bot Token (SLACK_BOT_TOKEN)**:
   - Go to https://api.slack.com/apps
   - Create a new app or select existing
   - Navigate to "OAuth & Permissions"
   - Install the app to your workspace
   - Copy the "Bot User OAuth Token"

2. **User Token (SLACK_USER_TOKEN)** - Optional:
   - In the same OAuth section
   - Copy the "User OAuth Token"
   - Required for: search_messages, search_files, search_all

### Required Bot Scopes

Ensure your Slack app has these OAuth scopes:

**Channel Management:**
- `channels:read` - View public channels
- `channels:write` - Create and archive channels
- `groups:read` - View private channels
- `groups:write` - Manage private channels

**Messaging:**
- `chat:write` - Post messages
- `chat:write.public` - Post to public channels
- `reactions:read` - View reactions
- `reactions:write` - Add/remove reactions

**Users:**
- `users:read` - View user profiles
- `users:read.email` - View user email addresses

**Search:**
- `search:read` - Search messages and files (requires user token)

## Usage

### Starting the Server

```bash
# Basic usage (read-only mode)
./slack-mcp

# Enable write operations
./slack-mcp -write-enabled
```

### As an MCP Bridge Backend

Register the Slack server in your MCP Bridge configuration:

```json
{
  "backends": [
    {
      "id": "slack",
      "name": "Slack Workspace",
      "type": "slack",
      "command": "./slack-mcp",
      "self_reporting": true
    }
  ]
}
```

## Available Tools

### Channel Management (8 tools)

| Tool | Description | Risk Level | Requires Approval |
|------|-------------|------------|-------------------|
| `list_channels` | List all accessible channels | Low | No |
| `get_channel_info` | Get detailed channel information | Low | No |
| `create_channel` | Create a new channel | High | Yes |
| `archive_channel` | Archive a channel | High | Yes |
| `get_channel_history` | Get message history | Low | No |
| `join_channel` | Join a channel | Med | No |
| `leave_channel` | Leave a channel | Med | No |
| `set_channel_topic` | Set channel topic | Med | No |

### Messaging (9 tools)

| Tool | Description | Risk Level | Requires Approval |
|------|-------------|------------|-------------------|
| `post_message` | Post a message to a channel | Med | No |
| `reply_to_thread` | Reply to a thread | Med | No |
| `get_thread_replies` | Get thread replies | Low | No |
| `add_reaction` | Add emoji reaction | Low | No |
| `remove_reaction` | Remove emoji reaction | Low | No |
| `get_message_reactions` | Get reactions on a message | Low | No |
| `update_message` | Edit a message | Med | No |
| `delete_message` | Delete a message | High | Yes |
| `send_dm` | Send a direct message | Med | No |

### Conversations (3 tools)

| Tool | Description | Risk Level | Requires Approval |
|------|-------------|------------|-------------------|
| `list_conversations` | List IM/MPIM conversations | Low | No |
| `get_dm_history` | Get DM history | Low | No |
| `open_dm` | Open/create a DM | Low | No |

### Search (3 tools)

**Note:** Requires `SLACK_USER_TOKEN`

| Tool | Description | Risk Level | Requires Approval |
|------|-------------|------------|-------------------|
| `search_messages` | Search messages | Low | No |
| `search_files` | Search files | Low | No |
| `search_all` | Search both messages and files | Low | No |

### User Management (7 tools)

| Tool | Description | Risk Level | Requires Approval |
|------|-------------|------------|-------------------|
| `list_users` | List all users | Low | No |
| `get_user_info` | Get user details | Low | No |
| `get_user_presence` | Get user presence status | Low | No |
| `lookup_user_by_email` | Find user by email | Low | No |
| `get_user_profile` | Get detailed profile | Low | No |
| `get_bot_info` | Get bot information | Low | No |
| `get_team_info` | Get workspace info | Low | No |

## Self-Reporting Capability

This server supports **EnforcerProfile self-reporting**, allowing the MCP Bridge to automatically discover tool capabilities and safety metadata without requiring API keys at startup.

### How It Works

1. **No Tokens Required for Registration**: The server registers all tools and their schemas even without Slack tokens
2. **Runtime Token Validation**: Tokens are only required when executing tools
3. **Safety Metadata**: Each tool reports its risk level, impact scope, resource cost, and approval requirements

### Example Self-Reported Metadata

```json
{
  "tool": "delete_message",
  "enforcer_profile": {
    "risk_level": "high",
    "impact_scope": "delete",
    "resource_cost": 3,
    "pii": true,
    "idempotent": false,
    "approval_required": true
  }
}
```

## Error Handling

The server provides clear error messages for common scenarios:

- **Missing Bot Token**: `SLACK_BOT_TOKEN is not set`
- **Missing User Token**: `SLACK_USER_TOKEN is not set (required for search)`
- **Invalid Channel**: `channel_not_found`
- **Permission Denied**: `not_in_channel` or `missing_scope`
- **Rate Limiting**: Automatically retries with exponential backoff

## Testing

Run the test suite:

```bash
cd mcp-framework
go test ./slack/... -v
```

### Test Coverage

The test suite includes:
- Server creation and initialization
- Tool registration verification
- EnforcerProfile validation
- Schema correctness
- Argument parsing helpers
- Client requirement checks
- Self-reporting capability tests

## Architecture

### Components

```
slack/
├── client.go         # Slack API client wrapper
├── slack.go          # Server initialization and tool registration
├── slack_tools.go    # All 30 tool implementations
├── slack_test.go     # Comprehensive test suite
└── README.md         # This file
```

### Key Design Decisions

1. **Token Flexibility**: Server starts without tokens for self-reporting, fails gracefully at runtime
2. **Enforcer Profiles**: Every tool reports safety metadata for automated policy enforcement
3. **Conditional Registration**: Search tools require user token but are always registered for discovery
4. **Write Protection**: All write operations are flagged with appropriate risk levels

## Security Considerations

### Data Handling
- No credentials are logged
- Tokens are read from environment only
- All API responses are properly sanitized

### Access Control
- Read operations: Low risk, no approval needed
- Write operations: Medium/High risk, may require approval
- Delete operations: High risk, always requires approval

### Rate Limiting
- Respects Slack's rate limits
- Implements automatic retry with backoff
- Reports rate limit errors clearly

## Troubleshooting

### Common Issues

**"channel_not_found" error**
- Ensure the bot is invited to the channel
- Check that the channel ID is correct

**"not_in_channel" error**
- The bot must join the channel first using `join_channel`

**"missing_scope" error**
- Your Slack app is missing required OAuth scopes
- Reinstall the app with updated scopes

**Search tools not working**
- Verify `SLACK_USER_TOKEN` is set
- User token must have `search:read` scope

### Debug Mode

Enable verbose logging:

```bash
export SLACK_DEBUG=1
./slack-mcp
```

## Contributing

When adding new tools:

1. Implement the `framework.ToolHandler` interface
2. Define appropriate JSON schema
3. Set correct EnforcerProfile metadata
4. Add comprehensive tests
5. Update this README

## License

This project is part of the MCP Bridge framework. See the main repository for license information.

## See Also

- [MCP Bridge Documentation](../README.md)
- [Slack API Documentation](https://api.slack.com/)
- [mcp-go Library](https://github.com/mark3labs/mcp-go)
