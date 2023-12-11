package util

import "fmt"

const (
	ColorReset  = "\033[0m"
	ColorBlack  = "\033[30m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
)

func PrintColored(message string, colorCode string) {
	fmt.Printf("%s%s%s\n", colorCode, message, ColorReset)
}
