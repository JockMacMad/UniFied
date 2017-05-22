package shell

import (
	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
)

func NewUnifiedShell() *ishell.Shell {
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()
	color.Set(color.FgGreen)
	shell.Print("Unified Interactive Shell\nType")
	color.Set(color.FgWhite)
	shell.Print(" help ")
	color.Set(color.FgGreen)
	shell.Print("for more info\n")
	color.Set(color.FgWhite)

	return shell
}