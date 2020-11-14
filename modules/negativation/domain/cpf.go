package domain

import (
	"errors"
	"strconv"
	"strings"
)

type CPF string

var (
	NonValidCharacter        = errors.New("non valid character")
	InvalidNumberOfCharacter = errors.New("invalid number of character")
)

const maxNumberOfCharacters = 11

func NewCPF(cpf string) (CPF, error) {
	cpf = removeSpecialCharacters(cpf)
	if haveNonNumericCharacters(cpf) {
		return "", NonValidCharacter
	}
	if haveInvalidNumberOfCharacters(cpf) {
		return "", InvalidNumberOfCharacter
	}
	return CPF(cpf), nil
}

func removeSpecialCharacters(cpf string) string {
	cpf = strings.ReplaceAll(cpf, ".", "")
	cpf = strings.ReplaceAll(cpf, "-", "")
	cpf = strings.ReplaceAll(cpf, " ", "")
	return cpf
}

func haveNonNumericCharacters(cpf string) bool {
	_, err := strconv.Atoi(cpf)
	return err != nil
}

func haveInvalidNumberOfCharacters(cpf string) bool {
	return len(cpf) != maxNumberOfCharacters
}
