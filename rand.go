package gah

import "hash/fnv"

//TODO use PCG32

// Hash64 returns a uint64 hash of the input string
func Hash64(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return uint64(h.Sum64())
}
