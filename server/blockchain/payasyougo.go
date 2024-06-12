// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

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

// PayAsYouGoAgreement is an auto generated low-level Go binding around an user-defined struct.
type PayAsYouGoAgreement struct {
	ContractId         [32]byte
	ClientId           string
	UnpaidBill         *big.Int
	LastUsageFetchTime uint64
	Rate               uint8
	Transactions       []PayAsYouGoTransaction
}

// PayAsYouGoTransaction is an auto generated low-level Go binding around an user-defined struct.
type PayAsYouGoTransaction struct {
	Amount          *big.Int
	Timestamp       uint64
	TransactionType uint8
}

// PayAsYouGoMetaData contains all meta data concerning the PayAsYouGo contract.
var PayAsYouGoMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_transactionAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BillAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BillPaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"}],\"name\":\"ContractCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"addBillToContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_transactionAddress\",\"type\":\"address\"}],\"name\":\"changeTransactionAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"contractIds\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"contracts\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"unpaidBill\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"lastUsageFetchTime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rate\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"}],\"name\":\"getContractById\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"unpaidBill\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"lastUsageFetchTime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rate\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"enumPayAsYouGo.TransactionType\",\"name\":\"transactionType\",\"type\":\"uint8\"}],\"internalType\":\"structPayAsYouGo.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"internalType\":\"structPayAsYouGo.Agreement\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getContracts\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"unpaidBill\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"lastUsageFetchTime\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rate\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"enumPayAsYouGo.TransactionType\",\"name\":\"transactionType\",\"type\":\"uint8\"}],\"internalType\":\"structPayAsYouGo.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"internalType\":\"structPayAsYouGo.Agreement[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"rate\",\"type\":\"uint8\"},{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"}],\"name\":\"newContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"contractId\",\"type\":\"bytes32\"}],\"name\":\"payBill\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// PayAsYouGoABI is the input ABI used to generate the binding from.
// Deprecated: Use PayAsYouGoMetaData.ABI instead.
var PayAsYouGoABI = PayAsYouGoMetaData.ABI

// PayAsYouGo is an auto generated Go binding around an Ethereum contract.
type PayAsYouGo struct {
	PayAsYouGoCaller     // Read-only binding to the contract
	PayAsYouGoTransactor // Write-only binding to the contract
	PayAsYouGoFilterer   // Log filterer for contract events
}

// PayAsYouGoCaller is an auto generated read-only Go binding around an Ethereum contract.
type PayAsYouGoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayAsYouGoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PayAsYouGoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayAsYouGoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PayAsYouGoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PayAsYouGoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PayAsYouGoSession struct {
	Contract     *PayAsYouGo       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PayAsYouGoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PayAsYouGoCallerSession struct {
	Contract *PayAsYouGoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// PayAsYouGoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PayAsYouGoTransactorSession struct {
	Contract     *PayAsYouGoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// PayAsYouGoRaw is an auto generated low-level Go binding around an Ethereum contract.
type PayAsYouGoRaw struct {
	Contract *PayAsYouGo // Generic contract binding to access the raw methods on
}

// PayAsYouGoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PayAsYouGoCallerRaw struct {
	Contract *PayAsYouGoCaller // Generic read-only contract binding to access the raw methods on
}

// PayAsYouGoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PayAsYouGoTransactorRaw struct {
	Contract *PayAsYouGoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPayAsYouGo creates a new instance of PayAsYouGo, bound to a specific deployed contract.
func NewPayAsYouGo(address common.Address, backend bind.ContractBackend) (*PayAsYouGo, error) {
	contract, err := bindPayAsYouGo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PayAsYouGo{PayAsYouGoCaller: PayAsYouGoCaller{contract: contract}, PayAsYouGoTransactor: PayAsYouGoTransactor{contract: contract}, PayAsYouGoFilterer: PayAsYouGoFilterer{contract: contract}}, nil
}

// NewPayAsYouGoCaller creates a new read-only instance of PayAsYouGo, bound to a specific deployed contract.
func NewPayAsYouGoCaller(address common.Address, caller bind.ContractCaller) (*PayAsYouGoCaller, error) {
	contract, err := bindPayAsYouGo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoCaller{contract: contract}, nil
}

// NewPayAsYouGoTransactor creates a new write-only instance of PayAsYouGo, bound to a specific deployed contract.
func NewPayAsYouGoTransactor(address common.Address, transactor bind.ContractTransactor) (*PayAsYouGoTransactor, error) {
	contract, err := bindPayAsYouGo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoTransactor{contract: contract}, nil
}

// NewPayAsYouGoFilterer creates a new log filterer instance of PayAsYouGo, bound to a specific deployed contract.
func NewPayAsYouGoFilterer(address common.Address, filterer bind.ContractFilterer) (*PayAsYouGoFilterer, error) {
	contract, err := bindPayAsYouGo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoFilterer{contract: contract}, nil
}

// bindPayAsYouGo binds a generic wrapper to an already deployed contract.
func bindPayAsYouGo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PayAsYouGoMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayAsYouGo *PayAsYouGoRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayAsYouGo.Contract.PayAsYouGoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayAsYouGo *PayAsYouGoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.PayAsYouGoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayAsYouGo *PayAsYouGoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.PayAsYouGoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PayAsYouGo *PayAsYouGoCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PayAsYouGo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PayAsYouGo *PayAsYouGoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PayAsYouGo *PayAsYouGoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.contract.Transact(opts, method, params...)
}

// ContractIds is a free data retrieval call binding the contract method 0x438c3aff.
//
// Solidity: function contractIds(uint256 ) view returns(bytes32)
func (_PayAsYouGo *PayAsYouGoCaller) ContractIds(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _PayAsYouGo.contract.Call(opts, &out, "contractIds", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ContractIds is a free data retrieval call binding the contract method 0x438c3aff.
//
// Solidity: function contractIds(uint256 ) view returns(bytes32)
func (_PayAsYouGo *PayAsYouGoSession) ContractIds(arg0 *big.Int) ([32]byte, error) {
	return _PayAsYouGo.Contract.ContractIds(&_PayAsYouGo.CallOpts, arg0)
}

// ContractIds is a free data retrieval call binding the contract method 0x438c3aff.
//
// Solidity: function contractIds(uint256 ) view returns(bytes32)
func (_PayAsYouGo *PayAsYouGoCallerSession) ContractIds(arg0 *big.Int) ([32]byte, error) {
	return _PayAsYouGo.Contract.ContractIds(&_PayAsYouGo.CallOpts, arg0)
}

// Contracts is a free data retrieval call binding the contract method 0xec56a373.
//
// Solidity: function contracts(bytes32 ) view returns(bytes32 contractId, string clientId, uint256 unpaidBill, uint64 lastUsageFetchTime, uint8 rate)
func (_PayAsYouGo *PayAsYouGoCaller) Contracts(opts *bind.CallOpts, arg0 [32]byte) (struct {
	ContractId         [32]byte
	ClientId           string
	UnpaidBill         *big.Int
	LastUsageFetchTime uint64
	Rate               uint8
}, error) {
	var out []interface{}
	err := _PayAsYouGo.contract.Call(opts, &out, "contracts", arg0)

	outstruct := new(struct {
		ContractId         [32]byte
		ClientId           string
		UnpaidBill         *big.Int
		LastUsageFetchTime uint64
		Rate               uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ContractId = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.ClientId = *abi.ConvertType(out[1], new(string)).(*string)
	outstruct.UnpaidBill = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LastUsageFetchTime = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.Rate = *abi.ConvertType(out[4], new(uint8)).(*uint8)

	return *outstruct, err

}

// Contracts is a free data retrieval call binding the contract method 0xec56a373.
//
// Solidity: function contracts(bytes32 ) view returns(bytes32 contractId, string clientId, uint256 unpaidBill, uint64 lastUsageFetchTime, uint8 rate)
func (_PayAsYouGo *PayAsYouGoSession) Contracts(arg0 [32]byte) (struct {
	ContractId         [32]byte
	ClientId           string
	UnpaidBill         *big.Int
	LastUsageFetchTime uint64
	Rate               uint8
}, error) {
	return _PayAsYouGo.Contract.Contracts(&_PayAsYouGo.CallOpts, arg0)
}

// Contracts is a free data retrieval call binding the contract method 0xec56a373.
//
// Solidity: function contracts(bytes32 ) view returns(bytes32 contractId, string clientId, uint256 unpaidBill, uint64 lastUsageFetchTime, uint8 rate)
func (_PayAsYouGo *PayAsYouGoCallerSession) Contracts(arg0 [32]byte) (struct {
	ContractId         [32]byte
	ClientId           string
	UnpaidBill         *big.Int
	LastUsageFetchTime uint64
	Rate               uint8
}, error) {
	return _PayAsYouGo.Contract.Contracts(&_PayAsYouGo.CallOpts, arg0)
}

// GetContractById is a free data retrieval call binding the contract method 0xa9c9a918.
//
// Solidity: function getContractById(bytes32 contractId) view returns((bytes32,string,uint256,uint64,uint8,(uint256,uint64,uint8)[]))
func (_PayAsYouGo *PayAsYouGoCaller) GetContractById(opts *bind.CallOpts, contractId [32]byte) (PayAsYouGoAgreement, error) {
	var out []interface{}
	err := _PayAsYouGo.contract.Call(opts, &out, "getContractById", contractId)

	if err != nil {
		return *new(PayAsYouGoAgreement), err
	}

	out0 := *abi.ConvertType(out[0], new(PayAsYouGoAgreement)).(*PayAsYouGoAgreement)

	return out0, err

}

// GetContractById is a free data retrieval call binding the contract method 0xa9c9a918.
//
// Solidity: function getContractById(bytes32 contractId) view returns((bytes32,string,uint256,uint64,uint8,(uint256,uint64,uint8)[]))
func (_PayAsYouGo *PayAsYouGoSession) GetContractById(contractId [32]byte) (PayAsYouGoAgreement, error) {
	return _PayAsYouGo.Contract.GetContractById(&_PayAsYouGo.CallOpts, contractId)
}

// GetContractById is a free data retrieval call binding the contract method 0xa9c9a918.
//
// Solidity: function getContractById(bytes32 contractId) view returns((bytes32,string,uint256,uint64,uint8,(uint256,uint64,uint8)[]))
func (_PayAsYouGo *PayAsYouGoCallerSession) GetContractById(contractId [32]byte) (PayAsYouGoAgreement, error) {
	return _PayAsYouGo.Contract.GetContractById(&_PayAsYouGo.CallOpts, contractId)
}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns((bytes32,string,uint256,uint64,uint8,(uint256,uint64,uint8)[])[])
func (_PayAsYouGo *PayAsYouGoCaller) GetContracts(opts *bind.CallOpts) ([]PayAsYouGoAgreement, error) {
	var out []interface{}
	err := _PayAsYouGo.contract.Call(opts, &out, "getContracts")

	if err != nil {
		return *new([]PayAsYouGoAgreement), err
	}

	out0 := *abi.ConvertType(out[0], new([]PayAsYouGoAgreement)).(*[]PayAsYouGoAgreement)

	return out0, err

}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns((bytes32,string,uint256,uint64,uint8,(uint256,uint64,uint8)[])[])
func (_PayAsYouGo *PayAsYouGoSession) GetContracts() ([]PayAsYouGoAgreement, error) {
	return _PayAsYouGo.Contract.GetContracts(&_PayAsYouGo.CallOpts)
}

// GetContracts is a free data retrieval call binding the contract method 0xc3a2a93a.
//
// Solidity: function getContracts() view returns((bytes32,string,uint256,uint64,uint8,(uint256,uint64,uint8)[])[])
func (_PayAsYouGo *PayAsYouGoCallerSession) GetContracts() ([]PayAsYouGoAgreement, error) {
	return _PayAsYouGo.Contract.GetContracts(&_PayAsYouGo.CallOpts)
}

// AddBillToContract is a paid mutator transaction binding the contract method 0xba281ab7.
//
// Solidity: function addBillToContract(bytes32 contractId, uint64 amount, uint64 timestamp) returns()
func (_PayAsYouGo *PayAsYouGoTransactor) AddBillToContract(opts *bind.TransactOpts, contractId [32]byte, amount uint64, timestamp uint64) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "addBillToContract", contractId, amount, timestamp)
}

// AddBillToContract is a paid mutator transaction binding the contract method 0xba281ab7.
//
// Solidity: function addBillToContract(bytes32 contractId, uint64 amount, uint64 timestamp) returns()
func (_PayAsYouGo *PayAsYouGoSession) AddBillToContract(contractId [32]byte, amount uint64, timestamp uint64) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.AddBillToContract(&_PayAsYouGo.TransactOpts, contractId, amount, timestamp)
}

// AddBillToContract is a paid mutator transaction binding the contract method 0xba281ab7.
//
// Solidity: function addBillToContract(bytes32 contractId, uint64 amount, uint64 timestamp) returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) AddBillToContract(contractId [32]byte, amount uint64, timestamp uint64) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.AddBillToContract(&_PayAsYouGo.TransactOpts, contractId, amount, timestamp)
}

// ChangeTransactionAddress is a paid mutator transaction binding the contract method 0xf4217f2a.
//
// Solidity: function changeTransactionAddress(address _transactionAddress) returns()
func (_PayAsYouGo *PayAsYouGoTransactor) ChangeTransactionAddress(opts *bind.TransactOpts, _transactionAddress common.Address) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "changeTransactionAddress", _transactionAddress)
}

// ChangeTransactionAddress is a paid mutator transaction binding the contract method 0xf4217f2a.
//
// Solidity: function changeTransactionAddress(address _transactionAddress) returns()
func (_PayAsYouGo *PayAsYouGoSession) ChangeTransactionAddress(_transactionAddress common.Address) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.ChangeTransactionAddress(&_PayAsYouGo.TransactOpts, _transactionAddress)
}

// ChangeTransactionAddress is a paid mutator transaction binding the contract method 0xf4217f2a.
//
// Solidity: function changeTransactionAddress(address _transactionAddress) returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) ChangeTransactionAddress(_transactionAddress common.Address) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.ChangeTransactionAddress(&_PayAsYouGo.TransactOpts, _transactionAddress)
}

// NewContract is a paid mutator transaction binding the contract method 0x05c23118.
//
// Solidity: function newContract(uint8 rate, string clientId) returns()
func (_PayAsYouGo *PayAsYouGoTransactor) NewContract(opts *bind.TransactOpts, rate uint8, clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "newContract", rate, clientId)
}

// NewContract is a paid mutator transaction binding the contract method 0x05c23118.
//
// Solidity: function newContract(uint8 rate, string clientId) returns()
func (_PayAsYouGo *PayAsYouGoSession) NewContract(rate uint8, clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.NewContract(&_PayAsYouGo.TransactOpts, rate, clientId)
}

// NewContract is a paid mutator transaction binding the contract method 0x05c23118.
//
// Solidity: function newContract(uint8 rate, string clientId) returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) NewContract(rate uint8, clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.NewContract(&_PayAsYouGo.TransactOpts, rate, clientId)
}

// PayBill is a paid mutator transaction binding the contract method 0xf0d65ec5.
//
// Solidity: function payBill(bytes32 contractId) payable returns()
func (_PayAsYouGo *PayAsYouGoTransactor) PayBill(opts *bind.TransactOpts, contractId [32]byte) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "payBill", contractId)
}

// PayBill is a paid mutator transaction binding the contract method 0xf0d65ec5.
//
// Solidity: function payBill(bytes32 contractId) payable returns()
func (_PayAsYouGo *PayAsYouGoSession) PayBill(contractId [32]byte) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.PayBill(&_PayAsYouGo.TransactOpts, contractId)
}

// PayBill is a paid mutator transaction binding the contract method 0xf0d65ec5.
//
// Solidity: function payBill(bytes32 contractId) payable returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) PayBill(contractId [32]byte) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.PayBill(&_PayAsYouGo.TransactOpts, contractId)
}

// PayAsYouGoBillAddedIterator is returned from FilterBillAdded and is used to iterate over the raw logs and unpacked data for BillAdded events raised by the PayAsYouGo contract.
type PayAsYouGoBillAddedIterator struct {
	Event *PayAsYouGoBillAdded // Event containing the contract specifics and raw log

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
func (it *PayAsYouGoBillAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayAsYouGoBillAdded)
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
		it.Event = new(PayAsYouGoBillAdded)
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
func (it *PayAsYouGoBillAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayAsYouGoBillAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayAsYouGoBillAdded represents a BillAdded event raised by the PayAsYouGo contract.
type PayAsYouGoBillAdded struct {
	ContractId [32]byte
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBillAdded is a free log retrieval operation binding the contract event 0xbde0d7a98e165595b50c03cf55b0293b83d8f4b3223e8fe0edf0c1ac95ecf916.
//
// Solidity: event BillAdded(bytes32 contractId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) FilterBillAdded(opts *bind.FilterOpts) (*PayAsYouGoBillAddedIterator, error) {

	logs, sub, err := _PayAsYouGo.contract.FilterLogs(opts, "BillAdded")
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoBillAddedIterator{contract: _PayAsYouGo.contract, event: "BillAdded", logs: logs, sub: sub}, nil
}

// WatchBillAdded is a free log subscription operation binding the contract event 0xbde0d7a98e165595b50c03cf55b0293b83d8f4b3223e8fe0edf0c1ac95ecf916.
//
// Solidity: event BillAdded(bytes32 contractId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) WatchBillAdded(opts *bind.WatchOpts, sink chan<- *PayAsYouGoBillAdded) (event.Subscription, error) {

	logs, sub, err := _PayAsYouGo.contract.WatchLogs(opts, "BillAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayAsYouGoBillAdded)
				if err := _PayAsYouGo.contract.UnpackLog(event, "BillAdded", log); err != nil {
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

// ParseBillAdded is a log parse operation binding the contract event 0xbde0d7a98e165595b50c03cf55b0293b83d8f4b3223e8fe0edf0c1ac95ecf916.
//
// Solidity: event BillAdded(bytes32 contractId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) ParseBillAdded(log types.Log) (*PayAsYouGoBillAdded, error) {
	event := new(PayAsYouGoBillAdded)
	if err := _PayAsYouGo.contract.UnpackLog(event, "BillAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayAsYouGoBillPaidIterator is returned from FilterBillPaid and is used to iterate over the raw logs and unpacked data for BillPaid events raised by the PayAsYouGo contract.
type PayAsYouGoBillPaidIterator struct {
	Event *PayAsYouGoBillPaid // Event containing the contract specifics and raw log

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
func (it *PayAsYouGoBillPaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayAsYouGoBillPaid)
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
		it.Event = new(PayAsYouGoBillPaid)
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
func (it *PayAsYouGoBillPaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayAsYouGoBillPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayAsYouGoBillPaid represents a BillPaid event raised by the PayAsYouGo contract.
type PayAsYouGoBillPaid struct {
	ContractId [32]byte
	Amount     *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBillPaid is a free log retrieval operation binding the contract event 0xa5b2fe6b6b3708f6c4805b27758d1745626d759eea73c3572b0da8b3a96db97b.
//
// Solidity: event BillPaid(bytes32 contractId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) FilterBillPaid(opts *bind.FilterOpts) (*PayAsYouGoBillPaidIterator, error) {

	logs, sub, err := _PayAsYouGo.contract.FilterLogs(opts, "BillPaid")
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoBillPaidIterator{contract: _PayAsYouGo.contract, event: "BillPaid", logs: logs, sub: sub}, nil
}

// WatchBillPaid is a free log subscription operation binding the contract event 0xa5b2fe6b6b3708f6c4805b27758d1745626d759eea73c3572b0da8b3a96db97b.
//
// Solidity: event BillPaid(bytes32 contractId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) WatchBillPaid(opts *bind.WatchOpts, sink chan<- *PayAsYouGoBillPaid) (event.Subscription, error) {

	logs, sub, err := _PayAsYouGo.contract.WatchLogs(opts, "BillPaid")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayAsYouGoBillPaid)
				if err := _PayAsYouGo.contract.UnpackLog(event, "BillPaid", log); err != nil {
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

// ParseBillPaid is a log parse operation binding the contract event 0xa5b2fe6b6b3708f6c4805b27758d1745626d759eea73c3572b0da8b3a96db97b.
//
// Solidity: event BillPaid(bytes32 contractId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) ParseBillPaid(log types.Log) (*PayAsYouGoBillPaid, error) {
	event := new(PayAsYouGoBillPaid)
	if err := _PayAsYouGo.contract.UnpackLog(event, "BillPaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayAsYouGoContractCreatedIterator is returned from FilterContractCreated and is used to iterate over the raw logs and unpacked data for ContractCreated events raised by the PayAsYouGo contract.
type PayAsYouGoContractCreatedIterator struct {
	Event *PayAsYouGoContractCreated // Event containing the contract specifics and raw log

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
func (it *PayAsYouGoContractCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayAsYouGoContractCreated)
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
		it.Event = new(PayAsYouGoContractCreated)
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
func (it *PayAsYouGoContractCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayAsYouGoContractCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayAsYouGoContractCreated represents a ContractCreated event raised by the PayAsYouGo contract.
type PayAsYouGoContractCreated struct {
	ContractId [32]byte
	ClientId   string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterContractCreated is a free log retrieval operation binding the contract event 0x102ef7a35251c7335cc8d842cacf8ddfa4e69191f4908aeab4327976ada485a8.
//
// Solidity: event ContractCreated(bytes32 contractId, string clientId)
func (_PayAsYouGo *PayAsYouGoFilterer) FilterContractCreated(opts *bind.FilterOpts) (*PayAsYouGoContractCreatedIterator, error) {

	logs, sub, err := _PayAsYouGo.contract.FilterLogs(opts, "ContractCreated")
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoContractCreatedIterator{contract: _PayAsYouGo.contract, event: "ContractCreated", logs: logs, sub: sub}, nil
}

// WatchContractCreated is a free log subscription operation binding the contract event 0x102ef7a35251c7335cc8d842cacf8ddfa4e69191f4908aeab4327976ada485a8.
//
// Solidity: event ContractCreated(bytes32 contractId, string clientId)
func (_PayAsYouGo *PayAsYouGoFilterer) WatchContractCreated(opts *bind.WatchOpts, sink chan<- *PayAsYouGoContractCreated) (event.Subscription, error) {

	logs, sub, err := _PayAsYouGo.contract.WatchLogs(opts, "ContractCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayAsYouGoContractCreated)
				if err := _PayAsYouGo.contract.UnpackLog(event, "ContractCreated", log); err != nil {
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

// ParseContractCreated is a log parse operation binding the contract event 0x102ef7a35251c7335cc8d842cacf8ddfa4e69191f4908aeab4327976ada485a8.
//
// Solidity: event ContractCreated(bytes32 contractId, string clientId)
func (_PayAsYouGo *PayAsYouGoFilterer) ParseContractCreated(log types.Log) (*PayAsYouGoContractCreated, error) {
	event := new(PayAsYouGoContractCreated)
	if err := _PayAsYouGo.contract.UnpackLog(event, "ContractCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
