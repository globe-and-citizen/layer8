// SPDX-License-Identifier: MIT

pragma solidity ^0.8.28;

error TransferFailed();

contract PayWithCrypto {
    address private contractOwner;
    address payable public layer8WalletAddress;

    event TrafficPaid(string, address, uint);

    constructor(address payable _layer8WalletAddress) {
        layer8WalletAddress = _layer8WalletAddress;
        contractOwner = msg.sender;
    }

    modifier onlyOwner() {
        require(msg.sender == contractOwner, "Only owner is allowed to execute this function");
        _;
    }

    function setLayer8WalletAddress(address payable _newLayer8WalletAddress) public onlyOwner {
        layer8WalletAddress = _newLayer8WalletAddress;
    }

    function pay(string calldata clientId) external payable {
        // transfer funds to layer8WalletAddress
        (bool success, ) = layer8WalletAddress.call{value: msg.value}("");
        if (!success) {
            revert TransferFailed();
        }

        emit TrafficPaid(clientId, msg.sender, msg.value);
    }
}
