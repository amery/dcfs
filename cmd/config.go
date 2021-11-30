package cmd

type Config struct{}

func (c *Config) SetDefaults() {}

func (c *Config) ReadInFile(filename string) error {
	return nil
}

func NewConfig() *Config {
	c := &Config{}
	c.SetDefaults()
	return c
}
