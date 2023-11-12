package config

import (
	"fmt"
	"myproject/internal/apperrors"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	ErrInternal = apperrors.New(apperrors.ErrInternal, "App-Config-")
)

type Config struct {
	Redis      *Redis      `yaml:"redis"`
	Psql       *Postgres   `yaml:"postgres"`
	HttpServer *HttpServer `yaml:"server"`
	Kafka      *Kafka      `yaml:"kafka"`
	Prometheus *Prometheus `yaml:"prometheus"`
}

func New() (*Config, error) {
	f, err := os.Open("./config.yml")
	if err != nil {
		ErrInternal.AddLocation("New-os.Open")
		ErrInternal.SetErr(err)
		return &Config{}, ErrInternal
	}
	c := &Config{}

	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		ErrInternal.AddLocation("New-DecodeYaml")
		ErrInternal.SetErr(err)
		return &Config{}, ErrInternal
	}
	return c, nil
}

type Redis struct {
	Port string
	Host string
}

func (r *Redis) HostPort() string {
	return fmt.Sprintf("%s:%s", r.Host, r.Port)
}

type Postgres struct {
	DbDriver      string
	Host          string
	Port          string
	UserName      string
	Password      string
	Schema        string
	SslMode       string
	MigrationPath string
}

func (p *Postgres) ConnString() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s", p.DbDriver, p.UserName, p.Password, p.Host, p.Port, p.Schema, p.SslMode)
}

type HttpServer struct {
	Host         string
	Port         string
	writeTimeout time.Duration `yaml:"writetimeout"`
	readTimeout  time.Duration `yaml:"readtimeout"`
}

func (s *HttpServer) WriteTimeout() time.Duration {
	return s.writeTimeout * time.Second
}
func (s *HttpServer) ReadTimeout() time.Duration {
	return s.readTimeout * time.Second
}
func (s *HttpServer) HostPort() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

type Kafka struct {
	Protocol string
	Host     string
	Port     string
	Topic    string
	GroupId  string
}

func (k *Kafka) HostPort() string {
	return fmt.Sprintf("%s:%s", k.Host, k.Port)
}

type Prometheus struct {
	Host string
	Port string
}

func (p *Prometheus) Addr() string {
	return fmt.Sprintf("%s:%s", p.Host, p.Port)
}
