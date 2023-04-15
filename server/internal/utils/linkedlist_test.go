package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLinkedList(t *testing.T) {
	cond := func(current, next int) bool {
		return current < next
	}

	type test struct {
		name string
		args []int
		want []int
	}
	tests := []test{
		{
			name: "insert in ascending order",
			args: []int{1, 2, 3, 4, 5},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "insert in descending order",
			args: []int{5, 4, 3, 2, 1},
			want: []int{1, 2, 3, 4, 5},
		},
		{
			name: "insert in random order",
			args: []int{3, 2, 1, 4, 5},
			want: []int{1, 2, 3, 4, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ll := NewLinkedList[int](cond)
			ll.InsertMany(tt.args)

			var got []int
			iter := ll.Iterator()
			for iter.HasNext() {
				v, err := iter.Next()
				if err != nil {
					break
				}
				got = append(got, v)
			}
			assert.Equal(t, tt.want, got)
			assert.False(t, iter.HasNext())
		})
	}
}
