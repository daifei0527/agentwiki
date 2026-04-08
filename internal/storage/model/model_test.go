package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNowMillis(t *testing.T) {
	millis := NowMillis()
	assert.True(t, millis > 0)
}

func TestGetLevelWeight(t *testing.T) {
	tests := []struct {
		level    int32
		expected float64
	}{
		{UserLevelLv0, 0.0},
		{UserLevelLv1, 1.0},
		{UserLevelLv2, 1.2},
		{UserLevelLv3, 1.5},
		{UserLevelLv4, 2.0},
		{UserLevelLv5, 3.0},
		{6, 0.0}, // 无效级别
		{-1, 0.0}, // 负数级别
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Level%d", tt.level), func(t *testing.T) {
			weight := GetLevelWeight(tt.level)
			assert.Equal(t, tt.expected, weight)
		})
	}
}

func TestKnowledgeEntry(t *testing.T) {
	testTitle := "Test Entry"
	testContent := "# Test Content"
	testCategory := "test/category"
	testCreatedBy := "test-creator"

	// 测试创建新条目
	entry := NewKnowledgeEntry(testTitle, testContent, testCategory, testCreatedBy)
	assert.NotNil(t, entry)
	assert.NotEmpty(t, entry.ID)
	assert.Equal(t, testTitle, entry.Title)
	assert.Equal(t, testContent, entry.Content)
	assert.Equal(t, testCategory, entry.Category)
	assert.Equal(t, testCreatedBy, entry.CreatedBy)
	assert.Equal(t, int64(1), entry.Version)
	assert.NotZero(t, entry.CreatedAt)
	assert.NotZero(t, entry.UpdatedAt)
	assert.Equal(t, entry.CreatedAt, entry.UpdatedAt)
	assert.Equal(t, 0.0, entry.Score)
	assert.Equal(t, int32(0), entry.ScoreCount)
	assert.Equal(t, EntryStatusDraft, entry.Status)
	assert.Equal(t, "CC-BY-SA-4.0", entry.License)
	assert.NotEmpty(t, entry.ContentHash)

	// 测试计算内容哈希
	originalHash := entry.ContentHash
	entry.Content = "Updated content"
	newHash := entry.ComputeContentHash()
	assert.NotEqual(t, originalHash, newHash)

	// 测试JSON序列化和反序列化
	jsonData, err := entry.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	newEntry := &KnowledgeEntry{}
	err = newEntry.FromJSON(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, entry.ID, newEntry.ID)
	assert.Equal(t, entry.Title, newEntry.Title)
	assert.Equal(t, entry.Content, newEntry.Content)
	assert.Equal(t, entry.Category, newEntry.Category)
	assert.Equal(t, entry.CreatedBy, newEntry.CreatedBy)
}

func TestUser(t *testing.T) {
	user := &User{
		PublicKey:       "test-public-key",
		AgentName:       "Test Agent",
		UserLevel:       UserLevelLv1,
		Email:           "test@example.com",
		EmailVerified:   true,
		Phone:           "1234567890",
		RegisteredAt:    1234567890,
		LastActive:      1234567890,
		ContributionCnt: 10,
		RatingCnt:       5,
		NodeId:          "test-node",
		Status:          UserStatusActive,
	}

	// 测试JSON序列化和反序列化
	jsonData, err := user.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	newUser := &User{}
	err = newUser.FromJSON(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, user.PublicKey, newUser.PublicKey)
	assert.Equal(t, user.AgentName, newUser.AgentName)
	assert.Equal(t, user.UserLevel, newUser.UserLevel)
	assert.Equal(t, user.Email, newUser.Email)
	assert.Equal(t, user.EmailVerified, newUser.EmailVerified)
	assert.Equal(t, user.Phone, newUser.Phone)
	assert.Equal(t, user.RegisteredAt, newUser.RegisteredAt)
	assert.Equal(t, user.LastActive, newUser.LastActive)
	assert.Equal(t, user.ContributionCnt, newUser.ContributionCnt)
	assert.Equal(t, user.RatingCnt, newUser.RatingCnt)
	assert.Equal(t, user.NodeId, newUser.NodeId)
	assert.Equal(t, user.Status, newUser.Status)
}

func TestRating(t *testing.T) {
	rating := &Rating{
		ID:           "test-rating-id",
		EntryId:      "test-entry-id",
		RaterPubkey:  "test-rater-pubkey",
		Score:        4.5,
		Weight:       1.0,
		WeightedScore: 4.5,
		RatedAt:      1234567890,
		Comment:      "Great entry!",
	}

	// 测试JSON序列化和反序列化
	jsonData, err := rating.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	newRating := &Rating{}
	err = newRating.FromJSON(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, rating.ID, newRating.ID)
	assert.Equal(t, rating.EntryId, newRating.EntryId)
	assert.Equal(t, rating.RaterPubkey, newRating.RaterPubkey)
	assert.Equal(t, rating.Score, newRating.Score)
	assert.Equal(t, rating.Weight, newRating.Weight)
	assert.Equal(t, rating.WeightedScore, newRating.WeightedScore)
	assert.Equal(t, rating.RatedAt, newRating.RatedAt)
	assert.Equal(t, rating.Comment, newRating.Comment)
}

func TestCategory(t *testing.T) {
	category := &Category{
		ID:          "test-category-id",
		Path:        "test/category",
		Name:        "Test Category",
		ParentId:    "parent-id",
		Level:       2,
		SortOrder:   10,
		IsBuiltin:   true,
		MaintainedBy: "maintainer-pubkey",
		CreatedAt:   1234567890,
	}

	// 测试JSON序列化和反序列化
	jsonData, err := category.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	newCategory := &Category{}
	err = newCategory.FromJSON(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, category.ID, newCategory.ID)
	assert.Equal(t, category.Path, newCategory.Path)
	assert.Equal(t, category.Name, newCategory.Name)
	assert.Equal(t, category.ParentId, newCategory.ParentId)
	assert.Equal(t, category.Level, newCategory.Level)
	assert.Equal(t, category.SortOrder, newCategory.SortOrder)
	assert.Equal(t, category.IsBuiltin, newCategory.IsBuiltin)
	assert.Equal(t, category.MaintainedBy, newCategory.MaintainedBy)
	assert.Equal(t, category.CreatedAt, newCategory.CreatedAt)
}

func TestNodeInfo(t *testing.T) {
	nodeInfo := &NodeInfo{
		NodeId:         "test-node-id",
		NodeType:       NodeTypeFull,
		PeerId:         "test-peer-id",
		PublicKey:      "test-public-key",
		Addresses:      []string{"/ip4/127.0.0.1/tcp/12345"},
		Version:        "v1.0.0",
		EntryCount:     1000,
		CategoryMirror: []string{"test/category"},
		LastSync:       1234567890,
		Uptime:         3600,
		AllowMirror:    true,
		BandwidthLimit: 1048576,
	}

	// 测试JSON序列化和反序列化
	jsonData, err := nodeInfo.ToJSON()
	assert.NoError(t, err)
	assert.NotEmpty(t, jsonData)

	newNodeInfo := &NodeInfo{}
	err = newNodeInfo.FromJSON(jsonData)
	assert.NoError(t, err)
	assert.Equal(t, nodeInfo.NodeId, newNodeInfo.NodeId)
	assert.Equal(t, nodeInfo.NodeType, newNodeInfo.NodeType)
	assert.Equal(t, nodeInfo.PeerId, newNodeInfo.PeerId)
	assert.Equal(t, nodeInfo.PublicKey, newNodeInfo.PublicKey)
	assert.Equal(t, nodeInfo.Addresses, newNodeInfo.Addresses)
	assert.Equal(t, nodeInfo.Version, newNodeInfo.Version)
	assert.Equal(t, nodeInfo.EntryCount, newNodeInfo.EntryCount)
	assert.Equal(t, nodeInfo.CategoryMirror, newNodeInfo.CategoryMirror)
	assert.Equal(t, nodeInfo.LastSync, newNodeInfo.LastSync)
	assert.Equal(t, nodeInfo.Uptime, newNodeInfo.Uptime)
	assert.Equal(t, nodeInfo.AllowMirror, newNodeInfo.AllowMirror)
	assert.Equal(t, nodeInfo.BandwidthLimit, newNodeInfo.BandwidthLimit)
}

func TestIsCJK(t *testing.T) {
	tests := []struct {
		char     rune
		expected bool
	}{
		{'a', false},      // 英文
		{'1', false},      // 数字
		{'中', true},      // 中文
		{'日', true},      // 日文
		{'韩', true},      // 韩文
		{'，', true},      // 中文标点
		{'。', true},      // 中文句号
		{'！', true},      // 中文感叹号
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result := IsCJK(tt.char)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestContainsCJK(t *testing.T) {
	tests := []struct {
		str      string
		expected bool
	}{
		{"Hello", false},
		{"12345", false},
		{"Hello 世界", true},
		{"中文测试", true},
		{"日本語", true},
		{"한국어", true},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			result := ContainsCJK(tt.str)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNormalizeKey(t *testing.T) {
	tests := []struct {
		str      string
		expected string
	}{
		{"Hello World", "hello world"},
		{"  TEST  ", "test"},
		{"UPPERCASE", "uppercase"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			result := NormalizeKey(tt.str)
			assert.Equal(t, tt.expected, result)
		})
	}
}
