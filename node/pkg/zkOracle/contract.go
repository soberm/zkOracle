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

// ZKOracleContractABI is the input ABI used to generate the binding from.
const ZKOracleContractABI = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

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
