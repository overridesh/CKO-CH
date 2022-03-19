package model

import "testing"

func TestNewStatus(t *testing.T) {
	tests := []struct {
		name   string
		input  Status
		expect Status
	}{
		{
			name:   "NewStatus_failed",
			input:  NewStatus("failed"),
			expect: Failed,
		},
		{
			name:   "NewStatus_Approved",
			input:  NewStatus("aPprOved"),
			expect: Approved,
		},
		{
			name:   "NewStatus_RANDOM",
			input:  NewStatus("RANDOM"),
			expect: Pending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", tt.input, tt.expect)
			}
		})
	}
}

func TestIsApproved(t *testing.T) {
	tests := []struct {
		name   string
		input  Status
		expect bool
	}{
		{
			name:   "Approved_IsApproved",
			input:  Approved,
			expect: true,
		},
		{
			name:   "Failed_IsApproved",
			input:  Failed,
			expect: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isApproved := tt.input.IsApproved()
			if isApproved != tt.expect {
				t.Errorf("expect values are equals, but got diferent, output: %v, expect: %v", isApproved, tt.expect)
			}
		})
	}
}
