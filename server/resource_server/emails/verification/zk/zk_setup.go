package zk

import (
	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"globe-and-citizen/layer8/server/resource_server/emails/verification/zk/circuit"
	"log"
)

func RunZkSnarksSetup() (constraint.ConstraintSystem, groth16.ProvingKey, groth16.VerifyingKey) {
	cs := GenerateConstraintSystem()

	provingKey, verifyingKey, err := groth16.Setup(cs)
	if err != nil {
		log.Fatalf("Error happened during the groth16 setup: %e", err)
	}

	return cs, provingKey, verifyingKey
}

func GenerateConstraintSystem() constraint.ConstraintSystem {
	zkCircuit := circuit.NewMimcCircuit()

	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, zkCircuit)
	if err != nil {
		log.Fatalf("Error while generating zk-snarks constraint system: %e", err)
	}

	return cs
}
