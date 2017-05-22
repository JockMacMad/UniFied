package shell

import "github.com/abiosoft/ishell"

func NewUnifiedShell() *ishell.Shell {
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	shell.Println("Unified Interactive Shell\nType help for more info")

	return shell
}