package zk

import (
	"bytes"
	"fmt"
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk/circuit"
	"globe-and-citizen/layer8/server/resource_server/utils"
)

type IProofProcessor interface {
	GenerateProof(emailAddress string, salt string, verificationCode string) ([]byte, error)
	VerifyProof(verificationCode string, salt string, proofBytes []byte) error
}

type ProofProcessor struct {
	cs constraint.ConstraintSystem

	provingKey      groth16.ProvingKey
	verificationKey groth16.VerifyingKey
}

func NewProofProcessor(
	cs constraint.ConstraintSystem,
	provingKey groth16.ProvingKey,
	verificationKey groth16.VerifyingKey,
) *ProofProcessor {
	g := new(ProofProcessor)

	g.cs = cs
	g.provingKey = provingKey
	g.verificationKey = verificationKey

	return g
}

func (pv *ProofProcessor) GenerateProof(
	emailAddress string,
	salt string,
	verificationCode string,
) ([]byte, error) {
	emailAsCircuitVariables, err := utils.StringToCircuitVariables(emailAddress)
	if err != nil {
		return []byte{}, err
	}
	saltAsCircuitVariables, err := utils.StringToCircuitVariables(salt)
	if err != nil {
		return []byte{}, err
	}

	codeAsCircuitVariables, err := utils.ConvertCodeToCircuitVariables(verificationCode)
	if err != nil {
		return []byte{}, err
	}

	circ := &circuit.MimcCircuit{
		EmailAsVariables: emailAsCircuitVariables, /* secret */
		SaltAsVariables:  saltAsCircuitVariables,  /* public */
		VerificationCode: codeAsCircuitVariables,  /* public */
	}

	witness, err := frontend.NewWitness(
		circ,
		ecc.BN254.ScalarField(),
	)
	if err != nil {
		return []byte{}, fmt.Errorf("error while generating zk-snarks witness: %e", err)
	}

	proof, err := groth16.Prove(pv.cs, pv.provingKey, witness)
	if err != nil {
		return []byte{}, err
	}

	var byteBuffer bytes.Buffer
	_, err = proof.WriteTo(&byteBuffer)
	if err != nil {
		return []byte{}, fmt.Errorf("error while writing proof to byte buffer: %e", err)
	}

	return byteBuffer.Bytes(), nil
}

func (pv *ProofProcessor) VerifyProof(
	verificationCode string, salt string, proofBytes []byte,
) error {
	codeAsCircuitVariables, err := utils.ConvertCodeToCircuitVariables(verificationCode)
	saltAsCircuitVariables, err := utils.StringToCircuitVariables(salt)

	var proof = groth16.NewProof(ecc.BN254)
	_, err = proof.ReadFrom(bytes.NewReader(proofBytes))
	if err != nil {
		return fmt.Errorf("error while reading proof bytes: %e", err)
	}

	emailAsVariables := [utils.EmailFrRepresentationSize]frontend.Variable{}
	for i := 0; i < utils.EmailFrRepresentationSize; i++ {
		emailAsVariables[i] = 0
	}

	witness, err := frontend.NewWitness(
		&circuit.MimcCircuit{
			EmailAsVariables: emailAsVariables,
			SaltAsVariables:  saltAsCircuitVariables,
			VerificationCode: codeAsCircuitVariables,
		},
		ecc.BN254.ScalarField(),
	)
	if err != nil {
		return fmt.Errorf("error while constructing a witness for zk proof verification: %e", err)
	}

	publicWitness, err := witness.Public()
	if err != nil {
		return fmt.Errorf("error while retrieving public witness: %e", err)
	}

	err = groth16.Verify(proof, pv.verificationKey, publicWitness)

	if err != nil {
		return fmt.Errorf("could not verify proof: %e", err)
	}

	return nil
}
