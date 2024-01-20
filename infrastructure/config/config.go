package config

type AppConfig struct {
	MySQLCfg  MySQLConfig       `yaml:"mysql"`
	LogCfg    LogConfig         `yaml:"log"`
	ServerCfg ServerConfig      `yaml:"server"`
	TaskTime  ScheduledTaskTime `yaml:"time"`
	AuthCfg   AuthConfig        `yaml:"auth"`
	InfluxCfg InfluxConfig      `yaml:"influx"`
}

type MySQLConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type LogConfig struct {
	Level      string `yaml:"level"`
	SysLogPath string `yaml:"sys_log_path"`
	ErrLogPath string `yaml:"err_log_path"`
	Rotate     int    `yaml:"rotate"`
	MaxSize    int    `yaml:"max_size"`
	MaxBackups int    `yaml:"max_backups"`
	MaxAge     int    `yaml:"max_age_day"`
	Compress   bool   `yaml:"is_compress"`
	StdOut     bool   `yaml:"std_out"`
}

type ServerConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	GoroutineNum int    `yaml:"goroutine_num"`
}

type ScheduledTaskTime struct {
	QueryTime string `yaml:"query_time"`
}

type AuthConfig struct {
	SecretKey string `yaml:"secret_key"`
}

type InfluxConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
}
