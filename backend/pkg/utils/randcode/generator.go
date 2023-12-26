package randcode

import (
	"math"
	"math/rand"
)

type GenType int

const (
	TYPE_DEFAULT = iota
	TYPE_DIGIT
	TYPE_LETTER
	TYPE_MIXED
)

func GenVerifyCode(length uint32, codeType GenType) string {
	switch codeType {
	case TYPE_DEFAULT:
		fallthrough
	case TYPE_DIGIT:
		return genVerifyCode([]byte("0123456789"), int(length))
	case TYPE_LETTER:
		return genVerifyCode([]byte("abcdefghijklmnopqrstuvwxyz"), int(length))
	case TYPE_MIXED:
		return genVerifyCode([]byte("0123456789abcdefghijklmnopqrstuvwxyz"), int(length))
	default:
		return "UNKNOWN_TYPE"
	}
}

func genVerifyCode(chars []byte, length int) string {
	charNums := len(chars)
	maskBitLength := 0
	if charNums >= 1 && charNums <= 2 {
		maskBitLength = 1
	} else if charNums >= 3 && charNums <= 4 {
		maskBitLength = 2
	} else if charNums >= 5 && charNums <= 8 {
		maskBitLength = 3
	} else if charNums >= 9 && charNums <= 16 {
		maskBitLength = 4
	} else {
		for i := 5; i < 16; i++ {
			if int(math.Pow(float64(2), float64(i))) >= charNums {
				maskBitLength = i
				break
			}
		}
	}
	maskBits := 1<<maskBitLength - 1
	res := make([]byte, length)
	for i, cache, times := 0, rand.Int63(), 64/maskBitLength; i < length; {
		if times <= 0 {
			cache = rand.Int63()
			times = 64 / maskBitLength
		}
		idx := cache & int64(maskBits)
		cache = cache >> maskBitLength
		times--
		if idx >= int64(len(chars)) {
			continue
		}
		res[i] = chars[idx]
		i++
	}
	return string(res)
}
