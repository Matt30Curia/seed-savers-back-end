package auth

import (
	"backend/seed-savers/config"
	"fmt"
	

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

var opts = SessionOptions{
	CookiesKey: "dffhdjfdsjfsdh",
	MaxAge:     86400 * 30,
	HttpOnly:   true,
	Secure:     false,
}

func NewOauth() *AuthStore{

	authStore := NewCookieStore(opts)

	gothic.Store = authStore.Store
	goth.UseProviders(
		google.New(config.Envs.GoogleClientID, config.Envs.GoogleClientSecretId,buildCallbackURL("google")))

	return authStore

}

func buildCallbackURL(provider string) string {
	return fmt.Sprintf("%s:%s/auth/%s/callback", config.Envs.PublicHost, config.Envs.Port, provider)
}
