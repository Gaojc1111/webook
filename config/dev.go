//go:build !k8s

package config

var Config = config{
	DB: DBConfig{
		DSN: "root:root@tcp(localhost:3306)/webook",
	},
	Redis: RedisConfig{
		Addr: "192.168.136.135:6379",
	},
}
