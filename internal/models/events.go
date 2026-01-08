package models

import (
	"time"

	"github.com/gagliardetto/solana-go"
)

type EventType string

const (
	EventTypeTokensMinted         EventType = "TokensMintedEvent"
	EventTypeTokensTransferred    EventType = "TokensTransferredEvent"
	EventTypeTokensBurned         EventType = "TokensBurnedEvent"
	EventTypeDelegateApproved     EventType = "DelegateApprovedEvent"
	EventTypeDelegateRevoked      EventType = "DelegateRevokedEvent"
	EventTypeTokenAccountClosed   EventType = "TokenAccountClosedEvent"
	EventTypeTokenAccountFrozen   EventType = "TokenAccountFrozenEvent"
	EventTypeTokenAccountThawed   EventType = "TokenAccountThawedEvent"
	EventTypeUserAccountCreated   EventType = "UserAccountCreatedEvent"
	EventTypeUserAccountUpdated   EventType = "UserAccountUpdatedEvent"
	EventTypeUserAccountClosed    EventType = "UserAccountClosedEvent"
	EventTypeConfigUpdated        EventType = "ConfigUpdatedEvent"
	EventTypeProgramPaused        EventType = "ProgramPausedEvent"
	EventTypeNftCollectionCreated EventType = "NftCollectionCreatedEvent"
	EventTypeNftMinted            EventType = "NftMintedEvent"
	EventTypeNftListed            EventType = "NftListedEvent"
	EventTypeNftSold              EventType = "NftSoldEvent"
	EventTypeNftListingCancelled  EventType = "NftListingCancelledEvent"
	EventTypeNftOfferCreated      EventType = "NftOfferCreatedEvent"
	EventTypeNftOfferAccepted     EventType = "NftOfferAcceptedEvent"
)

type BaseEvent struct {
	ID        string           `bson:"_id,omitempty" json:"id,omitempty"`
	EventType EventType        `bson:"event_type" json:"event_type"`
	Signature string           `bson:"signature" json:"signature"`
	Slot      uint64           `bson:"slot" json:"slot"`
	BlockTime time.Time        `bson:"block_time" json:"block_time"`
	ProgramID solana.PublicKey `bson:"program_id" json:"program_id"`
	CreatedAt time.Time        `bson:"created_at" json:"created_at"`
	RawData   []byte           `bson:"raw_data,omitempty" json:"raw_data,omitempty"`
}

type TokensMintedEvent struct {
	BaseEvent `bson:",inline"`
	Mint      solana.PublicKey `bson:"mint" json:"mint"`
	Recipient solana.PublicKey `bson:"recipient" json:"recipient"`
	Amount    uint64           `bson:"amount" json:"amount"`
	Timestamp int64            `bson:"timestamp" json:"timestamp"`
}

type TokensTransferredEvent struct {
	BaseEvent `bson:",inline"`
	Mint      solana.PublicKey `bson:"mint" json:"mint"`
	From      solana.PublicKey `bson:"from" json:"from"`
	To        solana.PublicKey `bson:"to" json:"to"`
	Amount    uint64           `bson:"amount" json:"amount"`
	Timestamp int64            `bson:"timestamp" json:"timestamp"`
}

type TokensBurnedEvent struct {
	BaseEvent `bson:",inline"`
	Mint      solana.PublicKey `bson:"mint" json:"mint"`
	Owner     solana.PublicKey `bson:"owner" json:"owner"`
	Amount    uint64           `bson:"amount" json:"amount"`
	Timestamp int64            `bson:"timestamp" json:"timestamp"`
}

type UserAccountCreatedEvent struct {
	BaseEvent `bson:",inline"`
	User      solana.PublicKey `bson:"user" json:"user"`
	Authority solana.PublicKey `bson:"authority" json:"authority"`
	Timestamp int64            `bson:"timestamp" json:"timestamp"`
}

type UserAccountUpdatedEvent struct {
	BaseEvent `bson:",inline"`
	User      solana.PublicKey `bson:"user" json:"user"`
	OldPoints uint64           `bson:"old_points" json:"old_points"`
	NewPoints uint64           `bson:"new_points" json:"new_points"`
	Timestamp int64            `bson:"timestamp" json:"timestamp"`
}

type ConfigUpdatedEvent struct {
	BaseEvent `bson:",inline"`
	Admin     solana.PublicKey `bson:"admin" json:"admin"`
	OldFee    uint64           `bson:"old_fee" json:"old_fee"`
	NewFee    uint64           `bson:"new_fee" json:"new_fee"`
	Timestamp int64            `bson:"timestamp" json:"timestamp"`
}

type NftMintedEvent struct {
	BaseEvent  `bson:",inline"`
	NftMint    solana.PublicKey `bson:"nft_mint" json:"nft_mint"`
	Collection solana.PublicKey `bson:"collection" json:"collection"`
	Owner      solana.PublicKey `bson:"owner" json:"owner"`
	Name       string           `bson:"name" json:"name"`
	Uri        string           `bson:"uri" json:"uri"`
	Timestamp  int64            `bson:"timestamp" json:"timestamp"`
}
