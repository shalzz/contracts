package kyber_test

import (
	"context"
	"os"
    "fmt"
	"testing"
    "math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
    "github.com/tokencard/contracts/v2/pkg/bindings/mocks"
	"github.com/tokencard/contracts/v2/pkg/bindings/mocks/kyber"
    // "github.com/tokencard/contracts/v2/pkg/bindings/mocks/makerDAO"
    "github.com/tokencard/contracts/v2/pkg/bindings"
    "github.com/tokencard/ethertest"
    . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
    . "github.com/tokencard/contracts/v2/test/shared"
)


var Wallet *bindings.Wallet
var WalletAddress common.Address

var KyberNetworkProxy *kyber.KyberNetworkProxy
var KyberNetworkProxyAddress common.Address

var KyberNetwork *kyber.KyberNetwork
var KyberNetworkAddress common.Address

var TKNReserve *kyber.KyberReserve
var TKNReserveAddress common.Address

var DAIReserve *kyber.KyberReserve
var DAIReserveAddress common.Address

var TKNLiquidityConversionRates *kyber.LiquidityConversionRates
var TKNLiquidityConversionRatesAddress common.Address

var DAILiquidityConversionRates *kyber.LiquidityConversionRates
var DAILiquidityConversionRatesAddress common.Address

var FeeBurner *kyber.FeeBurner
var FeeBurnerAddress common.Address

var ExpectedRate *kyber.ExpectedRate
var ExpectedRateAddress common.Address

var SanityRates *kyber.SanityRates
var SanityRatesAddress common.Address

var KNCBurner *mocks.BurnerToken
var KNCBurnerAddress common.Address

var DAI *mocks.BurnerToken
var DAIAddress common.Address

var KNCWallet *ethertest.Account

var TKNWallet *ethertest.Account

func TestTokenWhitelistSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "kyber Suite")
}

var _ = BeforeEach(func() {
	err := InitializeBackend()
	Expect(err).ToNot(HaveOccurred())

    var tx *types.Transaction

    //initialize wallets
    KNCWallet = ethertest.NewAccount()
    TKNWallet = ethertest.NewAccount()

    // deploy contracts
    KyberNetworkProxyAddress, tx, KyberNetworkProxy, err = kyber.DeployKyberNetworkProxy(BankAccount.TransactOpts(), Backend, Owner.Address())
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    KyberNetworkAddress, tx, KyberNetwork, err = kyber.DeployKyberNetwork(BankAccount.TransactOpts(), Backend, Owner.Address())
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    DAIAddress, tx, DAI, err = mocks.DeployBurnerToken(Owner.TransactOpts(), Backend)
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    KNCBurnerAddress, tx, KNCBurner, err = mocks.DeployBurnerToken(Owner.TransactOpts(), Backend)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    kncEthPrec := new(big.Int)
    kncEthPrec.SetString("1273294871580578838478",10)
    FeeBurnerAddress, tx, FeeBurner, err = kyber.DeployFeeBurner(BankAccount.TransactOpts(), Backend, Owner.Address(), KNCBurnerAddress, KyberNetworkAddress, kncEthPrec)
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    ExpectedRateAddress, tx, ExpectedRate, err = kyber.DeployExpectedRate(BankAccount.TransactOpts(), Backend, KyberNetworkAddress, KNCBurnerAddress, Owner.Address())
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    TKNLiquidityConversionRatesAddress, tx, TKNLiquidityConversionRates, err = kyber.DeployLiquidityConversionRates(BankAccount.TransactOpts(), Backend, Owner.Address(), TKNBurnerAddress)
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    DAILiquidityConversionRatesAddress, tx, DAILiquidityConversionRates, err = kyber.DeployLiquidityConversionRates(BankAccount.TransactOpts(), Backend, Owner.Address(), DAIAddress)
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    SanityRatesAddress, tx, SanityRates, err = kyber.DeploySanityRates(BankAccount.TransactOpts(), Backend, Owner.Address())
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    TKNReserveAddress, tx, TKNReserve, err = kyber.DeployKyberReserve(BankAccount.TransactOpts(), Backend, KyberNetworkAddress, TKNLiquidityConversionRatesAddress, Owner.Address())
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    DAIReserveAddress, tx, DAIReserve, err = kyber.DeployKyberReserve(BankAccount.TransactOpts(), Backend, KyberNetworkAddress, DAILiquidityConversionRatesAddress, Owner.Address())
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    WalletAddress, tx, Wallet, err = bindings.DeployWallet(BankAccount.TransactOpts(), Backend, Owner.Address(), true, ENSRegistryAddress, TokenWhitelistName, ControllerName, LicenceName, EthToWei(1))
	Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    //set network addresses and params
    tx, err = KyberNetworkProxy.SetKyberNetworkContract(Owner.TransactOpts(), KyberNetworkAddress)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.SetKyberProxy(Owner.TransactOpts(), KyberNetworkProxyAddress)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.SetFeeBurner(Owner.TransactOpts(), FeeBurnerAddress)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.SetExpectedRate(Owner.TransactOpts(), ExpectedRateAddress)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.SetParams(Owner.TransactOpts(), big.NewInt(100000000000), big.NewInt(20))
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.SetEnable(Owner.TransactOpts(), true)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.AddOperator(Owner.TransactOpts(), Owner.Address())
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.AddReserve(Owner.TransactOpts(), TKNReserveAddress, false)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.AddReserve(Owner.TransactOpts(), DAIReserveAddress, false)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.ListPairForReserve(Owner.TransactOpts(), TKNReserveAddress, TKNBurnerAddress, true, true, true)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KyberNetwork.ListPairForReserve(Owner.TransactOpts(), DAIReserveAddress, DAIAddress, true, true, true)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = TKNReserve.SetTokenWallet(Owner.TransactOpts(), TKNBurnerAddress, TKNWallet.Address())
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = TKNLiquidityConversionRates.SetReserveAddress(Owner.TransactOpts(), TKNReserveAddress)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = DAILiquidityConversionRates.SetReserveAddress(Owner.TransactOpts(), DAIReserveAddress)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    maxCapBuyInWei := new(big.Int)
    maxCapBuyInWei.SetString("10000000000000000000",10)
    tx, err = TKNLiquidityConversionRates.SetLiquidityParams(Owner.TransactOpts(),
    big.NewInt(7696581394),
    big.NewInt(1374389534),
    big.NewInt(40),
    maxCapBuyInWei,
    maxCapBuyInWei,
    big.NewInt(25),
    big.NewInt(5000000000000000),
    big.NewInt(1250000000000000))
    Expect(err).ToNot(HaveOccurred())
    Backend.Commit()
    Expect(isSuccessful(tx)).To(BeTrue())

    maxCapBuyInWei.SetString("10000000000000000000",10)
    tx, err = DAILiquidityConversionRates.SetLiquidityParams(Owner.TransactOpts(),
    big.NewInt(7696581394),
    big.NewInt(1374389534),
    big.NewInt(40),
    maxCapBuyInWei,
    maxCapBuyInWei,
    big.NewInt(25),
    big.NewInt(5000000000000000),
    big.NewInt(1250000000000000))
    Expect(err).ToNot(HaveOccurred())
    Backend.Commit()
    Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = FeeBurner.AddOperator(Owner.TransactOpts(), Owner.Address())
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = FeeBurner.SetReserveData(Owner.TransactOpts(), TKNReserveAddress, big.NewInt(25), KNCWallet.Address())
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = FeeBurner.SetReserveData(Owner.TransactOpts(), DAIReserveAddress, big.NewInt(25), KNCWallet.Address())
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = FeeBurner.SetWalletFees(Owner.TransactOpts(), KNCWallet.Address(), big.NewInt(25))
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = SanityRates.AddOperator(Owner.TransactOpts(), Owner.Address())
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())


    //deposit ETH and TKN to the reserve Contracts
    BankAccount.MustTransfer(Backend, TKNReserveAddress, EthToWei(100))
    BankAccount.MustTransfer(Backend, DAIReserveAddress, EthToWei(100))
    BankAccount.MustTransfer(Backend, KNCWallet.Address(), EthToWei(1))
    BankAccount.MustTransfer(Backend, TKNWallet.Address(), EthToWei(1))

    tx, err = TKNBurner.Mint(Owner.TransactOpts(), TKNWallet.Address(), big.NewInt(3800000000000))
    Expect(err).ToNot(HaveOccurred())
    Backend.Commit()
    Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = TKNBurner.Approve(TKNWallet.TransactOpts(), TKNReserveAddress, big.NewInt(3700000000000))
    Expect(err).ToNot(HaveOccurred())
    Backend.Commit()
    Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = DAI.Mint(Owner.TransactOpts(), DAIReserveAddress, EthToWei(38000))
    Expect(err).ToNot(HaveOccurred())
    Backend.Commit()
    Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = DAIReserve.ApproveWithdrawAddress(Owner.TransactOpts(), DAIAddress, Owner.Address(), true)
    Expect(err).ToNot(HaveOccurred())
	Backend.Commit()
	Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KNCBurner.Mint(Owner.TransactOpts(), KNCWallet.Address(), EthToWei(38000))
    Expect(err).ToNot(HaveOccurred())
    Backend.Commit()
    Expect(isSuccessful(tx)).To(BeTrue())

    tx, err = KNCBurner.Approve(KNCWallet.TransactOpts(), FeeBurnerAddress, EthToWei(38000))
    Expect(err).ToNot(HaveOccurred())
    Backend.Commit()
    Expect(isSuccessful(tx)).To(BeTrue())
})

var _ = AfterEach(func() {
    td := CurrentGinkgoTestDescription()

	if td.Failed {
		fmt.Fprintf(GinkgoWriter, "\nLast Executed Smart Contract Line for %s:%d\n", td.FileName, td.LineNumber)
		fmt.Fprintln(GinkgoWriter, TestRig.LastExecuted())
	}
	err := Backend.Close()
	Expect(err).ToNot(HaveOccurred())

})

var _ = AfterSuite(func() {
	// TestRig.ExpectMinimumCoverage("mocks/kyber/KyberNetworkProxy.sol", 0.0)
	TestRig.PrintGasUsage(os.Stdout)
})

func isSuccessful(tx *types.Transaction) bool {
	r, err := Backend.TransactionReceipt(context.Background(), tx.Hash())
	Expect(err).ToNot(HaveOccurred())
	return r.Status == types.ReceiptStatusSuccessful
}
