package jwttoken

import (
	"fmt"
	"kuncenduit-backend/shared/createerror"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const jwtIssuer = "kuncenduit-backend"    //for jwt registered claim: iss
const jwtAudience = "kuncenduit-frontend" //for jwt registered claim: aud
//to be set dynamically: exp (expiration time): Time after which the JWT expires
//to be set dynamically: nbf (not before time): Time before which the JWT must not be accepted for processing
//to be set dynamically: iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
//to be set dynamically: jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
//unset because already user data will be set in a custom claim: sub (subject): Subject of the JWT (the user id)
const KDCustomClaimName = "kdcc"

func getPrivateKey() (jwk.Key, error) {
	//generate jwk using https://mkjwk.org/
	//  use oct, key size 2048, key use signature, algorithm HS512,
	//  key id is your custom secret string (32 bytes = 2 * uuid_v4), show X.509 Yes
	//then minify the public and private keypair and put it below
	//  https://www.webtoolkitonline.com/json-minifier.html
	jsonRSAPrivateKey := []byte(`{"kty":"oct","use":"sig","kid":"8bbe07de-1867-42f1-8b2b-c0833522c375-2ba6e9a7-d10c-4bcf-b8a2-5cda28afdc00","k":"bIGBLO8z6tjdcyqb-tVZud-n5ZoKWZDqsO3MhMzXI7F9_cYM7NHsoFL90s_dguIaASt5M4DtSItjtuab_dv5evJxoggdtS1Su_a1_g5FkPnjsbvFzOJr4ixEqUv-aGJzo8GVxuecjngKOQX422sdvnTjVDBchEHWooceaqGNtl3sLcAI4w_EQBIYAQNVKDoeZEZZvE3ppT5yGHKAaBQLSzP2-Cqo0xBTw_Sc3FtrFFXBZy7DHahF8New-sKt80ICilPr6R1WGsBCLRAevguAbLLD2YV9YaCDjoIMXhs7A7Q2JZ2y7fRqikKlhaGPJtL9DW-Z2CWvEsdqJX-NHCg_ag","alg":"HS512"}`)
	keyForSigningJwt, err := jwk.ParseKey(jsonRSAPrivateKey)
	if err != nil {
		err = fmt.Errorf("error parsing private key: %v", err)
		return nil, err
	}
	return keyForSigningJwt, nil
}

func BuildAndSign(claimValue interface{}, expiration time.Time) (string, error) {
	jwtToken, err := jwt.NewBuilder().
		Issuer(jwtIssuer).Audience([]string{jwtAudience}).IssuedAt(time.Now()).
		Expiration(expiration).Claim(KDCustomClaimName, claimValue).
		Build()
	if err != nil {
		err = fmt.Errorf("error building jwt token: %v", err)
		return "", err
	}

	keyForSigningJwt, err := getPrivateKey()
	if err != nil {
		return "", err
	}
	signedJwtToken, err := jwt.Sign(jwtToken, jwt.WithKey(jwa.HS512, keyForSigningJwt))
	if err != nil {
		err = fmt.Errorf("error signing jwt token: %v", err)
		return "", err
	}

	return string(signedJwtToken), nil
}

func VerifySignedToken(signedJwtToken string) (jwt.Token, error) {
	keyForSigningJwt, err := getPrivateKey()
	if err != nil {
		return nil, err
	}

	jwtToken, err := jwt.Parse([]byte(signedJwtToken), jwt.WithKey(jwa.HS512, keyForSigningJwt))
	if err != nil {
		err = fmt.Errorf("error parsing jwt token: %v", err)
		return nil, err
	}

	if jwtToken.Expiration().Unix() < time.Now().Unix() {
		return nil, fmt.Errorf(createerror.UserCredentialExpiredErrorCode)
	}
	if jwtToken.Issuer() != jwtIssuer {
		return nil, fmt.Errorf(createerror.UserCredentialIncorrectErrorCode)
	}

	return jwtToken, nil
}
