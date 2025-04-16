// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rpcclient

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// PayWithCryptoSmartContractClientMetaData contains all meta data concerning the PayWithCryptoSmartContractClient contract.
var PayWithCryptoSmartContractClientMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_layer8WalletAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"TransferFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"TrafficPaid\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"layer8WalletAddress\",\"outputs\":[{\"internalType\":\"addresspayable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pay\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_newLayer8WalletAddress\",\"type\":\"address\"}],\"name\":\"setLayer8WalletAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b506040516105f63803806105f683398181016040528101906100319190610114565b8060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550335f5f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055505061013f565b5f5ffd5b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6100e3826100ba565b9050919050565b6100f3816100d9565b81146100fd575f5ffd5b50565b5f8151905061010e816100ea565b92915050565b5f60208284031215610129576101286100b6565b5b5f61013684828501610100565b91505092915050565b6104aa8061014c5f395ff3fe608060405260043610610033575f3560e01c80631b9265b814610037578063c882265514610041578063fdb497ff1461006b575b5f5ffd5b61003f610093565b005b34801561004c575f5ffd5b5061005561018f565b60405161006291906102c4565b60405180910390f35b348015610076575f5ffd5b50610091600480360381019061008c919061030b565b6101b4565b005b5f60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16346040516100d990610363565b5f6040518083038185875af1925050503d805f8114610113576040519150601f19603f3d011682016040523d82523d5f602084013e610118565b606091505b5050905080610153576040517f90b8ec1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f4c4cdf51b6d4aa2fe87567c8a994cc56f75dad4da01a749a995bd4fe7ae6185133346040516101849291906103af565b60405180910390a150565b60015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b5f5f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614610242576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161023990610456565b60405180910390fd5b8060015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050565b5f73ffffffffffffffffffffffffffffffffffffffff82169050919050565b5f6102ae82610285565b9050919050565b6102be816102a4565b82525050565b5f6020820190506102d75f8301846102b5565b92915050565b5f5ffd5b6102ea816102a4565b81146102f4575f5ffd5b50565b5f81359050610305816102e1565b92915050565b5f602082840312156103205761031f6102dd565b5b5f61032d848285016102f7565b91505092915050565b5f81905092915050565b50565b5f61034e5f83610336565b915061035982610340565b5f82019050919050565b5f61036d82610343565b9150819050919050565b5f61038182610285565b9050919050565b61039181610377565b82525050565b5f819050919050565b6103a981610397565b82525050565b5f6040820190506103c25f830185610388565b6103cf60208301846103a0565b9392505050565b5f82825260208201905092915050565b7f4f6e6c79206f776e657220697320616c6c6f77656420746f20657865637574655f8201527f20746869732066756e6374696f6e000000000000000000000000000000000000602082015250565b5f610440602e836103d6565b915061044b826103e6565b604082019050919050565b5f6020820190508181035f83015261046d81610434565b905091905056fea2646970667358221220f0d264dea93ff1f4b325f04b779394f5a7a401a504b719c0fdd6d2c774ca2e4a64736f6c634300081c0033",
}

// PayWithCryptoSmartContractClientABI is the input ABI used to generate the binding from.
// Deprecated: Use PayWithCryptoSmartContractClientMetaData.ABI instead.
var PayWithCryptoSmartContractClientABI = PayWithCryptoSmartContractClientMetaData.ABI

// PayWithCryptoSmartContractClientBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PayWithCryptoSmartContractClientMetaData.Bin instead.
var PayWithCryptoSmartContractClientBin = PayWithCryptoSmartContractClientMetaData.Bin

// DeployPayWithCryptoSmartContractClient deploys a new Ethereum contract, binding an instance of PayWithCryptoSmartContractClient to it.
func DeployPayWithCryptoSmartContractClient(auth *bind.TransactOpts, backend bind.ContractBackend, _layer8WalletAddress common.Address) (common.Address, *types.Transaction, *PayWithCryptoSmartContractClient, error) {
	parsed, err := PayWithCryptoSmartContractClientMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PayWithCryptoSmartContractClientBin), backend, _layer8WalletAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &PayWithCryptoSmartContractClient{PayWithCryptoSmartContractClientCaller: PayWithCryptoSmartContractClientCaller{contract: contract}, PayWithCryptoSmartContractClientTransactor: PayWithCryptoSmartContractClientTransactor{contract: contract}, PayWithCryptoSmartContractClientFilterer: PayWithCryptoSmartContractClientFilterer{contract: contract}}, nil
}

// PayWithCryptoSmartContractClient is an auto generated Go binding around an Ethereum contract.
type PayWithCryptoSmartContractClient struct {
	PayWithCryptoSmartContractClientCaller     // Read-only binding to the contract
	PayWithCryptoSmartContractClientTransactor // Write-only binding to the contract
	PayWithCryptoSmartContractClientFilterer   // Log filterer for contract events
}

// PayWithCryptoSmartContractClientCaller is an auto generated read-only Go binding around an Ethereum contract.
type PayWithCryptoSmartContractClientCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayWithCryptoSmartContractClientTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PayWithCryptoSmartContractClientTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayWithCryptoSmartContractClientFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PayWithCryptoSmartContractClientFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayWithCryptoSmartContractClientSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PayWithCryptoSmartContractClientSession struct {
	Contract     *PayWithCryptoSmartContractClient // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// PayWithCryptoSmartContractClientCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PayWithCryptoSmartContractClientCallerSession struct {
	Contract *PayWithCryptoSmartContractClientCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// PayWithCryptoSmartContractClientTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PayWithCryptoSmartContractClientTransactorSession struct {
	Contract     *PayWithCryptoSmartContractClientTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// PayWithCryptoSmartContractClientRaw is an auto generated low-level Go binding around an Ethereum contract.
type PayWithCryptoSmartContractClientRaw struct {
	Contract *PayWithCryptoSmartContractClient // Generic contract binding to access the raw methods on
}

// PayWithCryptoSmartContractClientCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PayWithCryptoSmartContractClientCallerRaw struct {
	Contract *PayWithCryptoSmartContractClientCaller // Generic read-only contract binding to access the raw methods on
}

// PayWithCryptoSmartContractClientTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PayWithCryptoSmartContractClientTransactorRaw struct {
	Contract *PayWithCryptoSmartContractClientTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPayWithCryptoSmartContractClient creates a new instance of PayWithCryptoSmartContractClient, bound to a specific deployed contract.
func NewPayWithCryptoSmartContractClient(address common.Address, backend bind.ContractBackend) (*PayWithCryptoSmartContractClient, error) {
	contract, err := bindPayWithCryptoSmartContractClient(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PayWithCryptoSmartContractClient{PayWithCryptoSmartContractClientCaller: PayWithCryptoSmartContractClientCaller{contract: contract}, PayWithCryptoSmartContractClientTransactor: PayWithCryptoSmartContractClientTransactor{contract: contract}, PayWithCryptoSmartContractClientFilterer: PayWithCryptoSmartContractClientFilterer{contract: contract}}, nil
}

// NewPayWithCryptoSmartContractClientCaller creates a new read-only instance of PayWithCryptoSmartContractClient, bound to a specific deployed contract.
func NewPayWithCryptoSmartContractClientCaller(address common.Address, caller bind.ContractCaller) (*PayWithCryptoSmartContractClientCaller, error) {
	contract, err := bindPayWithCryptoSmartContractClient(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PayWithCryptoSmartContractClientCaller{contract: contract}, nil
}

// NewPayWithCryptoSmartContractClientTransactor creates a new write-only instance of PayWithCryptoSmartContractClient, bound to a specific deployed contract.
func NewPayWithCryptoSmartContractClientTransactor(address common.Address, transactor bind.ContractTransactor) (*PayWithCryptoSmartContractClientTransactor, error) {
	contract, err := bindPayWithCryptoSmartContractClient(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PayWithCryptoSmartContractClientTransactor{contract: contract}, nil
}

// NewPayWithCryptoSmartContractClientFilterer creates a new log filterer instance of PayWithCryptoSmartContractClient, bound to a specific deployed contract.
func NewPayWithCryptoSmartContractClientFilterer(address common.Address, filterer bind.ContractFilterer) (*PayWithCryptoSmartContractClientFilterer, error) {
	contract, err := bindPayWithCryptoSmartContractClient(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PayWithCryptoSmartContractClientFilterer{contract: contract}, nil
}

// bindPayWithCryptoSmartContractClient binds a generic wrapper to an already deployed contract.
func bindPayWithCryptoSmartContractClient(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PayWithCryptoSmartContractClientMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayWithCryptoSmartContractClient.Contract.PayWithCryptoSmartContractClientCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.Contract.PayWithCryptoSmartContractClientTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.Contract.PayWithCryptoSmartContractClientTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayWithCryptoSmartContractClient.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.Contract.contract.Transact(opts, method, params...)
}

// Layer8WalletAddress is a free data retrieval call binding the contract method 0xc8822655.
//
// Solidity: function layer8WalletAddress() view returns(address)
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientCaller) Layer8WalletAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PayWithCryptoSmartContractClient.contract.Call(opts, &out, "layer8WalletAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Layer8WalletAddress is a free data retrieval call binding the contract method 0xc8822655.
//
// Solidity: function layer8WalletAddress() view returns(address)
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientSession) Layer8WalletAddress() (common.Address, error) {
	return _PayWithCryptoSmartContractClient.Contract.Layer8WalletAddress(&_PayWithCryptoSmartContractClient.CallOpts)
}

// Layer8WalletAddress is a free data retrieval call binding the contract method 0xc8822655.
//
// Solidity: function layer8WalletAddress() view returns(address)
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientCallerSession) Layer8WalletAddress() (common.Address, error) {
	return _PayWithCryptoSmartContractClient.Contract.Layer8WalletAddress(&_PayWithCryptoSmartContractClient.CallOpts)
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientTransactor) Pay(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.contract.Transact(opts, "pay")
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientSession) Pay() (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.Contract.Pay(&_PayWithCryptoSmartContractClient.TransactOpts)
}

// Pay is a paid mutator transaction binding the contract method 0x1b9265b8.
//
// Solidity: function pay() payable returns()
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientTransactorSession) Pay() (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.Contract.Pay(&_PayWithCryptoSmartContractClient.TransactOpts)
}

// SetLayer8WalletAddress is a paid mutator transaction binding the contract method 0xfdb497ff.
//
// Solidity: function setLayer8WalletAddress(address _newLayer8WalletAddress) returns()
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientTransactor) SetLayer8WalletAddress(opts *bind.TransactOpts, _newLayer8WalletAddress common.Address) (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.contract.Transact(opts, "setLayer8WalletAddress", _newLayer8WalletAddress)
}

// SetLayer8WalletAddress is a paid mutator transaction binding the contract method 0xfdb497ff.
//
// Solidity: function setLayer8WalletAddress(address _newLayer8WalletAddress) returns()
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientSession) SetLayer8WalletAddress(_newLayer8WalletAddress common.Address) (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.Contract.SetLayer8WalletAddress(&_PayWithCryptoSmartContractClient.TransactOpts, _newLayer8WalletAddress)
}

// SetLayer8WalletAddress is a paid mutator transaction binding the contract method 0xfdb497ff.
//
// Solidity: function setLayer8WalletAddress(address _newLayer8WalletAddress) returns()
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientTransactorSession) SetLayer8WalletAddress(_newLayer8WalletAddress common.Address) (*types.Transaction, error) {
	return _PayWithCryptoSmartContractClient.Contract.SetLayer8WalletAddress(&_PayWithCryptoSmartContractClient.TransactOpts, _newLayer8WalletAddress)
}

// PayWithCryptoSmartContractClientTrafficPaidIterator is returned from FilterTrafficPaid and is used to iterate over the raw logs and unpacked data for TrafficPaid events raised by the PayWithCryptoSmartContractClient contract.
type PayWithCryptoSmartContractClientTrafficPaidIterator struct {
	Event *PayWithCryptoSmartContractClientTrafficPaid // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PayWithCryptoSmartContractClientTrafficPaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayWithCryptoSmartContractClientTrafficPaid)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PayWithCryptoSmartContractClientTrafficPaid)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PayWithCryptoSmartContractClientTrafficPaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayWithCryptoSmartContractClientTrafficPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayWithCryptoSmartContractClientTrafficPaid represents a TrafficPaid event raised by the PayWithCryptoSmartContractClient contract.
type PayWithCryptoSmartContractClientTrafficPaid struct {
	Arg0 common.Address
	Arg1 *big.Int
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterTrafficPaid is a free log retrieval operation binding the contract event 0x4c4cdf51b6d4aa2fe87567c8a994cc56f75dad4da01a749a995bd4fe7ae61851.
//
// Solidity: event TrafficPaid(address arg0, uint256 arg1)
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientFilterer) FilterTrafficPaid(opts *bind.FilterOpts) (*PayWithCryptoSmartContractClientTrafficPaidIterator, error) {

	logs, sub, err := _PayWithCryptoSmartContractClient.contract.FilterLogs(opts, "TrafficPaid")
	if err != nil {
		return nil, err
	}
	return &PayWithCryptoSmartContractClientTrafficPaidIterator{contract: _PayWithCryptoSmartContractClient.contract, event: "TrafficPaid", logs: logs, sub: sub}, nil
}

// WatchTrafficPaid is a free log subscription operation binding the contract event 0x4c4cdf51b6d4aa2fe87567c8a994cc56f75dad4da01a749a995bd4fe7ae61851.
//
// Solidity: event TrafficPaid(address arg0, uint256 arg1)
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientFilterer) WatchTrafficPaid(opts *bind.WatchOpts, sink chan<- *PayWithCryptoSmartContractClientTrafficPaid) (event.Subscription, error) {

	logs, sub, err := _PayWithCryptoSmartContractClient.contract.WatchLogs(opts, "TrafficPaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayWithCryptoSmartContractClientTrafficPaid)
				if err := _PayWithCryptoSmartContractClient.contract.UnpackLog(event, "TrafficPaid", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTrafficPaid is a log parse operation binding the contract event 0x4c4cdf51b6d4aa2fe87567c8a994cc56f75dad4da01a749a995bd4fe7ae61851.
//
// Solidity: event TrafficPaid(address arg0, uint256 arg1)
func (_PayWithCryptoSmartContractClient *PayWithCryptoSmartContractClientFilterer) ParseTrafficPaid(log types.Log) (*PayWithCryptoSmartContractClientTrafficPaid, error) {
	event := new(PayWithCryptoSmartContractClientTrafficPaid)
	if err := _PayWithCryptoSmartContractClient.contract.UnpackLog(event, "TrafficPaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
