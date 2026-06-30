package calc

import "testing"

// 問題1: テーブル駆動テストを書いてください

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		// ここにテストケースを追加してください
		// {"正の数", 2, 3, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// TestSubtract, TestMultiply, TestDivide も実装してください
