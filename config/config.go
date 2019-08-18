package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	_SERVICE_VERSION = "ServiceVersion"
	_ENV_CONF_DIR    = "GIN_CONF_DIR"
	_ENV_PSM         = "ENV_PSM"
)

var (
	webConfig Config
)

type Config struct {
	ConfigDir string
	YamlConfig
}

type YamlConfig struct {
	Mysql         *MySql     `yaml:"mysql"`
	PSM           string     `yaml:"psm"`
	Port          int        `yaml:"port"`
	Addr          string     `yaml:"addr"`
	LoggerConf    LoggerConf `yaml:"loggerconf"`
	AlgorithmHost string     `yaml:"algoritmhost"`
}

type LoggerConf struct {
	Level       int    `yaml:"level"`
	FileLog     bool   `yaml:"filelog"`
	LogDir      string `yaml:"logdir"`
	LogInterval string `yaml:"loginterval"`
	ConsoleLog  bool   `yaml:"consolelog"`
}

func GetAlgorithmHost() string {
	return webConfig.AlgorithmHost
}

func (c *Config) SetPSM(psm string) {
	c.PSM = psm
}

func GetAddr() string {
	return webConfig.Addr
}

func GetPort() int {
	return webConfig.Port
}

func Level() int {
	return webConfig.LoggerConf.Level
}
func PSM() string {
	return webConfig.PSM
}

func LogDir() string {
	return webConfig.LoggerConf.LogDir
}

func FileLog() bool {
	return webConfig.LoggerConf.FileLog
}

func LogInterval() string {
	return webConfig.LoggerConf.LogInterval
}

func ConsoleLog() bool {
	return webConfig.LoggerConf.ConsoleLog
}

func GetMysql() *MySql {
	return webConfig.Mysql
}

type MySql struct {
	DSN         string        `yaml:"dsn"`
	Active      int           `yaml:"active"`
	Idle        int           `yaml:"idle"`
	IdleTimeout time.Duration `yaml:"idle_timeout"`
}

func (mysql *MySql) GetDSN() string {
	return mysql.DSN
}

func (mysql *MySql) GetActive() int {
	return mysql.Active
}

func (mysql *MySql) GetIdle() int {
	return mysql.Idle
}

func (mysql *MySql) GetIdleTimeout() time.Duration {
	return mysql.IdleTimeout
}

func LoadConf() {
	parseFlags()
	parseConf()
	fmt.Fprintf(os.Stdout, "Web config: %#v\n", webConfig)
}

func parseFlags() {
	flag.StringVar(&webConfig.ConfigDir, "conf", "", "support config file.")
	flag.Parse()
	if webConfig.ConfigDir == "" {
		webConfig.ConfigDir = os.Getenv(_ENV_CONF_DIR)
	}
	if webConfig.ConfigDir == "" {
		fmt.Fprintf(os.Stderr, "Conf dir is not specified, use -conf option or %s environment\n", _ENV_CONF_DIR)
		usage()
	}
	psm := os.Getenv(_ENV_PSM)
	webConfig.SetPSM(psm)
}

func usage() {
	flag.Usage()
	os.Exit(-1)
}

func ConfigDir() string {
	return webConfig.ConfigDir
}

func parseConf() {
	v := viper.New()
	v.SetEnvPrefix("GIN")
	confFile := filepath.Join(ConfigDir(), strings.Replace(PSM(), ".", "_", -1)+".yaml")
	v.SetConfigFile(confFile)
	if err := v.ReadInConfig(); err != nil {
		msg := fmt.Sprintf("Failed to load web config: %s, %s", confFile, err)
		fmt.Fprintf(os.Stderr, "%s\n", msg)
		panic(msg)
	}
	yamlConfig := &webConfig.YamlConfig
	if err := v.Unmarshal(yamlConfig); err != nil {
		msg := fmt.Sprintf("Failed to unmarshal app config: %s", err)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		panic(msg)
	}
}
