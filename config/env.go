// config/env.go
package config

import (
	"fmt"
	"os"
)

type Env struct {
	Domain      string
	Secure      bool
	HttpOnly    bool
	RedirectURI string
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
		Domain:      domain,
		Secure:      secure,
		HttpOnly:    httpOnly,
		RedirectURI: redirectURI,
	}

	return EnvInfo
}
