package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/bits"
	"globe-and-citizen/layer8/server/resource_server/utils"
)

type MimcCircuit struct {
	SaltAsVariables  [utils.InputFrRepresentationSize]frontend.Variable `gnark:",public"`
	InputAsVariables [utils.InputFrRepresentationSize]frontend.Variable `gnark:",secret"`

	VerificationCode [utils.VerificationCodeSize]frontend.Variable `gnark:",public"`
}

func NewMimcCircuit() *MimcCircuit {
	return new(MimcCircuit)
}

func (c *MimcCircuit) Define(api frontend.API) error {
	mimcInstance, _ := mimc.NewMiMC(api)

	for i := 0; i < utils.InputFrRepresentationSize; i++ {
		inputFrVariableBits := bits.ToBinary(api, c.InputAsVariables[i])
		saltFrVariableBits := bits.ToBinary(api, c.SaltAsVariables[i])

		currentVariable := frontend.Variable(0)
		powerOfTwo := frontend.Variable(1)
		for j := 0; j < len(inputFrVariableBits); j++ {
			xoredBit := api.Xor(inputFrVariableBits[j], saltFrVariableBits[j])
			currentVariable = api.Add(currentVariable, api.Mul(xoredBit, powerOfTwo))
			powerOfTwo = api.Mul(powerOfTwo, 2)
		}

		mimcInstance.Write(currentVariable)
	}

	mimcHash := mimcInstance.Sum()

	mimcBits := bits.ToBinary(api, mimcHash)

	code := make([]frontend.Variable, utils.VerificationCodeSize)

	ind := 0
	for i := 23; i >= 0; i -= 4 {
		a0 := api.Mul(mimcBits[i], 8)
		a1 := api.Mul(mimcBits[i-1], 4)
		a2 := api.Mul(mimcBits[i-2], 2)
		a3 := mimcBits[i-3]

		code[ind] = api.Add(a0, a1, a2, a3)

		ind++
	}

	for i := 0; i < utils.VerificationCodeSize; i++ {
		api.AssertIsEqual(code[i], c.VerificationCode[i])
	}

	return nil
}
