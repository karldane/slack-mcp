package slack

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/karldane/mcp-framework/framework"
	"github.com/mark3labs/mcp-go/mcp"
)

// ============================================================================
// CHANNEL TOOLS
// ============================================================================

type ListChannelsTool struct{ client *Client }

func (t *ListChannelsTool) Name() string { return "list_channels" }
func (t *ListChannelsTool) Description() string {
	return "List all accessible Slack channels in the workspace."
}
func (t *ListChannelsTool) Schema() mcp.ToolInputSchema { return listChannelsSchema() }
func (t *ListChannelsTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	types := getString(args, "types", "public_channel,private_channel")
	excludeArchived := getBool(args, "exclude_archived", true)
	limit := getInt(args, "limit", 100)
	channels, _, err := t.client.ListChannels(ctx, types, excludeArchived, limit)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(channels)
	return framework.TextResult(string(b)), nil
}
func (t *ListChannelsTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

type GetChannelInfoTool struct{ client *Client }

func (t *GetChannelInfoTool) Name() string { return "get_channel_info" }
func (t *GetChannelInfoTool) Description() string {
	return "Get detailed information about a specific Slack channel."
}
func (t *GetChannelInfoTool) Schema() mcp.ToolInputSchema { return getChannelInfoSchema() }
func (t *GetChannelInfoTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	ch, err := t.client.GetChannelInfo(ctx, channelID)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(ch)
	return framework.TextResult(string(b)), nil
}
func (t *GetChannelInfoTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

type CreateChannelTool struct{ client *Client }

func (t *CreateChannelTool) Name() string                { return "create_channel" }
func (t *CreateChannelTool) Description() string         { return "Create a new Slack channel." }
func (t *CreateChannelTool) Schema() mcp.ToolInputSchema { return createChannelSchema() }
func (t *CreateChannelTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	name := getRequiredString(args, "name")
	isPrivate := getBool(args, "is_private", false)
	ch, err := t.client.CreateChannel(ctx, name, isPrivate)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(ch)
	return framework.TextResult(string(b)), nil
}
func (t *CreateChannelTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskHigh),
		framework.WithImpact(framework.ImpactWrite),
		framework.WithResourceCost(6),
		framework.WithPII(false),
		framework.WithIdempotent(false),
		framework.WithApprovalReq(true),
	)
}

type ArchiveChannelTool struct{ client *Client }

func (t *ArchiveChannelTool) Name() string                { return "archive_channel" }
func (t *ArchiveChannelTool) Description() string         { return "Archive a Slack channel." }
func (t *ArchiveChannelTool) Schema() mcp.ToolInputSchema { return archiveChannelSchema() }
func (t *ArchiveChannelTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	if err := t.client.ArchiveChannel(ctx, channelID); err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true, "channel_id": %q}`, channelID)), nil
}
func (t *ArchiveChannelTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskHigh),
		framework.WithImpact(framework.ImpactDelete),
		framework.WithResourceCost(5),
		framework.WithPII(false),
		framework.WithIdempotent(false),
		framework.WithApprovalReq(true),
	)
}

type GetChannelHistoryTool struct{ client *Client }

func (t *GetChannelHistoryTool) Name() string { return "get_channel_history" }
func (t *GetChannelHistoryTool) Description() string {
	return "Fetch message history from a Slack channel."
}
func (t *GetChannelHistoryTool) Schema() mcp.ToolInputSchema { return getChannelHistorySchema() }
func (t *GetChannelHistoryTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	limit := getInt(args, "limit", 20)
	oldest := getString(args, "oldest", "")
	latest := getString(args, "latest", "")
	history, err := t.client.GetChannelHistory(ctx, channelID, limit, oldest, latest)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(history)
	return framework.TextResult(string(b)), nil
}
func (t *GetChannelHistoryTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(3),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type JoinChannelTool struct{ client *Client }

func (t *JoinChannelTool) Name() string                { return "join_channel" }
func (t *JoinChannelTool) Description() string         { return "Join a Slack channel." }
func (t *JoinChannelTool) Schema() mcp.ToolInputSchema { return joinChannelSchema() }
func (t *JoinChannelTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	ch, _, err := t.client.JoinChannel(ctx, channelID)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(ch)
	return framework.TextResult(string(b)), nil
}
func (t *JoinChannelTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskMed),
		framework.WithImpact(framework.ImpactWrite),
		framework.WithResourceCost(3),
		framework.WithPII(false),
		framework.WithIdempotent(false),
	)
}

type LeaveChannelTool struct{ client *Client }

func (t *LeaveChannelTool) Name() string                { return "leave_channel" }
func (t *LeaveChannelTool) Description() string         { return "Leave a Slack channel." }
func (t *LeaveChannelTool) Schema() mcp.ToolInputSchema { return leaveChannelSchema() }
func (t *LeaveChannelTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	ok, err := t.client.LeaveChannel(ctx, channelID)
	if err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true, "left": %v}`, ok)), nil
}
func (t *LeaveChannelTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskMed),
		framework.WithImpact(framework.ImpactWrite),
		framework.WithResourceCost(3),
		framework.WithPII(false),
		framework.WithIdempotent(false),
	)
}

type SetChannelTopicTool struct{ client *Client }

func (t *SetChannelTopicTool) Name() string                { return "set_channel_topic" }
func (t *SetChannelTopicTool) Description() string         { return "Set the topic for a Slack channel." }
func (t *SetChannelTopicTool) Schema() mcp.ToolInputSchema { return setChannelTopicSchema() }
func (t *SetChannelTopicTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	topic := getRequiredString(args, "topic")
	ch, err := t.client.SetChannelTopic(ctx, channelID, topic)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(ch)
	return framework.TextResult(string(b)), nil
}
func (t *SetChannelTopicTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskMed),
		framework.WithImpact(framework.ImpactWrite),
		framework.WithResourceCost(3),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

// ============================================================================
// MESSAGE TOOLS
// ============================================================================

type PostMessageTool struct{ client *Client }

func (t *PostMessageTool) Name() string                { return "post_message" }
func (t *PostMessageTool) Description() string         { return "Post a message to a Slack channel." }
func (t *PostMessageTool) Schema() mcp.ToolInputSchema { return postMessageSchema() }
func (t *PostMessageTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	text := getRequiredString(args, "text")
	threadTS := getString(args, "thread_ts", "")
	unfurlLinks := getBool(args, "unfurl_links", true)
	unfurlMedia := getBool(args, "unfurl_media", true)
	chanID, ts, err := t.client.PostMessage(ctx, channelID, text, threadTS, unfurlLinks, unfurlMedia, false)
	if err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true, "channel": %q, "ts": %q}`, chanID, ts)), nil
}
func (t *PostMessageTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskMed),
		framework.WithImpact(framework.ImpactWrite),
		framework.WithResourceCost(3),
		framework.WithPII(false),
		framework.WithIdempotent(false),
	)
}

type ReplyToThreadTool struct{ client *Client }

func (t *ReplyToThreadTool) Name() string                { return "reply_to_thread" }
func (t *ReplyToThreadTool) Description() string         { return "Reply to a message thread in Slack." }
func (t *ReplyToThreadTool) Schema() mcp.ToolInputSchema { return replyToThreadSchema() }
func (t *ReplyToThreadTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	threadTS := getRequiredString(args, "thread_ts")
	text := getRequiredString(args, "text")
	broadcast := getBool(args, "broadcast", false)
	chanID, ts, err := t.client.PostMessage(ctx, channelID, text, threadTS, false, false, broadcast)
	if err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true, "channel": %q, "ts": %q}`, chanID, ts)), nil
}
func (t *ReplyToThreadTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskMed),
		framework.WithImpact(framework.ImpactWrite),
		framework.WithResourceCost(3),
		framework.WithPII(false),
		framework.WithIdempotent(false),
	)
}

type GetThreadRepliesTool struct{ client *Client }

func (t *GetThreadRepliesTool) Name() string                { return "get_thread_replies" }
func (t *GetThreadRepliesTool) Description() string         { return "Get all replies in a message thread." }
func (t *GetThreadRepliesTool) Schema() mcp.ToolInputSchema { return getThreadRepliesSchema() }
func (t *GetThreadRepliesTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	threadTS := getRequiredString(args, "thread_ts")
	limit := getInt(args, "limit", 20)
	msgs, _, _, err := t.client.GetThreadReplies(ctx, channelID, threadTS, limit)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(msgs)
	return framework.TextResult(string(b)), nil
}
func (t *GetThreadRepliesTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(3),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type AddReactionTool struct{ client *Client }

func (t *AddReactionTool) Name() string                { return "add_reaction" }
func (t *AddReactionTool) Description() string         { return "Add an emoji reaction to a message." }
func (t *AddReactionTool) Schema() mcp.ToolInputSchema { return addReactionSchema() }
func (t *AddReactionTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	timestamp := getRequiredString(args, "timestamp")
	emoji := getRequiredString(args, "emoji")
	if err := t.client.AddReaction(ctx, channelID, timestamp, emoji); err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true}`)), nil
}
func (t *AddReactionTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

type RemoveReactionTool struct{ client *Client }

func (t *RemoveReactionTool) Name() string                { return "remove_reaction" }
func (t *RemoveReactionTool) Description() string         { return "Remove an emoji reaction from a message." }
func (t *RemoveReactionTool) Schema() mcp.ToolInputSchema { return removeReactionSchema() }
func (t *RemoveReactionTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	timestamp := getRequiredString(args, "timestamp")
	emoji := getRequiredString(args, "emoji")
	if err := t.client.RemoveReaction(ctx, channelID, timestamp, emoji); err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true}`)), nil
}
func (t *RemoveReactionTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(false),
		framework.WithIdempotent(false),
	)
}

type GetMessageReactionsTool struct{ client *Client }

func (t *GetMessageReactionsTool) Name() string { return "get_message_reactions" }
func (t *GetMessageReactionsTool) Description() string {
	return "Get all reactions on a specific message."
}
func (t *GetMessageReactionsTool) Schema() mcp.ToolInputSchema { return getMessageReactionsSchema() }
func (t *GetMessageReactionsTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	timestamp := getRequiredString(args, "timestamp")
	reactions, err := t.client.GetReactions(ctx, channelID, timestamp, true)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(reactions)
	return framework.TextResult(string(b)), nil
}
func (t *GetMessageReactionsTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

type UpdateMessageTool struct{ client *Client }

func (t *UpdateMessageTool) Name() string                { return "update_message" }
func (t *UpdateMessageTool) Description() string         { return "Update an existing message." }
func (t *UpdateMessageTool) Schema() mcp.ToolInputSchema { return updateMessageSchema() }
func (t *UpdateMessageTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	timestamp := getRequiredString(args, "timestamp")
	text := getRequiredString(args, "text")
	_, ts, _, err := t.client.UpdateMessage(ctx, channelID, timestamp, text)
	if err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true, "ts": %q}`, ts)), nil
}
func (t *UpdateMessageTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskMed),
		framework.WithImpact(framework.ImpactWrite),
		framework.WithResourceCost(3),
		framework.WithPII(false),
		framework.WithIdempotent(false),
	)
}

type DeleteMessageTool struct{ client *Client }

func (t *DeleteMessageTool) Name() string                { return "delete_message" }
func (t *DeleteMessageTool) Description() string         { return "Delete a message from a channel." }
func (t *DeleteMessageTool) Schema() mcp.ToolInputSchema { return deleteMessageSchema() }
func (t *DeleteMessageTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	timestamp := getRequiredString(args, "timestamp")
	_, _, err := t.client.DeleteMessage(ctx, channelID, timestamp)
	if err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true}`)), nil
}
func (t *DeleteMessageTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskHigh),
		framework.WithImpact(framework.ImpactDelete),
		framework.WithResourceCost(3),
		framework.WithPII(true),
		framework.WithIdempotent(false),
		framework.WithApprovalReq(true),
	)
}

type SendDMTool struct{ client *Client }

func (t *SendDMTool) Name() string                { return "send_dm" }
func (t *SendDMTool) Description() string         { return "Send a direct message to a user." }
func (t *SendDMTool) Schema() mcp.ToolInputSchema { return sendDMSchema() }
func (t *SendDMTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	userID := getRequiredString(args, "user_id")
	text := getRequiredString(args, "text")
	ch, _, _, err := t.client.OpenDM(ctx, []string{userID})
	if err != nil {
		return framework.TextResult(""), err
	}
	chanID, ts, err := t.client.PostMessage(ctx, ch.ID, text, "", false, false, false)
	if err != nil {
		return framework.TextResult(""), err
	}
	return framework.TextResult(fmt.Sprintf(`{"ok": true, "channel": %q, "ts": %q}`, chanID, ts)), nil
}
func (t *SendDMTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskMed),
		framework.WithImpact(framework.ImpactWrite),
		framework.WithResourceCost(3),
		framework.WithPII(true),
		framework.WithIdempotent(false),
	)
}

type ListConversationsTool struct{ client *Client }

func (t *ListConversationsTool) Name() string { return "list_conversations" }
func (t *ListConversationsTool) Description() string {
	return "List direct message and group DM conversations."
}
func (t *ListConversationsTool) Schema() mcp.ToolInputSchema { return listConversationsSchema() }
func (t *ListConversationsTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	types := getString(args, "types", "im,mpim")
	limit := getInt(args, "limit", 50)
	chans, _, err := t.client.ListConversations(ctx, types, limit)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(chans)
	return framework.TextResult(string(b)), nil
}
func (t *ListConversationsTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

type GetDMHistoryTool struct{ client *Client }

func (t *GetDMHistoryTool) Name() string { return "get_dm_history" }
func (t *GetDMHistoryTool) Description() string {
	return "Fetch message history from a DM or group DM conversation."
}
func (t *GetDMHistoryTool) Schema() mcp.ToolInputSchema { return getDMHistorySchema() }
func (t *GetDMHistoryTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	channelID := getRequiredString(args, "channel_id")
	limit := getInt(args, "limit", 20)
	oldest := getString(args, "oldest", "")
	latest := getString(args, "latest", "")
	history, err := t.client.GetDMHistory(ctx, channelID, limit, oldest, latest)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(history)
	return framework.TextResult(string(b)), nil
}
func (t *GetDMHistoryTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(3),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type OpenDMTool struct{ client *Client }

func (t *OpenDMTool) Name() string                { return "open_dm" }
func (t *OpenDMTool) Description() string         { return "Open a direct message or group DM conversation." }
func (t *OpenDMTool) Schema() mcp.ToolInputSchema { return openDMSchema() }
func (t *OpenDMTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	userIDsStr := getRequiredString(args, "user_ids")
	ch, _, _, err := t.client.OpenDM(ctx, []string{userIDsStr})
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(ch)
	return framework.TextResult(string(b)), nil
}
func (t *OpenDMTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

// ============================================================================
// SEARCH TOOLS
// ============================================================================

type SearchMessagesTool struct{ client *Client }

func (t *SearchMessagesTool) Name() string { return "search_messages" }
func (t *SearchMessagesTool) Description() string {
	return "Search messages across the Slack workspace."
}
func (t *SearchMessagesTool) Schema() mcp.ToolInputSchema { return searchMessagesSchema() }
func (t *SearchMessagesTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	query := getRequiredString(args, "query")
	sort := getString(args, "sort", "timestamp")
	sortDir := getString(args, "sort_dir", "desc")
	count := getInt(args, "count", 20)
	results, err := t.client.SearchMessages(ctx, query, sort, sortDir, count)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(results)
	return framework.TextResult(string(b)), nil
}
func (t *SearchMessagesTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(4),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type SearchFilesTool struct{ client *Client }

func (t *SearchFilesTool) Name() string                { return "search_files" }
func (t *SearchFilesTool) Description() string         { return "Search files across the Slack workspace." }
func (t *SearchFilesTool) Schema() mcp.ToolInputSchema { return searchFilesSchema() }
func (t *SearchFilesTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	query := getRequiredString(args, "query")
	sort := getString(args, "sort", "timestamp")
	sortDir := getString(args, "sort_dir", "desc")
	count := getInt(args, "count", 20)
	results, err := t.client.SearchFiles(ctx, query, sort, sortDir, count)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(results)
	return framework.TextResult(string(b)), nil
}
func (t *SearchFilesTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(4),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type SearchAllTool struct{ client *Client }

func (t *SearchAllTool) Name() string { return "search_all" }
func (t *SearchAllTool) Description() string {
	return "Search both messages and files across the Slack workspace."
}
func (t *SearchAllTool) Schema() mcp.ToolInputSchema { return searchAllSchema() }
func (t *SearchAllTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	query := getRequiredString(args, "query")
	sort := getString(args, "sort", "timestamp")
	sortDir := getString(args, "sort_dir", "desc")
	count := getInt(args, "count", 10)
	msgs, files, err := t.client.SearchAll(ctx, query, sort, sortDir, count)
	if err != nil {
		return framework.TextResult(""), err
	}
	result := map[string]interface{}{"messages": msgs, "files": files}
	b, _ := json.Marshal(result)
	return framework.TextResult(string(b)), nil
}
func (t *SearchAllTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(4),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

// ============================================================================
// USER TOOLS
// ============================================================================

type ListUsersTool struct{ client *Client }

func (t *ListUsersTool) Name() string                { return "list_users" }
func (t *ListUsersTool) Description() string         { return "List all users in the Slack workspace." }
func (t *ListUsersTool) Schema() mcp.ToolInputSchema { return listUsersSchema() }
func (t *ListUsersTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	limit := getInt(args, "limit", 100)
	includeLocale := getBool(args, "include_locale", false)
	users, err := t.client.ListUsers(ctx, limit, includeLocale)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(users)
	return framework.TextResult(string(b)), nil
}
func (t *ListUsersTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(3),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type GetUserInfoTool struct{ client *Client }

func (t *GetUserInfoTool) Name() string { return "get_user_info" }
func (t *GetUserInfoTool) Description() string {
	return "Get detailed information about a specific user."
}
func (t *GetUserInfoTool) Schema() mcp.ToolInputSchema { return getUserInfoSchema() }
func (t *GetUserInfoTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	userID := getRequiredString(args, "user_id")
	user, err := t.client.GetUserInfo(ctx, userID)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(user)
	return framework.TextResult(string(b)), nil
}
func (t *GetUserInfoTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type GetUserPresenceTool struct{ client *Client }

func (t *GetUserPresenceTool) Name() string { return "get_user_presence" }
func (t *GetUserPresenceTool) Description() string {
	return "Get the presence/online status of a user."
}
func (t *GetUserPresenceTool) Schema() mcp.ToolInputSchema { return getUserPresenceSchema() }
func (t *GetUserPresenceTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	userID := getRequiredString(args, "user_id")
	presence, err := t.client.GetUserPresence(ctx, userID)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(presence)
	return framework.TextResult(string(b)), nil
}
func (t *GetUserPresenceTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

type LookupUserByEmailTool struct{ client *Client }

func (t *LookupUserByEmailTool) Name() string                { return "lookup_user_by_email" }
func (t *LookupUserByEmailTool) Description() string         { return "Find a user by their email address." }
func (t *LookupUserByEmailTool) Schema() mcp.ToolInputSchema { return lookupUserByEmailSchema() }
func (t *LookupUserByEmailTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	email := getRequiredString(args, "email")
	user, err := t.client.LookupUserByEmail(ctx, email)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(user)
	return framework.TextResult(string(b)), nil
}
func (t *LookupUserByEmailTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type GetUserProfileTool struct{ client *Client }

func (t *GetUserProfileTool) Name() string                { return "get_user_profile" }
func (t *GetUserProfileTool) Description() string         { return "Get the profile information for a user." }
func (t *GetUserProfileTool) Schema() mcp.ToolInputSchema { return getUserProfileSchema() }
func (t *GetUserProfileTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	userID := getRequiredString(args, "user_id")
	profile, err := t.client.GetUserProfile(ctx, userID)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(profile)
	return framework.TextResult(string(b)), nil
}
func (t *GetUserProfileTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(2),
		framework.WithPII(true),
		framework.WithIdempotent(true),
	)
}

type GetBotInfoTool struct{ client *Client }

func (t *GetBotInfoTool) Name() string { return "get_bot_info" }
func (t *GetBotInfoTool) Description() string {
	return "Get information about the authenticated bot user."
}
func (t *GetBotInfoTool) Schema() mcp.ToolInputSchema { return getBotInfoSchema() }
func (t *GetBotInfoTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	info, err := t.client.GetBotInfo(ctx)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(info)
	return framework.TextResult(string(b)), nil
}
func (t *GetBotInfoTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(1),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

type GetTeamInfoTool struct{ client *Client }

func (t *GetTeamInfoTool) Name() string { return "get_team_info" }
func (t *GetTeamInfoTool) Description() string {
	return "Get information about the Slack workspace/team."
}
func (t *GetTeamInfoTool) Schema() mcp.ToolInputSchema { return getTeamInfoSchema() }
func (t *GetTeamInfoTool) Handle(ctx context.Context, args map[string]interface{}) (framework.ToolResult, error) {
	info, err := t.client.GetTeamInfo(ctx)
	if err != nil {
		return framework.TextResult(""), err
	}
	b, _ := json.Marshal(info)
	return framework.TextResult(string(b)), nil
}
func (t *GetTeamInfoTool) GetEnforcerProfile() *framework.EnforcerProfile {
	return framework.NewEnforcerProfile(
		framework.WithRisk(framework.RiskLow),
		framework.WithImpact(framework.ImpactRead),
		framework.WithResourceCost(1),
		framework.WithPII(false),
		framework.WithIdempotent(true),
	)
}

// ============================================================================
// SCHEMA HELPERS
// ============================================================================
// SCHEMA HELPERS
// ============================================================================

func listChannelsSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"types":            map[string]interface{}{"type": "string", "description": "Comma-separated channel types: public_channel,private_channel (default: public_channel,private_channel)"},
			"exclude_archived": map[string]interface{}{"type": "boolean", "description": "Exclude archived channels (default: true)"},
			"limit":            map[string]interface{}{"type": "number", "description": "Maximum number of channels to return (default: 100)"},
		},
	}
}

func getChannelInfoSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
		},
		Required: []string{"channel_id"},
	}
}

func createChannelSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"name":       map[string]interface{}{"type": "string", "description": "Name of the channel to create"},
			"is_private": map[string]interface{}{"type": "boolean", "description": "Create as private channel (default: false)"},
		},
		Required: []string{"name"},
	}
}

func archiveChannelSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID to archive"},
		},
		Required: []string{"channel_id"},
	}
}

func getChannelHistorySchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
			"limit":      map[string]interface{}{"type": "number", "description": "Number of messages to return (default: 20)"},
			"oldest":     map[string]interface{}{"type": "string", "description": "Start of time range (Unix timestamp)"},
			"latest":     map[string]interface{}{"type": "string", "description": "End of time range (Unix timestamp)"},
		},
		Required: []string{"channel_id"},
	}
}

func joinChannelSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID to join"},
		},
		Required: []string{"channel_id"},
	}
}

func leaveChannelSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID to leave"},
		},
		Required: []string{"channel_id"},
	}
}

func setChannelTopicSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
			"topic":      map[string]interface{}{"type": "string", "description": "The new topic"},
		},
		Required: []string{"channel_id", "topic"},
	}
}

func postMessageSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id":   map[string]interface{}{"type": "string", "description": "The channel ID to post to"},
			"text":         map[string]interface{}{"type": "string", "description": "Message text"},
			"thread_ts":    map[string]interface{}{"type": "string", "description": "Thread timestamp to reply to"},
			"unfurl_links": map[string]interface{}{"type": "boolean", "description": "Unfurl links (default: true)"},
			"unfurl_media": map[string]interface{}{"type": "boolean", "description": "Unfurl media (default: true)"},
		},
		Required: []string{"channel_id", "text"},
	}
}

func replyToThreadSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
			"thread_ts":  map[string]interface{}{"type": "string", "description": "Thread timestamp to reply to"},
			"text":       map[string]interface{}{"type": "string", "description": "Reply text"},
			"broadcast":  map[string]interface{}{"type": "boolean", "description": "Broadcast to channel (default: false)"},
		},
		Required: []string{"channel_id", "thread_ts", "text"},
	}
}

func getThreadRepliesSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
			"thread_ts":  map[string]interface{}{"type": "string", "description": "Thread timestamp"},
			"limit":      map[string]interface{}{"type": "number", "description": "Number of replies to return (default: 20)"},
		},
		Required: []string{"channel_id", "thread_ts"},
	}
}

func addReactionSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
			"timestamp":  map[string]interface{}{"type": "string", "description": "Message timestamp"},
			"emoji":      map[string]interface{}{"type": "string", "description": "Emoji name (e.g. thumbsup)"},
		},
		Required: []string{"channel_id", "timestamp", "emoji"},
	}
}

func removeReactionSchema() mcp.ToolInputSchema {
	return addReactionSchema()
}

func getMessageReactionsSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
			"timestamp":  map[string]interface{}{"type": "string", "description": "Message timestamp"},
		},
		Required: []string{"channel_id", "timestamp"},
	}
}

func updateMessageSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
			"timestamp":  map[string]interface{}{"type": "string", "description": "Message timestamp to update"},
			"text":       map[string]interface{}{"type": "string", "description": "New message text"},
		},
		Required: []string{"channel_id", "timestamp", "text"},
	}
}

func deleteMessageSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The channel ID"},
			"timestamp":  map[string]interface{}{"type": "string", "description": "Message timestamp to delete"},
		},
		Required: []string{"channel_id", "timestamp"},
	}
}

func sendDMSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"user_id": map[string]interface{}{"type": "string", "description": "The user ID to send a DM to"},
			"text":    map[string]interface{}{"type": "string", "description": "DM text"},
		},
		Required: []string{"user_id", "text"},
	}
}

func listConversationsSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"types": map[string]interface{}{"type": "string", "description": "Conversation types: im,mpim (default: im,mpim)"},
			"limit": map[string]interface{}{"type": "number", "description": "Maximum results (default: 50)"},
		},
	}
}

func getDMHistorySchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"channel_id": map[string]interface{}{"type": "string", "description": "The DM channel ID"},
			"limit":      map[string]interface{}{"type": "number", "description": "Number of messages (default: 20)"},
			"oldest":     map[string]interface{}{"type": "string", "description": "Start of time range"},
			"latest":     map[string]interface{}{"type": "string", "description": "End of time range"},
		},
		Required: []string{"channel_id"},
	}
}

func openDMSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"user_ids": map[string]interface{}{"type": "string", "description": "Comma-separated user IDs"},
		},
		Required: []string{"user_ids"},
	}
}

func searchMessagesSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"query":    map[string]interface{}{"type": "string", "description": "Search query"},
			"sort":     map[string]interface{}{"type": "string", "description": "Sort by: score or timestamp (default: timestamp)"},
			"sort_dir": map[string]interface{}{"type": "string", "description": "Sort direction: asc or desc (default: desc)"},
			"count":    map[string]interface{}{"type": "number", "description": "Number of results (default: 20)"},
		},
		Required: []string{"query"},
	}
}

func searchFilesSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"query":    map[string]interface{}{"type": "string", "description": "Search query"},
			"sort":     map[string]interface{}{"type": "string", "description": "Sort by: score or timestamp (default: timestamp)"},
			"sort_dir": map[string]interface{}{"type": "string", "description": "Sort direction: asc or desc (default: desc)"},
			"count":    map[string]interface{}{"type": "number", "description": "Number of results (default: 20)"},
		},
		Required: []string{"query"},
	}
}

func searchAllSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"query":    map[string]interface{}{"type": "string", "description": "Search query"},
			"sort":     map[string]interface{}{"type": "string", "description": "Sort by: score or timestamp (default: timestamp)"},
			"sort_dir": map[string]interface{}{"type": "string", "description": "Sort direction: asc or desc (default: desc)"},
			"count":    map[string]interface{}{"type": "number", "description": "Number of results (default: 10)"},
		},
		Required: []string{"query"},
	}
}

func listUsersSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"limit":          map[string]interface{}{"type": "number", "description": "Maximum users to return (default: 100)"},
			"include_locale": map[string]interface{}{"type": "boolean", "description": "Include user locale (default: false)"},
		},
	}
}

func getUserInfoSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"user_id": map[string]interface{}{"type": "string", "description": "The user ID"},
		},
		Required: []string{"user_id"},
	}
}

func getUserPresenceSchema() mcp.ToolInputSchema {
	return getUserInfoSchema()
}

func lookupUserByEmailSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"email": map[string]interface{}{"type": "string", "description": "User email address"},
		},
		Required: []string{"email"},
	}
}

func getUserProfileSchema() mcp.ToolInputSchema {
	return getUserInfoSchema()
}

func getBotInfoSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{Type: "object"}
}

func getTeamInfoSchema() mcp.ToolInputSchema {
	return mcp.ToolInputSchema{Type: "object"}
}

// ============================================================================
// ARGUMENT PARSING HELPERS
// ============================================================================

func getString(args map[string]interface{}, key, defaultVal string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return defaultVal
}

func getBool(args map[string]interface{}, key string, defaultVal bool) bool {
	if v, ok := args[key].(bool); ok {
		return v
	}
	return defaultVal
}

func getInt(args map[string]interface{}, key string, defaultVal int) int {
	if v, ok := args[key].(float64); ok {
		return int(v)
	}
	return defaultVal
}

func getRequiredString(args map[string]interface{}, key string) string {
	if v, ok := args[key].(string); ok {
		return v
	}
	return ""
}
