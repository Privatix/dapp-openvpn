// Package main is installer main package.
package main

import (
	"flag"

	"github.com/privatix/dapp-openvpn/cmd/installer/common/cli/agent"
	"github.com/privatix/dapp-openvpn/cmd/installer/common/cli/connstr"
	"github.com/privatix/dapp-openvpn/cmd/installer/common/cli/rootdir"
	"github.com/privatix/dapp-openvpn/common/cli/version"
	"github.com/privatix/dapp-openvpn/common/database/sql"
)

// Values for versioning.
var (
	Commit  string
	Version string
)

func main() {
	verFlag := version.NewVersionFlag(Commit, Version)
	agentFlag := agent.NewWhetherAgentFlag()
	connStrFlag := connstr.NewConnFlag()
	rootDirFlag := rootdir.NewFlagRootDir()

	flag.Parse()

	verFlag.Process()

	db, err := sql.NewDBFromConnStr(connStrFlag.Value().(*string))
	if err != nil {
		panic(err)
	}
	defer sql.CloseDB(db)

	err = agentFlag.Process()
	if err != nil {
		panic(err)
	}

	err = rootDirFlag.AdditionalParams(db, agentFlag.Value().(*bool))
	if err != nil {
		panic(err)
	}

	err = rootDirFlag.Process()
	if err != nil {
		panic(err)
	}
}
