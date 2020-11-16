package domain

import (
	"encoding/base64"
	"errors"
	"github.com/cossacklabs/themis/gothemis/cell"
	"github.com/cossacklabs/themis/gothemis/keys"
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

func (cpf CPF) Encrypt(symmetricKey string, encryptionContext string) (CPF, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(symmetricKey)
	if err != nil {
		return "", err
	}
	cell, err := cell.ContextImprintWithKey(&keys.SymmetricKey{Value: keyBytes})
	if err != nil || cell == nil {
		return "", err
	}
	encrypted, err := cell.Encrypt([]byte(cpf), []byte(encryptionContext))
	if err != nil {
		return "", err
	}
	encryptedCpf := base64.StdEncoding.EncodeToString(encrypted)
	return CPF(encryptedCpf), nil
}

func (cpf CPF) Decrypt(symmetricKey string, encryptionContext string) (CPF, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(symmetricKey)
	if err != nil {
		return "", err
	}
	cell, err := cell.ContextImprintWithKey(&keys.SymmetricKey{Value: keyBytes})
	if err != nil || cell == nil {
		return "", err
	}
	keyBytes, err = base64.StdEncoding.DecodeString(string(cpf))
	if err != nil {
		return "", err
	}
	decrypted, err := cell.Decrypt(keyBytes, []byte(encryptionContext))
	if err != nil {
		return "", err
	}
	decryptedCpf, err := NewCPF(string(decrypted))
	if err != nil {
		return "", err
	}
	return decryptedCpf, nil
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
