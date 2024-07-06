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

    struct Client {
        string clientId;
        uint256 unpaidBill;
        uint64 lastUsageFetchTime;
        uint64 rate;
        Transaction[] transactions;
    }

    struct Transaction {
        uint256 amount;
        uint64 timestamp;
        TransactionType transactionType;
    }

    struct BillingInput {
        string clientId;
        uint64 amount;
        uint64 timestamp;
    }

    string[] public clientIDs;
    mapping(string => Client) public clients;

    constructor(
        address payable _transactionAddress
    ) {
        owner = msg.sender;
        transactionAddress = _transactionAddress;
    }

    event ClientCreated(string clientId);
    event BillAdded(string clientId, uint256 amount);
    event BillPaid(string clientId, uint256 amount);

    function changeTransactionAddress(
        address payable _transactionAddress
    ) external onlyOwner {
        transactionAddress = _transactionAddress;
    }

    function newClient(
        uint64 rate,
        string memory clientId
    ) external onlyOwner {
        Client storage clientToBeStored = clients[clientId];

        clientToBeStored.clientId = clientId;
        clientToBeStored.unpaidBill = 0;
        clientToBeStored.rate = rate;
        clientToBeStored.lastUsageFetchTime = uint64(block.timestamp);

        clientIDs.push(clientId);

        emit ClientCreated(clientId); 
    }

    function addBillToClient(
        string memory clientId,
        uint64 amount,
        uint64 timestamp
    ) public onlyOwner {
        Client storage updatedClient = clients[clientId];
        uint256 amountToBePaid = amount * updatedClient.rate;

        Transaction memory transaction = Transaction({
            amount: amountToBePaid,
            timestamp: timestamp,
            transactionType: TransactionType.BILL
        });

        updatedClient.transactions.push(transaction);
        updatedClient.unpaidBill += amountToBePaid;
        updatedClient.lastUsageFetchTime = timestamp;

        emit BillAdded(clientId, amountToBePaid);
    }

    function bulkAddBillToClient(
        BillingInput[] memory billings
    ) external onlyOwner {
        for (uint256 i = 0; i < billings.length; i++) {
            addBillToClient(billings[i].clientId, billings[i].amount, billings[i].timestamp);
        }
    }

    function payBill(string memory clientId) external payable {
        Client storage updatedClient = clients[clientId];
        require(msg.value <= updatedClient.unpaidBill, "Overpaying the bill");

        updatedClient.unpaidBill -= msg.value;

        Transaction memory transaction = Transaction({
            amount: msg.value,
            timestamp: uint64(block.timestamp),
            transactionType: TransactionType.PAYMENT
        });

        updatedClient.transactions.push(transaction);

        (bool sent, ) = transactionAddress.call{value: msg.value}("");
        require(sent, "Failed to send payment to contract owner");
        
        emit BillPaid(clientId, msg.value);
    }

    function changeRate(string memory clientId, uint64 rate) external onlyOwner {
        Client storage updatedClient = clients[clientId];
        updatedClient.rate = rate;
    }

    function changeAllClientRates(uint64 rate) external onlyOwner {
        for (uint256 i = 0; i < clientIDs.length; i++) {
            clients[clientIDs[i]].rate = rate;
        }
    }

    function getClientById(string memory clientId)
        external
        view
        returns (Client memory)
    {
        return clients[clientId];
    }

    function getClients() external view returns (Client[] memory) {
        Client[] memory clientToResponse = new Client[](clientIDs.length);

        for (uint256 i = 0; i < clientIDs.length; i++) {
            clientToResponse[i] = clients[clientIDs[i]];
        }

        return clientToResponse;
    }
}
