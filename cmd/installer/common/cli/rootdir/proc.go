package rootdir

import (
	"fmt"
	"path/filepath"

	"gopkg.in/reform.v1"

	"github.com/privatix/dapp-openvpn/cmd/adapter/common/cfg"
	"github.com/privatix/dapp-openvpn/common/database"
	"github.com/privatix/dapp-openvpn/common/path"
)

const (
	templatePath = "templates"
	productPath  = "products"

	offeringTemplate = "offering.json"
	accessTemplate   = "access.json"

	serverProduct = "server.json"
	clientProduct = "client.json"

	adapterConfig = "adapter.config"
)

// Rootdir flag processor errors.
var (
	ErrNotAssociated = fmt.Errorf("product is not associated " +
		"with the template")
	ErrNotFile     = fmt.Errorf("object is not file")
	ErrNotAllItems = fmt.Errorf("some required items not found")
)

func processor(dir *string, adjust *bool,
	db *reform.DB, agent bool) error {
	srvProduct, cliProduct, err := handler(*dir, db)
	if err != nil {
		return err
	}

	if *adjust {
		var product *database.Product

		if agent {
			product = srvProduct
		} else {
			product = cliProduct
		}

		configFile := filepath.Join(*dir, adapterConfig)

		err = adjustment(product, configFile)
		if err != nil {
			return err
		}
	}

	for _, product := range []*database.Product{srvProduct, cliProduct} {
		err = importProduct(db, product)
		if err != nil {
			return err
		}
	}

	return nil
}

func handler(dir string, db *reform.DB) (srvProduct,
	cliProduct *database.Product, err error) {
	err = validateRoot(dir)
	if err != nil {
		return nil, nil, err
	}

	offerTplFile := filepath.Join(dir, templatePath, offeringTemplate)
	accessTplFile := filepath.Join(dir, templatePath, accessTemplate)

	offerTpl, _, err := templates(db, offerTplFile, accessTplFile)
	if err != nil {
		return nil, nil, err
	}

	serverProductFile := filepath.Join(dir, productPath, serverProduct)
	clientProductFile := filepath.Join(dir, productPath, clientProduct)

	return products(serverProductFile, clientProductFile, offerTpl.ID)
}

func adjustment(product *database.Product, configFile string) error {
	pass, err := database.SetProductAuth(product)
	if err != nil {
		return err
	}

	config := cfg.DefaultConfig()

	err = path.ReadJSONFile(configFile, &config)
	if err != nil {
		return err
	}

	config.Connector.Username = product.ID
	config.Connector.Password = pass

	return path.WriteJSONFile(configFile, &config)
}

func templates(db *reform.DB, offer,
	access string) (offerTpl, accessTpl *database.Template, err error) {
	offerTpl, err = importTemplate(offer, db)
	if err != nil {
		return nil, nil, err
	}

	accessTpl, err = importTemplate(access, db)
	if err != nil {
		return nil, nil, err
	}
	return offerTpl, accessTpl, err
}

func products(serverFile, clientFile, templateID string) (srvProduct,
	cliProduct *database.Product, err error) {
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
	items, err := path.ReadDir(name)
	if err != nil {
		return err
	}

	for _, v := range items {
		if expect[v] && !path.IsFile(filepath.Join(name, v)) {
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
	err := path.IsDirWriteable(dir)
	if err != nil {
		return err
	}

	rootItems := map[string]bool{
		templatePath:  false,
		productPath:   false,
		adapterConfig: true,
	}

	err = validateDir(dir, rootItems)
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

func productConcord(product *database.Product, tplID string) bool {
	if product.OfferTplID == nil {
		return false
	}
	return *product.OfferTplID == tplID
}

func importTemplate(file string, db *reform.DB) (*database.Template, error) {
	var template *database.Template

	err := path.ReadJSONFile(file, &template)
	if err != nil {
		return nil, err
	}

	err = db.Insert(template)
	if err != nil {
		return nil, err
	}
	return template, err
}

func importProduct(db *reform.DB, product *database.Product) error {
	return db.Insert(product)
}

func productFromFile(file string) (product *database.Product, err error) {
	err = path.ReadJSONFile(file, &product)
	return product, err
}
