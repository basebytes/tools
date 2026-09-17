package tools

import "regexp"

var (
	phoneMaskReg = regexp.MustCompile(`(\d{3})(\d{4})(\d+)`)
	cerNoMaskReg = regexp.MustCompile(`(\d{6})(\d+)(.{4})`)
)

func MaskPhone(raw string) (result string) {
	if raw != "" {
		result = phoneMaskReg.ReplaceAllString(raw, "$1****$3")
	}
	return
}

func MaskCerNo(raw string) (result string) {
	if raw != "" {
		result = cerNoMaskReg.ReplaceAllString(raw, "$1********$3")
	}
	return
}
