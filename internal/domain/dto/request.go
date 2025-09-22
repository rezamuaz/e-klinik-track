package dto

import (
	"e-klinik/internal/domain/entity"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreatedArtistParams struct {
	Name        string             `json:"name"`
	KanjiName   string             `json:"kanji_name"`
	Thumbnails  []entity.Thumbnail `json:"thumbnails"`
	Description pgtype.Text        `json:"description"`
}

type CreatedGenreParams struct {
	Title string `json:"title"`
}

type FindArtistsParams struct {
	Name  string `json:"name"`
	Limit int32  `json:"limit"`
	Page  int32  `json:"page"`
}

type FindTermParams struct {
	Taxonomy string `json:"taxonomy"`
	Initial  string `json:"initial"`
}

type CreateTermParams struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type FindPostsParams struct {
	Name  string `json:"name"`
	Limit int32  `json:"limit"`
	Page  int32  `json:"page"`
}
type CreateThirdParty struct {
	ConnectionUriDomain   string   `json:"connection_uri_domain"`
	AppID                 string   `json:"app_id"`
	TenantID              string   `json:"tenant_id"`
	ThirdPartyID          string   `json:"third_party_id"`
	Name                  *string  `json:"name"`
	OidcDiscoveryEndpoint *string  `json:"oidc_discovery_endpoint"`
	RequireEmail          *bool    `json:"require_email"`
	ClientType            string   `json:"client_type"`
	ClientID              string   `json:"client_id"`
	ClientSecret          *string  `json:"client_secret"`
	Scope                 []string `json:"scope"`
}

type CreateApp struct {
	AppID                string `json:"app_id"`
	TenantID             string `json:"tenant_id"`
	AppDisplayName       string `json:"app_display_name"`
	ConnectionUriDomain  string `json:"connection_uri_domain"`
	EmailPasswordEnabled *bool  `json:"email_password_enabled"`
	PasswordlessEnabled  *bool  `json:"passwordless_enabled"`
	ThirdPartyEnabled    *bool  `json:"third_party_enabled"`
}
