package blockchain

import (
	"bytes"
	"errors"
	"fmt"
	"math"

	"github.com/ioeXNetwork/ioeX.MainChain/config"
	. "github.com/ioeXNetwork/ioeX.MainChain/core"
	. "github.com/ioeXNetwork/ioeX.MainChain/errors"
	"github.com/ioeXNetwork/ioeX.MainChain/log"

	. "github.com/ioeXNetwork/ioeX.Utility/common"
	. "github.com/ioeXNetwork/ioeX.Utility/crypto"
)

// CheckTransactionSanity verifys received single transaction
func CheckTransactionSanity(version uint32, txn *Transaction) ErrCode {
	if err := CheckTransactionSize(txn); err != nil {
		log.Warn("[CheckTransactionSize],", err)
		return ErrTransactionSize
	}

	if err := CheckTransactionInput(txn); err != nil {
		log.Warn("[CheckTransactionInput],", err)
		return ErrInvalidInput
	}

	if err := CheckTransactionOutput(version, txn); err != nil {
		log.Warn("[CheckTransactionOutput],", err)
		return ErrInvalidOutput
	}

	if err := CheckAssetPrecision(txn); err != nil {
		log.Warn("[CheckAssetPrecesion],", err)
		return ErrAssetPrecision
	}

	if err := CheckAttributeProgram(txn); err != nil {
		log.Warn("[CheckAttributeProgram],", err)
		return ErrAttributeProgram
	}

	if err := CheckTransactionPayload(txn); err != nil {
		log.Warn("[CheckTransactionPayload],", err)
		return ErrTransactionPayload
	}

	// check iterms above for Coinbase transaction
	if txn.IsCoinBaseTx() {
		return Success
	}

	return Success
}

// CheckTransactionContext verifys a transaction with history transaction in ledger
func CheckTransactionContext(txn *Transaction) ErrCode {
	// check if duplicated with transaction in ledger
	if exist := DefaultLedger.Store.IsTxHashDuplicate(txn.Hash()); exist {
		log.Warn("[CheckTransactionContext] duplicate transaction check failed.")
		return ErrTransactionDuplicate
	}

	if txn.IsCoinBaseTx() {
		return Success
	}

	// check double spent transaction
	if DefaultLedger.IsDoubleSpend(txn) {
		log.Warn("[CheckTransactionContext] IsDoubleSpend check faild.")
		return ErrDoubleSpend
	}

	references, err := DefaultLedger.Store.GetTxReference(txn)
	if err != nil {
		log.Warn("[CheckTransactionContext] get transaction reference failed")
		return ErrUnknownReferredTx
	}

	if err := CheckTransactionUTXOLock(txn, references); err != nil {
		log.Warn("[CheckTransactionUTXOLock],", err)
		return ErrUTXOLocked
	}

	if err := CheckTransactionFee(txn, references); err != nil {
		log.Warn("[CheckTransactionFee],", err)
		return ErrTransactionBalance
	}
	if err := CheckDestructionAddress(references); err != nil {
		log.Warn("[CheckDestructionAddress], ", err)
		return ErrInvalidInput
	}
	if err := CheckTransactionSignature(txn, references); err != nil {
		log.Warn("[CheckTransactionSignature],", err)
		return ErrTransactionSignature
	}

	if err := CheckTransactionCoinbaseOutputLock(txn); err != nil {
		log.Warn("[CheckTransactionCoinbaseLock]", err)
		return ErrIneffectiveCoinbase
	}
	return Success
}

func CheckDestructionAddress(references map[*Input]*Output) error {
	for _, output := range references {
		// this uint168 code
		// is the program hash of the IOEX foundation destruction address ELANULLXXXXXXXXXXXXXXXXXXXXXYvs3rr
		// we allow no output from destruction address.
		// So add a check here in case someone crack the private key of this address.
		if output.ProgramHash == Uint168([21]uint8{33, 32, 254, 229, 215, 235, 62, 92, 125, 49, 151, 254, 207, 108, 13, 227, 15, 136, 154, 206, 247}) {
			return errors.New("cannot use utxo in the IOEX foundation destruction address")
		}
	}
	return nil
}

func CheckTransactionCoinbaseOutputLock(txn *Transaction) error {
	type lockTxInfo struct {
		isCoinbaseTx bool
		locktime     uint32
	}
	transactionCache := make(map[Uint256]lockTxInfo)
	currentHeight := DefaultLedger.Store.GetHeight()
	var referTxn *Transaction
	for _, input := range txn.Inputs {
		var lockHeight uint32
		var isCoinbase bool
		referHash := input.Previous.TxID
		if _, ok := transactionCache[referHash]; ok {
			lockHeight = transactionCache[referHash].locktime
			isCoinbase = transactionCache[referHash].isCoinbaseTx
		} else {
			var err error
			referTxn, _, err = DefaultLedger.Store.GetTransaction(referHash)
			// TODO
			// we have executed DefaultLedger.Store.GetTxReference(txn) before.
			//So if we can't find referTxn here, there must be a data inconsistent problem,
			// because we do not add lock correctly.This problem will be fixed later on.
			if err != nil {
				return errors.New("[CheckTransactionCoinbaseOutputLock] get tx reference failed:" + err.Error())
			}
			lockHeight = referTxn.LockTime
			isCoinbase = referTxn.IsCoinBaseTx()
			transactionCache[referHash] = lockTxInfo{isCoinbase, lockHeight}
		}

		if isCoinbase && currentHeight-lockHeight < config.Parameters.ChainParam.CoinbaseLockTime {
			return errors.New("cannot unlock coinbase transaction output")
		}
	}
	return nil
}

//validate the transaction of duplicate UTXO input
func CheckTransactionInput(txn *Transaction) error {
	if txn.IsCoinBaseTx() {
		if len(txn.Inputs) != 1 {
			return errors.New("coinbase must has only one input")
		}
		coinbaseInputHash := txn.Inputs[0].Previous.TxID
		coinbaseInputIndex := txn.Inputs[0].Previous.Index
		//TODO :check sequence
		if !coinbaseInputHash.IsEqual(EmptyHash) || coinbaseInputIndex != math.MaxUint16 {
			return errors.New("invalid coinbase input")
		}

		return nil
	}

	if len(txn.Inputs) <= 0 {
		return errors.New("transaction has no inputs")
	}
	existingTxInputs := make(map[string]struct{})
	for _, input := range txn.Inputs {
		if input.Previous.TxID.IsEqual(EmptyHash) && (input.Previous.Index == math.MaxUint16) {
			return errors.New("invalid transaction input")
		}
		if _, exists := existingTxInputs[input.ReferKey()]; exists {
			return errors.New("duplicated transaction inputs")
		} else {
			existingTxInputs[input.ReferKey()] = struct{}{}
		}
	}

	return nil
}

func CheckTransactionOutput(version uint32, txn *Transaction) error {
	if len(txn.Outputs) > math.MaxUint16 {
		return errors.New("output count should not be greater than 65535(MaxUint16)")
	}

	if txn.IsCoinBaseTx() {
		if len(txn.Outputs) < 2 {
			return errors.New("coinbase output is not enough, at least 2")
		}

		var totalReward = Fixed64(0)
		var foundationReward = Fixed64(0)
		for _, output := range txn.Outputs {
			if output.AssetID != DefaultLedger.Blockchain.AssetID {
				return errors.New("asset ID in coinbase is invalid")
			}
			totalReward += output.Value

			if output.ProgramHash.IsEqual(FoundationAddress11) {
				foundationReward += output.Value
			}
		}

		if Fixed64(foundationReward) < FoundationRewardAmountPerBlock {
			return errors.New("Reward to foundation in coinbase < 20")
		}

		return nil
	}

	if len(txn.Outputs) < 1 {
		return errors.New("transaction has no outputs")
	}
	// check if output address is valid
	for _, output := range txn.Outputs {
		if output.AssetID != DefaultLedger.Blockchain.AssetID {
			return errors.New("asset ID in output is invalid")
		}

		// output value must >= 0
		if output.Value < Fixed64(0) {
			return errors.New("Invalide transaction UTXO output")
		}
		if version&CheckTxOut == CheckTxOut {
			if !CheckOutputProgramHash(output.ProgramHash) {
				return errors.New("output address is invalid")
			}
		}

	}

	return nil
}

func CheckOutputProgramHash(programHash Uint168) bool {
	prefix := programHash[0]
	return prefix == PrefixStandard || prefix == PrefixMultisig || prefix == PrefixCrossChain

}

func CheckTransactionUTXOLock(txn *Transaction, references map[*Input]*Output) error {
	if txn.IsCoinBaseTx() {
		return nil
	}
	for input, output := range references {

		if output.OutputLock == 0 {
			//check next utxo
			continue
		}
		if input.Sequence != math.MaxUint32-1 {
			return errors.New("Invalid input sequence")
		}
		if txn.LockTime < output.OutputLock {
			return errors.New("UTXO output locked")
		}
	}
	return nil
}

func CheckTransactionSize(txn *Transaction) error {
	size := txn.GetSize()
	if size <= 0 || size > config.Parameters.MaxBlockSize {
		return fmt.Errorf("Invalid transaction size: %d bytes", size)
	}

	return nil
}

func CheckAssetPrecision(txn *Transaction) error {
	if len(txn.Outputs) == 0 {
		return nil
	}
	assetOutputs := make(map[Uint256][]*Output)

	for _, v := range txn.Outputs {
		assetOutputs[v.AssetID] = append(assetOutputs[v.AssetID], v)
	}
	for k, outputs := range assetOutputs {
		asset, err := DefaultLedger.GetAsset(k)
		if err != nil {
			return errors.New("The asset not exist in local blockchain.")
		}
		precision := asset.Precision
		for _, output := range outputs {
			if !checkAmountPrecise(output.Value, precision) {
				return errors.New("The precision of asset is incorrect.")
			}
		}
	}
	return nil
}

func CheckTransactionFee(tx *Transaction, references map[*Input]*Output) error {
	var outputValue Fixed64
	var inputValue Fixed64
	for _, output := range tx.Outputs {
		outputValue += output.Value
	}
	for _, reference := range references {
		inputValue += reference.Value
	}
	if inputValue < Fixed64(config.Parameters.PowConfiguration.MinTxFee)+outputValue {
		return fmt.Errorf("transaction fee not enough")
	}
	return nil
}

func CheckAttributeProgram(tx *Transaction) error {
	// Coinbase transaction do not check attribute and program
	if tx.IsCoinBaseTx() {
		return nil
	}

	// Check attributes
	for _, attr := range tx.Attributes {
		if !IsValidAttributeType(attr.Usage) {
			return fmt.Errorf("invalid attribute usage %v", attr.Usage)
		}
	}

	// Check programs
	if len(tx.Programs) == 0 {
		return fmt.Errorf("no programs found in transaction")
	}
	for _, program := range tx.Programs {
		if program.Code == nil {
			return fmt.Errorf("invalid program code nil")
		}
		if program.Parameter == nil {
			return fmt.Errorf("invalid program parameter nil")
		}
		_, err := ToProgramHash(program.Code)
		if err != nil {
			return fmt.Errorf("invalid program code %x", program.Code)
		}
	}
	return nil
}

func CheckTransactionSignature(tx *Transaction, references map[*Input]*Output) error {
	hashes, err := GetTxProgramHashes(tx, references)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	tx.SerializeUnsigned(buf)

	// Sort first
	SortProgramHashes(hashes)
	SortPrograms(tx.Programs)

	return RunPrograms(buf.Bytes(), hashes, tx.Programs)
}

func checkAmountPrecise(amount Fixed64, precision byte) bool {
	return amount.IntValue()%int64(math.Pow(10, float64(8-precision))) == 0
}

func CheckTransactionPayload(txn *Transaction) error {
	switch pld := txn.Payload.(type) {
	case *PayloadRegisterAsset:
		if pld.Asset.Precision < MinPrecision || pld.Asset.Precision > MaxPrecision {
			return errors.New("Invalide asset Precision.")
		}
		if !checkAmountPrecise(pld.Amount, pld.Asset.Precision) {
			return errors.New("Invalide asset value,out of precise.")
		}
	case *PayloadTransferAsset:
	case *PayloadRecord:
	case *PayloadCoinBase:
	default:
		return errors.New("[txValidator],invalidate transaction payload type.")
	}
	return nil
}
