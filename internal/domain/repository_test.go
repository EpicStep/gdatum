// Copyright 2025 Stepan Rabotkin.
// SPDX-License-Identifier: Apache-2.0.

package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestListServerSummariesParams_Validate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		params  ListServerSummariesParams
		wantErr bool
	}{
		{
			name: "Valid",
			params: ListServerSummariesParams{
				Limit:  100,
				Offset: 100,
			},
			wantErr: false,
		},
		{
			name: "InvalidLimitIsZero",
			params: ListServerSummariesParams{
				Limit:  0,
				Offset: 100,
			},
			wantErr: true,
		},
		{
			name: "InvalidLimitLessThenZero",
			params: ListServerSummariesParams{
				Limit:  -1,
				Offset: 100,
			},
			wantErr: true,
		},
		{
			name: "InvalidOffset",
			params: ListServerSummariesParams{
				Limit:  100,
				Offset: -1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tt.wantErr, tt.params.Validate() != nil)
		})
	}
}

func TestTimeRange_Validate(t *testing.T) {
	t.Parallel()

	testTime := time.Now().Truncate(time.Hour)

	tests := []struct {
		name      string
		timeRange TimeRange
		maxDelta  time.Duration
		wantErr   bool
	}{
		{
			name: "Valid",
			timeRange: TimeRange{
				From: testTime,
				To:   testTime.Add(time.Hour),
			},
			wantErr: false,
		},
		{
			name: "ValidWithMaxDelta",
			timeRange: TimeRange{
				From: testTime,
				To:   testTime.Add(time.Hour),
			},
			maxDelta: time.Hour * 2,
			wantErr:  false,
		},
		{
			name: "InvalidToBeforeFrom",
			timeRange: TimeRange{
				From: testTime.Add(time.Hour),
				To:   testTime,
			},
			wantErr: true,
		},
		{
			name: "InvalidMaxDeltaOverflow",
			timeRange: TimeRange{
				From: testTime,
				To:   testTime.Add(time.Hour),
			},
			maxDelta: time.Minute,
			wantErr:  true,
		},
		{
			name: "InvalidFromIsZero",
			timeRange: TimeRange{
				From: time.Time{},
				To:   testTime,
			},
			wantErr: true,
		},
		{
			name: "InvalidToIsZero",
			timeRange: TimeRange{
				From: testTime,
				To:   time.Time{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.wantErr, tt.timeRange.Validate(tt.maxDelta) != nil)
		})
	}
}

func TestListServerStatisticsParams_Validate(t *testing.T) {
	t.Parallel()

	testTime := time.Now().Truncate(time.Hour)

	tests := []struct {
		name    string
		params  ListServerStatisticsParams
		wantErr bool
	}{
		{
			name: "Valid",
			params: ListServerStatisticsParams{
				TimeRange: TimeRange{
					From: testTime,
					To:   testTime.Add(time.Hour),
				},
			},
			wantErr: false,
		},
		{
			name:    "Invalid",
			params:  ListServerStatisticsParams{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.wantErr, tt.params.Validate() != nil)
		})
	}
}
