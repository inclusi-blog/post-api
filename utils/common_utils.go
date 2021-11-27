package utils

import (
	"crypto/rand"
	"math/big"
)

func GenRandNum(min, max int64) int64 {
	bg := big.NewInt(max - min)
	//using crypto rand for better unique random number based on system heat, performance, process count, thread pool count etc.,
	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}

	return n.Int64() + min
}
