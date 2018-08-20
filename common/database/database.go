// Package database provides the implementations for database objects.
package database

import (
	"encoding/json"
	"math/big"

	"github.com/privatix/dapp-openvpn/common/crypto"
)

const (
	passwordLength = 12
	saltLength     = 9 * 1e18
)

//go:generate reform

// Template is a user defined structures.
// It can be an offer or access template.
//reform:templates
type Template struct {
	ID   string          `json:"id" reform:"id,pk"`
	Hash string          `json:"hash" reform:"hash"`
	Raw  json.RawMessage `json:"raw" reform:"raw"`
	Kind string          `json:"kind" reform:"kind"`
}

// Product stores configuration settings
// for specific product and authentication for adapter.
//reform:products
type Product struct {
	ID                     string          `json:"id" reform:"id,pk"`
	Name                   string          `json:"name" reform:"name"`
	OfferTplID             *string         `json:"offerTplID" reform:"offer_tpl_id"`
	OfferAccessID          *string         `json:"offerAccessID" reform:"offer_access_id"`
	UsageRepType           string          `json:"usageRepType" reform:"usage_rep_type"`
	IsServer               bool            `json:"isServer" reform:"is_server"`
	Salt                   uint64          `json:"-" reform:"salt"`
	Password               string          `json:"-" reform:"password"`
	ClientIdent            string          `json:"clientIdent" reform:"client_ident"`
	Config                 json.RawMessage `json:"config" reform:"config"`
	ServiceEndpointAddress *string         `json:"serviceEndpointAddress" reform:"service_endpoint_address"`
}

// SetProductAuth sets random password and salt to the product.
func SetProductAuth(product *Product) (string, error) {
	product.ID = crypto.NewUUID()

	salt, err := crypto.RandNumber(big.NewInt(saltLength))
	if err != nil {
		return "", err
	}

	pass := crypto.RandPass(passwordLength)

	passwordHash, err := crypto.HashPassword(pass, string(salt.Uint64()))
	if err != nil {
		return "", err
	}

	product.Password = passwordHash
	product.Salt = salt.Uint64()

	return pass, nil
}
