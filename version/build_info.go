package version

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	// Returns the raw build information set by the linker during build
	RawBuildInfo func() (compileDate, version string)
)

// Returns the prepared build information
func BuildInfo() (compileDate, version string) {
	c, v := RawBuildInfo()
	compileTime := "unknown"

	if seconds, err := strconv.ParseInt(c, 10, 64); err == nil {
		compileTime = time.Unix(seconds, 0).Format(time.RFC1123)
	}

	return compileTime, v
}

// Returns the version string
func String() string {
	compileDate, version := BuildInfo()

	sb := new(strings.Builder)
	fmt.Fprintln(sb, "ScrutZone "+version)
	fmt.Fprintf(sb, "Compiled at %s with %s\n\n", compileDate, runtime.Version())

	return sb.String()
}
