package model_test

import (
	"testing"

	"github.com/Goboolean/fetch-system.IaC/pkg/model"
	"github.com/stretchr/testify/assert"
)



func TestFormat(t *testing.T) {
	test := []struct {
		name  string
		str   string
		match bool
	}{
		{
			str: "stock.google.usa.1m",
			match: true,
		},
		{
			str: "stock.amazon.usa.1s",
			match: true,
		},
		{
			str: "test.goboolean.kor.r",
			match: false,
		},
		{
			str: "test.goboolean.kor.t",
			match: true,
		},
		{
			str: "stock.goboolean.kr.r",
			match: false,
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			got := model.IsSymbolValid(tt.str)
			assert.Equal(t, tt.match, got)
		})
	}
}