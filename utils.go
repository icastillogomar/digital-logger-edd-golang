package eddlogger

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	ColorRed    = "\033[91m"
	ColorYellow = "\033[93m"
	ColorGreen  = "\033[92m"
	ColorCyan   = "\033[96m"
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
)

func supportsColor() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	if os.Getenv("FORCE_COLOR") != "" {
		return true
	}
	return true
}

func colorize(text, color string) string {
	if supportsColor() {
		return color + text + ColorReset
	}
	return text
}

func LogError(message string) {
	prefix := colorize("[digital-edd-logger] ERROR:", ColorRed+ColorBold)
	fmt.Fprintf(os.Stderr, "%s %s\n", prefix, message)
}

func LogWarning(message string) {
	prefix := colorize("[digital-edd-logger] WARNING:", ColorYellow+ColorBold)
	fmt.Fprintf(os.Stderr, "%s %s\n", prefix, message)
}

func LogInfo(message string) {
	prefix := colorize("[digital-edd-logger]", ColorCyan)
	fmt.Printf("%s %s\n", prefix, message)
}

func IsProduction() bool {
	env := strings.ToLower(os.Getenv("ENV"))
	return env == "prod" || env == "production" || env == "qas" || env == "qa"
}

func GetMexicoTimeAsUTC() string {
	loc := time.FixedZone("CST", -6*60*60)
	now := time.Now().In(loc)
	return now.Format("2006-01-02T15:04:05.000Z")
}
