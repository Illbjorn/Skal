package pprint

type color string

const (
	colorReset   color = "\033[0m"
	colorRed     color = "\033[31m"
	colorGreen   color = "\033[32m"
	colorYellow  color = "\033[33m"
	colorBlue    color = "\033[34m"
	colorMagenta color = "\033[35m"
	colorCyan    color = "\033[36m"
	colorGray    color = "\033[37m"
	colorWhite   color = "\033[97m"
)
