package minhash

import (
	"github.com/dgryski/go-farm"
	"github.com/dgryski/go-spooky"
	"github.com/shawnohare/go-minhash"
)

func CompareByMinHash(set1, set2 []string) float64 {
	h1 := spooky.Hash64
	h2 := farm.Hash64
	size := 100

	mw1 := minhash.New(h1, h2, size)
	mw2 := minhash.New(h1, h2, size)

	for _, x := range set1 {
		mw1.Push(x)
	}

	for _, x := range set2 {
		mw2.Push(x)
	}

	return mw1.Similarity(mw2)
}
