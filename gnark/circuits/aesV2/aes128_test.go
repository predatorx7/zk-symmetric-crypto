package aes_v2

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/test"
)

func TestAES128(t *testing.T) {
	assert := test.NewAssert(t)

	key := "7E24067817FAE0D743D6CE1F32539163"
	plaintext := "000102030405060708090A0B0C0D0E0F101112131415161718191A1B1C1D1E1F000102030405060708090A0B0C0D0E0F101112131415161718191A10A0B0C0D0E0F10111213141516171819B1C1D1E1F"
	// ciphertext := "5104A106168A72D9790D41EE8EDAD388EB2E1EFC46DA57C8FCE630DF9141BE28"
	Nonce := "006CB6DBC0543B59DA48D90B"
	Counter := 1

	// calculate ciphertext ourselves
	block, err := aes.NewCipher(mustHex(key))
	if err != nil {
		panic(err)
	}
	ctr := cipher.NewCTR(block, append(mustHex(Nonce), binary.BigEndian.AppendUint32(nil, uint32(Counter))...))
	ciphertext := make([]byte, len(mustHex(plaintext)))
	ctr.XORKeyStream(ciphertext, mustHex(plaintext))

	keyAssign := StrToIntSlice(key, true)
	ptAssign := StrToIntSlice(plaintext, true)
	ctAssign := ciphertext // StrToIntSlice(ciphertext, true)
	nonceAssign := StrToIntSlice(Nonce, true)

	// witness values preparation
	assignment := AESCircuit{
		AESBaseCircuit: AESBaseCircuit{
			Key:     make([]frontend.Variable, 16),
			Counter: Counter,
			Nonce:   [12]frontend.Variable{},
			In:      [BLOCKS * 16]frontend.Variable{},
		},
		Out: [BLOCKS * 16]frontend.Variable{},
	}

	// assign values here because required to use make in assignment
	for i := 0; i < len(keyAssign); i++ {
		assignment.Key[i] = keyAssign[i]
	}
	for i := 0; i < len(ptAssign); i++ {
		assignment.In[i] = ptAssign[i]
	}
	for i := 0; i < len(ctAssign); i++ {
		assignment.Out[i] = ctAssign[i]
	}

	for i := 0; i < len(nonceAssign); i++ {
		assignment.Nonce[i] = nonceAssign[i]
	}

	assert.CheckCircuit(&AESCircuit{
		AESBaseCircuit: AESBaseCircuit{
			Key:     make([]frontend.Variable, 16),
			Counter: Counter,
			Nonce:   [12]frontend.Variable{},
			In:      [BLOCKS * 16]frontend.Variable{},
		},
		Out: [BLOCKS * 16]frontend.Variable{},
	}, test.WithValidAssignment(&assignment))
}

func StrToIntSlice(inputData string, hexRepresentation bool) []int {
	var byteSlice []byte
	if hexRepresentation {
		hexBytes, _ := hex.DecodeString(inputData)
		byteSlice = hexBytes
	} else {
		byteSlice = []byte(inputData)
	}

	var data []int
	for i := 0; i < len(byteSlice); i++ {
		data = append(data, int(byteSlice[i]))
	}
	return data
}

func mustHex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return b
}

func TestCompile(t *testing.T) {
	curve := ecc.BN254.ScalarField()

	witness := AESCircuit{
		AESBaseCircuit: AESBaseCircuit{
			Key: make([]frontend.Variable, 16),
		},
	}

	r1css, err := frontend.Compile(curve, r1cs.NewBuilder, &witness)
	if err != nil {
		panic(err)
	}

	fmt.Printf("constraints: %d\n", r1css.GetNbConstraints())

}
