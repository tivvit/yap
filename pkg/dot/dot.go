package dot

import "strings"

func EscapeHtml(label string) string {
	label = strings.ReplaceAll(label, "\"", "&quot;")
	label = strings.ReplaceAll(label, "&", "&amp;")
	label = strings.ReplaceAll(label, "<", "&lt;")
	label = strings.ReplaceAll(label, ">", "&gt;")
	return label
}
