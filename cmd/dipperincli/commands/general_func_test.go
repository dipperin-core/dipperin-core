// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.


package commands

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/consts"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestMoneyValueToCSCoin(t *testing.T) {

	moneyValue1 := "0.001"
	moneyValue2 := "7.890765"

	moneyValue3 := "a.2234343454"

	moneyValue4 := "300"

	moneyValue5 := "0.0000000001"

	value, err := MoneyValueToCSCoin(moneyValue1)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(0.001*consts.DIP), value)

	value, err = MoneyValueToCSCoin(moneyValue2)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(7.890765*consts.DIP), value)

	value, err = MoneyValueToCSCoin(moneyValue3)
	assert.Error(t, err)

	value, err = MoneyValueToCSCoin(moneyValue4)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(300*consts.DIP), value)

	value, err = MoneyValueToCSCoin(moneyValue5)
	assert.Error(t, err)

	/*	decimalValue, err := strconv.ParseInt("000000101",10,64)
		log.Info("the decimalValue is:","decimalValue",decimalValue)*/
}

func TestCSCoinToMoneyValue(t *testing.T) {
	csCoinValue1 := (*hexutil.Big)(big.NewInt(10000000))

	csCoinValue2 := (*hexutil.Big)(big.NewInt(34545000))

	csCoinValue3 := (*hexutil.Big)(big.NewInt(600000000))

	csCoinValue4 := (*hexutil.Big)(big.NewInt(897878600000000))

	csCoinValue5 := (*hexutil.Big)(big.NewInt(3069000000000))

	moneyValue1, err := CSCoinToMoneyValue(csCoinValue1)
	assert.NoError(t, err)
	assert.Equal(t, "0.01", moneyValue1)

	moneyValue2, err := CSCoinToMoneyValue(csCoinValue2)
	assert.NoError(t, err)
	assert.Equal(t, "0.034545", moneyValue2)

	moneyValue3, err := CSCoinToMoneyValue(csCoinValue3)
	assert.NoError(t, err)
	assert.Equal(t, "0.6", moneyValue3)

	moneyValue4, err := CSCoinToMoneyValue(csCoinValue4)
	assert.NoError(t, err)
	assert.Equal(t, "897878.6", moneyValue4)

	moneyValue5, err := CSCoinToMoneyValue(csCoinValue5)
	assert.NoError(t, err)
	assert.Equal(t, "3069", moneyValue5)
}

func TestDecimalToInter(t *testing.T) {
	moneyValue1 := "0.001"
	moneyValue2 := "7.89"

	moneyValue3 := "a.2234"

	moneyValue4 := "300"

	moneyValue5 := "0.0001"
	unitBit := 3

	unit := 1
	for i := 0; i < unitBit; i++ {
		unit *= 10
	}
	assert.EqualValues(t, 1000, unit)

	value, err := DecimalToInter(moneyValue1, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(1), value)

	value, err = DecimalToInter(moneyValue2, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(7890), value)

	value, err = DecimalToInter(moneyValue3, unitBit)
	assert.Error(t, err)

	value, err = DecimalToInter(moneyValue4, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, big.NewInt(300000), value)

	value, err = DecimalToInter(moneyValue5, unitBit)
	assert.Error(t, err)
}

func TestInterToDecimal(t *testing.T) {
	csCoinValue1 := (*hexutil.Big)(big.NewInt(100))

	csCoinValue2 := (*hexutil.Big)(big.NewInt(34545))

	csCoinValue3 := (*hexutil.Big)(big.NewInt(6000))

	csCoinValue4 := (*hexutil.Big)(big.NewInt(8978786))

	csCoinValue5 := (*hexutil.Big)(big.NewInt(30690000))

	unitBit := 3

	moneyValue1, err := InterToDecimal(csCoinValue1, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "0.1", moneyValue1)

	moneyValue2, err := InterToDecimal(csCoinValue2, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "34.545", moneyValue2)

	moneyValue3, err := InterToDecimal(csCoinValue3, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "6", moneyValue3)

	moneyValue4, err := InterToDecimal(csCoinValue4, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "8978.786", moneyValue4)

	moneyValue5, err := InterToDecimal(csCoinValue5, unitBit)
	assert.NoError(t, err)
	assert.Equal(t, "30690", moneyValue5)
}

/*func closeWallet(){
	log.Info("close wallet~~~~~~~~~~~~~~~~~~")
}

func receiveExitSignal(exit chan int) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(
		sigCh,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)

	log.Info("receiveExitSignal start ~~~~~~~~~~~~~~`")
	for {
		select {
		case <-exit:
			fmt.Printf("exit ~~~~~~~~~~~~~~~~~")
			log.Info("exit ~~~~~~~~~~~~~~")
			closeWallet()
			return
		case s := <-sigCh:
			log.Info("receive signal", "signal", s)
			fmt.Printf("receive signal")
			closeWallet()
			return
		}
	}
}*/

func TestCheckAndChangeHexToAddress(t *testing.T) {
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    common.Address
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CheckAndChangeHexToAddress(tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckAndChangeHexToAddress() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CheckAndChangeHexToAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseWalletPathAndName(t *testing.T) {
	type args struct {
		inputPath string
	}
	tests := []struct {
		name     string
		args     args
		wantPath string
		wantName string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotName := ParseWalletPathAndName(tt.args.inputPath)
			if gotPath != tt.wantPath {
				t.Errorf("ParseWalletPathAndName() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotName != tt.wantName {
				t.Errorf("ParseWalletPathAndName() gotName = %v, want %v", gotName, tt.wantName)
			}
		})
	}
}
