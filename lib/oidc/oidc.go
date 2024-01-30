package oidc

import (
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3"
	"github.com/rokoucha/go-oidc-idp-example/lib/keychain"
	"github.com/rokoucha/go-oidc-idp-example/lib/user"
)

const (
	SigningKeyKid = "oidc-signing-key"
)

type Client struct {
	Id          string
	Name        string
	RedirectUri string
}

type Oidc struct {
	baseUrl  string
	clients  []Client
	keychain *keychain.Keychain
	signer   jose.Signer
}

type Config struct {
	BaseUrl  string
	Clients  []Client
	Keychain *keychain.Keychain
}

func New(config Config) (*Oidc, error) {
	signingKey, err := config.Keychain.Create(SigningKeyKid)
	if err != nil {
		return nil, err
	}

	signer, err := jose.NewSigner(jose.SigningKey{
		Algorithm: jose.RS256,
		Key:       signingKey,
	}, nil)
	if err != nil {
		return nil, err
	}

	return &Oidc{
		baseUrl:  config.BaseUrl,
		clients:  config.Clients,
		keychain: config.Keychain,
		signer:   signer,
	}, nil
}

type AuthenticationRequest struct {
	Scope        string
	ResponseType string
	ClientID     string
	RedirectUri  string
	Nonce        string
}

type IDTokenPayload struct {
	Issuer     string `json:"iss"`
	Subject    string `json:"sub"`
	Audience   string `json:"aud"`
	Expiration int64  `json:"exp"`
	IssuedAt   int64  `json:"iat"`
	Nonce      string `json:"nonce"`
	Name       string `json:"name"`
}

type OpenIDProviderMetadata struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	JWKsUri                          string   `json:"jwks_uri"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	SubjectTypesSupported            []string `json:"subject_types_supported"`
	IdTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
}

func (o *Oidc) GetOpenIDProviderMetadata() OpenIDProviderMetadata {
	return OpenIDProviderMetadata{
		Issuer:                o.baseUrl,
		AuthorizationEndpoint: o.baseUrl + "/oidc/auth",
		JWKsUri:               o.baseUrl + "/oidc/jwks",
		ResponseTypesSupported: []string{
			"id_token",
		},
		SubjectTypesSupported: []string{
			"public",
		},
		IdTokenSigningAlgValuesSupported: []string{
			string(jose.RS256),
		},
	}
}

func (o *Oidc) getClient(clientId string) (Client, bool) {
	idx := slices.IndexFunc(o.clients, func(c Client) bool {
		return c.Id == clientId
	})

	if idx == -1 {
		return Client{}, false
	}

	return o.clients[idx], true
}

func (o *Oidc) GetPublicKeys() []jose.JSONWebKey {
	jwks := o.keychain.GetAll()
	publicKeys := make([]jose.JSONWebKey, len(jwks))
	for i, jwk := range jwks {
		publicKeys[i] = jwk.Public()
	}
	return publicKeys
}

func (o *Oidc) ValidateAuthenticationRequest(req AuthenticationRequest) error {
	if !strings.Contains(req.Scope, "openid") {
		return errors.New("invalid_scope")
	}

	if req.ResponseType != "id_token" {
		return errors.New("unsupported_response_type")
	}

	_, ok := o.getClient(req.ClientID)
	if !ok {
		return errors.New("access_denied")
	}

	return nil
}

func (o *Oidc) GenerateIDToken(user user.UserInfo, clientID string, nonce string) (string, error) {
	payload, err := json.Marshal(IDTokenPayload{
		Issuer:     o.baseUrl,
		Subject:    user.ID.String(),
		Audience:   clientID,
		Expiration: time.Now().Add(time.Hour * 24).Unix(),
		IssuedAt:   time.Now().Unix(),
		Nonce:      nonce,
		Name:       user.Username,
	})
	if err != nil {
		return "", err
	}

	jws, err := o.signer.Sign(payload)
	if err != nil {
		return "", err
	}

	return jws.CompactSerialize()
}
