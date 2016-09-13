// +build !windows

package zlog

const (
	linesSep      = "\n"
	prefixWarning = "  [\033[35mwarning\033[0m]: "
	prefixError   = "  [ \033[31merror\033[0m ]: "
	prefixInfo    = "     [info]:"
	suffixOK      = "[\033[32mOK\033[0m]"
	suffixWarning = "[\033[35mWARNING\033[0m]" + linesSep
	suffixError   = "[ \033[31mERROR\033[0m ]" + linesSep
	prefixStep    = "Step: "
	endOutputLine = linesSep + "################ Zlog session ######### %s"
)
