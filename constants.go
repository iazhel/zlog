package zlog

const (
	// If step has linesToFreeOsMem msgs, Step() return memory to OS.
	linesToFreeOsMem = 100000
	// Linux canvas
	prefixWarning = lineSep + "  [\033[35mwarning\033[0m]: "
	prefixError   = lineSep + "  [ \033[31merror\033[0m ]: "
	suffixOK      = "[\033[32mOK\033[0m]"
	suffixWarning = "[\033[35mWARNING\033[0m]"
	suffixError   = "[ \033[31mERROR\033[0m ]"

	// Windows canvas
	prefixWarning_Win = lineSep + "  [warning]: "
	prefixError_Win   = lineSep + "  [ error ]: "
	suffixOK_Win      = "[OK]"
	suffixWarning_Win = "[WARNING]"
	suffixError_Win   = "[ ERROR ]"

	// common const
	lineSep         = "\r\n"
	prefixInfo      = lineSep + "     [info]: "
	endOutputLine   = "\r\n############### End of session.############## %s\r\n"
	unknownStepName = "Unknown step%s"
	firstInfoMsg    = ", first msg: "
	prefixStep      = "" //lineSep

	// The number '65' is [Ok],[ERROR],[WANING] distance.
	suffixFormat = "%-65s %s" // suffix format

)
