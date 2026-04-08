package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, "local", cfg.Node.Type)
	assert.Equal(t, "agentwiki-node-1", cfg.Node.Name)
	assert.Equal(t, "./data", cfg.Node.DataDir)
	assert.Equal(t, "./logs", cfg.Node.LogDir)
	assert.Equal(t, "info", cfg.Node.LogLevel)

	assert.Equal(t, 18530, cfg.Network.ListenPort)
	assert.Equal(t, 18531, cfg.Network.APIPort)
	assert.Empty(t, cfg.Network.SeedNodes)
	assert.True(t, cfg.Network.DHTEnabled)
	assert.True(t, cfg.Network.MDNSEnabled)

	assert.True(t, cfg.Sync.AutoSync)
	assert.Equal(t, 300, cfg.Sync.IntervalSeconds)
	assert.Empty(t, cfg.Sync.MirrorCategories)
	assert.Equal(t, 1024, cfg.Sync.MaxLocalSizeMB)
	assert.Equal(t, "gzip", cfg.Sync.Compression)

	assert.True(t, cfg.Sharing.AllowMirror)
	assert.Equal(t, 100, cfg.Sharing.BandwidthLimitMB)
	assert.Equal(t, 10, cfg.Sharing.MaxConcurrent)

	assert.Equal(t, "./data/keys", cfg.User.PrivateKeyPath)
	assert.Empty(t, cfg.User.Email)
	assert.True(t, cfg.User.AutoRegister)

	assert.False(t, cfg.SMTP.Enabled)
	assert.Empty(t, cfg.SMTP.Host)
	assert.Equal(t, 587, cfg.SMTP.Port)
	assert.Empty(t, cfg.SMTP.Username)
	assert.Empty(t, cfg.SMTP.Password)
	assert.Empty(t, cfg.SMTP.From)

	assert.True(t, cfg.API.Enabled)
	assert.True(t, cfg.API.CORS)
}

func TestLoadConfig(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 写入测试配置
	testConfig := `{
		"node": {
			"type": "seed",
			"name": "test-seed",
			"data_dir": "/tmp/data",
			"log_dir": "/tmp/logs",
			"log_level": "debug"
		},
		"network": {
			"listen_port": 12345,
			"api_port": 12346,
			"seed_nodes": ["/ip4/127.0.0.1/tcp/54321"],
			"dht_enabled": false,
			"mdns_enabled": false
		}
	}`

	err := os.WriteFile(configPath, []byte(testConfig), 0644)
	require.NoError(t, err)

	// 加载配置
	cfg, err := Load(configPath)
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	assert.Equal(t, "seed", cfg.Node.Type)
	assert.Equal(t, "test-seed", cfg.Node.Name)
	assert.Equal(t, "/tmp/data", cfg.Node.DataDir)
	assert.Equal(t, "/tmp/logs", cfg.Node.LogDir)
	assert.Equal(t, "debug", cfg.Node.LogLevel)

	assert.Equal(t, 12345, cfg.Network.ListenPort)
	assert.Equal(t, 12346, cfg.Network.APIPort)
	assert.Len(t, cfg.Network.SeedNodes, 1)
	assert.Equal(t, "/ip4/127.0.0.1/tcp/54321", cfg.Network.SeedNodes[0])
	assert.False(t, cfg.Network.DHTEnabled)
	assert.False(t, cfg.Network.MDNSEnabled)
}

func TestLoadWithEnv(t *testing.T) {
	// 设置环境变量
	err := os.Setenv("AGENTWIKI_NODE_TYPE", "seed")
	require.NoError(t, err)
	err = os.Setenv("AGENTWIKI_NODE_NAME", "env-test-node")
	require.NoError(t, err)
	err = os.Setenv("AGENTWIKI_NETWORK_LISTEN_PORT", "54321")
	require.NoError(t, err)
	err = os.Setenv("AGENTWIKI_SYNC_AUTO_SYNC", "false")
	require.NoError(t, err)

	// 加载默认配置并应用环境变量
	cfg := DefaultConfig()
	cfg = LoadWithEnv(cfg)

	assert.Equal(t, "seed", cfg.Node.Type)
	assert.Equal(t, "env-test-node", cfg.Node.Name)
	assert.Equal(t, 54321, cfg.Network.ListenPort)
	assert.False(t, cfg.Sync.AutoSync)

	// 清理环境变量
	os.Unsetenv("AGENTWIKI_NODE_TYPE")
	os.Unsetenv("AGENTWIKI_NODE_NAME")
	os.Unsetenv("AGENTWIKI_NETWORK_LISTEN_PORT")
	os.Unsetenv("AGENTWIKI_SYNC_AUTO_SYNC")
}

func TestValidate(t *testing.T) {
	// 测试有效配置
	validCfg := DefaultConfig()
	err := Validate(validCfg)
	assert.NoError(t, err)

	// 测试无效的节点类型
	invalidTypeCfg := DefaultConfig()
	invalidTypeCfg.Node.Type = "invalid"
	err = Validate(invalidTypeCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的节点类型")

	// 测试无效的端口
	invalidPortCfg := DefaultConfig()
	invalidPortCfg.Network.ListenPort = 0
	err = Validate(invalidPortCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的监听端口")

	// 测试相同的端口
	samePortCfg := DefaultConfig()
	samePortCfg.Network.ListenPort = 18530
	samePortCfg.Network.APIPort = 18530
	err = Validate(samePortCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "监听端口和 API 端口不能相同")

	// 测试无效的日志级别
	invalidLogLevelCfg := DefaultConfig()
	invalidLogLevelCfg.Node.LogLevel = "invalid"
	err = Validate(invalidLogLevelCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的日志级别")

	// 测试无效的压缩算法
	invalidCompressionCfg := DefaultConfig()
	invalidCompressionCfg.Sync.Compression = "invalid"
	err = Validate(invalidCompressionCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "无效的压缩算法")

	// 测试启用SMTP但缺少配置
	smtpCfg := DefaultConfig()
	smtpCfg.SMTP.Enabled = true
	err = Validate(smtpCfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SMTP 已启用但未配置主机地址")
}

func TestSave(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// 创建测试配置
	cfg := DefaultConfig()
	cfg.Node.Name = "test-save-node"
	cfg.Network.ListenPort = 98765

	// 保存配置
	err := Save(cfg, configPath)
	require.NoError(t, err)

	// 验证文件存在
	_, err = os.Stat(configPath)
	require.NoError(t, err)

	// 重新加载配置
	loadedCfg, err := Load(configPath)
	require.NoError(t, err)

	assert.Equal(t, "test-save-node", loadedCfg.Node.Name)
	assert.Equal(t, 98765, loadedCfg.Network.ListenPort)
}

func TestToJSON(t *testing.T) {
	cfg := DefaultConfig()
	jsonStr, err := cfg.ToJSON()
	require.NoError(t, err)
	assert.NotEmpty(t, jsonStr)
}

func TestString(t *testing.T) {
	cfg := DefaultConfig()
	str := cfg.String()
	assert.NotEmpty(t, str)
}
