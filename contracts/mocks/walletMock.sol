// SPDX-License-Identifier: GPLv3

pragma solidity ^0.6.11;

contract WalletMock {
    /// @dev Ether can be deposited from any source, so this contract must be payable by anyone.
    receive() external payable {}

    function transfer(address payable _to, uint256 _amount) external {
        _to.transfer(_amount);
    }

    function sendValue(address payable _to, uint256 _amount) external {
        (bool success, ) = _to.call{value: _amount}("");
        require(success, "sendValue failed");
    }
}
