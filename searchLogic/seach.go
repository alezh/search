package searchLogic

import (

	"unicode"
)

func init()  {
	cryptTable = make(map[uint64]uint64)
}
//hash索引库
type HashTable struct {
	NHashA  uint64
	NHashB  uint64
	Ids     []int
	BExists int
}

const (
	HASH_OFFSET = iota     //0
	HASH_A                 //1
	HASH_B                 //2
)

var cryptTable map[uint64]uint64

var Table map[uint64]*HashTable

var nTableSize uint64



func MPQHashTableInit()  {
	var i uint64 = 0

	initCryptTable()

	Table = make(map[uint64]*HashTable)

	nTableSize = uint64(len(cryptTable))

	for i = 0; i < nTableSize; i++{
		hash := HashTable{
			NHashA: 0,
			NHashB: 0,
			Ids:make([]int,0),
			BExists:0,
		}
		Table[i]=&hash
	}
}

// 暴雪hash算法 长度0x100 * 5  0x500    int 1280
func initCryptTable() {
	var seed, index1, index2 uint64 = 0x00100001, 0, 0
	i := 0
	for index1 = 0; index1 < 0x100; index1 += 1 {
		for index2, i = index1, 0; i < 10; index2 += 0x100 {
			seed = (seed*125 + 3) % 0x2aaaab
			temp1 := (seed & 0xffff) << 0x10
			seed = (seed*125 + 3) % 0x2aaaab
			temp2 := seed & 0xffff
			cryptTable[index2] = temp1 | temp2
			i += 1
		}
	}
}
// 计算字符串的哈希值
func HashString(lpszString string, dwHashType uint64) uint64 {
	var i = 0
	var ch uint64 = 0
	var seed1, seed2 uint64 = 0x7FED7FED, 0xEEEEEEEE
	var key uint8
	strLen := len(lpszString)
	for i < strLen {
		key = lpszString[i]
		ch = uint64(unicode.ToUpper(rune(key)))
		i += 1
		seed1 = cryptTable[(dwHashType << 8) + ch ] ^ (seed1 + seed2)
		seed2 = uint64(ch) + seed1 + seed2 + (seed2 << 5) + 3
	}
	return uint64(seed1)
}

//获取是否存在
func GetHashTableIsExist(lpszString string) ([]int , bool)  {
	nHash := HashString( lpszString, HASH_OFFSET )
	nHashA := HashString( lpszString, HASH_A )
	nHashB := HashString( lpszString, HASH_B )
	nHashStart := nHash % nTableSize
	nHashPos := nHashStart

	if Table[nHashPos].BExists >0 {
		if Table[nHashPos].NHashA == nHashA && Table[nHashPos].NHashB == nHashB {
			return Rm_duplicate(Table[nHashPos].Ids),true
		}else {
			nHashPos = (nHashPos +1) % nTableSize
		}
		if nHashPos == nHashStart{
			return nil , false
		}
	}
	return nil , false
}

func Rm_duplicate(list []int) []int {
	var x []int = []int{}
	for _, i := range list {
		if len(x) == 0 {
			x = append(x, i)
		} else {
			for k, v := range x {
				if i == v {
					break
				}
				if k == len(x)-1 {
					x = append(x, i)
				}
			}
		}
	}
	return x
}

func InsertString(lpszString string,id int) {

	nHash := HashString( lpszString, HASH_OFFSET )
	nHashA := HashString( lpszString, HASH_A )
	nHashB := HashString( lpszString, HASH_B )
	nHashStart := nHash % nTableSize
	nHashPos := nHashStart

	if Table[nHashPos].BExists >0 {
		//插入ID
		hash := Table[nHashPos]
		hash.Ids = append(Table[nHashPos].Ids,id)
		Table[nHashPos] = hash

		nHashPos = (nHashPos + 1) % nTableSize
		if nHashPos == nHashStart{
			return
		}
	}

	hash := HashTable{
		NHashA: nHashA,
		NHashB: nHashB,
		BExists:1,
	}
	hash.Ids = append(hash.Ids,id)

	Table[nHashPos] = &hash
	return
}
