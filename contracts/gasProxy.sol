/**
 *  Copyright (C) 2019 The Contract Wallet Company Limited
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

pragma solidity ^0.5.17;

import "./internals/controllable.sol";
import "./internals/gasRefundable.sol";


contract GasProxy is Controllable, GasRefundable {
    /// @notice Emits the transaction executed by the controller.
    event ExecutedTransaction(address _destination, uint256 _value, bytes _data, bytes _returnData);

    /// @param _ens_ is the address of the ENS registry.
    /// @param _controllerNode_ ENS node of the controller contract.
    /// @param _gasTokenAddress ENS node of the gas token contract.
    constructor(address _ens_, bytes32 _controllerNode_, address _gasTokenAddress) public GasRefundable(_gasTokenAddress) {
        _initializeENSResolvable(_ens_);
        _initializeControllable(_controllerNode_);
    }

    /// @param _gasTokenAddress Address of the gas token used to refund gas.
    function setGasToken(address _gasTokenAddress) external onlyController {
        _setGasToken(_gasTokenAddress);
    }

    /// @param _gasCost Gas cost of the gas token free method call.
    function setFreeCallGasCost(uint256 _gasCost) external onlyController {
        _setFreeCallGasCost(_gasCost);
    }

    /// @param _gasRefund Amount of gas refunded per unit of gas token.
    function setGasRefundPerUnit(uint256 _gasRefund) external onlyController {
        _setGasRefundPerUnit(_gasRefund);
    }

    /// @notice Executes a controller operation and refunds gas using gas tokens.
    /// @param _destination Destination address of the executed transaction.
    /// @param _value Amount of ETH (wei) to be sent together with the transaction.
    /// @param _data Data payload of the controller transaction.
    function executeTransaction(address _destination, uint256 _value, bytes calldata _data) external onlyController refundGas {
        (bool success, bytes memory returnData) = _destination.call.value(_value)(_data);
        require(success, "external call failed");
        emit ExecutedTransaction(_destination, _value, _data, returnData);
    }
}
