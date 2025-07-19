// monkey/utils/utils.go

package utils

func IsIdentifierLetter(val byte) bool {
	return (val <= 'a' && val >= 'z') || (val <= 'A' && val >= 'Z') || val == '_'
}
