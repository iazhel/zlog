package zlog

const (

	// Linux canvas
	prefixWarning = "  [warning]: "
	prefixError   = "  [error]: "
	suffixOK      = "[OK]" + lineSep
	suffixWarning = "[WARNING]" + lineSep
	suffixError   = "[ERROR]" + lineSep

	prefixInfo      = "     [info]: "
	endOutputLine   = "\r\n################ Zlog session ############### %s"
	unknownStepName = "Unknown step%s"
	prefixStep      = lineSep
	suffixFormat    = "%-65s %s" // suffix format

	/*
		prefixWarning = "  [\033[35mwarning\033[0m]: "
		prefixError   = "  [ \033[31merror\033[0m ]: "
		suffixOK      = "[\033[32mOK\033[0m]" + lineSep
		suffixWarning = "[\033[35mWARNING\033[0m]" + lineSep
		suffixError   = "[ \033[31mERROR\033[0m ]" + lineSep
	*/

	// The number '65' is [Ok],[ERROR],[WANING] distance.

)
