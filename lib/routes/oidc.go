package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-jose/go-jose/v3"
	"github.com/rokoucha/go-oidc-idp-example/lib/oidc"
)

func (r *Routes) WellKnownOpenIdConfiguration(res http.ResponseWriter, req *http.Request) {
	json, err := json.Marshal(r.oidc.GetOpenIDProviderMetadata())
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(json)
	return
}

func (r *Routes) OidcAuth(res http.ResponseWriter, req *http.Request) {
	var authReq oidc.AuthenticationRequest
	switch req.Method {
	case "GET":
		authReq = oidc.AuthenticationRequest{
			Scope:        req.URL.Query().Get("scope"),
			ResponseType: req.URL.Query().Get("response_type"),
			ClientID:     req.URL.Query().Get("client_id"),
			RedirectUri:  req.URL.Query().Get("redirect_uri"),
			Nonce:        req.URL.Query().Get("nonce"),
		}
	case "POST":
		err := req.ParseForm()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		authReq = oidc.AuthenticationRequest{
			Scope:        req.Form.Get("scope"),
			ResponseType: req.Form.Get("response_type"),
			ClientID:     req.Form.Get("client_id"),
			RedirectUri:  req.Form.Get("redirect_uri"),
			Nonce:        req.Form.Get("nonce"),
		}
	default:
		http.Redirect(res, req, fmt.Sprintf("%s?error=invalid_request&error_description=method not allowed", authReq.RedirectUri), http.StatusFound)
		return
	}

	fmt.Printf("%#v\n", authReq)

	err := r.oidc.ValidateAuthenticationRequest(authReq)
	if err != nil {
		http.Redirect(res, req, fmt.Sprintf("%s?error=%s", authReq.RedirectUri, err.Error()), http.StatusFound)
		return
	}

	user, err := r.getUserFromSession(req)
	if err != nil {
		http.Redirect(res, req, fmt.Sprintf("%s?error=login_required", authReq.RedirectUri), http.StatusFound)
		return
	}

	idToken, err := r.oidc.GenerateIDToken(user, authReq.ClientID, authReq.Nonce)
	if err != nil {
		http.Redirect(res, req, fmt.Sprintf("%s?error=server_error", authReq.RedirectUri), http.StatusFound)
		return
	}

	http.Redirect(res, req, fmt.Sprintf("%s?id_token=%s", authReq.RedirectUri, idToken), http.StatusFound)
	return
}

func (r *Routes) OidcJwks(res http.ResponseWriter, req *http.Request) {
	keys := r.oidc.GetPublicKeys()
	jwks := struct {
		Keys []jose.JSONWebKey `json:"keys"`
	}{Keys: keys}

	json, err := json.Marshal(jwks)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(json)
	return
}
