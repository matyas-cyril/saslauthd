package main

// https://www.rfc-editor.org/rfc/rfc7519#section-4.1
type jwtStruct struct {
	Iss string `json:"iss"` // Emetteur
	Sub string `json:"sub"` // Sujet
	Aud string `json:"aud"` // Audience
	Exp uint64 `json:"exp"` // Expiration (epoch)
	Nbf uint64 `json:"nbf"` // Pas avant (epoch)
	Iat uint64 `json:"iat"` // Date de délivrance (epoch)
	Jti string `json:"jti"` // Identifiant unique du token
	Usr string `json:"usr"` // Utilisateur
	Dom string `json:"dom"` // Domaine
	Uid string `json:"uid"` // Identifiant utilisateur. Si défini alors on n'utilise pas Usr et Dom
}

type jwtCredent struct {
	Aud     []string
	Pwd     []byte
	VirtDom bool // Utilise-t-on les domaines pour l'authentification si disponible user@dom
}
