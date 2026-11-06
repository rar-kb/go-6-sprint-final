package service

import (
	"errors"
	"strings"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

var ErrEmptyInput = errors.New("service: пустая строка")

// Convert определяет что передано, возвращает обратную конвертацию
func Convert(input string) (string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", ErrEmptyInput
	}

	if isMorse(input) {
		return morse.ToText(input), nil
	}

	return morse.ToMorse(input), nil
}

// isMorse возвращает true, если строка состоит из символов которые допускаются в азбуке морзе
func isMorse(s string) bool {
	if !strings.ContainsAny(s, ".-") {
		return false
	}

	return strings.Trim(s, ".-/ \t\n\r\v\f") == ""
}
