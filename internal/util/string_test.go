package util

import "testing"

func TestGenKey(t *testing.T) {
	tt := []struct {
		keys []string
		want string
	}{
		{
			keys: []string{"key1", "key2"},
			want: "key1:key2",
		},
		{
			keys: []string{"key1", "key2", "key3"},
			want: "key1:key2:key3",
		},
	}

	for _, tc := range tt {
		got := GenKey(tc.keys...)
		if got != tc.want {
			t.Errorf("GenKey(%v) = %v, want %v", tc.keys, got, tc.want)
		}
	}
}
