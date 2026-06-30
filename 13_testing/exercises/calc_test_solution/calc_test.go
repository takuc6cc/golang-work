package calc

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"正の数", 2, 3, 5},
		{"負の数", -1, -2, -3},
		{"ゼロ", 0, 0, 0},
		{"正と負", 5, -3, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{5, 3, 2},
		{0, 5, -5},
		{-3, -2, -1},
	}
	for _, tt := range tests {
		if got := Subtract(tt.a, tt.b); got != tt.want {
			t.Errorf("Subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestMultiply(t *testing.T) {
	tests := []struct{ a, b, want int }{
		{3, 4, 12},
		{-2, 5, -10},
		{0, 100, 0},
	}
	for _, tt := range tests {
		if got := Multiply(tt.a, tt.b); got != tt.want {
			t.Errorf("Multiply(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestDivide(t *testing.T) {
	t.Run("正常系", func(t *testing.T) {
		tests := []struct{ a, b, want int }{
			{10, 2, 5},
			{9, 3, 3},
			{7, 2, 3},
		}
		for _, tt := range tests {
			got, err := Divide(tt.a, tt.b)
			if err != nil {
				t.Fatalf("予期しないエラー: %v", err)
			}
			if got != tt.want {
				t.Errorf("Divide(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		}
	})

	t.Run("ゼロ除算", func(t *testing.T) {
		_, err := Divide(10, 0)
		if err == nil {
			t.Error("エラーが返るはずですが、nil が返りました")
		}
	})
}
