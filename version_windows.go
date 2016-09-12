package zlog

const (
	linesSep      = "\r\n"
	prefixWarning = "  [warning]: "
	prefixError   = "  [error]  : "
	prefixInfo    = "     [info]:"
	suffixOK      = "[OK]"
	suffixWarning = "[WARNING]" + linesSep
	suffixError   = "[ERROR]" + linesSep
	prefixStep    = "Step: "
	endOutputLine = "\r\n################ Zlog session ######### %s"
)
