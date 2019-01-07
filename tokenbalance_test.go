package tokenbalance

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"math/big"
	"os"
	"os/user"
	"path"
	"testing"

	"github.com/vitelabs/go-vite/log15"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestFailedConnection(t *testing.T) {
	c := &Config{
		GethLocation: "https://google.com",
		Logs:         true,
	}
	err := c.Connect()
	assert.Error(t, err)
}

func TestFailingNoConfig(t *testing.T) {
	_, err := New("0xd26114cd6EE289AccF82350c8d8487fedB8A0C07", "0x42d4722b804585cdf6406fa7739e794b0aa8b1ff")
	assert.Error(t, err)
}

func TestConnection(t *testing.T) {
	c := &Config{
		GethLocation: os.Getenv("ETH"),
		Logs:         true,
	}
	err := c.Connect()
	assert.Nil(t, err)
}

func TestZeroDecimal(t *testing.T) {
	number := big.NewInt(123456789)
	tokenCorrected := bigIntString(number, 0)
	assert.Equal(t, "123456789", tokenCorrected)
}

func TestZeroBalance(t *testing.T) {
	number := big.NewInt(0)
	tokenCorrected := bigIntString(number, 18)
	assert.Equal(t, "0.0", tokenCorrected)
}

func TestFormatDecimal(t *testing.T) {
	number := big.NewInt(0)
	number.SetString("72094368689712", 10)
	tokenCorrected := bigIntString(number, 18)
	assert.Equal(t, "0.000072094368689712", tokenCorrected)
}

func TestFormatSmallDecimal(t *testing.T) {
	number := big.NewInt(0)
	number.SetString("123", 10)
	tokenCorrected := bigIntString(number, 18)
	assert.Equal(t, "0.000000000000000123", tokenCorrected)
}

func TestFormatVerySmallDecimal(t *testing.T) {
	number := big.NewInt(0)
	number.SetString("1142400000000001", 10)
	tokenCorrected := bigIntString(number, 18)
	assert.Equal(t, "0.001142400000000001", tokenCorrected)
}

func TestFailedNewTokenBalance(t *testing.T) {
	_, err := New("0x42D4722B804585CDf6406fa7739e794b0Aa8b1FF", "0x42d4722b804585cdf6406fa7739e794b0aa8b1ff")
	assert.Error(t, err)
}

func TestSymbolFix(t *testing.T) {
	symbol := symbolFix("0x86Fa049857E0209aa7D9e616F7eb3b3B78ECfdb0")
	assert.Equal(t, "EOS", symbol)
}

func TestTokenBalance_ToJSON(t *testing.T) {
	symbol := symbolFix("0x86Fa049857E0209aa7D9e616F7eb3b3B78ECfdb0")
	assert.Equal(t, "EOS", symbol)
}

func TestNewTokenBalance(t *testing.T) {
	c := &Config{
		GethLocation: "http://192.168.31.49:8545",
		Logs:         true,
	}
	err := c.Connect()
	assert.Nil(t, err)
	tb, err := New("0x1b793e49237758dbd8b752afc9eb4b329d5da016", "0x4bFa8b23EfEd0CaBDC7A7bBE575Ea0110792b73E")
	assert.Nil(t, err)

	t.Log(tb.Symbol)
	t.Log(tb.Balance)
	t.Log(tb.Decimals)
	//assert.Equal(t, "0x42D4722B804585CDf6406fa7739e794b0Aa8b1FF", tb.Wallet.String())
	//assert.Equal(t, "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07", tb.Contract.String())
	//assert.Equal(t, "600000.0", tb.BalanceString())
	//assert.Equal(t, "1.020095885777777767", tb.ETHString())
	//assert.Equal(t, int64(18), tb.Decimals)
	//assert.Equal(t, "OMG", tb.Symbol)

}

func TestTokenAllBalance(t *testing.T) {
	current, err := user.Current()
	if err != nil {
		panic(err)
	}
	log15.Root().SetHandler(
		log15.LvlFilterHandler(log15.LvlInfo, log15.Must.FileHandler(path.Join(current.HomeDir, "log/balances.log"), log15.JsonFormat())),
	)
	c := &Config{
		GethLocation: "http://192.168.31.49:8545",
		Logs:         true,
	}
	err = c.Connect()
	assert.Nil(t, err)

	tb := &TokenBalance{
		Contract: common.HexToAddress("0x1b793e49237758dbd8b752afc9eb4b329d5da016"),
		Wallet:   common.HexToAddress("0x4bFa8b23EfEd0CaBDC7A7bBE575Ea0110792b73E"),
		Decimals: 0,
		Balance:  big.NewInt(0),
		ctx:      context.TODO(),
	}

	tk, err := newTokenCaller(tb.Contract, Geth)
	if err != nil {
		panic(err)
	}
	fi, err := os.Open("/Users/jie/log/address_final.log")
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		balance, err := tk.BalanceOf(nil, common.HexToAddress(string(a)))
		if err != nil {
			continue
		}
		log15.Info(fmt.Sprintf("%s,%s", string(a), balance.String()))
		//fmt.Println(fmt.Sprintf("%s,%s", w, balance.String()))
	}

	//assert.Equal(t, "0x42D4722B804585CDf6406fa7739e794b0Aa8b1FF", tb.Wallet.String())
	//assert.Equal(t, "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07", tb.Contract.String())
	//assert.Equal(t, "600000.0", tb.BalanceString())
	//assert.Equal(t, "1.020095885777777767", tb.ETHString())
	//assert.Equal(t, int64(18), tb.Decimals)
	//assert.Equal(t, "OMG", tb.Symbol)

}
