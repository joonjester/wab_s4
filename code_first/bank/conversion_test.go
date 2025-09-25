package bank

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConverCurrency(t *testing.T) {
	tests := map[string]struct {
		base       Currency
		target     Currency
		amount     float64
		mockRate   float64
		wantAmount float64
		wantErr    bool
	}{
		"Happy Path: EUR->USD": {
			base:       EUR,
			target:     USD,
			amount:     100,
			mockRate:   1.2,
			wantAmount: 120,
			wantErr:    false,
		},
		"Happy Path: EUR->JPN": {
			base:       EUR,
			target:     JPN,
			amount:     100,
			mockRate:   0.7,
			wantAmount: 70,
			wantErr:    false,
		},
		"Happy Path: EUR->GBP": {
			base:       EUR,
			target:     GBP,
			amount:     100,
			mockRate:   1.3,
			wantAmount: 130,
			wantErr:    false,
		},
		"Unhappy Path: API return empty rates": {
			base:       EUR,
			target:     USD,
			amount:     100,
			mockRate:   0,
			wantAmount: 0,
			wantErr:    true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.wantErr {
					_ = json.NewEncoder(w).Encode(RatesResponse{
						Rates: map[string]float64{},
					})
					return
				}

				resp := RatesResponse{
					Rates: map[string]float64{
						string(tt.target): tt.mockRate * tt.amount,
					},
					Base: string(tt.base),
				}
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			oldApi := frankfurterAPI
			frankfurterAPI = server.URL
			defer func() { frankfurterAPI = oldApi }()

			got, err := ConvertCurrency(tt.amount, tt.base, tt.target)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexcepted error: %v", err)
			}

			if got == nil || *got != tt.wantAmount {
				t.Errorf("got %v, want %v", got, tt.wantAmount)
			}

		})
	}
}
