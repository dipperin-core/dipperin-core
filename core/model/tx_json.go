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


package model

import (
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/hexutil"
	"math/big"
	"encoding/json"
	"errors"
	"github.com/dipperin/dipperin-core/third-party/crypto/cs-crypto"
)

type TransactionJSON struct {
	TxData txData
	Wit    witness
}

func (tx Transaction) MarshalJSON() ([]byte, error) {
	tJson := TransactionJSON{
		TxData: tx.data,
		Wit:    tx.wit,
	}
	return json.Marshal(&tJson)
}

func (tx *Transaction) UnmarshalJSON(input []byte) error {
	tJson := TransactionJSON{}
	err := json.Unmarshal(input, &tJson)
	tx.data = tJson.TxData
	tx.wit = tJson.Wit
	id := deriveChainId(tx.wit.V)
	temp := big.NewInt(0).Sub(tx.wit.V, big.NewInt(0).Mul(id, big.NewInt(2)))
	v := big.NewInt(0).Sub(temp, big.NewInt(54))
	if !cs_crypto.ValidSigValue(tx.wit.R, tx.wit.S, v) {
		return errors.New("UnmarshalJSON invalid transaction v, r, s values")
	}

	return err
}

func (t txData) MarshalJSON() ([]byte, error) {
	type txdata struct {
		AccountNonce hexutil.Uint64  `json:"nonce"    gencodec:"required"`
		//Version      hexutil.Uint64  `json:"version" gencodec:"required"`
		Recipient    *common.Address `json:"to"       rlp:"nil"`
		HashLock     *common.Hash    `json:"hashlock" rlp:"nil"`
		TimeLock     *hexutil.Big    `json:"timelock" gencodec:"required"`
		Amount       *hexutil.Big    `json:"value"    gencodec:"required"`
		Fee          *hexutil.Big    `json:"fee"      gencodec:"required"`
		ExtraData    hexutil.Bytes   `json:"extradata"    gencodec:"required"`
	}
	var enc txdata
	enc.AccountNonce = hexutil.Uint64(t.AccountNonce)
	//enc.Version = hexutil.Uint64(t.Version)
	enc.Recipient = t.Recipient
	enc.HashLock = t.HashLock
	enc.TimeLock = (*hexutil.Big)(t.TimeLock)
	enc.Amount = (*hexutil.Big)(t.Amount)
	enc.Fee = (*hexutil.Big)(t.Fee)
	enc.ExtraData = t.ExtraData
	return json.Marshal(&enc)
}

func (t *txData) UnmarshalJSON(input []byte) error {
	type txdata struct {
		AccountNonce *hexutil.Uint64 `json:"nonce"    gencodec:"required"`
		//Version      *hexutil.Uint64 `json:"version" gencodec:"required"`
		Recipient    *common.Address `json:"to"       rlp:"nil"`
		HashLock     *common.Hash    `json:"hashlock" rlp:"nil"`
		TimeLock     *hexutil.Big    `json:"timelock" gencodec:"required"`
		Amount       *hexutil.Big    `json:"value"    gencodec:"required"`
		Fee          *hexutil.Big    `json:"fee"      gencodec:"required"`
		ExtraData    *hexutil.Bytes  `json:"extradata"    gencodec:"required"`
	}
	var dec txdata
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.AccountNonce == nil {
		return errors.New("missing required field 'nonce' for txData")

	}
	t.AccountNonce = uint64(*dec.AccountNonce)
	//if dec.Version == nil {
	//	return errors.New("missing required field 'version' for txData")
	//}
	//t.Version = uint64(*dec.Version)
	if dec.Recipient != nil {
		t.Recipient = dec.Recipient
	}
	if dec.HashLock != nil {
		t.HashLock = dec.HashLock
	}
	if dec.TimeLock != nil {
		t.TimeLock = (*big.Int)(dec.TimeLock)
		//return errors.New("missing required field 'timelock' for txData")
	}
	if dec.Amount == nil {
		return errors.New("missing required field 'amount' for txData")
	}
	t.Amount = (*big.Int)(dec.Amount)
	if dec.Fee == nil {
		return errors.New("missing required field 'fee' for txData")
	}
	t.Fee = (*big.Int)(dec.Fee)
	if dec.ExtraData == nil {
		return errors.New("missing required field 'extradata' for txData")
	}
	t.ExtraData = *dec.ExtraData
	return nil
}

func (t witness) MarshalJSON() ([]byte, error) {
	type wit struct {
		R       *hexutil.Big  `json:"r" gencodec:"required"`
		S       *hexutil.Big  `json:"s" gencodec:"required"`
		V       *hexutil.Big  `json:"v" gencodec:"required"`
		HashKey hexutil.Bytes `json:"hashkey"    gencodec:"required"`
	}
	var enc wit
	enc.R = (*hexutil.Big)(t.R)
	enc.S = (*hexutil.Big)(t.S)
	enc.V = (*hexutil.Big)(t.V)
	enc.HashKey = t.HashKey
	return json.Marshal(&enc)
}

func (t *witness) UnmarshalJSON(input []byte) error {
	type wit struct {
		R *hexutil.Big `json:"r" gencodec:"required"`
		S *hexutil.Big `json:"s" gencodec:"required"`
		V *hexutil.Big `json:"v" gencodec:"required"`
		// hash_key
		HashKey *hexutil.Bytes `json:"hashkey"    gencodec:"required"`
	}
	var dec wit
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.HashKey != nil {
		t.HashKey = *dec.HashKey
	}
	if dec.R == nil {
		return errors.New("missing required field 'R' for witness")
	}
	t.R = (*big.Int)(dec.R)
	if dec.S == nil {
		return errors.New("missing required field 'S' for witness")
	}
	t.S = (*big.Int)(dec.S)
	if dec.V == nil {
		return errors.New("missing required field 'V' for witness")
	}
	t.V = (*big.Int)(dec.V)
	return nil
}
