package decoder

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gagliardetto/solana-go"
	"github.com/lugondev/go-indexer-solana-starter/internal/models"
)

type CounterLogParser struct {
	programID solana.PublicKey
}

func NewCounterLogParser(programID solana.PublicKey) *CounterLogParser {
	return &CounterLogParser{
		programID: programID,
	}
}

func (p *CounterLogParser) ParseLogs(logs []string, accounts []solana.PublicKey) ([]CounterAction, error) {
	var actions []CounterAction

	for _, log := range logs {
		if !strings.Contains(log, "Program log:") {
			continue
		}

		action := p.parseLogMessage(log, accounts)
		if action != nil {
			actions = append(actions, *action)
		}
	}

	return actions, nil
}

type CounterAction struct {
	Type         models.EventType
	Counter      solana.PublicKey
	Authority    *solana.PublicKey
	OldValue     *uint64
	NewValue     *uint64
	AddedValue   *uint64
	Payer        *solana.PublicKey
	FeeCollector *solana.PublicKey
	Payment      *uint64
}

func (p *CounterLogParser) parseLogMessage(log string, accounts []solana.PublicKey) *CounterAction {
	msgPrefix := "Program log: "
	if !strings.Contains(log, msgPrefix) {
		return nil
	}

	msg := strings.TrimSpace(strings.Split(log, msgPrefix)[1])

	var counter solana.PublicKey
	if len(accounts) > 0 {
		counter = accounts[0]
	}

	if msg == "Counter initialized" {
		return &CounterAction{
			Type:     models.EventTypeCounterInitialized,
			Counter:  counter,
			NewValue: uint64Ptr(0),
		}
	}

	if strings.HasPrefix(msg, "Counter incremented to: ") {
		newValue := p.extractNumber(msg, "Counter incremented to: ")
		if newValue != nil {
			oldValue := *newValue - 1
			return &CounterAction{
				Type:     models.EventTypeCounterIncremented,
				Counter:  counter,
				OldValue: &oldValue,
				NewValue: newValue,
			}
		}
	}

	if strings.HasPrefix(msg, "Counter decremented to: ") {
		newValue := p.extractNumber(msg, "Counter decremented to: ")
		if newValue != nil {
			oldValue := *newValue + 1
			return &CounterAction{
				Type:     models.EventTypeCounterDecremented,
				Counter:  counter,
				OldValue: &oldValue,
				NewValue: newValue,
			}
		}
	}

	if strings.HasPrefix(msg, "Added ") && strings.Contains(msg, "to counter") {
		re := regexp.MustCompile(`Added (\d+) to counter\. New value: (\d+)`)
		matches := re.FindStringSubmatch(msg)
		if len(matches) == 3 {
			added, _ := strconv.ParseUint(matches[1], 10, 64)
			newVal, _ := strconv.ParseUint(matches[2], 10, 64)
			oldVal := newVal - added
			return &CounterAction{
				Type:       models.EventTypeCounterAdded,
				Counter:    counter,
				OldValue:   &oldVal,
				AddedValue: &added,
				NewValue:   &newVal,
			}
		}
	}

	if msg == "Counter reset" {
		return &CounterAction{
			Type:     models.EventTypeCounterReset,
			Counter:  counter,
			OldValue: nil,
			NewValue: uint64Ptr(0),
		}
	}

	if strings.HasPrefix(msg, "Payment of ") && strings.Contains(msg, "lamports received") {
		re := regexp.MustCompile(`Payment of (\d+) lamports received\. Counter incremented to: (\d+)`)
		matches := re.FindStringSubmatch(msg)
		if len(matches) == 3 {
			payment, _ := strconv.ParseUint(matches[1], 10, 64)
			newCount, _ := strconv.ParseUint(matches[2], 10, 64)

			var payer, feeCollector *solana.PublicKey
			if len(accounts) > 1 {
				payer = &accounts[1]
			}
			if len(accounts) > 2 {
				feeCollector = &accounts[2]
			}

			return &CounterAction{
				Type:         models.EventTypeCounterPaymentReceived,
				Counter:      counter,
				Payment:      &payment,
				NewValue:     &newCount,
				Payer:        payer,
				FeeCollector: feeCollector,
			}
		}
	}

	return nil
}

func (p *CounterLogParser) extractNumber(msg, prefix string) *uint64 {
	numStr := strings.TrimPrefix(msg, prefix)
	numStr = strings.TrimSpace(numStr)

	num, err := strconv.ParseUint(numStr, 10, 64)
	if err != nil {
		return nil
	}

	return &num
}

func uint64Ptr(v uint64) *uint64 {
	return &v
}

func (p *CounterLogParser) ExtractCounterAccounts(tx interface{}) []solana.PublicKey {
	return nil
}

func IsCounterProgramLog(log string, programID solana.PublicKey) bool {
	programStr := fmt.Sprintf("Program %s invoke", programID.String())
	return strings.Contains(log, programStr) || strings.Contains(log, "Program log:")
}
