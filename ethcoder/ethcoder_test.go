package ethcoder

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
)

func TestAbiCoder(t *testing.T) {
	packed, err := AbiCoder([]string{"uint256[]", "uint256[]"}, []interface{}{[]*big.Int{big.NewInt(44)}, []*big.Int{big.NewInt(22)}})
	assert.NoError(t, err)

	packedHex := hexutil.Encode(packed)

	// ethers.utils.defaultAbiCoder.encode(['uint256[]','uint256[]'], [[44],[22]])
	// "0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002c00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000016"
	expected := `0x000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002c00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000016`

	assert.Equal(t, expected, packedHex)
	assert.Equal(t, len(expected), len(packedHex))
}

func TestSolidityPack(t *testing.T) {
	// string
	{
		// ethers.utils.solidityPack(['string'], ['peϣer'])
		// "0x7065cfa36572"
		h, err := solidityArgumentPackHex("string", "peϣer", false)
		assert.NoError(t, err)
		assert.Equal(t, "0x7065cfa36572", h)
	}

	// address
	{
		// ethers.utils.solidityPack(['address'], ['0x39d28D4c4191a584acabe021F5B905887a6B5247'])
		// "0x39d28d4c4191a584acabe021f5b905887a6b5247"
		h, err := solidityArgumentPackHex("address", common.HexToAddress("0x39d28D4c4191a584acabe021F5B905887a6B5247"), false)
		assert.NoError(t, err)
		assert.Equal(t, "0x39d28d4c4191a584acabe021f5b905887a6b5247", h)
	}

	// bytes
	{
		// ethers.utils.solidityPack(['bytes'], [[0,1,2,3]])
		// "0x00010203"
		h, err := solidityArgumentPackHex("bytes", []byte{0, 1, 2, 3}, false)
		assert.NoError(t, err)
		assert.Equal(t, "0x00010203", h)
	}

	// bool
	{
		// ethers.utils.solidityPack(['bool'], [true])
		// "0x01"
		h, err := solidityArgumentPackHex("bool", true, false)
		assert.NoError(t, err)
		assert.Equal(t, "0x01", h)

		h, err = solidityArgumentPackHex("bool", false, false)
		assert.NoError(t, err)
		assert.Equal(t, "0x00", h)
	}

	// uint256
	{
		// ethers.utils.solidityPack(['uint256'], [55])
		// "0x0000000000000000000000000000000000000000000000000000000000000037"
		h, err := solidityArgumentPackHex("uint256", big.NewInt(55), false)
		assert.NoError(t, err)
		assert.Equal(t, "0x0000000000000000000000000000000000000000000000000000000000000037", h)
	}

	// bytes8
	{
		// ethers.utils.solidityPack(['bytes8'], [[0,1,2,3,4,5,6,7]])
		// "0x0001020304050607"
		h, err := solidityArgumentPackHex("bytes8", [8]byte{0, 1, 2, 3, 4, 5, 6, 7}, false)
		assert.NoError(t, err)
		assert.Equal(t, "0x0001020304050607", h)
	}

	// address[]
	{
		// ethers.utils.solidityPack(['address[]'], [['0x39d28D4c4191a584acabe021F5B905887a6B5247']])
		// "0x00000000000000000000000039d28d4c4191a584acabe021f5b905887a6b5247"
		h, err := solidityArgumentPackHex("address[]", []common.Address{common.HexToAddress("0x39d28D4c4191a584acabe021F5B905887a6B5247")}, false)
		assert.NoError(t, err)
		assert.Equal(t, "0x00000000000000000000000039d28d4c4191a584acabe021f5b905887a6b5247", h)
	}

	// string[]
	{
		// ethers.utils.solidityPack(['string[]'], [['sup','eth']])
		// "0x737570657468"
		h, err := solidityArgumentPackHex("string[]", []string{"sup", "eth"}, false)
		assert.NoError(t, err)
		assert.Equal(t, "0x737570657468", h)
	}

	// bool[]
	{
		// ethers.utils.solidityPack(['bool[]'], [[true,true]])
		// "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001"
		h, err := solidityArgumentPackHex("bool[]", []bool{true, true}, false)
		assert.NoError(t, err)
		assert.Equal(t, "0x00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001", h)
	}
}
