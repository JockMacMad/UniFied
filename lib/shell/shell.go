package shell

import (
	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
)

func NewUnifiedShell() *ishell.Shell {
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	terminal := ishell.New()
	color.Set(color.FgGreen)
	terminal.Print("Unified Interactive Shell\nType")
	color.Set(color.FgWhite)
	terminal.Print(" help ")
	color.Set(color.FgGreen)
	terminal.Print("for more info\n")
	color.Set(color.FgWhite)

	return terminal
}
