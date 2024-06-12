// SPDX-License-Identifier: MIT
pragma solidity ^0.8.24;

contract PayAsYouGo {
    address owner;
    address payable transactionAddress;

    enum TransactionType {
        PAYMENT,
        BILL
    }

    modifier onlyOwner() {
        require(msg.sender == owner, "Not the contract owner");
        _;
    }

    struct Agreement {
        bytes32 contractId;
        string clientId;
        uint256 unpaidBill;
        uint64 lastUsageFetchTime;
        uint8 rate;
        Transaction[] transactions;
    }

    struct Transaction {
        uint256 amount;
        uint64 timestamp;
        TransactionType transactionType;
    }

    bytes32[] public contractIds;
    mapping(bytes32 => Agreement) public contracts;

    constructor(
        address payable _transactionAddress
    ) {
        owner = msg.sender;
        transactionAddress = _transactionAddress;
    }

    function changeTransactionAddress(
        address payable _transactionAddress
    ) external onlyOwner {
        transactionAddress = _transactionAddress;
    }

    function newContract(
        uint8 rate,
        string memory clientId
    ) external onlyOwner returns (bytes32) {
        bytes32 contractId = keccak256(
            abi.encodePacked(clientId, block.timestamp)
        );

        Agreement storage contractToBeStored = contracts[contractId];

        contractToBeStored.contractId = contractId;
        contractToBeStored.clientId = clientId;
        contractToBeStored.unpaidBill = 0;
        contractToBeStored.rate = rate;
        contractToBeStored.lastUsageFetchTime = uint64(block.timestamp);

        contractIds.push(contractId);

        return contractId;
    }

    function addBillToContract(
        bytes32 contractId,
        uint64 amount,
        uint64 timestamp
    ) external onlyOwner {
        Agreement storage updatedContract = contracts[contractId];
        uint256 amountToBePaid = amount * updatedContract.rate;

        Transaction memory transaction = Transaction({
            amount: amountToBePaid,
            timestamp: timestamp,
            transactionType: TransactionType.BILL
        });

        updatedContract.transactions.push(transaction);
        updatedContract.unpaidBill += amountToBePaid;
        updatedContract.lastUsageFetchTime = timestamp;
    }

    function payBill(bytes32 contractId) external payable {
        Agreement storage updatedContract = contracts[contractId];
        require(msg.value <= updatedContract.unpaidBill, "Overpaying the bill");

        updatedContract.unpaidBill -= msg.value;

        Transaction memory transaction = Transaction({
            amount: msg.value,
            timestamp: uint64(block.timestamp),
            transactionType: TransactionType.PAYMENT
        });

        updatedContract.transactions.push(transaction);

        (bool sent, ) = transactionAddress.call{value: msg.value}("");
        require(sent, "Failed to send payment to contract owner");
    }

    function getContractById(bytes32 contractId)
        external
        view
        returns (Agreement memory)
    {
        return contracts[contractId];
    }

    function getContracts() external view returns (Agreement[] memory) {
        Agreement[] memory agreements = new Agreement[](contractIds.length);

        for (uint256 i = 0; i < contractIds.length; i++) {
            agreements[i] = contracts[contractIds[i]];
        }

        return agreements;
    }
}
