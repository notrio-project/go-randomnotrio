package randomnotrio

import (
	"bytes"
	"encoding/hex"
	"runtime"
	"testing"
)

var tests = [][3][]byte{
	{
		[]byte("test key 000"),
		[]byte("This is a test"),
		[]byte("11832c48e712bb71b46112f8ca62b6cf926146acd751f6c523e9b31d8c766551"),
	},
	{
		[]byte("test key 000"),
		[]byte("Lorem ipsum dolor sit amet"),
		[]byte("f78a38ecbb737b68aafc38060dad3aa4251576d0d1bfe8dc851b39c6ace2d882"),
	},
	{
		[]byte("test key 000"),
		[]byte("sed do eiusmod tempor incididunt ut labore et dolore magna aliqua"),
		[]byte("14596e4dfa62c57ce7b32a3b2bab5c9056f77876078a8f4bb40597c3dbccab5a"),
	},
	{
		[]byte("test key 001"),
		[]byte("sed do eiusmod tempor incididunt ut labore et dolore magna aliqua"),
		[]byte("9c846d32bbeaa3df31a7cccba8d29d9a6d9a7a9a4233fb6c106c3018455b6ceb"),
	},
}

func TestMain(t *testing.T) {
	flags := GetFlags()
	runtime.GOMAXPROCS(runtime.NumCPU())

	for _, tp := range tests {
		cache, _ := AllocCache(flags)
		t.Log("AllocCache finished")
		InitCache(cache, tp[0])
		t.Log("InitCache finished")

		vm, _ := CreateLightVM(cache, flags)
		hash := CalculateHash(vm, tp[1])

		var hashCorrect = make([]byte, hex.DecodedLen(len(tp[2])))
		_, err := hex.Decode(hashCorrect, tp[2])
		if err != nil {
			t.Log(err)
		}

		if !bytes.Equal(hash[:], hashCorrect) {
			t.Logf("incorrect hash: expected %x, got %x", hashCorrect, hash)
			t.Fail()
		}

		DestroyVM(vm)
		ReleaseCache(cache)
	}
}

func BenchmarkCalculateHash(b *testing.B) {
	flags := GetFlags()
	cache, _ := AllocCache(flags | FlagFullMEM)
	ds, _ := AllocDataset(flags | FlagFullMEM)
	InitCache(cache, []byte("123"))
	var workerNum = uint64(runtime.NumCPU())
	InitDatasetMultithread(ds, cache, workerNum)
	vm, _ := CreateVM(cache, ds, flags|FlagFullMEM)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		CalculateHash(vm, []byte("123"))
	}

	DestroyVM(vm)
}
