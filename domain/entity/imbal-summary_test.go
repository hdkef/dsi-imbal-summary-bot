package entity

import (
	"reflect"
	"testing"
	"time"
)

func TestImbalSummary_GetImbal(t *testing.T) {
	type fields struct {
		imbals []Imbal
	}

	tests := []struct {
		name   string
		fields fields
		want   []MonthlySummaryResult
		want1  float64
	}{
		{
			name: "should be ok",
			fields: fields{
				imbals: []Imbal{
					{
						date:   time.Now().Add(62 * 24 * time.Hour),
						amount: 10000,
					},
					{
						date:   time.Now().Add(62 * 24 * time.Hour),
						amount: 10,
					},
					{
						date:   time.Now(),
						amount: 5000,
					},
					{
						date:   time.Now(),
						amount: 7000,
					},
					{
						date:   time.Now().Add(-62 * 24 * time.Hour),
						amount: 10000,
					},
					{
						date:   time.Now().Add(-62 * 24 * time.Hour),
						amount: 10,
					},
					{
						date:   time.Now().Add(31 * 24 * time.Hour),
						amount: 8000,
					},
					{
						date:   time.Now().Add(31 * 24 * time.Hour),
						amount: 9000,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ImbalSummary{}

			for _, v := range tt.fields.imbals {
				i.SetImbal(v)
			}

			got, got1 := i.GetImbal()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ImbalSummary.GetImbal() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ImbalSummary.GetImbal() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
