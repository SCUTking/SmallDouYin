package config

// GinConfig 瀹氫箟 Gin 閰嶇疆鏂囦欢鐨勭粨鏋勪綋
type GinConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// MySQLConfig 瀹氫箟 mysql 閰嶇疆鏂囦欢缁撴瀯浣�
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Username     string `mapstructure:"username"`
	Password     string `mapstructure:"password"`
	DBname       string `mapstructure:"db_name"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// RedisConfig 瀹氫箟 redis 閰嶇疆鏂囦欢缁撴瀯浣�
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// JWTConfig 瀹氫箟 jwt 閰嶇疆鏂囦欢缁撴瀯浣�
type JWTConfig struct {
	SigningKey string `mapstructure:"signing_key"`
}

// System 瀹氫箟椤圭洰閰嶇疆鏂囦欢缁撴瀯浣�
type System struct {
	GinConfig   *GinConfig   `mapstructure:"gin"`
	MySQLConfig *MySQLConfig `mapstructure:"mysql"`
	RedisConfig *RedisConfig `mapstructure:"redis"`
	JWTConfig   *JWTConfig   `mapstructure:"jwt"`
}
