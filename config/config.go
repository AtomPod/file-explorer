package config

import (
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

//Config 配置
type Config struct {
	Mode          string              `json:"mode" yaml:"mode" mapstructure:"mode"`
	ServerName    string              `json:"server_name" yaml:"server_name" mapstructure:"server_name"`
	HTTP          HTTPConfig          `json:"http" yaml:"http" mapstructure:"http"`
	Database      DatabaseConfig      `json:"database" yaml:"database" mapstructure:"database"`
	UserService   UserServiceConfig   `json:"user_service" yaml:"user_service" mapstructure:"user_service"`
	FileService   FileServiceConfig   `json:"file_service" yaml:"file_service" mapstructure:"file_service"`
	Prometheus    PromConfig          `json:"prometheus" yaml:"prometheus"  mapstructure:"prometheus"`
	Cache         CacheConfig         `json:"cache" yaml:"cache" mapstructure:"cache"`
	Email         EMailConfig         `json:"email" yaml:"email" mapstructure:"email"`
	EmailTemplate EmailTemplateConfig `json:"email_template" yaml:"email_template" mapstructure:"email_template"`
	Log           LogConfig           `json:"log" yaml:"log" mapstructure:"log"`
	Raw           *viper.Viper        `json:"-" yaml:"-" mapstructure:"-"`
}

//LogConfig 日志配置
type LogConfig struct {
	Level            string   `json:"level" yaml:"level"`
	Encoding         string   `json:"encoding" yaml:"encoding"`
	OutputPaths      []string `json:"output_paths" yaml:"output_paths" mapstructure:"output_paths"`
	ErrorOutputPaths []string `json:"error_output_paths" yaml:"error_output_paths" mapstructure:"error_output_paths"`
}

//HTTPConfig http配置
type HTTPConfig struct {
	Host     string `json:"host" yaml:"host" mapstructure:"host"`
	Port     string `json:"port" yaml:"port" mapstructure:"port"`
	SSL      bool   `json:"ssl" yaml:"ssl" mapstructure:"ssl"`
	CertPath string `json:"cert_path" yaml:"cert_path" mapstructure:"cert_path"`
	KeyPath  string `json:"key_path" yaml:"key_path" mapstructure:"key_path"`
	AutoTLS  struct {
		Enable   bool     `json:"enable" yaml:"enable" mapstructure:"enable"`
		CacheDir string   `json:"cache_dir" yaml:"cache_dir" mapstructure:"cache_dir"`
		Host     []string `json:"hosts" yaml:"hosts" mapstructure:"hosts"`
	} `json:"auto_tls" yaml:"auto_tls" mapstructure:"auto_tls"`
}

//DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Engine   string `json:"engine" yaml:"engine" mapstructure:"engine"`
	Host     string `json:"host" yaml:"host" mapstructure:"host"`
	Port     int    `json:"port" yaml:"port" mapstructure:"port"`
	DBName   string `json:"dbname" yaml:"dbname" mapstructure:"dbname"`
	User     string `json:"user" yaml:"user" mapstructure:"user"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
	SSLMode  string `json:"sslmode" yaml:"sslmode" mapstructure:"sslmode"`
}

//VerificationCodeConfig 验证码配置信息
type VerificationCodeConfig struct {
	TypeName      string        `json:"type_name" yaml:"type_name" mapstructure:"type_map"`
	Expiration    time.Duration `json:"expiration" yaml:"expiration" mapstructure:"expiration"`
	MailTemplName string        `json:"mail_templ_name" yaml:"mail_templ_name" mapstructure:"mail_templ_name"`
}

//UserServiceConfig 用户服务配置
type UserServiceConfig struct {
	JWT struct {
		SigningKey string        `json:"signing_key" yaml:"signing_key" mapstructure:"signing_key"`
		Expire     time.Duration `json:"expire" yaml:"expire" mapstructure:"expire"`
	} `json:"jwt" yaml:"jwt" mapstructure:"jwt"`

	VerificationCode struct {
		Email         VerificationCodeConfig `json:"email" yaml:"email" mapstructure:"email"`
		ResetPassword VerificationCodeConfig `json:"reset_password" yaml:"reset_password" mapstructure:"reset_password"`
	} `json:"verification_code" yaml:"verification_code" mapstructure:"verification_code"`
}

//FileServiceConfig 文件服务配置
type FileServiceConfig struct {
	BasePath     string `json:"basepath" yaml:"basepath" mapstructure:"basepath"`
	absolutePath string
}

//CacheConfig 缓存配置
type CacheConfig struct {
	Engine       string   `json:"engine" yaml:"engine" mapstructure:"engine"`
	Locations    []string `json:"locations" yaml:"locations" mapstructure:"locations"`
	KeyPrefix    string   `json:"key_prefix" yaml:"key_prefix" mapstructure:"key_prefix"`
	MaxCacheSize int64    `json:"max_cache_size" yaml:"max_cache_size" mapstructure:"max_cache_size"`
}

//EMailConfig email配置
type EMailConfig struct {
	Host         string `json:"host" yaml:"host" mapstructure:"host"`
	Port         int    `json:"port" yaml:"port" mapstructure:"port"`
	Username     string `json:"username" yaml:"username" mapstructure:"username"`
	Password     string `json:"password" yaml:"password" mapstructure:"password"`
	SSL          bool   `json:"ssl" yaml:"ssl" mapstructure:"ssl"`
	CertPath     string `json:"cert_path" yaml:"cert_path" mapstructure:"cert_path"`
	KeyPath      string `json:"key_path" yaml:"key_path" mapstructure:"key_path"`
	LocalName    string `json:"local_name" yaml:"local_name" mapstructure:"local_name"`
	MaxQueueSize int64  `json:"max_queue_size" yaml:"max_queue_size" mapstructure:"max_queue_size"`
}

//EmailTemplateItem Email模板项信息
type EmailTemplateItem struct {
	Subject     string `json:"subject" yaml:"subject" mapstructure:"subject"`
	ContentType string `json:"content_type" yaml:"content_type" mapstructure:"content_type"`
	Alias       string `json:"alias" yaml:"alias" mapstructure:"alias"`
	Filename    string `json:"filename" yaml:"filename" mapstructure:"filename"`
}

//EmailTemplateConfig Email模板配置
type EmailTemplateConfig struct {
	BasePath string              `json:"base_path" yaml:"base_path" mapstructure:"base_path"`
	Items    []EmailTemplateItem `json:"items" yaml:"items" mapstructure:"items"`
}

//PromConfig prometheus配置
type PromConfig struct {
	Enable      bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
	GatewayAddr string `json:"gateway_addr" yaml:"gateway_addr" mapstructure:"gateway_addr"`
	IntervalSec int    `json:"interval_sec" yaml:"interval_sec" mapstructure:"interval_sec"`
}

//FileAbsolutePath 获取绝对路径
func (fs *FileServiceConfig) FileAbsolutePath() string {
	return fs.absolutePath
}

//Load 加载配置信息
func Load(configPaths []string, confName, confType string) *Config {
	v := viper.New()
	v.SetConfigType(confType)
	v.SetConfigName(confName)
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	conf := &Config{}

	return initializeConfig(v, conf)
}

func initializeConfig(v *viper.Viper, conf *Config) *Config {
	if err := v.Unmarshal(conf); err != nil {
		panic(err)
	}

	wd, _ := os.Getwd()

	if conf.FileService.BasePath == "" {
		conf.FileService.BasePath = "files"
	}

	conf.FileService.absolutePath = filepath.Join(wd, conf.FileService.BasePath)

	if conf.UserService.JWT.Expire == time.Duration(0) {
		conf.UserService.JWT.Expire = time.Duration(2) * time.Hour
	}
	if conf.UserService.JWT.SigningKey == "" {
		conf.UserService.JWT.SigningKey = "golang-service"
	}

	verificationCodeConf := &conf.UserService.VerificationCode
	if verificationCodeConf.Email.TypeName == "" {
		verificationCodeConf.Email.TypeName = "register"
	}
	if verificationCodeConf.Email.Expiration == time.Duration(0) {
		verificationCodeConf.Email.Expiration = 2 * time.Hour
	}
	if verificationCodeConf.ResetPassword.TypeName == "" {
		verificationCodeConf.ResetPassword.TypeName = "reset_password"
	}
	if verificationCodeConf.ResetPassword.Expiration == time.Duration(0) {
		verificationCodeConf.ResetPassword.Expiration = 2 * time.Hour
	}

	if conf.Email.MaxQueueSize == 0 {
		conf.Email.MaxQueueSize = 1024
	}

	if conf.Prometheus.Enable {
		if conf.Prometheus.GatewayAddr == "" {
			conf.Prometheus.GatewayAddr = "http://localhost:9091"
		}

		if conf.Prometheus.IntervalSec == 0 {
			conf.Prometheus.IntervalSec = 15
		}
	}

	if conf.HTTP.Port == "" {
		conf.HTTP.Port = "8080"
	}

	if conf.HTTP.AutoTLS.Enable {
		if conf.HTTP.AutoTLS.CacheDir == "" {
			conf.HTTP.AutoTLS.CacheDir = "./tls_cache"
		}
	}

	conf.Raw = v
	return conf
}
