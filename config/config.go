package config

type Config struct {
	ApiKey    string
	AccessKey string
	DBName    string
	DBUser    string
	DBPass    string
}

func InitConfig() Config {
	return Config{
		ApiKey:    "$2a$10$iUWJxV84LRWQEkWZB/d/s.YMP1lCafcgp16S/7Q.dIR9BCP/Ahrvu",
		AccessKey: "$2a$10$hhjK7nJchphzgNzBRufFg.o9nIAQps07NVIomjXzP7Q28/hEVjflW",
		DBName:    "snaphub",
		DBUser:    "kranid",
		DBPass:    "lakrima54123",
	}
}
