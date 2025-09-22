package pkg

import (
	"context"
	"e-klinik/config"
	"log"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func NewOidcProvider(cfg *config.Config) *oidc.Provider {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	provider, err := oidc.NewProvider(ctx, cfg.Oidc.IssuerUrl)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	return provider

}

func NewOidcConfig(cfg *config.Config, provider *oidc.Provider) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cfg.Oidc.ClientId,
		ClientSecret: cfg.Oidc.ClientSecret,
		RedirectURL:  cfg.Oidc.RedirectUrl,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "openid"},
	}
}
