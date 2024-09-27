package pkg

import (
	"context"
	"log"
	"sysbitBroker/config"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func NewOpenIDProvider(cfg *config.Config) *oidc.Provider {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	provider, err := oidc.NewProvider(ctx, cfg.OpenID.IssuerUrl)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	return provider

}

func NewOpenIDConfig(cfg *config.Config, provider *oidc.Provider) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.OpenID.ClientId,
		ClientSecret: cfg.OpenID.ClientSecret,
		RedirectURL:  cfg.OpenID.RedirectUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "openid"},
	}
}
