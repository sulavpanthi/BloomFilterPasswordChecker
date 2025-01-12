package entity

import (
	"encoding/gob"
	"math"
	"os"

	"github.com/rs/zerolog"
	"github.com/sulavpanthi/BloomFilterPasswordChecker/pkg/appcontext"
)

type BloomFilter struct {
	BitArray          []bool
	BitArraySize      uint64
	HashFunctionCount uint64
}

var log *zerolog.Logger

func getBitArraySize(expectedElements uint64, falsePositiveProbability float64) uint64 {
	m := -(float64(expectedElements) * math.Log(falsePositiveProbability)) / math.Pow(math.Log(2), 2)
	return uint64(math.Ceil(m))
}

func getHashFunctionsCount(expectedElements, BitArraySize uint64) uint64 {
	HashFunctionCount := float64(BitArraySize) / float64(expectedElements) * math.Log(2)
	return uint64(math.Ceil(HashFunctionCount))
}

func New(expectedElements uint64, falsePositiveProbability float64) *BloomFilter {

	log := appcontext.Get().Logger
	log.Info().Msg("Bloom filter initializing now......")

	BitArraySize := getBitArraySize(expectedElements, falsePositiveProbability)
	HashFunctionCount := getHashFunctionsCount(expectedElements, BitArraySize)
	BitArray := make([]bool, BitArraySize)

	bloomFilter := &BloomFilter{
		BitArray:          BitArray,
		BitArraySize:      BitArraySize,
		HashFunctionCount: HashFunctionCount,
	}
	log.Info().Msg("Bloom filter initialized successfully")
	return bloomFilter
}

func fnv1Hash64(input string, seed uint64) uint64 {
	log := appcontext.Get().Logger
	// TODO: Move this prime value to config
	prime := uint64(0x100000001B3) // FNV-1 64-bit prime
	hash := seed

	for i := 0; i < len(input); i++ {
		hash ^= uint64(input[i])
		hash *= prime
	}
	log.Info().Uint64("hash before returning", hash).Msg("log := appcontext.Get().Logger")
	return hash
}

func (bf *BloomFilter) Add(word string) {

	log := appcontext.Get().Logger

	// TODO: Use strategy design pattern here
	var i uint64
	for ; i < bf.HashFunctionCount; i++ {
		hash := fnv1Hash64(word, i)
		bit := hash % bf.BitArraySize
		log.Info().Uint64("bit", bit).Msg("Bit to set is- ")
		bf.BitArray[bit] = true
	}
}

func (bf *BloomFilter) Check(word string) bool {

	// TODO: Use strategy design pattern here
	var i uint64
	for ; i < bf.HashFunctionCount; i++ {
		hash := fnv1Hash64(word, i)
		bit := hash % bf.BitArraySize
		log.Info().Uint64("bit", bit).Msg("Bit to set is- ")
		if !bf.BitArray[bit] {
			return false
		}
	}

	return true
}

func (bf *BloomFilter) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		log.Error().Err(err).Msg("Cannot create file: ")
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(bf)
	if err != nil {
		log.Error().Err(err).Msg("error encoding Bloom filter: ")
		return err
	}

	return nil
}

func Load(filename string) (*BloomFilter, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("error opening file")
		return nil, err
	}
	defer file.Close()

	b := &BloomFilter{}
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(b)
	if err != nil {
		log.Error().Err(err).Msg("error decoding Bloom filter:")
		return nil, err
	}

	return b, nil
}
