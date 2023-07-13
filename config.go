package asgardian

type Config struct {
  Port int
  Host string
  Env string
  Key string
}

func CreateConfig() Config {
  return Config{
    Port: 8080,
    Host: "localhost",
    Env: "development",
    Key: "1234",
  }
}
