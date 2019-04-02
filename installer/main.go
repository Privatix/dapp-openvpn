package main

import (
	"database/sql"
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
		return writeAppVersion(t)
	}); err != nil {
		panic(err)
	}
}

func writeAppVersion(tx *reform.TX) error {
	versionSetting := &data.Setting{}
	err := tx.FindOneTo(versionSetting, "key", data.SettingAppVersion)
	if err == sql.ErrNoRows {
		err = tx.Insert(&data.Setting{
			Key:         data.SettingAppVersion,
			Value:       appVersion,
			Permissions: data.ReadOnly,
			Name:        "App version",
		})
	} else if err == nil {
		fmt.Printf("%s before update: %v\n",
			data.SettingAppVersion, versionSetting.Value)

		versionSetting.Value = appVersion
		err = tx.Update(versionSetting)
	}

	if err != nil {
		fmt.Println("failed to write app version: ", err)
		return err
	}

	fmt.Printf("%s after update: %v\n",
		data.SettingAppVersion, appVersion)
	return err
}
