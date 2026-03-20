package slack

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/slack-go/slack"
)

type Client struct {
	botClient    *slack.Client
	userClient   *slack.Client
	hasBotToken  bool
	hasUserToken bool
}

func NewClient() *Client {
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	userToken := os.Getenv("SLACK_USER_TOKEN")

	var botClient *slack.Client
	if botToken != "" {
		botClient = slack.New(botToken)
	}

	var userClient *slack.Client
	if userToken != "" {
		userClient = slack.New(userToken)
	}

	return &Client{
		botClient:    botClient,
		userClient:   userClient,
		hasBotToken:  botToken != "",
		hasUserToken: userToken != "",
	}
}

func (c *Client) IsConnected() bool {
	return c.botClient != nil
}

func (c *Client) HasUserToken() bool {
	return c.hasUserToken
}

func (c *Client) requireBot() error {
	if c.botClient == nil {
		return fmt.Errorf("SLACK_BOT_TOKEN is not set")
	}
	return nil
}

func (c *Client) requireUser() error {
	if !c.hasUserToken {
		return fmt.Errorf("SLACK_USER_TOKEN is not set (required for search)")
	}
	return nil
}

// ---- Channels ----

func (c *Client) ListChannels(ctx context.Context, types string, excludeArchived bool, limit int) ([]slack.Channel, string, error) {
	if err := c.requireBot(); err != nil {
		return nil, "", err
	}
	params := slack.GetConversationsParameters{
		Types:           []string{types},
		ExcludeArchived: excludeArchived,
		Limit:           limit,
	}
	return c.botClient.GetConversationsContext(ctx, &params)
}

func (c *Client) GetChannelInfo(ctx context.Context, channelID string) (*slack.Channel, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	params := slack.GetConversationInfoInput{ChannelID: channelID}
	return c.botClient.GetConversationInfoContext(ctx, &params)
}

func (c *Client) CreateChannel(ctx context.Context, name string, isPrivate bool) (*slack.Channel, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	params := slack.CreateConversationParams{ChannelName: name, IsPrivate: isPrivate}
	return c.botClient.CreateConversationContext(ctx, params)
}

func (c *Client) ArchiveChannel(ctx context.Context, channelID string) error {
	if err := c.requireBot(); err != nil {
		return err
	}
	return c.botClient.ArchiveConversationContext(ctx, channelID)
}

func (c *Client) GetChannelHistory(ctx context.Context, channelID string, limit int, oldest, latest string) (*slack.GetConversationHistoryResponse, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	params := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
		Latest:    latest,
		Oldest:    oldest,
	}
	return c.botClient.GetConversationHistoryContext(ctx, &params)
}

func (c *Client) JoinChannel(ctx context.Context, channelID string) (*slack.Channel, string, error) {
	if err := c.requireBot(); err != nil {
		return nil, "", err
	}
	ch, ts, _, err := c.botClient.JoinConversationContext(ctx, channelID)
	return ch, ts, err
}

func (c *Client) LeaveChannel(ctx context.Context, channelID string) (bool, error) {
	if err := c.requireBot(); err != nil {
		return false, err
	}
	return c.botClient.LeaveConversationContext(ctx, channelID)
}

func (c *Client) SetChannelTopic(ctx context.Context, channelID, topic string) (*slack.Channel, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	return c.botClient.SetTopicOfConversationContext(ctx, channelID, topic)
}

// ---- Messages ----

func (c *Client) PostMessage(ctx context.Context, channelID, text, threadTS string, unfurlLinks, unfurlMedia, broadcast bool) (string, string, error) {
	if err := c.requireBot(); err != nil {
		return "", "", err
	}
	params := slack.PostMessageParameters{
		UnfurlLinks:     unfurlLinks,
		UnfurlMedia:     unfurlMedia,
		ThreadTimestamp: threadTS,
		ReplyBroadcast:  broadcast,
	}
	return c.botClient.PostMessageContext(ctx, channelID,
		slack.MsgOptionPostMessageParameters(params),
		slack.MsgOptionText(text, false),
	)
}

func (c *Client) UpdateMessage(ctx context.Context, channelID, timestamp, text string) (string, string, string, error) {
	if err := c.requireBot(); err != nil {
		return "", "", "", err
	}
	return c.botClient.UpdateMessageContext(ctx, channelID, timestamp,
		slack.MsgOptionText(text, false),
	)
}

func (c *Client) DeleteMessage(ctx context.Context, channelID, timestamp string) (string, string, error) {
	if err := c.requireBot(); err != nil {
		return "", "", err
	}
	return c.botClient.DeleteMessageContext(ctx, channelID, timestamp)
}

// ---- Reactions ----

func (c *Client) AddReaction(ctx context.Context, channelID, timestamp, emoji string) error {
	if err := c.requireBot(); err != nil {
		return err
	}
	item := slack.ItemRef{Channel: channelID, Timestamp: timestamp}
	return c.botClient.AddReactionContext(ctx, emoji, item)
}

func (c *Client) RemoveReaction(ctx context.Context, channelID, timestamp, emoji string) error {
	if err := c.requireBot(); err != nil {
		return err
	}
	item := slack.ItemRef{Channel: channelID, Timestamp: timestamp}
	return c.botClient.RemoveReactionContext(ctx, emoji, item)
}

func (c *Client) GetReactions(ctx context.Context, channelID, timestamp string, full bool) (slack.ReactedItem, error) {
	if err := c.requireBot(); err != nil {
		return slack.ReactedItem{}, err
	}
	item := slack.ItemRef{Channel: channelID, Timestamp: timestamp}
	params := slack.GetReactionsParameters{Full: full}
	return c.botClient.GetReactionsContext(ctx, item, params)
}

// ---- Conversations ----

func (c *Client) GetThreadReplies(ctx context.Context, channelID, threadTS string, limit int) ([]slack.Message, bool, string, error) {
	if err := c.requireBot(); err != nil {
		return nil, false, "", err
	}
	params := slack.GetConversationRepliesParameters{
		ChannelID: channelID,
		Timestamp: threadTS,
		Limit:     limit,
	}
	return c.botClient.GetConversationRepliesContext(ctx, &params)
}

func (c *Client) OpenDM(ctx context.Context, userIDs []string) (*slack.Channel, bool, bool, error) {
	if err := c.requireBot(); err != nil {
		return nil, false, false, err
	}
	params := slack.OpenConversationParameters{Users: userIDs}
	return c.botClient.OpenConversationContext(ctx, &params)
}

func (c *Client) ListConversations(ctx context.Context, types string, limit int) ([]slack.Channel, string, error) {
	if err := c.requireBot(); err != nil {
		return nil, "", err
	}
	params := slack.GetConversationsParameters{Types: []string{types}, Limit: limit}
	return c.botClient.GetConversationsContext(ctx, &params)
}

func (c *Client) GetDMHistory(ctx context.Context, channelID string, limit int, oldest, latest string) (*slack.GetConversationHistoryResponse, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	params := slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
		Latest:    latest,
		Oldest:    oldest,
	}
	return c.botClient.GetConversationHistoryContext(ctx, &params)
}

// ---- Users ----

func (c *Client) ListUsers(ctx context.Context, limit int, includeLocale bool) ([]slack.User, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	opts := []slack.GetUsersOption{slack.GetUsersOptionLimit(limit)}
	if includeLocale {
		opts = append(opts, slack.GetUsersOptionPresence(true))
	}
	users, err := c.botClient.GetUsersContext(ctx, opts...)
	return users, err
}

func (c *Client) GetUserInfo(ctx context.Context, userID string) (*slack.User, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	return c.botClient.GetUserInfoContext(ctx, userID)
}

func (c *Client) GetUserPresence(ctx context.Context, userID string) (*slack.UserPresence, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	return c.botClient.GetUserPresenceContext(ctx, userID)
}

func (c *Client) LookupUserByEmail(ctx context.Context, email string) (*slack.User, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	return c.botClient.GetUserByEmailContext(ctx, email)
}

func (c *Client) GetUserProfile(ctx context.Context, userID string) (*slack.UserProfile, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	params := slack.GetUserProfileParameters{UserID: userID}
	return c.botClient.GetUserProfileContext(ctx, &params)
}

func (c *Client) GetBotInfo(ctx context.Context) (*slack.AuthTestResponse, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	return c.botClient.AuthTestContext(ctx)
}

func (c *Client) GetTeamInfo(ctx context.Context) (*slack.TeamInfo, error) {
	if err := c.requireBot(); err != nil {
		return nil, err
	}
	return c.botClient.GetTeamInfoContext(ctx)
}

// ---- Search (requires user token) ----

func (c *Client) SearchMessages(ctx context.Context, query, sort, sortDir string, count int) (*slack.SearchMessages, error) {
	if err := c.requireUser(); err != nil {
		return nil, err
	}
	params := slack.SearchParameters{
		Sort:          sort,
		SortDirection: sortDir,
		Count:         count,
	}
	return c.userClient.SearchMessagesContext(ctx, query, params)
}

func (c *Client) SearchFiles(ctx context.Context, query, sort, sortDir string, count int) (*slack.SearchFiles, error) {
	if err := c.requireUser(); err != nil {
		return nil, err
	}
	params := slack.SearchParameters{
		Sort:          sort,
		SortDirection: sortDir,
		Count:         count,
	}
	return c.userClient.SearchFilesContext(ctx, query, params)
}

func (c *Client) SearchAll(ctx context.Context, query, sort, sortDir string, count int) (*slack.SearchMessages, *slack.SearchFiles, error) {
	if err := c.requireUser(); err != nil {
		return nil, nil, err
	}
	params := slack.SearchParameters{
		Sort:          sort,
		SortDirection: sortDir,
		Count:         count,
	}
	return c.userClient.SearchContext(ctx, query, params)
}

// ---- Helpers ----

func ParseTimestamp(ts string) (time.Time, error) {
	sec, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timestamp: %s", ts)
	}
	return time.Unix(sec, 0), nil
}
