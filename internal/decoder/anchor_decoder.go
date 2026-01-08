package decoder

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/lugondev/go-indexer-solana-starter/internal/models"
)

type EventDecoder struct {
	discriminators map[string]models.EventType
}

func NewEventDecoder() *EventDecoder {
	return &EventDecoder{
		discriminators: makeDiscriminatorMap(),
	}
}

func makeDiscriminatorMap() map[string]models.EventType {
	return map[string]models.EventType{
		eventDiscriminator("TokensMintedEvent"):         models.EventTypeTokensMinted,
		eventDiscriminator("TokensTransferredEvent"):    models.EventTypeTokensTransferred,
		eventDiscriminator("TokensBurnedEvent"):         models.EventTypeTokensBurned,
		eventDiscriminator("DelegateApprovedEvent"):     models.EventTypeDelegateApproved,
		eventDiscriminator("DelegateRevokedEvent"):      models.EventTypeDelegateRevoked,
		eventDiscriminator("TokenAccountClosedEvent"):   models.EventTypeTokenAccountClosed,
		eventDiscriminator("TokenAccountFrozenEvent"):   models.EventTypeTokenAccountFrozen,
		eventDiscriminator("TokenAccountThawedEvent"):   models.EventTypeTokenAccountThawed,
		eventDiscriminator("UserAccountCreatedEvent"):   models.EventTypeUserAccountCreated,
		eventDiscriminator("UserAccountUpdatedEvent"):   models.EventTypeUserAccountUpdated,
		eventDiscriminator("UserAccountClosedEvent"):    models.EventTypeUserAccountClosed,
		eventDiscriminator("ConfigUpdatedEvent"):        models.EventTypeConfigUpdated,
		eventDiscriminator("ProgramPausedEvent"):        models.EventTypeProgramPaused,
		eventDiscriminator("NftCollectionCreatedEvent"): models.EventTypeNftCollectionCreated,
		eventDiscriminator("NftMintedEvent"):            models.EventTypeNftMinted,
		eventDiscriminator("NftListedEvent"):            models.EventTypeNftListed,
		eventDiscriminator("NftSoldEvent"):              models.EventTypeNftSold,
		eventDiscriminator("NftListingCancelledEvent"):  models.EventTypeNftListingCancelled,
		eventDiscriminator("NftOfferCreatedEvent"):      models.EventTypeNftOfferCreated,
		eventDiscriminator("NftOfferAcceptedEvent"):     models.EventTypeNftOfferAccepted,
	}
}

func eventDiscriminator(name string) string {
	discriminatorPreimage := []byte(fmt.Sprintf("event:%s", name))
	hash := sha256.Sum256(discriminatorPreimage)
	return base64.StdEncoding.EncodeToString(hash[:8])
}

func (d *EventDecoder) DecodeEvent(data []byte) (models.EventType, interface{}, error) {
	if len(data) < 8 {
		return "", nil, fmt.Errorf("data too short for discriminator")
	}

	discriminator := base64.StdEncoding.EncodeToString(data[:8])
	eventType, ok := d.discriminators[discriminator]
	if !ok {
		return "", nil, fmt.Errorf("unknown discriminator: %s", discriminator)
	}

	eventData := data[8:]
	decoder := bin.NewBinDecoder(eventData)

	switch eventType {
	case models.EventTypeTokensMinted:
		event, err := decodeTokensMinted(decoder)
		return eventType, event, err
	case models.EventTypeTokensTransferred:
		event, err := decodeTokensTransferred(decoder)
		return eventType, event, err
	case models.EventTypeTokensBurned:
		event, err := decodeTokensBurned(decoder)
		return eventType, event, err
	case models.EventTypeUserAccountCreated:
		event, err := decodeUserAccountCreated(decoder)
		return eventType, event, err
	case models.EventTypeUserAccountUpdated:
		event, err := decodeUserAccountUpdated(decoder)
		return eventType, event, err
	case models.EventTypeConfigUpdated:
		event, err := decodeConfigUpdated(decoder)
		return eventType, event, err
	case models.EventTypeNftMinted:
		event, err := decodeNftMinted(decoder)
		return eventType, event, err
	default:
		return eventType, nil, fmt.Errorf("decoder not implemented for %s", eventType)
	}
}

func decodeTokensMinted(decoder *bin.Decoder) (*models.TokensMintedEvent, error) {
	event := &models.TokensMintedEvent{}
	if err := decoder.Decode(&event.Mint); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Recipient); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Amount); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Timestamp); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeTokensTransferred(decoder *bin.Decoder) (*models.TokensTransferredEvent, error) {
	event := &models.TokensTransferredEvent{}
	if err := decoder.Decode(&event.Mint); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.From); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.To); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Amount); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Timestamp); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeTokensBurned(decoder *bin.Decoder) (*models.TokensBurnedEvent, error) {
	event := &models.TokensBurnedEvent{}
	if err := decoder.Decode(&event.Mint); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Owner); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Amount); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Timestamp); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeUserAccountCreated(decoder *bin.Decoder) (*models.UserAccountCreatedEvent, error) {
	event := &models.UserAccountCreatedEvent{}
	if err := decoder.Decode(&event.User); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Authority); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Timestamp); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeUserAccountUpdated(decoder *bin.Decoder) (*models.UserAccountUpdatedEvent, error) {
	event := &models.UserAccountUpdatedEvent{}
	if err := decoder.Decode(&event.User); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.OldPoints); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.NewPoints); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Timestamp); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeConfigUpdated(decoder *bin.Decoder) (*models.ConfigUpdatedEvent, error) {
	event := &models.ConfigUpdatedEvent{}
	if err := decoder.Decode(&event.Admin); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.OldFee); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.NewFee); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Timestamp); err != nil {
		return nil, err
	}
	return event, nil
}

func decodeNftMinted(decoder *bin.Decoder) (*models.NftMintedEvent, error) {
	event := &models.NftMintedEvent{}
	if err := decoder.Decode(&event.NftMint); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Collection); err != nil {
		return nil, err
	}
	if err := decoder.Decode(&event.Owner); err != nil {
		return nil, err
	}

	var nameLen uint32
	if err := decoder.Decode(&nameLen); err != nil {
		return nil, err
	}
	nameBytes := make([]byte, nameLen)
	if err := decoder.Decode(&nameBytes); err != nil {
		return nil, err
	}
	event.Name = string(nameBytes)

	var uriLen uint32
	if err := decoder.Decode(&uriLen); err != nil {
		return nil, err
	}
	uriBytes := make([]byte, uriLen)
	if err := decoder.Decode(&uriBytes); err != nil {
		return nil, err
	}
	event.Uri = string(uriBytes)

	if err := decoder.Decode(&event.Timestamp); err != nil {
		return nil, err
	}
	return event, nil
}

func ParseProgramData(logs []string) [][]byte {
	var programData [][]byte

	for _, log := range logs {
		if len(log) < 14 {
			continue
		}

		if log[:13] == "Program data:" {
			dataStr := log[14:]
			data, err := base64.StdEncoding.DecodeString(dataStr)
			if err != nil {
				continue
			}
			programData = append(programData, data)
		}
	}

	return programData
}

func FilterByProgramID(programID solana.PublicKey, data []byte) bool {
	if len(data) < 8 {
		return false
	}

	decoder := bin.NewBinDecoder(data)

	var eventProgramID solana.PublicKey
	if err := decoder.Decode(&eventProgramID); err != nil {
		return false
	}

	return eventProgramID.Equals(programID)
}
