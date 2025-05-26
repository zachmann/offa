package model

import (
	"strings"
)

type Claim string

// UserClaims holds claims about a user
type UserClaims map[Claim]any

func (claims UserClaims) GetForHeader(claim Claim) (string, bool) {
	return claims.getAsString(claim, ",")
}
func (claims UserClaims) GetForMemCache(claim Claim) (string, bool) {
	return claims.getAsString(claim, ":")
}

func (claims UserClaims) getAsString(claim Claim, sliceSeparator string) (string, bool) {
	v, ok := claims.GetString(claim)
	if ok {
		return v, true
	}
	vs, ok := claims.GetStringSlice(claim)
	if ok {
		return strings.Join(vs, sliceSeparator), true
	}
	return "", false
}

func (claims UserClaims) GetString(claim Claim) (string, bool) {
	v, ok := claims[claim]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func (claims UserClaims) GetStringSlice(claim Claim) ([]string, bool) {
	v, ok := claims[claim]
	if !ok {
		return nil, false
	}
	s, ok := v.([]string)
	return s, ok
}
