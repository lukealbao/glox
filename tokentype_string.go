// Code generated by "stringer -type=TokenType"; DO NOT EDIT.

package main

import "strconv"

const _TokenType_name = "LEFT_PARENRIGHT_PARENLEFT_BRACERIGHT_BRACECOMMADOTMINUSPLUSSEMICOLONSLASHSTARBANGBANG_EQUALEQUALEQUAL_EQUALGREATERGREATER_EQUALLESSLESS_EQUALIDENTIFIERSTRINGNUMBERANDCLASSELSEFALSEFUNFORIFNILORPRINTRETURNSUPERTHISTRUEVARWHILEEOF"

var _TokenType_index = [...]uint8{0, 10, 21, 31, 42, 47, 50, 55, 59, 68, 73, 77, 81, 91, 96, 107, 114, 127, 131, 141, 151, 157, 163, 166, 171, 175, 180, 183, 186, 188, 191, 193, 198, 204, 209, 213, 217, 220, 225, 228}

func (i TokenType) String() string {
	if i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}
