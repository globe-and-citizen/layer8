package utils

import (
	"encoding/binary"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark/frontend"
	"unicode/utf8"
)

const VerificationCodeSize = 6
const InputFrRepresentationSize = 38

const runesPerElement = 7
const bytes = 4
const elementByteSize = 32

func StringToFrElements(input string) ([InputFrRepresentationSize]fr.Element, error) {
	runeCount := utf8.RuneCountInString(input)

	runes := make([]rune, runeCount)

	ind := 0
	for _, runeValue := range input {
		runes[ind] = runeValue
		ind++
	}

	var elements [InputFrRepresentationSize]fr.Element

	ind = 0

	for i := 0; i < runeCount; i += runesPerElement {
		var elementBytes [elementByteSize]byte

		offset := 0
		for j := i; j < i+runesPerElement; j++ {
			if j == runeCount {
				break
			}

			binary.LittleEndian.PutUint32(elementBytes[offset:offset+bytes], uint32(runes[j]))

			offset += bytes
		}

		element, err := fr.LittleEndian.Element(&elementBytes)
		if err != nil {
			return [38]fr.Element{}, fmt.Errorf("error while converting input to FR elements: %e", err)
		}

		elements[ind] = element
		ind++
	}

	for ; ind < InputFrRepresentationSize; ind++ {
		elements[ind] = fr.NewElement(0)
	}

	return elements, nil
}

func StringToCircuitVariables(
	input string,
) ([InputFrRepresentationSize]frontend.Variable, error) {
	elements, err := StringToFrElements(input)
	if err != nil {
		return [InputFrRepresentationSize]frontend.Variable{}, err
	}

	var circuitVariables [InputFrRepresentationSize]frontend.Variable

	for i := 0; i < InputFrRepresentationSize; i++ {
		circuitVariables[i] = elements[i]
	}

	return circuitVariables, nil
}

func ConvertCodeToCircuitVariables(code string) ([6]frontend.Variable, error) {
	codeAsCircuitVariables := [6]frontend.Variable{}

	for i := 0; i < len(code); i++ {
		currByte := code[i]

		if currByte >= '0' && currByte <= '9' {
			currByte -= '0'
		} else if currByte >= 'a' && currByte <= 'f' {
			currByte = currByte - 'a' + 10
		} else {
			return codeAsCircuitVariables, fmt.Errorf("invalid character at index %d of the verification code", i)
		}

		codeAsCircuitVariables[i] = currByte
	}

	return codeAsCircuitVariables, nil
}
