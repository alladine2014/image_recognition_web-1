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
	PSM       string
	ConfigDir string
	Port      int
	Addr      string
	YamlConfig
}

type YamlConfig struct {
	Mysql *MySql `yaml:"mysql"`
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
	flag.StringVar(&webConfig.Addr, "addr", "", "service addr.")
	flag.StringVar(&webConfig.ConfigDir, "conf", "", "support config file.")
	flag.IntVar(&webConfig.Port, "port", 0, "service port.")
	flag.StringVar(&webConfig.PSM, "psm", "", "service port.")
	flag.Parse()
	if webConfig.ConfigDir == "" {
		webConfig.ConfigDir = os.Getenv(_ENV_CONF_DIR)
	}
	if webConfig.ConfigDir == "" {
		fmt.Fprintf(os.Stderr, "Conf dir is not specified, use -conf option or %s environment\n", _ENV_CONF_DIR)
		usage()
	}
	if webConfig.PSM == "" {
		webConfig.PSM = os.Getenv(_ENV_PSM)
	}
	if webConfig.PSM == "" {
		fmt.Fprintf(os.Stderr, "PSM is not specified use -psm option or %s environment\n", _ENV_PSM)
	}
}

func usage() {
	flag.Usage()
	os.Exit(-1)
}

func ConfigDir() string {
	return webConfig.ConfigDir
}

func PSM() string {
	return webConfig.PSM
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
