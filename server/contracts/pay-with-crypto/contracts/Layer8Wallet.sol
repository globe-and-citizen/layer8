// SPDX-License-Identifier: MIT

pragma solidity ^0.8.28;

contract Layer8Wallet {
    event PaymentReceived(address, uint);

    function getBalance() public view returns (uint) {
        return address(this).balance;
    }

    receive() external payable {
        emit PaymentReceived(msg.sender, msg.value);
    }
}
