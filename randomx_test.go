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
		[]byte("7266bdd4668d0718f9f8817af9a251a7b8ed49468a00717c74584ba2ac7e3e09"),
	},
	{
		[]byte("test key 000"),
		[]byte("Lorem ipsum dolor sit amet"),
		[]byte("bcdc7bce23f615e00d4369bec3c05cd7add07fc25c49529e7671582133d025f0"),
	},
	{
		[]byte("test key 000"),
		[]byte("sed do eiusmod tempor incididunt ut labore et dolore magna aliqua"),
		[]byte("c64944ed33e4abd159fb8692c81a923bbe7a7c346fd39349f76fb53fa284bc87"),
	},
	{
		[]byte("test key 001"),
		[]byte("sed do eiusmod tempor incididunt ut labore et dolore magna aliqua"),
		[]byte("6b339ad028e772e91290df111e61fa82af0eef1dfc53dfbcb04a98945b01dca6"),
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
