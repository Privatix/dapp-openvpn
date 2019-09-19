package main

import (
	"flag"
	"fmt"

	"gopkg.in/reform.v1"

	"github.com/privatix/dappctrl/data"
)

const appVersion = "0.12.0"

// Errors.
var (
	ErrNotAssociated = fmt.Errorf("product is not associated with the template")
	ErrNotFile       = fmt.Errorf("object is not file")
	ErrNotAllItems   = fmt.Errorf("some required items not found")
)

func main() {
	conn := flag.String("connstr",
		"user=postgres dbname=dappctrl sslmode=disable",
		"PostgreSQL connection string")
	dir := flag.String("rootdir", "", "Full path to root directory"+
		" of service adapter")
	setAuth := flag.Bool("setauth", false, "Generate authentication"+
		" credentials for service adapter")

	flag.Parse()

	db, err := data.NewDBFromConnStr(*conn)
	if err != nil {
		panic(err)
	}
	defer data.CloseDB(db)

	if err := db.InTransaction(func(t *reform.TX) error {
		if err := processor(*dir, *setAuth, t); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}
}
