package ethcoder

import (
	"math/big"
	"reflect"
	"regexp"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/pkg/errors"
)

// a port of ethers/utils/solidity.ts

var (
	regexArgBytes  = regexp.MustCompile(`^bytes([0-9]+)$`)
	regexArgNumber = regexp.MustCompile(`^(u?int)([0-9]*)$`)
	regexArgArray  = regexp.MustCompile(`^(.*)\[([0-9]*)\]$`)
)

func solidityArgumentPackHex(typ string, val interface{}, isArray bool) (string, error) {
	b, err := solidityArgumentPack(typ, val, isArray)
	if err != nil {
		return "", err
	}
	h := hexutil.Encode(b)
	return h, nil
}

func padZeros(array []byte, totalLength int) ([]byte, error) {
	if len(array) > totalLength {
		return nil, errors.New("array is larger than total expected length")
	}

	buf := make([]byte, totalLength)
	i := totalLength - 1
	for j := len(array) - 1; j >= 0; j-- {
		buf[i] = array[j]
		i--
	}
	return buf, nil
}

func solidityArgumentPack(typ string, val interface{}, isArray bool) ([]byte, error) {
	switch typ {
	case "address":
		v, ok := val.(common.Address)
		if !ok {
			return nil, errors.New("not an common.Address")
		}
		b := v.Bytes()
		if isArray {
			return padZeros(b, 32)
		}
		return b, nil

	case "string":
		v, ok := val.(string)
		if !ok {
			return nil, errors.New("not a string")
		}
		h := hexutil.Encode([]byte(v))
		b, err := hexutil.Decode(h)
		if err != nil {
			return nil, err
		}
		return b, nil

	case "bytes":
		b, ok := val.([]byte)
		if !ok {
			return nil, errors.New("not a []byte")
		}
		return b, nil

	case "bool":
		v, ok := val.(bool)
		if !ok {
			return nil, errors.New("not a bool")
		}
		var b []byte
		if v {
			b = []byte{1}
		} else {
			b = []byte{0}
		}
		if isArray {
			return padZeros(b, 32)
		}
		return b, nil
	}

	// numbers
	if match := regexArgNumber.FindStringSubmatch(typ); len(match) > 0 {
		size, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return nil, err
		}
		if (size%8 != 0) || size == 0 || size > 256 {
			return nil, errors.Errorf("invalid number type '%s'", typ)
		}
		if isArray {
			size = 256
		}
		v, ok := val.(*big.Int)
		if !ok {
			return nil, errors.New("not a *big.Int")
		}
		b := math.PaddedBigBytes(v, int(size/8))
		return b, nil
	}

	// bytes
	if match := regexArgBytes.FindStringSubmatch(typ); len(match) > 0 {
		size, err := strconv.ParseInt(match[1], 10, 64)
		if err != nil {
			return nil, err
		}
		if size == 0 || size > 32 {
			return nil, errors.Errorf("invalid number type '%s'", typ)
		}

		if isArray {
			// if (isArray) { return arrayify((value + Zeros).substring(0, 66)); }
			return nil, errors.New("unsupported, file ticket.")
		}

		rv := reflect.ValueOf(val)
		if rv.Type().Kind() != reflect.Array && rv.Type().Kind() != reflect.Slice {
			return nil, errors.New("not an array")
		}
		if rv.Type().Elem().Kind() != reflect.Uint8 {
			return nil, errors.New("not a byte array")
		}
		if rv.Len() != int(size) {
			return nil, errors.Errorf("not a [%d]byte", size)
		}

		v := make([]byte, size, size)
		var ok bool
		for i := 0; i < int(size); i++ {
			v[i], ok = rv.Index(i).Interface().(byte)
			if !ok {
				return nil, errors.New("unable to set byte")
			}
		}
		return v, nil
	}

	// arrays
	if match := regexArgArray.FindStringSubmatch(typ); len(match) > 0 {
		baseTyp := match[1]
		if match[2] == "" {
			match[2] = "0"
		}
		count, err := strconv.ParseInt(match[2], 10, 64)
		if err != nil {
			return nil, err
		}

		rv := reflect.ValueOf(val)
		if rv.Type().Kind() != reflect.Array && rv.Type().Kind() != reflect.Slice {
			return nil, errors.New("not an array")
		}
		size := rv.Len()
		if count > 0 && size != int(count) {
			return nil, errors.Errorf("array size does not match required size of %d", count)
		}

		buf := []byte{}
		for i := 0; i < size; i++ {
			b, err := solidityArgumentPack(baseTyp, rv.Index(i).Interface(), true)
			if err != nil {
				return nil, err
			}
			buf = append(buf, b...)
		}

		return buf, nil
	}

	return nil, errors.Errorf("unknown type '%s'", typ)
}
