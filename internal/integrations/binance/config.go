package binance

type BinanceConfig struct {
	ApiKey    string   `yaml:"apiKey"`
	ApiSecret string   `yaml:"apiSecret"`
	Basepaths []string `yaml:"basepaths"`
}
