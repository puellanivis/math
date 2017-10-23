package primative

import (
	"fmt"
)

type digitTransform []rune

func (t digitTransform) Transform(v interface{}) string {
	var s string

	switch v := v.(type) {
	case string:
		s = v
	default:
		s = fmt.Sprint(v)
	}

	runes := []rune(s)

	for i, r := range runes {
		d := r - '0'
		if d >= 0 && d < 10 {
			runes[i] = t[d]
		}
	}

	return string(runes)
}

var (
	superscript = digitTransform("⁰¹²³⁴⁵⁶⁷⁸⁹")
	subscript   = digitTransform("₀₁₂₃₄₅₆₇₈₉")
)
