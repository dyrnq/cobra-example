package cmd


type Config struct {
	Server struct {
		Address string `mapstructure:"address"`
		Port int `mapstructure:"port"`
	} `mapstructure:"server"`
}