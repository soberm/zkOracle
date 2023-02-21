// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zkOracle

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ZKOracleAccount is an auto generated low-level Go binding around an user-defined struct.
type ZKOracleAccount struct {
	Index   *big.Int
	PubKey  ZKOraclePublicKey
	Balance *big.Int
}

// ZKOraclePublicKey is an auto generated low-level Go binding around an user-defined struct.
type ZKOraclePublicKey struct {
	X *big.Int
	Y *big.Int
}

// ZKOracleContractABI is the input ABI used to generate the binding from.
const ZKOracleContractABI = "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"merkleTreeAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_seedX\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_seedY\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"request\",\"type\":\"uint256\"}],\"name\":\"BlockRequested\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"submitter\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"validators\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"request\",\"type\":\"uint256\"}],\"name\":\"BlockSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"Exiting\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"indexed\":false,\"internalType\":\"structZKOracle.PublicKey\",\"name\":\"pubkey\",\"type\":\"tuple\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Registered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"replaced\",\"type\":\"address\"}],\"name\":\"Replaced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"Withdrawn\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.PublicKey\",\"name\":\"pubKey\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.Account\",\"name\":\"account\",\"type\":\"tuple\"},{\"internalType\":\"uint256[]\",\"name\":\"path\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"helper\",\"type\":\"uint256[]\"}],\"name\":\"exit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"exitDelay\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAggregator\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"number\",\"type\":\"uint256\"}],\"name\":\"getBlockByNumber\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getExitTime\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getIPAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.PublicKey\",\"name\":\"pubKey\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.Account\",\"name\":\"account\",\"type\":\"tuple\"}],\"name\":\"hashAccount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.PublicKey\",\"name\":\"publicKey\",\"type\":\"tuple\"},{\"internalType\":\"string\",\"name\":\"ip\",\"type\":\"string\"}],\"name\":\"register\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.PublicKey\",\"name\":\"publicKey\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.PublicKey\",\"name\":\"pubKey\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.Account\",\"name\":\"toReplace\",\"type\":\"tuple\"},{\"internalType\":\"uint256[]\",\"name\":\"path\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"helper\",\"type\":\"uint256[]\"}],\"name\":\"replace\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"request\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validators\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"blockHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"postStateRoot\",\"type\":\"uint256\"},{\"internalType\":\"uint256[2]\",\"name\":\"a\",\"type\":\"uint256[2]\"},{\"internalType\":\"uint256[2][2]\",\"name\":\"b\",\"type\":\"uint256[2][2]\"},{\"internalType\":\"uint256[2]\",\"name\":\"c\",\"type\":\"uint256[2]\"}],\"name\":\"submitBlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"x\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"y\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.PublicKey\",\"name\":\"pubKey\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"internalType\":\"structZKOracle.Account\",\"name\":\"account\",\"type\":\"tuple\"},{\"internalType\":\"uint256[]\",\"name\":\"path\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"helper\",\"type\":\"uint256[]\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// ZKOracleContract is an auto generated Go binding around an Ethereum contract.
type ZKOracleContract struct {
	ZKOracleContractCaller     // Read-only binding to the contract
	ZKOracleContractTransactor // Write-only binding to the contract
	ZKOracleContractFilterer   // Log filterer for contract events
}

// ZKOracleContractCaller is an auto generated read-only Go binding around an Ethereum contract.
type ZKOracleContractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKOracleContractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ZKOracleContractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKOracleContractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ZKOracleContractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ZKOracleContractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ZKOracleContractSession struct {
	Contract     *ZKOracleContract // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ZKOracleContractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ZKOracleContractCallerSession struct {
	Contract *ZKOracleContractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ZKOracleContractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ZKOracleContractTransactorSession struct {
	Contract     *ZKOracleContractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ZKOracleContractRaw is an auto generated low-level Go binding around an Ethereum contract.
type ZKOracleContractRaw struct {
	Contract *ZKOracleContract // Generic contract binding to access the raw methods on
}

// ZKOracleContractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ZKOracleContractCallerRaw struct {
	Contract *ZKOracleContractCaller // Generic read-only contract binding to access the raw methods on
}

// ZKOracleContractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ZKOracleContractTransactorRaw struct {
	Contract *ZKOracleContractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewZKOracleContract creates a new instance of ZKOracleContract, bound to a specific deployed contract.
func NewZKOracleContract(address common.Address, backend bind.ContractBackend) (*ZKOracleContract, error) {
	contract, err := bindZKOracleContract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ZKOracleContract{ZKOracleContractCaller: ZKOracleContractCaller{contract: contract}, ZKOracleContractTransactor: ZKOracleContractTransactor{contract: contract}, ZKOracleContractFilterer: ZKOracleContractFilterer{contract: contract}}, nil
}

// NewZKOracleContractCaller creates a new read-only instance of ZKOracleContract, bound to a specific deployed contract.
func NewZKOracleContractCaller(address common.Address, caller bind.ContractCaller) (*ZKOracleContractCaller, error) {
	contract, err := bindZKOracleContract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractCaller{contract: contract}, nil
}

// NewZKOracleContractTransactor creates a new write-only instance of ZKOracleContract, bound to a specific deployed contract.
func NewZKOracleContractTransactor(address common.Address, transactor bind.ContractTransactor) (*ZKOracleContractTransactor, error) {
	contract, err := bindZKOracleContract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractTransactor{contract: contract}, nil
}

// NewZKOracleContractFilterer creates a new log filterer instance of ZKOracleContract, bound to a specific deployed contract.
func NewZKOracleContractFilterer(address common.Address, filterer bind.ContractFilterer) (*ZKOracleContractFilterer, error) {
	contract, err := bindZKOracleContract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractFilterer{contract: contract}, nil
}

// bindZKOracleContract binds a generic wrapper to an already deployed contract.
func bindZKOracleContract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ZKOracleContractABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZKOracleContract *ZKOracleContractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZKOracleContract.Contract.ZKOracleContractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZKOracleContract *ZKOracleContractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.ZKOracleContractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZKOracleContract *ZKOracleContractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.ZKOracleContractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ZKOracleContract *ZKOracleContractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ZKOracleContract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ZKOracleContract *ZKOracleContractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ZKOracleContract *ZKOracleContractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.contract.Transact(opts, method, params...)
}

// ExitDelay is a free data retrieval call binding the contract method 0x28388630.
//
// Solidity: function exitDelay() view returns(uint256)
func (_ZKOracleContract *ZKOracleContractCaller) ExitDelay(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZKOracleContract.contract.Call(opts, &out, "exitDelay")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ExitDelay is a free data retrieval call binding the contract method 0x28388630.
//
// Solidity: function exitDelay() view returns(uint256)
func (_ZKOracleContract *ZKOracleContractSession) ExitDelay() (*big.Int, error) {
	return _ZKOracleContract.Contract.ExitDelay(&_ZKOracleContract.CallOpts)
}

// ExitDelay is a free data retrieval call binding the contract method 0x28388630.
//
// Solidity: function exitDelay() view returns(uint256)
func (_ZKOracleContract *ZKOracleContractCallerSession) ExitDelay() (*big.Int, error) {
	return _ZKOracleContract.Contract.ExitDelay(&_ZKOracleContract.CallOpts)
}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(uint256)
func (_ZKOracleContract *ZKOracleContractCaller) GetAggregator(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _ZKOracleContract.contract.Call(opts, &out, "getAggregator")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(uint256)
func (_ZKOracleContract *ZKOracleContractSession) GetAggregator() (*big.Int, error) {
	return _ZKOracleContract.Contract.GetAggregator(&_ZKOracleContract.CallOpts)
}

// GetAggregator is a free data retrieval call binding the contract method 0x3ad59dbc.
//
// Solidity: function getAggregator() view returns(uint256)
func (_ZKOracleContract *ZKOracleContractCallerSession) GetAggregator() (*big.Int, error) {
	return _ZKOracleContract.Contract.GetAggregator(&_ZKOracleContract.CallOpts)
}

// GetExitTime is a free data retrieval call binding the contract method 0x3c8c4e75.
//
// Solidity: function getExitTime(address addr) view returns(uint256)
func (_ZKOracleContract *ZKOracleContractCaller) GetExitTime(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ZKOracleContract.contract.Call(opts, &out, "getExitTime", addr)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetExitTime is a free data retrieval call binding the contract method 0x3c8c4e75.
//
// Solidity: function getExitTime(address addr) view returns(uint256)
func (_ZKOracleContract *ZKOracleContractSession) GetExitTime(addr common.Address) (*big.Int, error) {
	return _ZKOracleContract.Contract.GetExitTime(&_ZKOracleContract.CallOpts, addr)
}

// GetExitTime is a free data retrieval call binding the contract method 0x3c8c4e75.
//
// Solidity: function getExitTime(address addr) view returns(uint256)
func (_ZKOracleContract *ZKOracleContractCallerSession) GetExitTime(addr common.Address) (*big.Int, error) {
	return _ZKOracleContract.Contract.GetExitTime(&_ZKOracleContract.CallOpts, addr)
}

// GetIPAddress is a free data retrieval call binding the contract method 0xefca2af7.
//
// Solidity: function getIPAddress(uint256 index) view returns(string)
func (_ZKOracleContract *ZKOracleContractCaller) GetIPAddress(opts *bind.CallOpts, index *big.Int) (string, error) {
	var out []interface{}
	err := _ZKOracleContract.contract.Call(opts, &out, "getIPAddress", index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetIPAddress is a free data retrieval call binding the contract method 0xefca2af7.
//
// Solidity: function getIPAddress(uint256 index) view returns(string)
func (_ZKOracleContract *ZKOracleContractSession) GetIPAddress(index *big.Int) (string, error) {
	return _ZKOracleContract.Contract.GetIPAddress(&_ZKOracleContract.CallOpts, index)
}

// GetIPAddress is a free data retrieval call binding the contract method 0xefca2af7.
//
// Solidity: function getIPAddress(uint256 index) view returns(string)
func (_ZKOracleContract *ZKOracleContractCallerSession) GetIPAddress(index *big.Int) (string, error) {
	return _ZKOracleContract.Contract.GetIPAddress(&_ZKOracleContract.CallOpts, index)
}

// HashAccount is a free data retrieval call binding the contract method 0xea368cff.
//
// Solidity: function hashAccount((uint256,(uint256,uint256),uint256) account) view returns(uint256)
func (_ZKOracleContract *ZKOracleContractCaller) HashAccount(opts *bind.CallOpts, account ZKOracleAccount) (*big.Int, error) {
	var out []interface{}
	err := _ZKOracleContract.contract.Call(opts, &out, "hashAccount", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// HashAccount is a free data retrieval call binding the contract method 0xea368cff.
//
// Solidity: function hashAccount((uint256,(uint256,uint256),uint256) account) view returns(uint256)
func (_ZKOracleContract *ZKOracleContractSession) HashAccount(account ZKOracleAccount) (*big.Int, error) {
	return _ZKOracleContract.Contract.HashAccount(&_ZKOracleContract.CallOpts, account)
}

// HashAccount is a free data retrieval call binding the contract method 0xea368cff.
//
// Solidity: function hashAccount((uint256,(uint256,uint256),uint256) account) view returns(uint256)
func (_ZKOracleContract *ZKOracleContractCallerSession) HashAccount(account ZKOracleAccount) (*big.Int, error) {
	return _ZKOracleContract.Contract.HashAccount(&_ZKOracleContract.CallOpts, account)
}

// Exit is a paid mutator transaction binding the contract method 0x3ea5f392.
//
// Solidity: function exit((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256[] helper) returns()
func (_ZKOracleContract *ZKOracleContractTransactor) Exit(opts *bind.TransactOpts, account ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.contract.Transact(opts, "exit", account, path, helper)
}

// Exit is a paid mutator transaction binding the contract method 0x3ea5f392.
//
// Solidity: function exit((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256[] helper) returns()
func (_ZKOracleContract *ZKOracleContractSession) Exit(account ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.Exit(&_ZKOracleContract.TransactOpts, account, path, helper)
}

// Exit is a paid mutator transaction binding the contract method 0x3ea5f392.
//
// Solidity: function exit((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256[] helper) returns()
func (_ZKOracleContract *ZKOracleContractTransactorSession) Exit(account ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.Exit(&_ZKOracleContract.TransactOpts, account, path, helper)
}

// GetBlockByNumber is a paid mutator transaction binding the contract method 0x7b31f51a.
//
// Solidity: function getBlockByNumber(uint256 number) payable returns()
func (_ZKOracleContract *ZKOracleContractTransactor) GetBlockByNumber(opts *bind.TransactOpts, number *big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.contract.Transact(opts, "getBlockByNumber", number)
}

// GetBlockByNumber is a paid mutator transaction binding the contract method 0x7b31f51a.
//
// Solidity: function getBlockByNumber(uint256 number) payable returns()
func (_ZKOracleContract *ZKOracleContractSession) GetBlockByNumber(number *big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.GetBlockByNumber(&_ZKOracleContract.TransactOpts, number)
}

// GetBlockByNumber is a paid mutator transaction binding the contract method 0x7b31f51a.
//
// Solidity: function getBlockByNumber(uint256 number) payable returns()
func (_ZKOracleContract *ZKOracleContractTransactorSession) GetBlockByNumber(number *big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.GetBlockByNumber(&_ZKOracleContract.TransactOpts, number)
}

// Register is a paid mutator transaction binding the contract method 0xb511e2dc.
//
// Solidity: function register((uint256,uint256) publicKey, string ip) payable returns()
func (_ZKOracleContract *ZKOracleContractTransactor) Register(opts *bind.TransactOpts, publicKey ZKOraclePublicKey, ip string) (*types.Transaction, error) {
	return _ZKOracleContract.contract.Transact(opts, "register", publicKey, ip)
}

// Register is a paid mutator transaction binding the contract method 0xb511e2dc.
//
// Solidity: function register((uint256,uint256) publicKey, string ip) payable returns()
func (_ZKOracleContract *ZKOracleContractSession) Register(publicKey ZKOraclePublicKey, ip string) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.Register(&_ZKOracleContract.TransactOpts, publicKey, ip)
}

// Register is a paid mutator transaction binding the contract method 0xb511e2dc.
//
// Solidity: function register((uint256,uint256) publicKey, string ip) payable returns()
func (_ZKOracleContract *ZKOracleContractTransactorSession) Register(publicKey ZKOraclePublicKey, ip string) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.Register(&_ZKOracleContract.TransactOpts, publicKey, ip)
}

// Replace is a paid mutator transaction binding the contract method 0xad917fd1.
//
// Solidity: function replace((uint256,uint256) publicKey, (uint256,(uint256,uint256),uint256) toReplace, uint256[] path, uint256[] helper) payable returns()
func (_ZKOracleContract *ZKOracleContractTransactor) Replace(opts *bind.TransactOpts, publicKey ZKOraclePublicKey, toReplace ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.contract.Transact(opts, "replace", publicKey, toReplace, path, helper)
}

// Replace is a paid mutator transaction binding the contract method 0xad917fd1.
//
// Solidity: function replace((uint256,uint256) publicKey, (uint256,(uint256,uint256),uint256) toReplace, uint256[] path, uint256[] helper) payable returns()
func (_ZKOracleContract *ZKOracleContractSession) Replace(publicKey ZKOraclePublicKey, toReplace ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.Replace(&_ZKOracleContract.TransactOpts, publicKey, toReplace, path, helper)
}

// Replace is a paid mutator transaction binding the contract method 0xad917fd1.
//
// Solidity: function replace((uint256,uint256) publicKey, (uint256,(uint256,uint256),uint256) toReplace, uint256[] path, uint256[] helper) payable returns()
func (_ZKOracleContract *ZKOracleContractTransactorSession) Replace(publicKey ZKOraclePublicKey, toReplace ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.Replace(&_ZKOracleContract.TransactOpts, publicKey, toReplace, path, helper)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xf9438379.
//
// Solidity: function submitBlock(uint256 index, uint256 request, uint256 validators, bytes32 blockHash, uint256 postStateRoot, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_ZKOracleContract *ZKOracleContractTransactor) SubmitBlock(opts *bind.TransactOpts, index *big.Int, request *big.Int, validators *big.Int, blockHash [32]byte, postStateRoot *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.contract.Transact(opts, "submitBlock", index, request, validators, blockHash, postStateRoot, a, b, c)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xf9438379.
//
// Solidity: function submitBlock(uint256 index, uint256 request, uint256 validators, bytes32 blockHash, uint256 postStateRoot, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_ZKOracleContract *ZKOracleContractSession) SubmitBlock(index *big.Int, request *big.Int, validators *big.Int, blockHash [32]byte, postStateRoot *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.SubmitBlock(&_ZKOracleContract.TransactOpts, index, request, validators, blockHash, postStateRoot, a, b, c)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xf9438379.
//
// Solidity: function submitBlock(uint256 index, uint256 request, uint256 validators, bytes32 blockHash, uint256 postStateRoot, uint256[2] a, uint256[2][2] b, uint256[2] c) returns()
func (_ZKOracleContract *ZKOracleContractTransactorSession) SubmitBlock(index *big.Int, request *big.Int, validators *big.Int, blockHash [32]byte, postStateRoot *big.Int, a [2]*big.Int, b [2][2]*big.Int, c [2]*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.SubmitBlock(&_ZKOracleContract.TransactOpts, index, request, validators, blockHash, postStateRoot, a, b, c)
}

// Withdraw is a paid mutator transaction binding the contract method 0x4175f9f9.
//
// Solidity: function withdraw((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256[] helper) returns()
func (_ZKOracleContract *ZKOracleContractTransactor) Withdraw(opts *bind.TransactOpts, account ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.contract.Transact(opts, "withdraw", account, path, helper)
}

// Withdraw is a paid mutator transaction binding the contract method 0x4175f9f9.
//
// Solidity: function withdraw((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256[] helper) returns()
func (_ZKOracleContract *ZKOracleContractSession) Withdraw(account ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.Withdraw(&_ZKOracleContract.TransactOpts, account, path, helper)
}

// Withdraw is a paid mutator transaction binding the contract method 0x4175f9f9.
//
// Solidity: function withdraw((uint256,(uint256,uint256),uint256) account, uint256[] path, uint256[] helper) returns()
func (_ZKOracleContract *ZKOracleContractTransactorSession) Withdraw(account ZKOracleAccount, path []*big.Int, helper []*big.Int) (*types.Transaction, error) {
	return _ZKOracleContract.Contract.Withdraw(&_ZKOracleContract.TransactOpts, account, path, helper)
}

// ZKOracleContractBlockRequestedIterator is returned from FilterBlockRequested and is used to iterate over the raw logs and unpacked data for BlockRequested events raised by the ZKOracleContract contract.
type ZKOracleContractBlockRequestedIterator struct {
	Event *ZKOracleContractBlockRequested // Event containing the contract specifics and raw log

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
func (it *ZKOracleContractBlockRequestedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKOracleContractBlockRequested)
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
		it.Event = new(ZKOracleContractBlockRequested)
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
func (it *ZKOracleContractBlockRequestedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKOracleContractBlockRequestedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKOracleContractBlockRequested represents a BlockRequested event raised by the ZKOracleContract contract.
type ZKOracleContractBlockRequested struct {
	Number  *big.Int
	Request *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBlockRequested is a free log retrieval operation binding the contract event 0x91ad2fc61c3d1a0988fc84b3e15d8a97480fadfe6743526d7bc0116257e842e9.
//
// Solidity: event BlockRequested(uint256 number, uint256 request)
func (_ZKOracleContract *ZKOracleContractFilterer) FilterBlockRequested(opts *bind.FilterOpts) (*ZKOracleContractBlockRequestedIterator, error) {

	logs, sub, err := _ZKOracleContract.contract.FilterLogs(opts, "BlockRequested")
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractBlockRequestedIterator{contract: _ZKOracleContract.contract, event: "BlockRequested", logs: logs, sub: sub}, nil
}

// WatchBlockRequested is a free log subscription operation binding the contract event 0x91ad2fc61c3d1a0988fc84b3e15d8a97480fadfe6743526d7bc0116257e842e9.
//
// Solidity: event BlockRequested(uint256 number, uint256 request)
func (_ZKOracleContract *ZKOracleContractFilterer) WatchBlockRequested(opts *bind.WatchOpts, sink chan<- *ZKOracleContractBlockRequested) (event.Subscription, error) {

	logs, sub, err := _ZKOracleContract.contract.WatchLogs(opts, "BlockRequested")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKOracleContractBlockRequested)
				if err := _ZKOracleContract.contract.UnpackLog(event, "BlockRequested", log); err != nil {
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

// ParseBlockRequested is a log parse operation binding the contract event 0x91ad2fc61c3d1a0988fc84b3e15d8a97480fadfe6743526d7bc0116257e842e9.
//
// Solidity: event BlockRequested(uint256 number, uint256 request)
func (_ZKOracleContract *ZKOracleContractFilterer) ParseBlockRequested(log types.Log) (*ZKOracleContractBlockRequested, error) {
	event := new(ZKOracleContractBlockRequested)
	if err := _ZKOracleContract.contract.UnpackLog(event, "BlockRequested", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKOracleContractBlockSubmittedIterator is returned from FilterBlockSubmitted and is used to iterate over the raw logs and unpacked data for BlockSubmitted events raised by the ZKOracleContract contract.
type ZKOracleContractBlockSubmittedIterator struct {
	Event *ZKOracleContractBlockSubmitted // Event containing the contract specifics and raw log

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
func (it *ZKOracleContractBlockSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKOracleContractBlockSubmitted)
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
		it.Event = new(ZKOracleContractBlockSubmitted)
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
func (it *ZKOracleContractBlockSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKOracleContractBlockSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKOracleContractBlockSubmitted represents a BlockSubmitted event raised by the ZKOracleContract contract.
type ZKOracleContractBlockSubmitted struct {
	Submitter  *big.Int
	Validators *big.Int
	Request    *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBlockSubmitted is a free log retrieval operation binding the contract event 0x7ccdeef917218144c4de8fcf67721096b04b8e55277a839612ff53d6e186cb5e.
//
// Solidity: event BlockSubmitted(uint256 submitter, uint256 validators, uint256 request)
func (_ZKOracleContract *ZKOracleContractFilterer) FilterBlockSubmitted(opts *bind.FilterOpts) (*ZKOracleContractBlockSubmittedIterator, error) {

	logs, sub, err := _ZKOracleContract.contract.FilterLogs(opts, "BlockSubmitted")
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractBlockSubmittedIterator{contract: _ZKOracleContract.contract, event: "BlockSubmitted", logs: logs, sub: sub}, nil
}

// WatchBlockSubmitted is a free log subscription operation binding the contract event 0x7ccdeef917218144c4de8fcf67721096b04b8e55277a839612ff53d6e186cb5e.
//
// Solidity: event BlockSubmitted(uint256 submitter, uint256 validators, uint256 request)
func (_ZKOracleContract *ZKOracleContractFilterer) WatchBlockSubmitted(opts *bind.WatchOpts, sink chan<- *ZKOracleContractBlockSubmitted) (event.Subscription, error) {

	logs, sub, err := _ZKOracleContract.contract.WatchLogs(opts, "BlockSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKOracleContractBlockSubmitted)
				if err := _ZKOracleContract.contract.UnpackLog(event, "BlockSubmitted", log); err != nil {
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

// ParseBlockSubmitted is a log parse operation binding the contract event 0x7ccdeef917218144c4de8fcf67721096b04b8e55277a839612ff53d6e186cb5e.
//
// Solidity: event BlockSubmitted(uint256 submitter, uint256 validators, uint256 request)
func (_ZKOracleContract *ZKOracleContractFilterer) ParseBlockSubmitted(log types.Log) (*ZKOracleContractBlockSubmitted, error) {
	event := new(ZKOracleContractBlockSubmitted)
	if err := _ZKOracleContract.contract.UnpackLog(event, "BlockSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKOracleContractExitingIterator is returned from FilterExiting and is used to iterate over the raw logs and unpacked data for Exiting events raised by the ZKOracleContract contract.
type ZKOracleContractExitingIterator struct {
	Event *ZKOracleContractExiting // Event containing the contract specifics and raw log

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
func (it *ZKOracleContractExitingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKOracleContractExiting)
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
		it.Event = new(ZKOracleContractExiting)
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
func (it *ZKOracleContractExitingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKOracleContractExitingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKOracleContractExiting represents a Exiting event raised by the ZKOracleContract contract.
type ZKOracleContractExiting struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterExiting is a free log retrieval operation binding the contract event 0x157b1813744b6e99f33bf9153540dd45bd711cbb7322dc3b9c43822687e94180.
//
// Solidity: event Exiting(address indexed sender)
func (_ZKOracleContract *ZKOracleContractFilterer) FilterExiting(opts *bind.FilterOpts, sender []common.Address) (*ZKOracleContractExitingIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ZKOracleContract.contract.FilterLogs(opts, "Exiting", senderRule)
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractExitingIterator{contract: _ZKOracleContract.contract, event: "Exiting", logs: logs, sub: sub}, nil
}

// WatchExiting is a free log subscription operation binding the contract event 0x157b1813744b6e99f33bf9153540dd45bd711cbb7322dc3b9c43822687e94180.
//
// Solidity: event Exiting(address indexed sender)
func (_ZKOracleContract *ZKOracleContractFilterer) WatchExiting(opts *bind.WatchOpts, sink chan<- *ZKOracleContractExiting, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ZKOracleContract.contract.WatchLogs(opts, "Exiting", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKOracleContractExiting)
				if err := _ZKOracleContract.contract.UnpackLog(event, "Exiting", log); err != nil {
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

// ParseExiting is a log parse operation binding the contract event 0x157b1813744b6e99f33bf9153540dd45bd711cbb7322dc3b9c43822687e94180.
//
// Solidity: event Exiting(address indexed sender)
func (_ZKOracleContract *ZKOracleContractFilterer) ParseExiting(log types.Log) (*ZKOracleContractExiting, error) {
	event := new(ZKOracleContractExiting)
	if err := _ZKOracleContract.contract.UnpackLog(event, "Exiting", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKOracleContractRegisteredIterator is returned from FilterRegistered and is used to iterate over the raw logs and unpacked data for Registered events raised by the ZKOracleContract contract.
type ZKOracleContractRegisteredIterator struct {
	Event *ZKOracleContractRegistered // Event containing the contract specifics and raw log

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
func (it *ZKOracleContractRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKOracleContractRegistered)
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
		it.Event = new(ZKOracleContractRegistered)
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
func (it *ZKOracleContractRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKOracleContractRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKOracleContractRegistered represents a Registered event raised by the ZKOracleContract contract.
type ZKOracleContractRegistered struct {
	Sender common.Address
	Index  *big.Int
	Pubkey ZKOraclePublicKey
	Value  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRegistered is a free log retrieval operation binding the contract event 0x195cbf3e686960cbbd84ae23e6bd39c44cf97ba7b9f1fbe6b06a84b2a7d7757c.
//
// Solidity: event Registered(address sender, uint256 index, (uint256,uint256) pubkey, uint256 value)
func (_ZKOracleContract *ZKOracleContractFilterer) FilterRegistered(opts *bind.FilterOpts) (*ZKOracleContractRegisteredIterator, error) {

	logs, sub, err := _ZKOracleContract.contract.FilterLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractRegisteredIterator{contract: _ZKOracleContract.contract, event: "Registered", logs: logs, sub: sub}, nil
}

// WatchRegistered is a free log subscription operation binding the contract event 0x195cbf3e686960cbbd84ae23e6bd39c44cf97ba7b9f1fbe6b06a84b2a7d7757c.
//
// Solidity: event Registered(address sender, uint256 index, (uint256,uint256) pubkey, uint256 value)
func (_ZKOracleContract *ZKOracleContractFilterer) WatchRegistered(opts *bind.WatchOpts, sink chan<- *ZKOracleContractRegistered) (event.Subscription, error) {

	logs, sub, err := _ZKOracleContract.contract.WatchLogs(opts, "Registered")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKOracleContractRegistered)
				if err := _ZKOracleContract.contract.UnpackLog(event, "Registered", log); err != nil {
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

// ParseRegistered is a log parse operation binding the contract event 0x195cbf3e686960cbbd84ae23e6bd39c44cf97ba7b9f1fbe6b06a84b2a7d7757c.
//
// Solidity: event Registered(address sender, uint256 index, (uint256,uint256) pubkey, uint256 value)
func (_ZKOracleContract *ZKOracleContractFilterer) ParseRegistered(log types.Log) (*ZKOracleContractRegistered, error) {
	event := new(ZKOracleContractRegistered)
	if err := _ZKOracleContract.contract.UnpackLog(event, "Registered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKOracleContractReplacedIterator is returned from FilterReplaced and is used to iterate over the raw logs and unpacked data for Replaced events raised by the ZKOracleContract contract.
type ZKOracleContractReplacedIterator struct {
	Event *ZKOracleContractReplaced // Event containing the contract specifics and raw log

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
func (it *ZKOracleContractReplacedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKOracleContractReplaced)
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
		it.Event = new(ZKOracleContractReplaced)
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
func (it *ZKOracleContractReplacedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKOracleContractReplacedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKOracleContractReplaced represents a Replaced event raised by the ZKOracleContract contract.
type ZKOracleContractReplaced struct {
	Sender   common.Address
	Replaced common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterReplaced is a free log retrieval operation binding the contract event 0x25f1f5bf76d36a31ec1b14caf64b580331299fd2037e3589426317b1cdfd4ecb.
//
// Solidity: event Replaced(address indexed sender, address indexed replaced)
func (_ZKOracleContract *ZKOracleContractFilterer) FilterReplaced(opts *bind.FilterOpts, sender []common.Address, replaced []common.Address) (*ZKOracleContractReplacedIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var replacedRule []interface{}
	for _, replacedItem := range replaced {
		replacedRule = append(replacedRule, replacedItem)
	}

	logs, sub, err := _ZKOracleContract.contract.FilterLogs(opts, "Replaced", senderRule, replacedRule)
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractReplacedIterator{contract: _ZKOracleContract.contract, event: "Replaced", logs: logs, sub: sub}, nil
}

// WatchReplaced is a free log subscription operation binding the contract event 0x25f1f5bf76d36a31ec1b14caf64b580331299fd2037e3589426317b1cdfd4ecb.
//
// Solidity: event Replaced(address indexed sender, address indexed replaced)
func (_ZKOracleContract *ZKOracleContractFilterer) WatchReplaced(opts *bind.WatchOpts, sink chan<- *ZKOracleContractReplaced, sender []common.Address, replaced []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}
	var replacedRule []interface{}
	for _, replacedItem := range replaced {
		replacedRule = append(replacedRule, replacedItem)
	}

	logs, sub, err := _ZKOracleContract.contract.WatchLogs(opts, "Replaced", senderRule, replacedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKOracleContractReplaced)
				if err := _ZKOracleContract.contract.UnpackLog(event, "Replaced", log); err != nil {
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

// ParseReplaced is a log parse operation binding the contract event 0x25f1f5bf76d36a31ec1b14caf64b580331299fd2037e3589426317b1cdfd4ecb.
//
// Solidity: event Replaced(address indexed sender, address indexed replaced)
func (_ZKOracleContract *ZKOracleContractFilterer) ParseReplaced(log types.Log) (*ZKOracleContractReplaced, error) {
	event := new(ZKOracleContractReplaced)
	if err := _ZKOracleContract.contract.UnpackLog(event, "Replaced", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ZKOracleContractWithdrawnIterator is returned from FilterWithdrawn and is used to iterate over the raw logs and unpacked data for Withdrawn events raised by the ZKOracleContract contract.
type ZKOracleContractWithdrawnIterator struct {
	Event *ZKOracleContractWithdrawn // Event containing the contract specifics and raw log

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
func (it *ZKOracleContractWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ZKOracleContractWithdrawn)
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
		it.Event = new(ZKOracleContractWithdrawn)
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
func (it *ZKOracleContractWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ZKOracleContractWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ZKOracleContractWithdrawn represents a Withdrawn event raised by the ZKOracleContract contract.
type ZKOracleContractWithdrawn struct {
	Sender common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdrawn is a free log retrieval operation binding the contract event 0xf45a04d08a70caa7eb4b747571305559ad9fdf4a093afd41506b35c8a306fa94.
//
// Solidity: event Withdrawn(address indexed sender)
func (_ZKOracleContract *ZKOracleContractFilterer) FilterWithdrawn(opts *bind.FilterOpts, sender []common.Address) (*ZKOracleContractWithdrawnIterator, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ZKOracleContract.contract.FilterLogs(opts, "Withdrawn", senderRule)
	if err != nil {
		return nil, err
	}
	return &ZKOracleContractWithdrawnIterator{contract: _ZKOracleContract.contract, event: "Withdrawn", logs: logs, sub: sub}, nil
}

// WatchWithdrawn is a free log subscription operation binding the contract event 0xf45a04d08a70caa7eb4b747571305559ad9fdf4a093afd41506b35c8a306fa94.
//
// Solidity: event Withdrawn(address indexed sender)
func (_ZKOracleContract *ZKOracleContractFilterer) WatchWithdrawn(opts *bind.WatchOpts, sink chan<- *ZKOracleContractWithdrawn, sender []common.Address) (event.Subscription, error) {

	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _ZKOracleContract.contract.WatchLogs(opts, "Withdrawn", senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ZKOracleContractWithdrawn)
				if err := _ZKOracleContract.contract.UnpackLog(event, "Withdrawn", log); err != nil {
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

// ParseWithdrawn is a log parse operation binding the contract event 0xf45a04d08a70caa7eb4b747571305559ad9fdf4a093afd41506b35c8a306fa94.
//
// Solidity: event Withdrawn(address indexed sender)
func (_ZKOracleContract *ZKOracleContractFilterer) ParseWithdrawn(log types.Log) (*ZKOracleContractWithdrawn, error) {
	event := new(ZKOracleContractWithdrawn)
	if err := _ZKOracleContract.contract.UnpackLog(event, "Withdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
