// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"time"

	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/mq/kafka"
)

const StructTagName = "yaml"

type Path string

type Index int

type CacheConfig struct {
	Topic         string `yaml:"topic"`
	SlotNum       int    `yaml:"slotNum"`
	SlotSize      int    `yaml:"slotSize"`
	SuccessExpire int    `yaml:"successExpire"`
	FailedExpire  int    `yaml:"failedExpire"`
}

type LocalCache struct {
	User         CacheConfig `yaml:"user"`
	Group        CacheConfig `yaml:"group"`
	Friend       CacheConfig `yaml:"friend"`
	Conversation CacheConfig `yaml:"conversation"`
}

type Log struct {
	StorageLocation     string `yaml:"storageLocation"`
	RotationTime        uint   `yaml:"rotationTime"`
	RemainRotationCount uint   `yaml:"remainRotationCount"`
	RemainLogLevel      int    `yaml:"remainLogLevel"`
	IsStdout            bool   `yaml:"isStdout"`
	IsJson              bool   `yaml:"isJson"`
	IsSimplify          bool   `yaml:"isSimplify"`
	WithStack           bool   `yaml:"withStack"`
}

type Mongo struct {
	URI         string   `yaml:"uri"`
	Address     []string `yaml:"address"`
	Database    string   `yaml:"database"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	AuthSource  string   `yaml:"authSource"`
	MaxPoolSize int      `yaml:"maxPoolSize"`
	MaxRetry    int      `yaml:"maxRetry"`
}
type Kafka struct {
	Username           string   `yaml:"username"`
	Password           string   `yaml:"password"`
	ProducerAck        string   `yaml:"producerAck"`
	CompressType       string   `yaml:"compressType"`
	Address            []string `yaml:"address"`
	ToRedisTopic       string   `yaml:"toRedisTopic"`
	ToMongoTopic       string   `yaml:"toMongoTopic"`
	ToPushTopic        string   `yaml:"toPushTopic"`
	ToOfflinePushTopic string   `yaml:"toOfflinePushTopic"`
	ToRedisGroupID     string   `yaml:"toRedisGroupID"`
	ToMongoGroupID     string   `yaml:"toMongoGroupID"`
	ToPushGroupID      string   `yaml:"toPushGroupID"`
	ToOfflineGroupID   string   `yaml:"toOfflinePushGroupID"`

	Tls TLSConfig `yaml:"tls"`
}
type TLSConfig struct {
	EnableTLS          bool   `yaml:"enableTLS"`
	CACrt              string `yaml:"caCrt"`
	ClientCrt          string `yaml:"clientCrt"`
	ClientKey          string `yaml:"clientKey"`
	ClientKeyPwd       string `yaml:"clientKeyPwd"`
	InsecureSkipVerify bool   `yaml:"insecureSkipVerify"`
}

type API struct {
	Api struct {
		ListenIP         string `yaml:"listenIP"`
		Ports            []int  `yaml:"ports"`
		CompressionLevel int    `yaml:"compressionLevel"`
	} `yaml:"api"`
	Prometheus struct {
		Enable       bool   `yaml:"enable"`
		AutoSetPorts bool   `yaml:"autoSetPorts"`
		Ports        []int  `yaml:"ports"`
		GrafanaURL   string `yaml:"grafanaURL"`
	} `yaml:"prometheus"`
}

type CronTask struct {
	CronExecuteTime   string   `yaml:"cronExecuteTime"`
	RetainChatRecords int      `yaml:"retainChatRecords"`
	FileExpireTime    int      `yaml:"fileExpireTime"`
	DeleteObjectType  []string `yaml:"deleteObjectType"`
}

type OfflinePushConfig struct {
	Enable bool   `yaml:"enable"`
	Title  string `yaml:"title"`
	Desc   string `yaml:"desc"`
	Ext    string `yaml:"ext"`
}

type NotificationConfig struct {
	IsSendMsg        bool              `yaml:"isSendMsg"`
	ReliabilityLevel int               `yaml:"reliabilityLevel"`
	UnreadCount      bool              `yaml:"unreadCount"`
	OfflinePush      OfflinePushConfig `yaml:"offlinePush"`
}

type Prometheus struct {
	Enable bool  `yaml:"enable"`
	Ports  []int `yaml:"ports"`
}

type MsgGateway struct {
	RPC         RPC        `yaml:"rpc"`
	Prometheus  Prometheus `yaml:"prometheus"`
	ListenIP    string     `yaml:"listenIP"`
	LongConnSvr struct {
		Ports               []int `yaml:"ports"`
		WebsocketMaxConnNum int   `yaml:"websocketMaxConnNum"`
		WebsocketMaxMsgLen  int   `yaml:"websocketMaxMsgLen"`
		WebsocketTimeout    int   `yaml:"websocketTimeout"`
	} `yaml:"longConnSvr"`
}

type MsgTransfer struct {
	Prometheus struct {
		Enable       bool  `yaml:"enable"`
		AutoSetPorts bool  `yaml:"autoSetPorts"`
		Ports        []int `yaml:"ports"`
	} `yaml:"prometheus"`
}

type Push struct {
	RPC                  RPC        `yaml:"rpc"`
	Prometheus           Prometheus `yaml:"prometheus"`
	MaxConcurrentWorkers int        `yaml:"maxConcurrentWorkers"`
	Enable               string     `yaml:"enable"`
	GeTui                struct {
		PushUrl      string `yaml:"pushUrl"`
		MasterSecret string `yaml:"masterSecret"`
		AppKey       string `yaml:"appKey"`
		Intent       string `yaml:"intent"`
		ChannelID    string `yaml:"channelID"`
		ChannelName  string `yaml:"channelName"`
	} `yaml:"geTui"`
	FCM struct {
		FilePath string `yaml:"filePath"`
		AuthURL  string `yaml:"authURL"`
	} `yaml:"fcm"`
	JPush struct {
		AppKey       string `yaml:"appKey"`
		MasterSecret string `yaml:"masterSecret"`
		PushURL      string `yaml:"pushURL"`
		PushIntent   string `yaml:"pushIntent"`
	} `yaml:"jpush"`
	IOSPush struct {
		PushSound  string `yaml:"pushSound"`
		BadgeCount bool   `yaml:"badgeCount"`
		Production bool   `yaml:"production"`
	} `yaml:"iosPush"`
	FullUserCache bool `yaml:"fullUserCache"`
}

type Auth struct {
	RPC         RPC        `yaml:"rpc"`
	Prometheus  Prometheus `yaml:"prometheus"`
	TokenPolicy struct {
		Expire int64 `yaml:"expire"`
	} `yaml:"tokenPolicy"`
}

type Conversation struct {
	RPC        RPC        `yaml:"rpc"`
	Prometheus Prometheus `yaml:"prometheus"`
}

type Friend struct {
	RPC        RPC        `yaml:"rpc"`
	Prometheus Prometheus `yaml:"prometheus"`
}

type Group struct {
	RPC                        RPC        `yaml:"rpc"`
	Prometheus                 Prometheus `yaml:"prometheus"`
	EnableHistoryForNewMembers bool       `yaml:"enableHistoryForNewMembers"`
}

type Msg struct {
	RPC          RPC        `yaml:"rpc"`
	Prometheus   Prometheus `yaml:"prometheus"`
	FriendVerify bool       `yaml:"friendVerify"`
}

type User struct {
	RPC        RPC        `yaml:"rpc"`
	Prometheus Prometheus `yaml:"prometheus"`
}

type RPC struct {
	RegisterIP string `yaml:"registerIP"`
	ListenIP   string `yaml:"listenIP"`
	Ports      []int  `yaml:"ports"`
	AutoPort   bool   `yaml:"autoPort"`
}

type Redis struct {
	Disable     bool     `yaml:"-"`
	Address     []string `yaml:"address"`
	Username    string   `yaml:"username"`
	Password    string   `yaml:"password"`
	ClusterMode bool     `yaml:"clusterMode"`
	DB          int      `yaml:"storage"`
	MaxRetry    int      `yaml:"maxRetry"`
	PoolSize    int      `yaml:"poolSize"`
}

type BeforeConfig struct {
	Enable         bool    `yaml:"enable"`
	Timeout        int     `yaml:"timeout"`
	FailedContinue bool    `yaml:"failedContinue"`
	DeniedTypes    []int32 `yaml:"deniedTypes"`
}

type AfterConfig struct {
	Enable       bool     `yaml:"enable"`
	Timeout      int      `yaml:"timeout"`
	AttentionIds []string `yaml:"attentionIds"`
	DeniedTypes  []int32  `yaml:"deniedTypes"`
}

type Share struct {
	Secret        string     `yaml:"secret"`
	IMAdminUserID []string   `yaml:"imAdminUserID"`
	MultiLogin    MultiLogin `yaml:"multiLogin"`
}

type MultiLogin struct {
	Policy       int `yaml:"policy"`
	MaxNumOneEnd int `yaml:"maxNumOneEnd"`
}

type RpcService struct {
	User           string `yaml:"user"`
	Friend         string `yaml:"friend"`
	Msg            string `yaml:"msg"`
	Push           string `yaml:"push"`
	MessageGateway string `yaml:"messageGateway"`
	Group          string `yaml:"group"`
	Auth           string `yaml:"auth"`
	Conversation   string `yaml:"conversation"`
	Third          string `yaml:"third"`
	Customer       string `yaml:"customer"`
}

func (r *RpcService) GetServiceNames() []string {
	return []string{
		r.User,
		r.Friend,
		r.Msg,
		r.Push,
		r.MessageGateway,
		r.Group,
		r.Auth,
		r.Conversation,
		r.Third,
		r.Customer,
	}
}

// FullConfig stores all configurations for before and after events
type Webhooks struct {
	URL                      string       `yaml:"url"`
	BeforeSendSingleMsg      BeforeConfig `yaml:"beforeSendSingleMsg"`
	BeforeUpdateUserInfoEx   BeforeConfig `yaml:"beforeUpdateUserInfoEx"`
	AfterUpdateUserInfoEx    AfterConfig  `yaml:"afterUpdateUserInfoEx"`
	AfterSendSingleMsg       AfterConfig  `yaml:"afterSendSingleMsg"`
	BeforeSendGroupMsg       BeforeConfig `yaml:"beforeSendGroupMsg"`
	BeforeMsgModify          BeforeConfig `yaml:"beforeMsgModify"`
	AfterSendGroupMsg        AfterConfig  `yaml:"afterSendGroupMsg"`
	AfterUserOnline          AfterConfig  `yaml:"afterUserOnline"`
	AfterUserOffline         AfterConfig  `yaml:"afterUserOffline"`
	AfterUserKickOff         AfterConfig  `yaml:"afterUserKickOff"`
	BeforeOfflinePush        BeforeConfig `yaml:"beforeOfflinePush"`
	BeforeOnlinePush         BeforeConfig `yaml:"beforeOnlinePush"`
	BeforeGroupOnlinePush    BeforeConfig `yaml:"beforeGroupOnlinePush"`
	BeforeAddFriend          BeforeConfig `yaml:"beforeAddFriend"`
	BeforeUpdateUserInfo     BeforeConfig `yaml:"beforeUpdateUserInfo"`
	AfterUpdateUserInfo      AfterConfig  `yaml:"afterUpdateUserInfo"`
	BeforeCreateGroup        BeforeConfig `yaml:"beforeCreateGroup"`
	AfterCreateGroup         AfterConfig  `yaml:"afterCreateGroup"`
	BeforeMemberJoinGroup    BeforeConfig `yaml:"beforeMemberJoinGroup"`
	BeforeSetGroupMemberInfo BeforeConfig `yaml:"beforeSetGroupMemberInfo"`
	AfterSetGroupMemberInfo  AfterConfig  `yaml:"afterSetGroupMemberInfo"`
	AfterQuitGroup           AfterConfig  `yaml:"afterQuitGroup"`
	AfterKickGroupMember     AfterConfig  `yaml:"afterKickGroupMember"`
	AfterDismissGroup        AfterConfig  `yaml:"afterDismissGroup"`
	BeforeApplyJoinGroup     BeforeConfig `yaml:"beforeApplyJoinGroup"`
	AfterGroupMsgRead        AfterConfig  `yaml:"afterGroupMsgRead"`
	AfterSingleMsgRead       AfterConfig  `yaml:"afterSingleMsgRead"`
	BeforeUserRegister       BeforeConfig `yaml:"beforeUserRegister"`
	AfterUserRegister        AfterConfig  `yaml:"afterUserRegister"`
	AfterTransferGroupOwner  AfterConfig  `yaml:"afterTransferGroupOwner"`
	BeforeSetFriendRemark    BeforeConfig `yaml:"beforeSetFriendRemark"`
	AfterSetFriendRemark     AfterConfig  `yaml:"afterSetFriendRemark"`
	AfterGroupMsgRevoke      AfterConfig  `yaml:"afterGroupMsgRevoke"`
	AfterJoinGroup           AfterConfig  `yaml:"afterJoinGroup"`
	BeforeInviteUserToGroup  BeforeConfig `yaml:"beforeInviteUserToGroup"`
	AfterSetGroupInfo        AfterConfig  `yaml:"afterSetGroupInfo"`
	BeforeSetGroupInfo       BeforeConfig `yaml:"beforeSetGroupInfo"`
	AfterSetGroupInfoEx      AfterConfig  `yaml:"afterSetGroupInfoEx"`
	BeforeSetGroupInfoEx     BeforeConfig `yaml:"beforeSetGroupInfoEx"`
	AfterRevokeMsg           AfterConfig  `yaml:"afterRevokeMsg"`
	BeforeAddBlack           BeforeConfig `yaml:"beforeAddBlack"`
	AfterAddFriend           AfterConfig  `yaml:"afterAddFriend"`
	BeforeAddFriendAgree     BeforeConfig `yaml:"beforeAddFriendAgree"`
	AfterAddFriendAgree      AfterConfig  `yaml:"afterAddFriendAgree"`
	AfterDeleteFriend        AfterConfig  `yaml:"afterDeleteFriend"`
	BeforeImportFriends      BeforeConfig `yaml:"beforeImportFriends"`
	AfterImportFriends       AfterConfig  `yaml:"afterImportFriends"`
	AfterRemoveBlack         AfterConfig  `yaml:"afterRemoveBlack"`
}

type ZooKeeper struct {
	Schema   string   `yaml:"schema"`
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
}

type Discovery struct {
	Enable     string     `yaml:"enable"`
	Etcd       Etcd       `yaml:"etcd"`
	Kubernetes Kubernetes `yaml:"kubernetes"`
	RpcService RpcService `yaml:"rpcService"`
}

type Kubernetes struct {
	Namespace string `yaml:"namespace"`
}

type Etcd struct {
	RootDirectory string   `yaml:"rootDirectory"`
	Address       []string `yaml:"address"`
	Username      string   `yaml:"username"`
	Password      string   `yaml:"password"`
}

func (m *Mongo) Build() *mongoutil.Config {
	return &mongoutil.Config{
		Uri:         m.URI,
		Address:     m.Address,
		Database:    m.Database,
		Username:    m.Username,
		Password:    m.Password,
		AuthSource:  m.AuthSource,
		MaxPoolSize: m.MaxPoolSize,
		MaxRetry:    m.MaxRetry,
	}
}

func (r *Redis) Build() *redisutil.Config {
	return &redisutil.Config{
		ClusterMode: r.ClusterMode,
		Address:     r.Address,
		Username:    r.Username,
		Password:    r.Password,
		DB:          r.DB,
		MaxRetry:    r.MaxRetry,
		PoolSize:    r.PoolSize,
	}
}

func (k *Kafka) Build() *kafka.Config {
	return &kafka.Config{
		Username:     k.Username,
		Password:     k.Password,
		ProducerAck:  k.ProducerAck,
		CompressType: k.CompressType,
		Addr:         k.Address,
		TLS: kafka.TLSConfig{
			EnableTLS:          k.Tls.EnableTLS,
			CACrt:              k.Tls.CACrt,
			ClientCrt:          k.Tls.ClientCrt,
			ClientKey:          k.Tls.ClientKey,
			ClientKeyPwd:       k.Tls.ClientKeyPwd,
			InsecureSkipVerify: k.Tls.InsecureSkipVerify,
		},
	}
}

func (l *CacheConfig) Failed() time.Duration {
	return time.Second * time.Duration(l.FailedExpire)
}

func (l *CacheConfig) Success() time.Duration {
	return time.Second * time.Duration(l.SuccessExpire)
}

func (l *CacheConfig) Enable() bool {
	return l.Topic != "" && l.SlotNum > 0 && l.SlotSize > 0
}

type AllConfig struct {
	Discovery    Discovery
	Kafka        Kafka
	LocalCache   LocalCache
	Log          Log
	Mongo        Mongo
	API          API
	CronTask     CronTask
	MsgGateway   MsgGateway
	MsgTransfer  MsgTransfer
	Push         Push
	Auth         Auth
	Conversation Conversation
	Friend       Friend
	Group        Group
	Msg          Msg
	User         User
	Redis        Redis
	Share        Share
	Webhooks     Webhooks
}

func (a *AllConfig) Name2Config(name string) any {
	switch name {
	case a.Discovery.GetConfigFileName():
		return a.Discovery
	case a.Kafka.GetConfigFileName():
		return a.Kafka
	case a.LocalCache.GetConfigFileName():
		return a.LocalCache
	case a.Log.GetConfigFileName():
		return a.Log
	case a.Mongo.GetConfigFileName():
		return a.Mongo
	case a.API.GetConfigFileName():
		return a.API
	case a.CronTask.GetConfigFileName():
		return a.CronTask
	case a.MsgGateway.GetConfigFileName():
		return a.MsgGateway
	case a.MsgTransfer.GetConfigFileName():
		return a.MsgTransfer
	case a.Push.GetConfigFileName():
		return a.Push
	case a.Auth.GetConfigFileName():
		return a.Auth
	case a.Conversation.GetConfigFileName():
		return a.Conversation
	case a.Friend.GetConfigFileName():
		return a.Friend
	case a.Group.GetConfigFileName():
		return a.Group
	case a.Msg.GetConfigFileName():
		return a.Msg
	case a.User.GetConfigFileName():
		return a.User
	case a.Redis.GetConfigFileName():
		return a.Redis
	case a.Share.GetConfigFileName():
		return a.Share
	case a.Webhooks.GetConfigFileName():
		return a.Webhooks
	default:
		return nil
	}
}

func (a *AllConfig) GetConfigNames() []string {
	return []string{
		a.Discovery.GetConfigFileName(),
		a.Kafka.GetConfigFileName(),
		a.LocalCache.GetConfigFileName(),
		a.Log.GetConfigFileName(),
		a.Mongo.GetConfigFileName(),
		a.API.GetConfigFileName(),
		a.CronTask.GetConfigFileName(),
		a.MsgGateway.GetConfigFileName(),
		a.MsgTransfer.GetConfigFileName(),
		a.Push.GetConfigFileName(),
		a.Auth.GetConfigFileName(),
		a.Conversation.GetConfigFileName(),
		a.Friend.GetConfigFileName(),
		a.Group.GetConfigFileName(),
		a.Msg.GetConfigFileName(),
		a.User.GetConfigFileName(),
		a.Redis.GetConfigFileName(),
		a.Share.GetConfigFileName(),
		a.Webhooks.GetConfigFileName(),
	}
}

const (
	FileName                         = "config.yaml"
	DiscoveryConfigFilename          = "discovery.yml"
	KafkaConfigFileName              = "kafka.yml"
	LocalCacheConfigFileName         = "local-cache.yml"
	LogConfigFileName                = "log.yml"
	MinioConfigFileName              = "minio.yml"
	MongodbConfigFileName            = "mongodb.yml"
	NotificationFileName             = "notification.yml"
	OpenIMAPICfgFileName             = "openim-api.yml"
	OpenIMCronTaskCfgFileName        = "openim-crontask.yml"
	OpenIMMsgGatewayCfgFileName      = "openim-msggateway.yml"
	OpenIMMsgTransferCfgFileName     = "openim-msgtransfer.yml"
	OpenIMPushCfgFileName            = "openim-push.yml"
	OpenIMRPCAuthCfgFileName         = "openim-rpc-auth.yml"
	OpenIMRPCConversationCfgFileName = "openim-rpc-conversation.yml"
	OpenIMRPCFriendCfgFileName       = "openim-rpc-friend.yml"
	OpenIMRPCGroupCfgFileName        = "openim-rpc-group.yml"
	OpenIMRPCMsgCfgFileName          = "openim-rpc-msg.yml"
	OpenIMRPCThirdCfgFileName        = "openim-rpc-third.yml"
	OpenIMRPCUserCfgFileName         = "openim-rpc-user.yml"
	RedisConfigFileName              = "redis.yml"
	ShareFileName                    = "share.yml"
	WebhooksConfigFileName           = "webhooks.yml"
)

func (d *Discovery) GetConfigFileName() string {
	return DiscoveryConfigFilename
}

func (k *Kafka) GetConfigFileName() string {
	return KafkaConfigFileName
}

func (lc *LocalCache) GetConfigFileName() string {
	return LocalCacheConfigFileName
}

func (l *Log) GetConfigFileName() string {
	return LogConfigFileName
}

func (m *Mongo) GetConfigFileName() string {
	return MongodbConfigFileName
}

func (a *API) GetConfigFileName() string {
	return OpenIMAPICfgFileName
}

func (ct *CronTask) GetConfigFileName() string {
	return OpenIMCronTaskCfgFileName
}

func (mg *MsgGateway) GetConfigFileName() string {
	return OpenIMMsgGatewayCfgFileName
}

func (mt *MsgTransfer) GetConfigFileName() string {
	return OpenIMMsgTransferCfgFileName
}

func (p *Push) GetConfigFileName() string {
	return OpenIMPushCfgFileName
}

func (a *Auth) GetConfigFileName() string {
	return OpenIMRPCAuthCfgFileName
}

func (c *Conversation) GetConfigFileName() string {
	return OpenIMRPCConversationCfgFileName
}

func (f *Friend) GetConfigFileName() string {
	return OpenIMRPCFriendCfgFileName
}

func (g *Group) GetConfigFileName() string {
	return OpenIMRPCGroupCfgFileName
}

func (m *Msg) GetConfigFileName() string {
	return OpenIMRPCMsgCfgFileName
}

func (u *User) GetConfigFileName() string {
	return OpenIMRPCUserCfgFileName
}

func (r *Redis) GetConfigFileName() string {
	return RedisConfigFileName
}

func (s *Share) GetConfigFileName() string {
	return ShareFileName
}

func (w *Webhooks) GetConfigFileName() string {
	return WebhooksConfigFileName
}

type Customer struct {
	RPC RPC `yaml:"rpc"`
}

func (x *Customer) GetConfigFileName() string {
	return "customer.yml"
}
