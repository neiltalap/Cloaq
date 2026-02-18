package main

func main() {
	cli := MakeCommandRegistry(
		&RunCommand{},
		&SettingsCommand{},
	)

	help := &HelpCommand{cli: cli}
	cli.commands[help.Name()] = help

	cli.Execute()
}
