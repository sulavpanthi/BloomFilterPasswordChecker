package usecase

import (
	"encoding/json"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/sulavpanthi/BloomFilterPasswordChecker/internal/entity"
	"github.com/sulavpanthi/BloomFilterPasswordChecker/pkg/appcontext"
)

type BloomFilterUseCase struct {
	BloomFilter *entity.BloomFilter
}

type BloomFilterJSON struct {
	BitArray          []int  `json:"bitArray"`
	BitArraySize      uint64 `json:"bitArraySize"`
	HashFunctionCount uint64 `json:"hashFunctionCount"`
}

func NewBloomFilterUseCase(size uint64, fp float64) *BloomFilterUseCase {
	bf := entity.New(size, fp)
	return &BloomFilterUseCase{
		BloomFilter: bf,
	}
}

func InitBloomFilterUseCase() *BloomFilterUseCase {
	var bloomFilterUseCase *BloomFilterUseCase
	config := appcontext.Get().Config
	log := appcontext.Get().Logger

	_, err := os.Open(config.BloomFilterFileName)
	if err == nil {
		log.Info().Msg("Bloom filter file exists")
		bloomFilter, _ := LoadJSON(config.BloomFilterFileName)
		bloomFilterUseCase = &BloomFilterUseCase{
			BloomFilter: bloomFilter,
		}
	} else {
		log.Info().Msg("Bloom filter file does not exist")
		bloomFilterUseCase = NewBloomFilterUseCase(config.ExpectedElements, config.FalsePositiveProbability)
	}

	log.Info().Interface("bloom_filter_use_case", *bloomFilterUseCase).Msg("Bloom filter being used")
	return bloomFilterUseCase
}

func (uc *BloomFilterUseCase) AddPassword(password string) {
	uc.BloomFilter.Add(password)
	config := appcontext.Get().Config
	uc.SaveAsJSON(config.BloomFilterFileName)
}

func (uc *BloomFilterUseCase) IsPasswordCommon(password string) bool {
	return uc.BloomFilter.Check(password)
}

func (uc *BloomFilterUseCase) SerializeAsJSON() *BloomFilterJSON {
	bitArrayInt := make([]int, len(uc.BloomFilter.BitArray))
	for i, bit := range uc.BloomFilter.BitArray {
		if bit {
			bitArrayInt[i] = 1
		} else {
			bitArrayInt[i] = 0
		}
	}

	jsonBF := BloomFilterJSON{
		BitArray:          bitArrayInt,
		BitArraySize:      uc.BloomFilter.BitArraySize,
		HashFunctionCount: uc.BloomFilter.HashFunctionCount,
	}
	return &jsonBF
}

func (uc *BloomFilterUseCase) SaveAsJSON(filename string) error {

	jsonBF := uc.SerializeAsJSON()
	file, err := os.Create(filename)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create file")
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(jsonBF); err != nil {
		log.Error().Err(err).Msg("Failed to encode bloom filter")
		return err
	}

	return nil
}

func LoadJSON(filename string) (*entity.BloomFilter, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Error().Err(err).Msg("failed to open file")
		return nil, err
	}
	defer file.Close()

	var jsonBF BloomFilterJSON
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&jsonBF); err != nil {
		log.Error().Err(err).Msg("failed to decode JSON")
		return nil, err
	}

	bitArray := make([]bool, len(jsonBF.BitArray))
	for i, bit := range jsonBF.BitArray {
		bitArray[i] = bit == 1
	}

	bf := &entity.BloomFilter{
		BitArray:          bitArray,
		BitArraySize:      jsonBF.BitArraySize,
		HashFunctionCount: jsonBF.HashFunctionCount,
	}

	return bf, nil
}
