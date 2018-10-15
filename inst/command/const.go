package command

const rootHelp = `
Usage:
  installer [command] [flags]
Available Commands:
  install     Install product package
  remove      Remove product package
Flags:
  --help      Display help information
  --version   Display the current version of this CLI
Use "installer [command] --help" for more information about a command.
`

const installHelp = `
Usage:
	installer install [flags]
Flags:
	--config	Configuration file
	--help		Display help information
`

const removeHelp = `
Usage:
	installer remove [flags]
Flags:
	--help		Display help information
	--workdir	Product install directory
`
