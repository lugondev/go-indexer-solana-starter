package solana

import (
	"context"
	"fmt"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type Client struct {
	rpc *rpc.Client
}

func NewClient(rpcURL, wsURL string) (*Client, error) {
	if rpcURL == "" {
		return nil, fmt.Errorf("rpcURL cannot be empty")
	}

	client := rpc.New(rpcURL)
	return &Client{
		rpc: client,
	}, nil
}

func (c *Client) GetSlot(ctx context.Context) (uint64, error) {
	slot, err := c.rpc.GetSlot(ctx, rpc.CommitmentConfirmed)
	if err != nil {
		return 0, fmt.Errorf("get slot: %w", err)
	}
	return slot, nil
}

func (c *Client) GetTransaction(ctx context.Context, signature solana.Signature) (*rpc.GetTransactionResult, error) {
	out, err := c.rpc.GetTransaction(
		ctx,
		signature,
		&rpc.GetTransactionOpts{
			Encoding:                       solana.EncodingBase64,
			Commitment:                     rpc.CommitmentConfirmed,
			MaxSupportedTransactionVersion: nil,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("get transaction: %w", err)
	}
	return out, nil
}

func (c *Client) GetSignaturesForAddress(ctx context.Context, address solana.PublicKey, limit int, before, until *solana.Signature) ([]*rpc.TransactionSignature, error) {
	opts := &rpc.GetSignaturesForAddressOpts{
		Limit: &limit,
	}
	if before != nil {
		opts.Before = *before
	}
	if until != nil {
		opts.Until = *until
	}

	sigs, err := c.rpc.GetSignaturesForAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("get signatures for address: %w", err)
	}
	return sigs, nil
}

func (c *Client) GetBlockTime(ctx context.Context, slot uint64) (int64, error) {
	blockTime, err := c.rpc.GetBlockTime(ctx, slot)
	if err != nil {
		return 0, fmt.Errorf("get block time: %w", err)
	}
	if blockTime == nil {
		return 0, fmt.Errorf("block time is nil")
	}
	return blockTime.Time().Unix(), nil
}

type Block struct {
	Slot              uint64
	Blockhash         string
	PreviousBlockhash string
	ParentSlot        uint64
	Transactions      []Transaction
}

type Transaction struct {
	Signature string
	Message   Message
	Meta      *TransactionMeta
}

type Message struct {
	AccountKeys     []string
	RecentBlockhash string
	Instructions    []Instruction
}

type Instruction struct {
	ProgramIDIndex int
	Accounts       []int
	Data           string
}

type TransactionMeta struct {
	Err               error
	Fee               uint64
	PreBalances       []uint64
	PostBalances      []uint64
	InnerInstructions []InnerInstruction
	LogMessages       []string
}

type InnerInstruction struct {
	Index        int
	Instructions []Instruction
}

func (c *Client) GetBlock(ctx context.Context, slot uint64) (*Block, error) {
	return nil, fmt.Errorf("not implemented")
}
