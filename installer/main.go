package main

import (
	"flag"
	"fmt"

	"github.com/privatix/dappctrl/data"
)

//go:generate go generate ../vendor/github.com/privatix/dappctrl/data/schema.go

// Errors.
var (
	ErrNotAssociated = fmt.Errorf("product is not associated with the template")
	ErrNotFile       = fmt.Errorf("object is not file")
	ErrNotAllItems   = fmt.Errorf("some required items not found")
)

func main() {
	agent := flag.Bool("agent", false, "Whether to install agent")
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

	err = processor(*dir, *setAuth, db, *agent)
	if err != nil {
		panic(err)
	}
}
