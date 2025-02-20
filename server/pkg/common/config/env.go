package config

import (
	"strings"
)

var (
	ShareFileName           = "share.yml"
	RedisConfigFileName     = "redis.yml"
	DiscoveryConfigFileName = "discovery.yml"
	MongodbConfigFileName   = "mongodb.yml"
	LogConfigFileName       = "log.yml"
	ChatAPIChatCfgFileName  = "wiseengage-api-customerservice.yml"
	ChatRPCChatCfgFileName  = "wiseengage-rpc-customerservice.yml"
)

var EnvPrefixMap map[string]string

func init() {
	EnvPrefixMap = make(map[string]string)
	fileNames := []string{
		ShareFileName,
		RedisConfigFileName,
		DiscoveryConfigFileName,
		MongodbConfigFileName,
		LogConfigFileName,
		ChatAPIChatCfgFileName,
		ChatRPCChatCfgFileName,
	}

	for _, fileName := range fileNames {
		envKey := strings.TrimSuffix(strings.TrimSuffix(fileName, ".yml"), ".yaml")
		envKey = "CHATENV_" + envKey
		envKey = strings.ToUpper(strings.ReplaceAll(envKey, "-", "_"))
		EnvPrefixMap[fileName] = envKey
	}
}

const (
	FlagConf          = "config_folder_path"
	FlagTransferIndex = "index"
)
