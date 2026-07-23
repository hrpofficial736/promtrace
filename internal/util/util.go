package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func GenerateID() string {
	b := make([]byte, 4)

	rand.Read(b)

	return hex.EncodeToString(b)
}

func ParseDuration(s string) (time.Duration, error) {
	numStr, found := strings.CutSuffix(s, "d")

	if !found {
		return time.ParseDuration(s)
	}

	days, err := strconv.Atoi(numStr)
	if err != nil {
		return 0, err
	}

	return time.Duration(days) * 24 * time.Hour, nil
}

func GetPromtraceHeadingAsciiText() string {
	return `
██████╗ ██████╗  ██████╗ ███╗   ███╗████████╗██████╗  █████╗  ██████╗███████╗
██╔══██╗██╔══██╗██╔═══██╗████╗ ████║╚══██╔══╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██████╔╝██████╔╝██║   ██║██╔████╔██║   ██║   ██████╔╝███████║██║     █████╗  
██╔═══╝ ██╔══██╗██║   ██║██║╚██╔╝██║   ██║   ██╔══██╗██╔══██║██║     ██╔══╝  
██║     ██║  ██║╚██████╔╝██║ ╚═╝ ██║   ██║   ██║  ██║██║  ██║╚██████╗███████╗
╚═╝     ╚═╝  ╚═╝ ╚═════╝ ╚═╝     ╚═╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝ ╚═════╝╚══════╝
`
}

func FmtCost(cost float64) string {
	return fmt.Sprintf("$%.6f", cost)
}
