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

// PayAsYouGoBillingInput is an auto generated low-level Go binding around an user-defined struct.
type PayAsYouGoBillingInput struct {
	ClientId  string
	Amount    uint64
	Timestamp uint64
}

// PayAsYouGoClient is an auto generated low-level Go binding around an user-defined struct.
type PayAsYouGoClient struct {
	ClientId           string
	UnpaidBill         *big.Int
	LastUsageFetchTime uint64
	Rate               uint64
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
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_transactionAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BillAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"BillPaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"}],\"name\":\"ClientCreated\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"addBillToClient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"amount\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"internalType\":\"structPayAsYouGo.BillingInput[]\",\"name\":\"billings\",\"type\":\"tuple[]\"}],\"name\":\"bulkAddBillToClient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"rate\",\"type\":\"uint64\"}],\"name\":\"changeAllClientRates\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint64\",\"name\":\"rate\",\"type\":\"uint64\"}],\"name\":\"changeRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_transactionAddress\",\"type\":\"address\"}],\"name\":\"changeTransactionAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"clientIDs\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"name\":\"clients\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"unpaidBill\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"lastUsageFetchTime\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rate\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"}],\"name\":\"getClientById\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"unpaidBill\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"lastUsageFetchTime\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rate\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"enumPayAsYouGo.TransactionType\",\"name\":\"transactionType\",\"type\":\"uint8\"}],\"internalType\":\"structPayAsYouGo.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"internalType\":\"structPayAsYouGo.Client\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getClients\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"unpaidBill\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"lastUsageFetchTime\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rate\",\"type\":\"uint64\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"enumPayAsYouGo.TransactionType\",\"name\":\"transactionType\",\"type\":\"uint8\"}],\"internalType\":\"structPayAsYouGo.Transaction[]\",\"name\":\"transactions\",\"type\":\"tuple[]\"}],\"internalType\":\"structPayAsYouGo.Client[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"rate\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"}],\"name\":\"newClient\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"clientId\",\"type\":\"string\"}],\"name\":\"payBill\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
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

// ClientIDs is a free data retrieval call binding the contract method 0x6af5ec39.
//
// Solidity: function clientIDs(uint256 ) view returns(string)
func (_PayAsYouGo *PayAsYouGoCaller) ClientIDs(opts *bind.CallOpts, arg0 *big.Int) (string, error) {
	var out []interface{}
	err := _PayAsYouGo.contract.Call(opts, &out, "clientIDs", arg0)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ClientIDs is a free data retrieval call binding the contract method 0x6af5ec39.
//
// Solidity: function clientIDs(uint256 ) view returns(string)
func (_PayAsYouGo *PayAsYouGoSession) ClientIDs(arg0 *big.Int) (string, error) {
	return _PayAsYouGo.Contract.ClientIDs(&_PayAsYouGo.CallOpts, arg0)
}

// ClientIDs is a free data retrieval call binding the contract method 0x6af5ec39.
//
// Solidity: function clientIDs(uint256 ) view returns(string)
func (_PayAsYouGo *PayAsYouGoCallerSession) ClientIDs(arg0 *big.Int) (string, error) {
	return _PayAsYouGo.Contract.ClientIDs(&_PayAsYouGo.CallOpts, arg0)
}

// Clients is a free data retrieval call binding the contract method 0x20ba1e9f.
//
// Solidity: function clients(string ) view returns(string clientId, uint256 unpaidBill, uint64 lastUsageFetchTime, uint64 rate)
func (_PayAsYouGo *PayAsYouGoCaller) Clients(opts *bind.CallOpts, arg0 string) (struct {
	ClientId           string
	UnpaidBill         *big.Int
	LastUsageFetchTime uint64
	Rate               uint64
}, error) {
	var out []interface{}
	err := _PayAsYouGo.contract.Call(opts, &out, "clients", arg0)

	outstruct := new(struct {
		ClientId           string
		UnpaidBill         *big.Int
		LastUsageFetchTime uint64
		Rate               uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ClientId = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.UnpaidBill = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.LastUsageFetchTime = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.Rate = *abi.ConvertType(out[3], new(uint64)).(*uint64)

	return *outstruct, err

}

// Clients is a free data retrieval call binding the contract method 0x20ba1e9f.
//
// Solidity: function clients(string ) view returns(string clientId, uint256 unpaidBill, uint64 lastUsageFetchTime, uint64 rate)
func (_PayAsYouGo *PayAsYouGoSession) Clients(arg0 string) (struct {
	ClientId           string
	UnpaidBill         *big.Int
	LastUsageFetchTime uint64
	Rate               uint64
}, error) {
	return _PayAsYouGo.Contract.Clients(&_PayAsYouGo.CallOpts, arg0)
}

// Clients is a free data retrieval call binding the contract method 0x20ba1e9f.
//
// Solidity: function clients(string ) view returns(string clientId, uint256 unpaidBill, uint64 lastUsageFetchTime, uint64 rate)
func (_PayAsYouGo *PayAsYouGoCallerSession) Clients(arg0 string) (struct {
	ClientId           string
	UnpaidBill         *big.Int
	LastUsageFetchTime uint64
	Rate               uint64
}, error) {
	return _PayAsYouGo.Contract.Clients(&_PayAsYouGo.CallOpts, arg0)
}

// GetClientById is a free data retrieval call binding the contract method 0x4c4d7c5a.
//
// Solidity: function getClientById(string clientId) view returns((string,uint256,uint64,uint64,(uint256,uint64,uint8)[]))
func (_PayAsYouGo *PayAsYouGoCaller) GetClientById(opts *bind.CallOpts, clientId string) (PayAsYouGoClient, error) {
	var out []interface{}
	err := _PayAsYouGo.contract.Call(opts, &out, "getClientById", clientId)

	if err != nil {
		return *new(PayAsYouGoClient), err
	}

	out0 := *abi.ConvertType(out[0], new(PayAsYouGoClient)).(*PayAsYouGoClient)

	return out0, err

}

// GetClientById is a free data retrieval call binding the contract method 0x4c4d7c5a.
//
// Solidity: function getClientById(string clientId) view returns((string,uint256,uint64,uint64,(uint256,uint64,uint8)[]))
func (_PayAsYouGo *PayAsYouGoSession) GetClientById(clientId string) (PayAsYouGoClient, error) {
	return _PayAsYouGo.Contract.GetClientById(&_PayAsYouGo.CallOpts, clientId)
}

// GetClientById is a free data retrieval call binding the contract method 0x4c4d7c5a.
//
// Solidity: function getClientById(string clientId) view returns((string,uint256,uint64,uint64,(uint256,uint64,uint8)[]))
func (_PayAsYouGo *PayAsYouGoCallerSession) GetClientById(clientId string) (PayAsYouGoClient, error) {
	return _PayAsYouGo.Contract.GetClientById(&_PayAsYouGo.CallOpts, clientId)
}

// GetClients is a free data retrieval call binding the contract method 0x6db80995.
//
// Solidity: function getClients() view returns((string,uint256,uint64,uint64,(uint256,uint64,uint8)[])[])
func (_PayAsYouGo *PayAsYouGoCaller) GetClients(opts *bind.CallOpts) ([]PayAsYouGoClient, error) {
	var out []interface{}
	err := _PayAsYouGo.contract.Call(opts, &out, "getClients")

	if err != nil {
		return *new([]PayAsYouGoClient), err
	}

	out0 := *abi.ConvertType(out[0], new([]PayAsYouGoClient)).(*[]PayAsYouGoClient)

	return out0, err

}

// GetClients is a free data retrieval call binding the contract method 0x6db80995.
//
// Solidity: function getClients() view returns((string,uint256,uint64,uint64,(uint256,uint64,uint8)[])[])
func (_PayAsYouGo *PayAsYouGoSession) GetClients() ([]PayAsYouGoClient, error) {
	return _PayAsYouGo.Contract.GetClients(&_PayAsYouGo.CallOpts)
}

// GetClients is a free data retrieval call binding the contract method 0x6db80995.
//
// Solidity: function getClients() view returns((string,uint256,uint64,uint64,(uint256,uint64,uint8)[])[])
func (_PayAsYouGo *PayAsYouGoCallerSession) GetClients() ([]PayAsYouGoClient, error) {
	return _PayAsYouGo.Contract.GetClients(&_PayAsYouGo.CallOpts)
}

// AddBillToClient is a paid mutator transaction binding the contract method 0xf0a537e4.
//
// Solidity: function addBillToClient(string clientId, uint64 amount, uint64 timestamp) returns()
func (_PayAsYouGo *PayAsYouGoTransactor) AddBillToClient(opts *bind.TransactOpts, clientId string, amount uint64, timestamp uint64) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "addBillToClient", clientId, amount, timestamp)
}

// AddBillToClient is a paid mutator transaction binding the contract method 0xf0a537e4.
//
// Solidity: function addBillToClient(string clientId, uint64 amount, uint64 timestamp) returns()
func (_PayAsYouGo *PayAsYouGoSession) AddBillToClient(clientId string, amount uint64, timestamp uint64) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.AddBillToClient(&_PayAsYouGo.TransactOpts, clientId, amount, timestamp)
}

// AddBillToClient is a paid mutator transaction binding the contract method 0xf0a537e4.
//
// Solidity: function addBillToClient(string clientId, uint64 amount, uint64 timestamp) returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) AddBillToClient(clientId string, amount uint64, timestamp uint64) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.AddBillToClient(&_PayAsYouGo.TransactOpts, clientId, amount, timestamp)
}

// BulkAddBillToClient is a paid mutator transaction binding the contract method 0x09b1f57a.
//
// Solidity: function bulkAddBillToClient((string,uint64,uint64)[] billings) returns()
func (_PayAsYouGo *PayAsYouGoTransactor) BulkAddBillToClient(opts *bind.TransactOpts, billings []PayAsYouGoBillingInput) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "bulkAddBillToClient", billings)
}

// BulkAddBillToClient is a paid mutator transaction binding the contract method 0x09b1f57a.
//
// Solidity: function bulkAddBillToClient((string,uint64,uint64)[] billings) returns()
func (_PayAsYouGo *PayAsYouGoSession) BulkAddBillToClient(billings []PayAsYouGoBillingInput) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.BulkAddBillToClient(&_PayAsYouGo.TransactOpts, billings)
}

// BulkAddBillToClient is a paid mutator transaction binding the contract method 0x09b1f57a.
//
// Solidity: function bulkAddBillToClient((string,uint64,uint64)[] billings) returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) BulkAddBillToClient(billings []PayAsYouGoBillingInput) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.BulkAddBillToClient(&_PayAsYouGo.TransactOpts, billings)
}

// ChangeAllClientRates is a paid mutator transaction binding the contract method 0x797bad26.
//
// Solidity: function changeAllClientRates(uint64 rate) returns()
func (_PayAsYouGo *PayAsYouGoTransactor) ChangeAllClientRates(opts *bind.TransactOpts, rate uint64) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "changeAllClientRates", rate)
}

// ChangeAllClientRates is a paid mutator transaction binding the contract method 0x797bad26.
//
// Solidity: function changeAllClientRates(uint64 rate) returns()
func (_PayAsYouGo *PayAsYouGoSession) ChangeAllClientRates(rate uint64) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.ChangeAllClientRates(&_PayAsYouGo.TransactOpts, rate)
}

// ChangeAllClientRates is a paid mutator transaction binding the contract method 0x797bad26.
//
// Solidity: function changeAllClientRates(uint64 rate) returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) ChangeAllClientRates(rate uint64) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.ChangeAllClientRates(&_PayAsYouGo.TransactOpts, rate)
}

// ChangeRate is a paid mutator transaction binding the contract method 0xa87d7ba7.
//
// Solidity: function changeRate(string clientId, uint64 rate) returns()
func (_PayAsYouGo *PayAsYouGoTransactor) ChangeRate(opts *bind.TransactOpts, clientId string, rate uint64) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "changeRate", clientId, rate)
}

// ChangeRate is a paid mutator transaction binding the contract method 0xa87d7ba7.
//
// Solidity: function changeRate(string clientId, uint64 rate) returns()
func (_PayAsYouGo *PayAsYouGoSession) ChangeRate(clientId string, rate uint64) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.ChangeRate(&_PayAsYouGo.TransactOpts, clientId, rate)
}

// ChangeRate is a paid mutator transaction binding the contract method 0xa87d7ba7.
//
// Solidity: function changeRate(string clientId, uint64 rate) returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) ChangeRate(clientId string, rate uint64) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.ChangeRate(&_PayAsYouGo.TransactOpts, clientId, rate)
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

// NewClient is a paid mutator transaction binding the contract method 0x03900e36.
//
// Solidity: function newClient(uint64 rate, string clientId) returns()
func (_PayAsYouGo *PayAsYouGoTransactor) NewClient(opts *bind.TransactOpts, rate uint64, clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "newClient", rate, clientId)
}

// NewClient is a paid mutator transaction binding the contract method 0x03900e36.
//
// Solidity: function newClient(uint64 rate, string clientId) returns()
func (_PayAsYouGo *PayAsYouGoSession) NewClient(rate uint64, clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.NewClient(&_PayAsYouGo.TransactOpts, rate, clientId)
}

// NewClient is a paid mutator transaction binding the contract method 0x03900e36.
//
// Solidity: function newClient(uint64 rate, string clientId) returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) NewClient(rate uint64, clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.NewClient(&_PayAsYouGo.TransactOpts, rate, clientId)
}

// PayBill is a paid mutator transaction binding the contract method 0x070abcc2.
//
// Solidity: function payBill(string clientId) payable returns()
func (_PayAsYouGo *PayAsYouGoTransactor) PayBill(opts *bind.TransactOpts, clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.contract.Transact(opts, "payBill", clientId)
}

// PayBill is a paid mutator transaction binding the contract method 0x070abcc2.
//
// Solidity: function payBill(string clientId) payable returns()
func (_PayAsYouGo *PayAsYouGoSession) PayBill(clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.PayBill(&_PayAsYouGo.TransactOpts, clientId)
}

// PayBill is a paid mutator transaction binding the contract method 0x070abcc2.
//
// Solidity: function payBill(string clientId) payable returns()
func (_PayAsYouGo *PayAsYouGoTransactorSession) PayBill(clientId string) (*types.Transaction, error) {
	return _PayAsYouGo.Contract.PayBill(&_PayAsYouGo.TransactOpts, clientId)
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
	ClientId string
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBillAdded is a free log retrieval operation binding the contract event 0x3d43d2c72e75437af21df25cfbf67c625914003c91122e46b6e3396434c1f7c2.
//
// Solidity: event BillAdded(string clientId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) FilterBillAdded(opts *bind.FilterOpts) (*PayAsYouGoBillAddedIterator, error) {

	logs, sub, err := _PayAsYouGo.contract.FilterLogs(opts, "BillAdded")
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoBillAddedIterator{contract: _PayAsYouGo.contract, event: "BillAdded", logs: logs, sub: sub}, nil
}

// WatchBillAdded is a free log subscription operation binding the contract event 0x3d43d2c72e75437af21df25cfbf67c625914003c91122e46b6e3396434c1f7c2.
//
// Solidity: event BillAdded(string clientId, uint256 amount)
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

// ParseBillAdded is a log parse operation binding the contract event 0x3d43d2c72e75437af21df25cfbf67c625914003c91122e46b6e3396434c1f7c2.
//
// Solidity: event BillAdded(string clientId, uint256 amount)
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
	ClientId string
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterBillPaid is a free log retrieval operation binding the contract event 0x85efa8140d8de83ee16bcc1d3f231311f1eaa6dacc4868130fc30689c5ab0b66.
//
// Solidity: event BillPaid(string clientId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) FilterBillPaid(opts *bind.FilterOpts) (*PayAsYouGoBillPaidIterator, error) {

	logs, sub, err := _PayAsYouGo.contract.FilterLogs(opts, "BillPaid")
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoBillPaidIterator{contract: _PayAsYouGo.contract, event: "BillPaid", logs: logs, sub: sub}, nil
}

// WatchBillPaid is a free log subscription operation binding the contract event 0x85efa8140d8de83ee16bcc1d3f231311f1eaa6dacc4868130fc30689c5ab0b66.
//
// Solidity: event BillPaid(string clientId, uint256 amount)
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

// ParseBillPaid is a log parse operation binding the contract event 0x85efa8140d8de83ee16bcc1d3f231311f1eaa6dacc4868130fc30689c5ab0b66.
//
// Solidity: event BillPaid(string clientId, uint256 amount)
func (_PayAsYouGo *PayAsYouGoFilterer) ParseBillPaid(log types.Log) (*PayAsYouGoBillPaid, error) {
	event := new(PayAsYouGoBillPaid)
	if err := _PayAsYouGo.contract.UnpackLog(event, "BillPaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PayAsYouGoClientCreatedIterator is returned from FilterClientCreated and is used to iterate over the raw logs and unpacked data for ClientCreated events raised by the PayAsYouGo contract.
type PayAsYouGoClientCreatedIterator struct {
	Event *PayAsYouGoClientCreated // Event containing the contract specifics and raw log

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
func (it *PayAsYouGoClientCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PayAsYouGoClientCreated)
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
		it.Event = new(PayAsYouGoClientCreated)
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
func (it *PayAsYouGoClientCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PayAsYouGoClientCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PayAsYouGoClientCreated represents a ClientCreated event raised by the PayAsYouGo contract.
type PayAsYouGoClientCreated struct {
	ClientId string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterClientCreated is a free log retrieval operation binding the contract event 0xeb98df470d17266538e4ee034952206621fad8d86ca38b090e92f64589108482.
//
// Solidity: event ClientCreated(string clientId)
func (_PayAsYouGo *PayAsYouGoFilterer) FilterClientCreated(opts *bind.FilterOpts) (*PayAsYouGoClientCreatedIterator, error) {

	logs, sub, err := _PayAsYouGo.contract.FilterLogs(opts, "ClientCreated")
	if err != nil {
		return nil, err
	}
	return &PayAsYouGoClientCreatedIterator{contract: _PayAsYouGo.contract, event: "ClientCreated", logs: logs, sub: sub}, nil
}

// WatchClientCreated is a free log subscription operation binding the contract event 0xeb98df470d17266538e4ee034952206621fad8d86ca38b090e92f64589108482.
//
// Solidity: event ClientCreated(string clientId)
func (_PayAsYouGo *PayAsYouGoFilterer) WatchClientCreated(opts *bind.WatchOpts, sink chan<- *PayAsYouGoClientCreated) (event.Subscription, error) {

	logs, sub, err := _PayAsYouGo.contract.WatchLogs(opts, "ClientCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PayAsYouGoClientCreated)
				if err := _PayAsYouGo.contract.UnpackLog(event, "ClientCreated", log); err != nil {
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

// ParseClientCreated is a log parse operation binding the contract event 0xeb98df470d17266538e4ee034952206621fad8d86ca38b090e92f64589108482.
//
// Solidity: event ClientCreated(string clientId)
func (_PayAsYouGo *PayAsYouGoFilterer) ParseClientCreated(log types.Log) (*PayAsYouGoClientCreated, error) {
	event := new(PayAsYouGoClientCreated)
	if err := _PayAsYouGo.contract.UnpackLog(event, "ClientCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
