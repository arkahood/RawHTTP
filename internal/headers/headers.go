package headers

import (
	"errors"
	"strings"
	"unicode"
)

const FieldLineSeperator = ":"
const crlf = "\r\n"

type Headers map[string]string

// there can be an unlimited amount of whitespace
// before and after the field-value (Header value). However, when parsing a field-name,
// there must be no spaces betwixt the colon and the field-name. In other words,
// these are valid:

// 'Host: localhost:42069'
// '          Host: localhost:42069    '

// But this is not:

// Host : localhost:42069

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	endIdx := strings.Index(string(data), crlf)
	if endIdx == -1 {
		return 0, false, nil
	}
	// if encounter \r\n in front of line that means we consumed all the headers/fieldLines
	if endIdx == 0 {
		return endIdx + 2, true, nil
	}

	currentLine := string(data[:endIdx])
	fieldParts := strings.SplitN(currentLine, FieldLineSeperator, 2)

	if len(fieldParts) != 2 {
		return 0, false, errors.New("field-line have wrong number of parts")
	}

	fieldName := strings.TrimLeft(fieldParts[0], " ")
	if len(fieldName) != len(strings.TrimSpace(fieldName)) {
		return 0, false, errors.New("error in field-name syntax, whitespace unexpected")
	}
	if !isValidFieldName(fieldName) {
		return 0, false, errors.New("invalid characters in field-name")
	}

	fieldValue := strings.TrimSpace(strings.Trim(fieldParts[1], crlf))
	// lowercase the fieldname while adding to the map
	val, exists := h[strings.ToLower(fieldName)]
	if exists {
		h[strings.ToLower(fieldName)] = val + "," + fieldValue
	} else {
		h[strings.ToLower(fieldName)] = fieldValue
	}

	return endIdx + 2, false, nil
}

func isValidFieldName(fieldName string) bool {
	allowedSpecials := map[rune]bool{
		'!': true, '#': true, '$': true, '%': true, '&': true, '\'': true,
		'*': true, '+': true, '-': true, '.': true, '^': true, '_': true,
		'`': true, '|': true, '~': true,
	}

	if len(fieldName) == 0 {
		return false
	}

	for _, ch := range fieldName {
		switch {
		case unicode.IsLetter(ch):
			continue
		case unicode.IsDigit(ch):
			continue
		case allowedSpecials[ch]:
			continue
		default:
			return false
		}
	}

	return true
}
