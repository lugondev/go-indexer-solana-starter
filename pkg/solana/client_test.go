package solana

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		rpcURL  string
		wsURL   string
		wantErr bool
	}{
		{
			name:    "valid URLs",
			rpcURL:  "https://api.mainnet-beta.solana.com",
			wsURL:   "wss://api.mainnet-beta.solana.com",
			wantErr: false,
		},
		{
			name:    "empty RPC URL",
			rpcURL:  "",
			wsURL:   "wss://api.mainnet-beta.solana.com",
			wantErr: true,
		},
		{
			name:    "empty WS URL is ok",
			rpcURL:  "https://api.mainnet-beta.solana.com",
			wsURL:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewClient(tt.rpcURL, tt.wsURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}
