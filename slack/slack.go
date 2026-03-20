package slack

import (
	"github.com/karldane/mcp-framework/framework"
)

type Server struct {
	*framework.Server
	client *Client
}

func NewServer() *Server {
	config := &framework.Config{
		Name:    "slack-mcp",
		Version: "1.0.0",
		Instructions: `Slack MCP Server providing access to Slack workspace management, messaging, user lookup, and search tools.

Authentication: Requires SLACK_BOT_TOKEN (xoxb-...) environment variable. Optional: SLACK_USER_TOKEN (xoxp-...) for search operations.

Channel management: list_channels, get_channel_info, create_channel, archive_channel, get_channel_history, join_channel, leave_channel, set_channel_topic.

Messaging: post_message, reply_to_thread, get_thread_replies, add_reaction, remove_reaction, get_message_reactions, update_message, delete_message, send_dm, get_dm_history.

Conversations: list_conversations, open_dm.

Search (requires SLACK_USER_TOKEN): search_messages, search_files, search_all.

User management: list_users, get_user_info, get_user_presence, lookup_user_by_email, get_user_profile, get_bot_info, get_team_info.`,
	}

	s := &Server{
		Server: framework.NewServerWithConfig(config),
		client: NewClient(),
	}

	s.registerTools()
	return s
}

func (s *Server) registerTools() {
	s.Server.RegisterTool(&ListChannelsTool{client: s.client})
	s.Server.RegisterTool(&GetChannelInfoTool{client: s.client})
	s.Server.RegisterTool(&CreateChannelTool{client: s.client})
	s.Server.RegisterTool(&ArchiveChannelTool{client: s.client})
	s.Server.RegisterTool(&GetChannelHistoryTool{client: s.client})
	s.Server.RegisterTool(&JoinChannelTool{client: s.client})
	s.Server.RegisterTool(&LeaveChannelTool{client: s.client})
	s.Server.RegisterTool(&SetChannelTopicTool{client: s.client})

	s.Server.RegisterTool(&PostMessageTool{client: s.client})
	s.Server.RegisterTool(&ReplyToThreadTool{client: s.client})
	s.Server.RegisterTool(&GetThreadRepliesTool{client: s.client})
	s.Server.RegisterTool(&AddReactionTool{client: s.client})
	s.Server.RegisterTool(&RemoveReactionTool{client: s.client})
	s.Server.RegisterTool(&GetMessageReactionsTool{client: s.client})
	s.Server.RegisterTool(&UpdateMessageTool{client: s.client})
	s.Server.RegisterTool(&DeleteMessageTool{client: s.client})
	s.Server.RegisterTool(&SendDMTool{client: s.client})
	s.Server.RegisterTool(&ListConversationsTool{client: s.client})
	s.Server.RegisterTool(&GetDMHistoryTool{client: s.client})
	s.Server.RegisterTool(&OpenDMTool{client: s.client})

	s.Server.RegisterTool(&SearchMessagesTool{client: s.client})
	s.Server.RegisterTool(&SearchFilesTool{client: s.client})
	s.Server.RegisterTool(&SearchAllTool{client: s.client})

	s.Server.RegisterTool(&ListUsersTool{client: s.client})
	s.Server.RegisterTool(&GetUserInfoTool{client: s.client})
	s.Server.RegisterTool(&GetUserPresenceTool{client: s.client})
	s.Server.RegisterTool(&LookupUserByEmailTool{client: s.client})
	s.Server.RegisterTool(&GetUserProfileTool{client: s.client})
	s.Server.RegisterTool(&GetBotInfoTool{client: s.client})
	s.Server.RegisterTool(&GetTeamInfoTool{client: s.client})
}

func (s *Server) Initialize() {
	s.Server.Initialize()
}

func (s *Server) Start() error {
	return s.Server.Start()
}

func (s *Server) SetWriteEnabled(enabled bool) {
	s.Server.SetWriteEnabled(enabled)
}

func (s *Server) IsWriteEnabled() bool {
	return s.Server.IsWriteEnabled()
}
