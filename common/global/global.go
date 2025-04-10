package global

import (
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// GlobalConfig 应用配置结构体
type GlobalConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	OSS      OSSConfig      `mapstructure:"oss"`
	Wechat   WechatConfig   `mapstructure:"wechat"`
	Shop     ShopConfig     `mapstructure:"shop"`
	Baidu    BaiduConfig    `mapstructure:"baidu"`
	Template TemplateConfig `mapstructure:"template"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string `mapstructure:"driver"`
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBName       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

// RedisConfig redis 配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	Database     int    `mapstructure:"database"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
	DialTimeout  int    `mapstructure:"dial_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
	Compress   bool   `mapstructure:"compress"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	AdminSecretKey string `mapstructure:"admin_secret_key"`
	AdminTTL       int    `mapstructure:"admin_ttl"`
	AdminTokenName string `mapstructure:"admin_token_name"`
	UserSecretKey  string `mapstructure:"user_secret_key"`
	UserTTL        int    `mapstructure:"user_ttl"`
	UserTokenName  string `mapstructure:"user_token_name"`
}

// OSSConfig Alibaba OSS配置
type OSSConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name"`
}

// WechatConfig 微信登录配置
type WechatConfig struct {
	AppID                 string `mapstructure:"app_id"`
	AppSecretKey          string `mapstructure:"app_secret_key"`
	MchID                 string `mapstructure:"mchid"`
	MchSerialNumber       string `mapstructure:"mch-serial-no"` // 商户证书序列号
	PrivateKeyFilePath    string `mapstructure:"private-key-file-path"`
	ApiV3Key              string `mapstructure:"api-v3-key"`
	WeChatPayCertFilePath string `mapstructure:"we-chat-pay-cert-file-path"`
	NotifyUrl             string `mapstructure:"notify-url"`
	RefundNotifyUrl       string `mapstructure:"refund-notify-url"`
}

// TemplateConfig xlsx模板文件
type TemplateConfig struct {
	Path string `mapstructure:"path"`
}

// ShopConfig 商店信息
type ShopConfig struct {
	Address string `mapstructure:"address"`
}

// BaiduConfig 百度地图配置
type BaiduConfig struct {
	AK string `mapstructure:"ak"`
}

// Config 全局配置实例
var Config GlobalConfig

// DB 全局数据库实例
var DB *gorm.DB

var Redis *redis.Client

// Logger 全局日志实例
var Logger *zap.Logger

// SugarLogger 全局Sugar日志实例
var SugarLogger *zap.SugaredLogger
