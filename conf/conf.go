package conf

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Server struct {
	System System `mapstructure:"system" json:"system" yaml:"system"`
	Logger Logger `mapstructure:"logger" json:"logger" yaml:"logger"`
	Jwt    Jwt    `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Mysql  Mysql  `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis  Redis  `mapstructure:"redis" json:"redis" yaml:"redis"`
	Gpt    Gpt    `mapstructure:"gpt" json:"gpt" yaml:"gpt"`
	File   File   `mapstructure:"file" json:"file" yaml:"file"`
}

// System 系统配置
type System struct {
	Env      string   `mapstructure:"env" json:"env" yaml:"env"`                // 环境
	Addr     string   `mapstructure:"addr" json:"addr" yaml:"addr"`             // 系统服务监听端口
	Name     string   `mapstructure:"name" json:"name" yaml:"name"`             // 系统服务名称
	Version  string   `mapstructure:"version" json:"version" yaml:"version"`    // 系统版本
	Http     Http     `mapstructure:"http" json:"http" yaml:"http"`             // HTTP 配置
	Security Security `mapstructure:"security" json:"security" yaml:"security"` // 安全配置
}

// Http HTTP配置
type Http struct {
	ReadTimeout  time.Duration `mapstructure:"readTimeout" json:"readTimeout" yaml:"readTimeout"`    // HTTP读取超时时间
	WriteTimeout time.Duration `mapstructure:"writeTimeout" json:"writeTimeout" yaml:"writeTimeout"` // HTTP写入超时时间
	IdleTimeout  time.Duration `mapstructure:"idleTimeout" json:"idleTimeout" yaml:"idleTimeout"`    // HTTP空闲超时时间
}

// Security 安全配置
type Security struct {
	Cors Cors `mapstructure:"cors" json:"cors" yaml:"cors"` // CORS 跨域配置
}

// Cors 跨域配置
type Cors struct {
	Enabled      bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`                // 是否启用跨域支持
	AllowOrigins string `mapstructure:"allowOrigins" json:"allowOrigins" yaml:"allowOrigins"` // 允许跨域的来源
	AllowMethods string `mapstructure:"allowMethods" json:"allowMethods" yaml:"allowMethods"` // 允许的HTTP方法
}

// Logger 用于配置日志
type Logger struct {
	LogLevel      string `mapstructure:"logLevel" json:"logLevel" yaml:"logLevel"`                // 日志级别（debug、info、warn、error、fatal、panic）
	LogFormat     string `mapstructure:"logFormat" json:"logFormat" yaml:"logFormat"`             // 日志格式（json、text、table）
	JSONFormatter bool   `mapstructure:"jsonFormatter" json:"jsonFormatter" yaml:"jsonFormatter"` // 是否格式化json
	Output        string `mapstructure:"output" json:"output" yaml:"output"`                      // 输出方式（console、file）
	LogPath       string `mapstructure:"logPath" json:"logPath" yaml:"logPath"`                   // 日志文件路径（当 Output 为 file 时有效）
	MaxSize       int    `mapstructure:"maxSize" json:"maxSize" yaml:"maxSize"`                   // 日志文件最大大小（MB）（当 Output 为 file 时有效）
	MaxBackups    int    `mapstructure:"maxBackups" json:"maxBackups" yaml:"maxBackups"`          // 最大保留日志文件数（当 Output 为 file 时有效）
	MaxAge        int    `mapstructure:"maxAge" json:"maxAge" yaml:"maxAge"`                      // 日志文件保留天数（当 Output 为 file 时有效）
}

// Mysql 数据库配置结构体
type Mysql struct {
	Host          string `mapstructure:"host" json:"host" yaml:"host"`                            // 数据库主机地址
	Port          int    `mapstructure:"port" json:"port" yaml:"port"`                            // 数据库端口
	User          string `mapstructure:"user" json:"user" yaml:"user"`                            // 数据库用户名
	Password      string `mapstructure:"password" json:"password" yaml:"password"`                // 数据库密码
	Db            string `mapstructure:"db" json:"db" yaml:"db"`                                  // 数据库名称
	Enabled       bool   `mapstructure:"enabled" json:"enabled" yaml:"enabled"`                   // 是否启用日志输出
	Level         string `mapstructure:"level" json:"level" yaml:"level"`                         // 日志级别
	SlowThreshold int    `mapstructure:"slowThreshold" json:"slowThreshold" yaml:"slowThreshold"` // 慢查询阈值
}

// Redis 缓存配置
type Redis struct {
	Host           string        `mapstructure:"host" json:"host" yaml:"host"`                               // Redis服务的IP地址
	Port           int           `mapstructure:"port" json:"port" yaml:"port"`                               // Redis服务的端口号
	Password       string        `mapstructure:"password" json:"password" yaml:"password"`                   // Redis认证的密码，如无密码则为空
	Db             int           `mapstructure:"db" json:"db" yaml:"db"`                                     // Redis数据库序号（默认0）
	MaxIdle        int           `mapstructure:"maxIdle" json:"maxIdle" yaml:"maxIdle"`                      // 最大空闲连接数，用于控制资源
	MaxActive      int           `mapstructure:"maxActive" json:"maxActive" yaml:"maxActive"`                // 最大活跃连接数（0 表示无限制）
	IdleTimeout    time.Duration `mapstructure:"idleTimeout" json:"idleTimeout" yaml:"idleTimeout"`          // 空闲连接超时时间，超时将关闭空闲连接
	ConnectTimeout time.Duration `mapstructure:"connectTimeout" json:"connectTimeout" yaml:"connectTimeout"` // Redis连接超时时间
	ReadTimeout    time.Duration `mapstructure:"readTimeout" json:"readTimeout" yaml:"readTimeout"`          // Redis读操作的超时时间
	WriteTimeout   time.Duration `mapstructure:"writeTimeout" json:"writeTimeout" yaml:"writeTimeout"`       // Redis写操作的超时时间
	PoolSize       int           `mapstructure:"poolSize" json:"poolSize" yaml:"poolSize"`                   // 连接池大小（默认为10*CPU核心数）
	LogEnabled     bool          `mapstructure:"logEnabled" json:"logEnabled" yaml:"logEnabled"`             // 是否记录Redis操作日志，默认true
}

// Jwt 鉴权
type Jwt struct {
	AccessSecret string `mapstructure:"accessSecret" json:"accessSecret" yaml:"accessSecret"`
	AccessExpire int64  `mapstructure:"accessExpire" json:"accessExpire" yaml:"accessExpire"`
}

// Gpt ai
type Gpt struct {
	ApiKey   string `mapstructure:"ApiKey" json:"apiKey" yaml:"ApiKey"`
	Endpoint string `mapstructure:"Endpoint" json:"endpoint" yaml:"Endpoint"`
	Model    string `mapstructure:"Model" json:"model" yaml:"Model"`
}

// File 文件配置
type File struct {
	TaskPath string `mapstructure:"taskPath" json:"taskPath" yaml:"taskPath"`
}

var GlobalConf *Server

// LoadConf 加载配置文件
func LoadConf(env string) {
	var server Server
	var confPath string

	if env == "pro" {
		confPath = "conf/conf.prod.yaml"
	} else if env == "test" {
		confPath = "conf/conf.test.yaml"
	} else {
		confPath = "conf/conf.dev.yaml"
	}
	vi := viper.New()

	vi.SetConfigFile(confPath)

	err := vi.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: #{err}"))
	}
	if err := vi.Unmarshal(&server); err != nil {
		fmt.Println(err)
	}
	server.System.Env = env
	// return &server
	GlobalConf = &server
}
