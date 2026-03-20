package slack

import (
	"context"
	"testing"

	"github.com/karldane/mcp-framework/framework"
)

// ============================================================================
// Test Suite 1: Server Creation and Initialization
// ============================================================================

func TestServerCreation(t *testing.T) {
	server := NewServer()
	if server == nil {
		t.Fatal("Expected server to be created")
	}
}

func TestServerInitializeWithoutTokens(t *testing.T) {
	server := NewServer()

	// Should not panic or error when tokens are missing
	server.Initialize()

	// Verify client is set up but may not be connected
	if server.client == nil {
		t.Error("Client should be initialized even without tokens")
	}

	// Check that tools are registered even without connection
	tools := server.ListTools()
	if len(tools) == 0 {
		t.Error("Tools should be registered even without tokens")
	}
}

// ============================================================================
// Test Suite 2: Tool Registration
// ============================================================================

func TestToolRegistrationCount(t *testing.T) {
	server := NewServer()
	server.Initialize()

	tools := server.ListTools()

	// Should have all 30 tools (search tools are registered but will fail at runtime if no user token)
	if len(tools) != 30 {
		t.Errorf("Expected 30 tools, got %d", len(tools))
	}
}

func TestRequiredToolsExist(t *testing.T) {
	server := NewServer()
	server.Initialize()

	requiredTools := []string{
		"list_channels",
		"get_channel_info",
		"create_channel",
		"archive_channel",
		"get_channel_history",
		"join_channel",
		"leave_channel",
		"set_channel_topic",
		"post_message",
		"reply_to_thread",
		"get_thread_replies",
		"add_reaction",
		"remove_reaction",
		"get_message_reactions",
		"update_message",
		"delete_message",
		"send_dm",
		"list_conversations",
		"get_dm_history",
		"open_dm",
		"search_messages",
		"search_files",
		"search_all",
		"list_users",
		"get_user_info",
		"get_user_presence",
		"lookup_user_by_email",
		"get_user_profile",
		"get_bot_info",
		"get_team_info",
	}

	registeredTools := make(map[string]bool)
	for _, tool := range server.ListTools() {
		registeredTools[tool] = true
	}

	for _, required := range requiredTools {
		if !registeredTools[required] {
			t.Errorf("Required tool '%s' not registered", required)
		}
	}
}

// ============================================================================
// Test Suite 3: EnforcerProfile Validation
// ============================================================================

func TestEnforcerProfile_ReadTools(t *testing.T) {
	tests := []struct {
		name string
		tool interface {
			GetEnforcerProfile() framework.EnforcerProfile
		}
		expectedRisk     framework.RiskLevel
		expectedImpact   framework.ImpactScope
		expectedApproval bool
	}{
		{
			name:             "ListChannelsTool",
			tool:             &ListChannelsTool{},
			expectedRisk:     framework.RiskLow,
			expectedImpact:   framework.ImpactRead,
			expectedApproval: false,
		},
		{
			name:             "GetChannelInfoTool",
			tool:             &GetChannelInfoTool{},
			expectedRisk:     framework.RiskLow,
			expectedImpact:   framework.ImpactRead,
			expectedApproval: false,
		},
		{
			name:             "GetChannelHistoryTool",
			tool:             &GetChannelHistoryTool{},
			expectedRisk:     framework.RiskLow,
			expectedImpact:   framework.ImpactRead,
			expectedApproval: false,
		},
		{
			name:             "ListUsersTool",
			tool:             &ListUsersTool{},
			expectedRisk:     framework.RiskLow,
			expectedImpact:   framework.ImpactRead,
			expectedApproval: false,
		},
		{
			name:             "GetUserInfoTool",
			tool:             &GetUserInfoTool{},
			expectedRisk:     framework.RiskLow,
			expectedImpact:   framework.ImpactRead,
			expectedApproval: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := tt.tool.GetEnforcerProfile()

			if profile.RiskLevel != tt.expectedRisk {
				t.Errorf("Expected risk %s, got %s", tt.expectedRisk, profile.RiskLevel)
			}

			if profile.ImpactScope != tt.expectedImpact {
				t.Errorf("Expected impact %s, got %s", tt.expectedImpact, profile.ImpactScope)
			}

			if profile.ApprovalReq != tt.expectedApproval {
				t.Errorf("Expected approval %v, got %v", tt.expectedApproval, profile.ApprovalReq)
			}
		})
	}
}

func TestEnforcerProfile_WriteTools(t *testing.T) {
	tests := []struct {
		name string
		tool interface {
			GetEnforcerProfile() framework.EnforcerProfile
		}
		expectedRisk     framework.RiskLevel
		expectedImpact   framework.ImpactScope
		expectedApproval bool
	}{
		{
			name:             "CreateChannelTool",
			tool:             &CreateChannelTool{},
			expectedRisk:     framework.RiskHigh,
			expectedImpact:   framework.ImpactWrite,
			expectedApproval: true,
		},
		{
			name:             "PostMessageTool",
			tool:             &PostMessageTool{},
			expectedRisk:     framework.RiskMed,
			expectedImpact:   framework.ImpactWrite,
			expectedApproval: false,
		},
		{
			name:             "DeleteMessageTool",
			tool:             &DeleteMessageTool{},
			expectedRisk:     framework.RiskHigh,
			expectedImpact:   framework.ImpactDelete,
			expectedApproval: true,
		},
		{
			name:             "ArchiveChannelTool",
			tool:             &ArchiveChannelTool{},
			expectedRisk:     framework.RiskHigh,
			expectedImpact:   framework.ImpactDelete,
			expectedApproval: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile := tt.tool.GetEnforcerProfile()

			if profile.RiskLevel != tt.expectedRisk {
				t.Errorf("Expected risk %s, got %s", tt.expectedRisk, profile.RiskLevel)
			}

			if profile.ImpactScope != tt.expectedImpact {
				t.Errorf("Expected impact %s, got %s", tt.expectedImpact, profile.ImpactScope)
			}

			if profile.ApprovalReq != tt.expectedApproval {
				t.Errorf("Expected approval %v, got %v", tt.expectedApproval, profile.ApprovalReq)
			}
		})
	}
}

func TestEnforcerProfile_HighRiskTools(t *testing.T) {
	highRiskTools := []struct {
		name string
		tool interface {
			GetEnforcerProfile() framework.EnforcerProfile
		}
	}{
		{"CreateChannelTool", &CreateChannelTool{}},
		{"ArchiveChannelTool", &ArchiveChannelTool{}},
		{"DeleteMessageTool", &DeleteMessageTool{}},
	}

	for _, tt := range highRiskTools {
		t.Run(tt.name, func(t *testing.T) {
			profile := tt.tool.GetEnforcerProfile()

			if profile.RiskLevel != framework.RiskHigh {
				t.Errorf("%s should have RiskHigh, got %s", tt.name, profile.RiskLevel)
			}

			if !profile.ApprovalReq {
				t.Errorf("%s should require approval", tt.name)
			}
		})
	}
}

// ============================================================================
// Test Suite 4: Tool Schemas
// ============================================================================

func TestToolSchemas(t *testing.T) {
	// Test that tools have valid schemas
	tool1 := &ListChannelsTool{}
	schema1 := tool1.Schema()
	if schema1.Type != "object" {
		t.Errorf("Expected schema type 'object', got '%s'", schema1.Type)
	}

	tool2 := &PostMessageTool{}
	schema2 := tool2.Schema()
	if schema2.Type != "object" {
		t.Errorf("Expected schema type 'object', got '%s'", schema2.Type)
	}
}

// ============================================================================
// Test Suite 5: Argument Parsing Helpers
// ============================================================================

func TestGetString(t *testing.T) {
	args := map[string]interface{}{
		"present": "value",
		"number":  42,
	}

	// Test present key
	if got := getString(args, "present", "default"); got != "value" {
		t.Errorf("Expected 'value', got '%s'", got)
	}

	// Test missing key returns default
	if got := getString(args, "missing", "default"); got != "default" {
		t.Errorf("Expected 'default', got '%s'", got)
	}

	// Test non-string value returns default
	if got := getString(args, "number", "default"); got != "default" {
		t.Errorf("Expected 'default' for non-string, got '%s'", got)
	}
}

func TestGetBool(t *testing.T) {
	args := map[string]interface{}{
		"true":   true,
		"false":  false,
		"string": "true",
	}

	// Test true value
	if got := getBool(args, "true", false); !got {
		t.Error("Expected true")
	}

	// Test false value
	if got := getBool(args, "false", true); got {
		t.Error("Expected false")
	}

	// Test missing key returns default
	if got := getBool(args, "missing", true); !got {
		t.Error("Expected default true for missing key")
	}

	// Test non-bool value returns default
	if got := getBool(args, "string", false); got {
		t.Error("Expected default false for non-bool")
	}
}

func TestGetInt(t *testing.T) {
	args := map[string]interface{}{
		"number": 42.0, // JSON numbers are float64
		"string": "42",
	}

	// Test number value
	if got := getInt(args, "number", 0); got != 42 {
		t.Errorf("Expected 42, got %d", got)
	}

	// Test missing key returns default
	if got := getInt(args, "missing", 100); got != 100 {
		t.Errorf("Expected 100, got %d", got)
	}

	// Test non-number value returns default
	if got := getInt(args, "string", 0); got != 0 {
		t.Errorf("Expected 0 for non-number, got %d", got)
	}
}

func TestGetRequiredString(t *testing.T) {
	args := map[string]interface{}{
		"present": "value",
	}

	// Test present key
	if got := getRequiredString(args, "present"); got != "value" {
		t.Errorf("Expected 'value', got '%s'", got)
	}

	// Test missing key returns empty string
	if got := getRequiredString(args, "missing"); got != "" {
		t.Errorf("Expected empty string for missing key, got '%s'", got)
	}
}

// ============================================================================
// Test Suite 6: Client Requirements
// ============================================================================

func TestClientRequireBot(t *testing.T) {
	// Note: requireBot checks the actual botClient field, not hasBotToken
	// In production, NewClient() initializes these properly from environment
	// For this test, we verify the behavior when client is nil vs not nil

	// Client without initialized botClient (simulating missing token)
	clientNoToken := &Client{}
	if err := clientNoToken.requireBot(); err == nil {
		t.Error("requireBot should error when botClient is nil")
	}
}

func TestClientRequireUser(t *testing.T) {
	// Client with user token
	clientWithToken := &Client{
		hasUserToken: true,
	}
	if err := clientWithToken.requireUser(); err != nil {
		t.Errorf("requireUser should not error with user token: %v", err)
	}

	// Client without user token
	clientNoToken := &Client{
		hasUserToken: false,
	}
	if err := clientNoToken.requireUser(); err == nil {
		t.Error("requireUser should error without user token")
	}
}

// ============================================================================
// Test Suite 7: Self-Reporting Capability
// ============================================================================

func TestSelfReportingWithoutTokens(t *testing.T) {
	// This tests that the server can self-report tool metadata
	// even without actual Slack tokens
	server := NewServer()
	server.Initialize()

	tools := server.ListTools()
	if len(tools) == 0 {
		t.Fatal("Server should be able to self-report tools without tokens")
	}

	// Verify all tools have valid metadata
	for _, toolName := range tools {
		ctx := context.Background()
		// Execute tool to trigger metadata check
		// Most tools will fail with "not connected" error, which is expected
		_, err := server.ExecuteTool(ctx, toolName, map[string]interface{}{})
		// We expect errors for most tools due to missing required args or no connection
		// but we should NOT get "tool not found"
		if err != nil && err.Error() == "tool '"+toolName+"' not found" {
			t.Errorf("Tool '%s' should be registered", toolName)
		}
	}
}

// ============================================================================
// Test Suite 8: Tool Names and Descriptions
// ============================================================================

func TestToolMetadata(t *testing.T) {
	tests := []struct {
		name string
		tool interface {
			Name() string
			Description() string
		}
		wantName    string
		wantDescHas string // Check that description contains this substring
	}{
		{
			name:        "ListChannelsTool",
			tool:        &ListChannelsTool{},
			wantName:    "list_channels",
			wantDescHas: "channels",
		},
		{
			name:        "PostMessageTool",
			tool:        &PostMessageTool{},
			wantName:    "post_message",
			wantDescHas: "message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tool.Name(); got != tt.wantName {
				t.Errorf("Name() = %v, want %v", got, tt.wantName)
			}
			if desc := tt.tool.Description(); !contains(desc, tt.wantDescHas) {
				t.Errorf("Description() = %v, should contain %v", desc, tt.wantDescHas)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(substr) == 0 || len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// ============================================================================
// Test Suite 9: Integration Tests
// ============================================================================

func TestServerWithWriteEnabled(t *testing.T) {
	// Create server and enable writes
	server := NewServer()
	server.SetWriteEnabled(true)
	server.Initialize()

	// Verify write tools are registered
	tools := server.ListTools()
	hasWriteTools := false
	for _, tool := range tools {
		if tool == "post_message" || tool == "create_channel" {
			hasWriteTools = true
			break
		}
	}

	if !hasWriteTools {
		t.Error("Write tools should be registered when write-enabled")
	}
}

func TestToolExecutionRequiresConnection(t *testing.T) {
	server := NewServer()
	server.Initialize()

	// Try to execute a tool without Slack connection
	ctx := context.Background()
	_, err := server.ExecuteTool(ctx, "list_channels", map[string]interface{}{})

	// Should fail because no Slack token is set
	if err == nil {
		t.Error("Expected error when executing tool without Slack connection")
	}
}
