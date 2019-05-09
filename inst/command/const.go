package command

const rootHelp = `
Usage:
  installer [command] [flags]
Available Commands:
  install     Install product package
  remove      Remove product package
  run         Run service
  start	      Start service
  stop	      Stop service
Flags:
  --help      Display help information
  --version   Display the current version of this CLI
Use "installer [command] --help" for more information about a command.
`

const installHelp = `
Usage:
	installer install [flags]
Flags:
  --config  Configuration file
  --help    Display help information
  --role    Product role
  --workdir Product install directory
`

const templateHelp = `
Usage:
  installer %s [flags]
Flags:
  --help    Display help information
  --workdir Product install directory
`

const envFile = "config/.env.config.json"
