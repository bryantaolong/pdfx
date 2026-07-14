package commands

import "strings"

func EnsurePDFExt(name string) string {
	if !strings.HasSuffix(strings.ToLower(name), ".pdf") {
		return name + ".pdf"
	}
	return name
}
