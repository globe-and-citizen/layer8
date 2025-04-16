# Pay with crypto

This module contains the smart contracts that enable pay with crypto functionality on the client portal, as well as
all the associated deployment scripts and unit tests.
Primary logic is contained within the `/contracts/PayWithCrypto.sol` smart contract, 
`/contracts/Layer8Wallet.sol` represents the Layer8 account used to receive 
payments (currently only for testing purposes), it's planned to use the account on CNC Portal instead when 
the functionality is properly tested and is ready for a production-level usage.
