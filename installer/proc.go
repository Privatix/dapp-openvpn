package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/sethvargo/go-password/password"
	"gopkg.in/reform.v1"

	"github.com/privatix/dappctrl/data"
	"github.com/privatix/dappctrl/util"

	"github.com/privatix/dapp-openvpn/adapter/config"
)

const (
	templatePath = "templates"
	productPath  = "products"

	offeringTemplate = "offering.json"
	accessTemplate   = "access.json"

	serverProduct = "server.json"
	clientProduct = "client.json"

	agentAdapterConfig  = "dappvpn.config.json"
	clientAdapterConfig = "dappvpn.config.json"

	jsonIdent = "    "

	passwordLength = 12
	saltLength     = 9 * 1e18
)

func processor(dir string, adjust bool,
	tx *reform.TX, agent bool) error {
	srvProduct, cliProduct, err := handler(dir, tx)
	if err != nil {
		return err
	}

	if adjust {
		var product *data.Product
		var adapterConfig string

		if agent {
			product = srvProduct
			adapterConfig = agentAdapterConfig
		} else {
			product = cliProduct
			adapterConfig = clientAdapterConfig
		}

		configFile := filepath.Join(dir, adapterConfig)

		err = adjustment(product, configFile)
		if err != nil {
			return err
		}
	}

	for _, product := range []*data.Product{srvProduct, cliProduct} {
		err = importProduct(tx, product)
		if err != nil {
			return err
		}
	}

	return nil
}

func handler(dir string, tx *reform.TX) (srvProduct,
	cliProduct *data.Product, err error) {
	err = validateRoot(dir)
	if err != nil {
		return nil, nil, err
	}

	offerTplFile := filepath.Join(dir, templatePath, offeringTemplate)
	accessTplFile := filepath.Join(dir, templatePath, accessTemplate)

	offerTpl, _, err := templates(tx, offerTplFile, accessTplFile)
	if err != nil {
		return nil, nil, err
	}

	serverProductFile := filepath.Join(dir, productPath, serverProduct)
	clientProductFile := filepath.Join(dir, productPath, clientProduct)

	return products(serverProductFile, clientProductFile, offerTpl.ID)
}

func adjustment(product *data.Product, configFile string) error {
	pass, err := setProductAuth(product)
	if err != nil {
		return err
	}

	cfg := config.NewConfig()

	err = util.ReadJSONFile(configFile, &cfg)
	if err != nil {
		return err
	}

	cfg.Connector.Username = product.ID
	cfg.Connector.Password = pass

	return util.WriteJSONFile(configFile, "", jsonIdent, &cfg)
}

func templates(tx *reform.TX, offer,
	access string) (offerTpl, accessTpl *data.Template, err error) {
	offerTpl, err = importTemplate(offer, tx)
	if err != nil {
		return nil, nil, err
	}

	accessTpl, err = importTemplate(access, tx)
	if err != nil {
		return nil, nil, err
	}
	return offerTpl, accessTpl, err
}

func products(serverFile, clientFile, templateID string) (srvProduct,
	cliProduct *data.Product, err error) {
	srvProduct, err = productFromFile(serverFile)
	if err != nil {
		return nil, nil, err
	}

	if !productConcord(srvProduct, templateID) {
		return nil, nil, ErrNotAssociated
	}

	cliProduct, err = productFromFile(clientFile)
	if err != nil {
		return nil, nil, err
	}

	if !productConcord(cliProduct, templateID) {
		return nil, nil, ErrNotAssociated
	}
	return srvProduct, cliProduct, err
}

func validateDir(name string, expect map[string]bool) error {
	info, err := os.Stat(name)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s - is not a directory", name)
	}

	dir, err := os.Open(name)
	if err != nil {
		return err
	}
	defer dir.Close()

	items, err := dir.Readdirnames(-1)
	if err != nil {
		return err
	}

	isFile := func(name string) bool {
		stat, err := os.Stat(name)
		if err != nil {
			return false
		}
		return !stat.IsDir()
	}

	for _, v := range items {
		if expect[v] && !isFile(filepath.Join(name, v)) {
			return ErrNotFile
		}

		delete(expect, v)
	}

	if len(expect) != 0 {
		return ErrNotAllItems
	}
	return nil
}

func validateRoot(dir string) error {
	rootItems := map[string]bool{
		templatePath:       false,
		productPath:        false,
		agentAdapterConfig: true,
	}

	err := validateDir(dir, rootItems)
	if err != nil {
		return err
	}

	tplItems := map[string]bool{
		offeringTemplate: true,
		accessTemplate:   true,
	}

	err = validateDir(filepath.Join(dir, templatePath), tplItems)
	if err != nil {
		return err
	}

	productItems := map[string]bool{
		serverProduct: true,
		clientProduct: true,
	}

	return validateDir(filepath.Join(dir, productPath), productItems)
}

func productConcord(product *data.Product, tplID string) bool {
	if product.OfferTplID == nil {
		return false
	}
	return *product.OfferTplID == tplID
}

func importTemplate(file string, tx *reform.TX) (*data.Template, error) {
	var template *data.Template

	err := util.ReadJSONFile(file, &template)
	if err != nil {
		return nil, err
	}

	err = tx.Insert(template)
	if err != nil {
		return nil, err
	}
	return template, err
}

func importProduct(tx *reform.TX, product *data.Product) error {
	return tx.Insert(product)
}

func productFromFile(file string) (product *data.Product, err error) {
	err = util.ReadJSONFile(file, &product)
	return product, err
}

func setProductAuth(product *data.Product) (string, error) {
	product.ID = util.NewUUID()

	salt, err := rand.Int(rand.Reader, big.NewInt(saltLength))
	if err != nil {
		return "", err
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(10))
	pass, _ := password.Generate(passwordLength, int(n.Int64()), 0, false, false)

	passwordHash, err := data.HashPassword(pass, string(salt.Uint64()))
	if err != nil {
		return "", err
	}

	product.Password = passwordHash
	product.Salt = salt.Uint64()

	return pass, nil
}
