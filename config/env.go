// config/env.go
package config

import (
	"fmt"
	"os"
)

type Env struct {
	Protocol           string
	ClinetDomain       string
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
	RedirectPath       string
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
	LineClientID       string
	LineClientSecret   string
	OutPutLoggerFile   string
}

var (
	GoogleSignInEnv Env
	GoogleSignUpEnv Env
	GoogleDeleteEnv Env
	LineSignInEnv   Env
	LineSignUpEnv   Env
	LineDeleteEnv   Env
	GlobalEnv       Env
)

func InitGoogleEnvs() {
	env := os.Getenv("ENV")

	GoogleSignInEnv = LeadEnv(env, "auth/google/signin/callback")
	GoogleSignUpEnv = LeadEnv(env, "auth/google/signup/callback")
	GoogleDeleteEnv = LeadEnv(env, "auth/google/delete/callback")
	LineSignInEnv = LeadEnv(env, "auth/line/signin/callback")
	LineSignUpEnv = LeadEnv(env, "auth/line/signup/callback")
	LineDeleteEnv = LeadEnv(env, "auth/line/delete/callback")
	GlobalEnv = LeadEnv(env, "")
}

const ENV = "ENV"

func LeadEnv(env string, path string) Env {
	var protocol string = "http"
	var secure bool = false
	var domain string = "localhost"
	var clinetDomain string = "localhost:3000"
	var httpOnly bool = false
	var redirectURI string = fmt.Sprintf("%s://%s/%s", protocol, domain, path)

	// ローカル以外の場合
	if env != "local" {
		protocol = "https"
		domain = os.Getenv("DOMAIN")
		clinetDomain = domain
		secure = true
		httpOnly = true
		RedirectURI = fmt.Sprintf("%s://%s/%s", protocol, domain, path)
	}

	EnvInfo := Env{
		Protocol:           protocol,
		ClinetDomain:       clinetDomain,
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
		RedirectPath:       os.Getenv("REDIRECT_PATH"),
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
		LineClientID:       os.Getenv("LINE_CLIENT_ID"),
		LineClientSecret:   os.Getenv("LINE_CLIENT_SECRET"),
		OutPutLoggerFile:   os.Getenv("OUT_PUT_LOGGER_FILE"),
	}

	return EnvInfo
}
