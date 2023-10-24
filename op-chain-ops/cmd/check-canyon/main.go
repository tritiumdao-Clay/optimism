package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum-optimism/optimism/op-service/sources"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

func CalcBaseFee(parent eth.BlockInfo, elasticity uint64, preCanyon bool) *big.Int {
	denomUint := uint64(250)
	if preCanyon {
		denomUint = uint64(50)
	}
	parentGasTarget := parent.GasLimit() / elasticity
	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.GasUsed() == parentGasTarget {
		return new(big.Int).Set(parent.BaseFee())
	}

	var (
		num   = new(big.Int)
		denom = new(big.Int)
	)

	if parent.GasUsed() > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		// max(1, parentBaseFee * gasUsedDelta / parentGasTarget / baseFeeChangeDenominator)
		num.SetUint64(parent.GasUsed() - parentGasTarget)
		num.Mul(num, parent.BaseFee())
		num.Div(num, denom.SetUint64(parentGasTarget))
		num.Div(num, denom.SetUint64(denomUint))
		baseFeeDelta := math.BigMax(num, common.Big1)

		return num.Add(parent.BaseFee(), baseFeeDelta)
	} else {
		// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
		// max(0, parentBaseFee * gasUsedDelta / parentGasTarget / baseFeeChangeDenominator)
		num.SetUint64(parentGasTarget - parent.GasUsed())
		num.Mul(num, parent.BaseFee())
		num.Div(num, denom.SetUint64(parentGasTarget))
		num.Div(num, denom.SetUint64(denomUint))
		baseFee := num.Sub(parent.BaseFee(), num)

		return math.BigMax(baseFee, common.Big0)
	}
}

func PreCanyonEncode(receipts types.Receipts) [][]byte {
	for _, receipt := range receipts {
		if receipt.Type == types.DepositTxType {
			receipt.DepositReceiptVersion = nil
		}
	}
	var out [][]byte
	for i := range receipts {
		var buf bytes.Buffer
		receipts.EncodeIndex(i, &buf)
		out = append(out, buf.Bytes())
	}
	return out
}

func PostCanyonEncode(receipts types.Receipts) [][]byte {
	v := uint64(1)
	for _, receipt := range receipts {
		if receipt.Type == types.DepositTxType {
			receipt.DepositReceiptVersion = &v
		}
	}
	var out [][]byte
	for i := range receipts {
		var buf bytes.Buffer
		receipts.EncodeIndex(i, &buf)
		out = append(out, buf.Bytes())
	}
	return out
}

func HashList(list [][]byte) common.Hash {
	t := trie.NewEmpty(trie.NewDatabase(rawdb.NewDatabase(memorydb.New()), nil))
	for i, value := range list {
		var index []byte
		val := make([]byte, len(value))
		copy(val, value)
		index = rlp.AppendUint64(index, uint64(i))
		if err := t.Update(index, val); err != nil {
			panic(err)
		}
	}
	return t.Hash()
}

type ReceiptFetcher interface {
	InfoByNumber(context.Context, uint64) (eth.BlockInfo, error)
	FetchReceipts(context.Context, common.Hash) (eth.BlockInfo, types.Receipts, error)
}

func ValidatePreCanyonReceipts(number uint64, client ReceiptFetcher) func() error {
	return func() error {

		block, err := client.InfoByNumber(context.Background(), number)
		if err != nil {
			return err
		}
		_, receipts, err := client.FetchReceipts(context.Background(), block.Hash())
		if err != nil {
			return err
		}

		have := block.ReceiptHash()
		want := HashList(PreCanyonEncode(receipts))
		if have != want {
			return fmt.Errorf("Receipts do not look correct as pre-canyon. have: %v, want: %v", have, want)
		}
		return nil
	}
}

func ValidatePreCanyon1559Params(number, elasticity uint64, client ReceiptFetcher) func() error {
	return func() error {

		block, err := client.InfoByNumber(context.Background(), number)
		if err != nil {
			return err
		}
		parent, err := client.InfoByNumber(context.Background(), number-1)
		if err != nil {
			return err
		}

		want := CalcBaseFee(parent, elasticity, true)
		have := block.BaseFee()
		if have.Cmp(want) != 0 {
			return fmt.Errorf("BaseFee does not match. have: %v. want: %v", have, want)
		}
		return nil
	}
}

func ValidatePostCanyonReceipts(number uint64, client ReceiptFetcher) func() error {
	return func() error {

		block, err := client.InfoByNumber(context.Background(), number)
		if err != nil {
			return err
		}
		_, receipts, err := client.FetchReceipts(context.Background(), block.Hash())
		if err != nil {
			return err
		}

		have := block.ReceiptHash()
		want := HashList(PostCanyonEncode(receipts))
		if have != want {
			return fmt.Errorf("Receipts do not look correct as post-canyon. have: %v, want: %v", have, want)
		}
		return nil
	}
}

func ValidatePostCanyon1559Params(number, elasticity uint64, client ReceiptFetcher) func() error {
	return func() error {
		block, err := client.InfoByNumber(context.Background(), number)
		if err != nil {
			return err
		}
		parent, err := client.InfoByNumber(context.Background(), number-1)
		if err != nil {
			return err
		}

		want := CalcBaseFee(parent, elasticity, false)
		have := block.BaseFee()
		if have.Cmp(want) != 0 {
			return fmt.Errorf("BaseFee does not match. have: %v. want: %v", have, want)
		}
		return nil
	}
}

func ValidatePair(pre, post func() error, preValid bool) {
	if preValid {
		if err := pre(); err != nil {
			log.Crit("Pre-state was invalid when it was expected to be valid", "err", err)
		}
		if err := post(); err == nil {
			log.Crit("Post-state was valid when it was expected to be invalid")
		}
	} else {
		if err := pre(); err == nil {
			log.Crit("Pre-state was valid when it was expected to be invalid")
		}
		if err := post(); err != nil {
			log.Crit("Post-state was invalid when it was expected to be valid", "err", err)
		}
	}
}

func main() {
	logger := log.New()

	// Define the flag variables
	var (
		preCanyon  bool
		number     uint64
		elasticity uint64
		rpcURL     string
	)

	// Define and parse the command-line flags
	flag.BoolVar(&preCanyon, "pre-canyon", true, "Set this flag to assert pre-canyon receipt hash behavior")
	flag.Uint64Var(&number, "number", 111253022, "block number to check")
	flag.Uint64Var(&elasticity, "elasticity", 6, "Specify the EIP-1559 elasticity. 6 on mainnet/sepolia. 10 on goerli")
	flag.StringVar(&rpcURL, "rpc-url", "https://mainnet.optimism.io", "Specify the RPC URL as a string")

	// Parse the command-line arguments
	flag.Parse()

	l1RPC, err := client.NewRPC(context.Background(), logger, rpcURL, client.WithDialBackoff(10))
	if err != nil {
		log.Crit("Error creating RPC", "err", err)
	}
	c := &rollup.Config{SeqWindowSize: 10}
	l1ClCfg := sources.L1ClientDefaultConfig(c, true, sources.RPCKindBasic)
	client, err := sources.NewL1Client(l1RPC, logger, nil, l1ClCfg)
	if err != nil {
		log.Crit("Error creating RPC", "err", err)
	}

	ValidatePair(ValidatePreCanyonReceipts(number, client), ValidatePostCanyonReceipts(number, client), preCanyon)
	ValidatePair(ValidatePreCanyon1559Params(number, elasticity, client), ValidatePostCanyon1559Params(number, elasticity, client), preCanyon)
}
