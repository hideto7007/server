// config/env.go
package config

import (
	"fmt"
	"os"
)

type Env struct {
	Domain             string
	Secure             bool
	HttpOnly           bool
	RedirectURI        string
	PsqlUser           string
	PsqlTestUser       string
	PsqlDbname         string
	PsqlPassword       string
	PsqlHost           string
	PsqlPort           string
	PsqlSslModel       string
	ReactClient        string
	VueClient          string
	DockerClient       string
	SwaggerClient      string
	GoogleAccounts     string
	GoogleApis         string
	JwtSecret          string
	DomainURL          string
	RedisPassword      string
	RedisDomain        string
	RedisPort          string
	SmtpHost           string
	SmtpPort           string
	FromEmail          string
	EmailPassword      string
	GoogleClientID     string
	GoogleClientSecret string
}

var (
	GoogleSignInEnv Env
	GoogleSignUpEnv Env
	GlobalEnv       Env
)

func InitGoogleEnvs() {
	env := os.Getenv("ENV")

	GoogleSignInEnv = LeadEnv(env, "auth/google/signin/callback")
	GoogleSignUpEnv = LeadEnv(env, "auth/google/signup/callback")
	GlobalEnv = LeadEnv(env, "")
}

const ENV = "ENV"

func LeadEnv(env string, path string) Env {
	var secure bool = false
	var domain string = "localhost"
	var httpOnly bool = false
	var redirectURI string = fmt.Sprintf("http://localhost:8080/%s", path)

	// ローカル以外の場合
	if env != "local" {
		domain = env
		secure = true
		httpOnly = true
		RedirectURI = fmt.Sprintf("%s/%s", env, path)
	}

	EnvInfo := Env{
		Domain:             domain,
		Secure:             secure,
		HttpOnly:           httpOnly,
		RedirectURI:        redirectURI,
		PsqlUser:           os.Getenv("PSQL_USER"),
		PsqlTestUser:       os.Getenv("PSQL_TEST_USER"),
		PsqlDbname:         os.Getenv("PSQL_DBNAME"),
		PsqlPassword:       os.Getenv("PSQL_PASSWORD"),
		PsqlHost:           os.Getenv("PSQL_HOST"),
		PsqlPort:           os.Getenv("PSQL_PORT"),
		PsqlSslModel:       os.Getenv("PSQL_SSLMODEL"),
		ReactClient:        os.Getenv("REACT_CLIENT"),
		VueClient:          os.Getenv("VUE_CLIENT"),
		DockerClient:       os.Getenv("DOCKER_CLIENT"),
		SwaggerClient:      os.Getenv("SWAGGER_CLIENT"),
		GoogleAccounts:     os.Getenv("GOOGLE_ACCOUNTS_CLIENT"),
		GoogleApis:         os.Getenv("GOOGLEAPIS_CLIENT"),
		JwtSecret:          os.Getenv("JWT_SECRET"),
		DomainURL:          os.Getenv("DOMAIN"),
		RedisPassword:      os.Getenv("REDIS_PASSWORD"),
		RedisDomain:        os.Getenv("REDIS_DOMAIN"),
		RedisPort:          os.Getenv("REDIS_PORT"),
		SmtpHost:           os.Getenv("SMTP_HOST"),
		SmtpPort:           os.Getenv("SMTP_PORT"),
		FromEmail:          os.Getenv("FROMEMAIL"),
		EmailPassword:      os.Getenv("PASSWORD"),
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	}

	return EnvInfo
}
