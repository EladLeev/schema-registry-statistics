package utils

import (
	"io/ioutil"
	"os"
	"reflect"
	"sync"
	"testing"
)

func TestDumpStats(t *testing.T) {
	testCases := []struct {
		name string
		stat ResultStats
		path string
		want string
	}{
		{
			name: "test case 1",
			stat: ResultStats{
				ResultStore: map[uint32][]int{
					1: {1, 2, 3},
					2: {4, 5, 6},
				},
			},
			path: "test1.json",
			want: `{"1":[1,2,3],"2":[4,5,6]}`,
		},
		{
			name: "test case 2",
			stat: ResultStats{
				ResultStore: map[uint32][]int{
					3: {7, 8, 9},
					4: {10, 11, 12},
				},
			},
			path: "test2.json",
			want: `{"3":[7,8,9],"4":[10,11,12]}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			DumpStats(tc.stat, tc.path)
			got, err := ioutil.ReadFile(tc.path)
			if err != nil {
				t.Fatalf("unexpected error reading file: %v", err)
			}
			if string(got) != tc.want {
				t.Errorf("got %q, want %q", string(got), tc.want)
			}
			// Clean up
			if err := os.Remove(tc.path); err != nil {
				t.Fatalf("unexpected error deleting file: %v", err)
			}
		})
	}
}

func TestAppendResult(t *testing.T) {
	l := sync.RWMutex{}
	testCases := []struct {
		name       string
		stat       ResultStats
		offset     int64
		schemaId   uint32
		wantResult ResultStats
	}{
		{
			name: "test case 1",
			stat: ResultStats{
				ResultStore: map[uint32][]int{
					1: {1, 2, 3},
				},
			},
			offset:   4,
			schemaId: 1,
			wantResult: ResultStats{
				ResultStore: map[uint32][]int{
					1: {1, 2, 3, 4},
				},
			},
		},
		{
			name: "test case 2",
			stat: ResultStats{
				ResultStore: map[uint32][]int{
					2: {5, 6, 7},
				},
			},
			offset:   8,
			schemaId: 2,
			wantResult: ResultStats{
				ResultStore: map[uint32][]int{
					2: {5, 6, 7, 8},
				},
			},
		},
		{
			name: "test case 3",
			stat: ResultStats{
				ResultStore: map[uint32][]int{
					3: {9, 10, 11},
				},
			},
			offset:   12,
			schemaId: 3,
			wantResult: ResultStats{
				ResultStore: map[uint32][]int{
					3: {9, 10, 11, 12},
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			AppendResult(tc.stat, tc.offset, tc.schemaId, &l)
			if !reflect.DeepEqual(tc.stat, tc.wantResult) {
				t.Errorf("got %+v, want %+v", tc.stat, tc.wantResult)
			}
		})
	}
}
