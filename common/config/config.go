package config

import (
	"bytes"
	"fmt"
	"github.com/shopspring/decimal"
	"os"
	"regexp"
	"strings"
	"takeout/common/database"
	"takeout/common/global"
	"takeout/common/logger"
	"takeout/common/redis"
	"takeout/internal/task"
	"time"

	"github.com/spf13/viper"
)

// 读取文件内容
func readFile(filePath string) []byte {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return []byte{}
	}
	return data
}

// 加载环境变量配置文件并返回map
func loadEnvConfig(envConfigPath string) (map[string]any, error) {
	if envConfigPath == "" {
		envConfigPath = "config-env.yaml"
	}

	// 检查环境变量配置文件是否存在
	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("env config file not found: %s", envConfigPath)
	}

	// 创建新的viper实例用于读取环境变量配置
	envViper := viper.New()
	envViper.SetConfigFile(envConfigPath)
	envViper.SetConfigType("yaml")

	// 读取环境变量配置文件
	if err := envViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read env config file: %w", err)
	}

	// 获取所有配置项
	envData := envViper.AllSettings()

	return envData, nil
}

// 合并配置文件并替换占位符
func mergeAndReplace(configFile, envFile string) ([]byte, error) {
	// 加载 config.yaml 模板
	configData := readFile(configFile)
	if len(configData) == 0 {
		return nil, fmt.Errorf("failed to read config file: %s", configFile)
	}

	// 加载环境变量配置
	envData, err := loadEnvConfig(envFile)
	if err != nil {
		return nil, err
	}

	// 替换占位符
	replaceConfig := replacePlaceholder(string(configData), envData)

	return []byte(replaceConfig), nil
}

// 替换 ${} 占位符的逻辑
func replacePlaceholder(content string, envMap map[string]any) string {
	reg := regexp.MustCompile(`\$\{([a-zA-Z0-9._]+)}`)
	// 返回替换后的文本
	return reg.ReplaceAllStringFunc(content, func(match string) string {
		// 找到要替换的词
		key := reg.FindStringSubmatch(match)[1]
		// 返回对应的内容
		return getValueFromMap(key, envMap)
	})
}

// 从 env.yaml 中获取嵌套字段的值
func getValueFromMap(key string, data map[string]any) string {
	keys := splitKey(key)
	value := data

	for _, k := range keys {
		if v, ok := value[k]; ok {
			// 存在嵌套
			if nestedMap, isMap := v.(map[string]any); isMap {
				value = nestedMap
			} else {
				return fmt.Sprintf("%v", v)
			}
		} else {
			return "${" + key + "}" // 没找到就原样返回
		}
	}

	return "${}" // 空值情况
}

// 分割嵌套路径的键
func splitKey(key string) []string {
	return strings.Split(key, ".")
}

// Init 初始化配置：直接对yaml文件进行替换占位符，viper直接从内存中加载配置文件
func Init(configPath string) error {
	// 设置默认配置文件路径
	if configPath == "" {
		configPath = "config.yaml" // 默认配置文件路径
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", configPath)
	}

	// 检查环境配置文件是否存在
	envConfigPath := "config-env.yaml"
	if _, err := os.Stat(envConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("env config file not found: %s", envConfigPath)
	}

	// 合并配置文件并替换占位符
	configData, err := mergeAndReplace(configPath, envConfigPath)
	if err != nil {
		return fmt.Errorf("failed to merge and replace config: %w", err)
	}

	// 从内存中读取替换好的配置
	viper.SetConfigType("yaml")
	if err = viper.ReadConfig(bytes.NewBuffer(configData)); err != nil {
		return fmt.Errorf("failed to read config from buffer: %w", err)
	}

	// 解析配置到结构体
	if err = viper.Unmarshal(&global.Config); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// 设置时域
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return fmt.Errorf("failed to initialize local: %w", err)
	}
	time.Local = loc

	// ☆ 设置全局变量，使 decimal.Decimal 序列化为 JSON 时不加引号
	decimal.MarshalJSONWithoutQuotes = true

	// 确保配置已经完全初始化后再初始化日志和数据库
	// 初始化日志
	if err = loggerInit(); err != nil {
		return fmt.Errorf("failed to initialize logger: %w", err)
	}

	// 初始化数据库
	if err = dbInit(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// 初始化缓存设置
	if err = redisInit(); err != nil {
		return fmt.Errorf("fail to initialize redis: %w", err)
	}

	// 初始化 Task
	if err = task.Init(); err != nil {
		return fmt.Errorf("fail to initialize task: %w", err)
	}

	return nil
}

// Close 关闭必要的连接
func Close() {
	closeDB()
	closeRedis()
}

// 初始化日志
func loggerInit() error {
	return logger.InitLogger()
}

// 初始化数据库
func dbInit() error {
	return database.InitDB()
}

func redisInit() error {
	return redis.InitRedis()
}

// 关闭数据库连接
func closeDB() {
	_ = database.Close()
}

func closeRedis() {
	_ = redis.Close()
}
