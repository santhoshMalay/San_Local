package keygen_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/zhuravlev-pe/course-watch/pkg/keygen"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	tc := map[string]struct {
		secret1     string
		context1    string
		bytesCount1 int
		secret2     string
		context2    string
		bytesCount2 int
		mustBeEqual bool
	}{
		"same, no context": {
			secret1:     "mySecret",
			context1:    "",
			bytesCount1: 16,
			secret2:     "mySecret",
			context2:    "",
			bytesCount2: 16,
			mustBeEqual: true,
		},
		"same, with context": {
			secret1:     "mySecret",
			context1:    "JWT signature key",
			bytesCount1: 16,
			secret2:     "mySecret",
			context2:    "JWT signature key",
			bytesCount2: 16,
			mustBeEqual: true,
		},
		"different secrets": {
			secret1:     "mySecret1",
			context1:    "",
			bytesCount1: 16,
			secret2:     "mySecret2",
			context2:    "",
			bytesCount2: 16,
			mustBeEqual: false,
		},
		"different contexts": {
			secret1:     "mySecret",
			context1:    "JWT signature key",
			bytesCount1: 16,
			secret2:     "mySecret",
			context2:    "encryption key",
			bytesCount2: 16,
			mustBeEqual: false,
		},
		"different sizes": {
			secret1:     "mySecret",
			context1:    "JWT signature key",
			bytesCount1: 16,
			secret2:     "mySecret",
			context2:    "JWT signature key",
			bytesCount2: 8,
			mustBeEqual: false,
		},
	}
	for name, c := range tc {
		t.Run(name, func(t *testing.T) {
			t.Log("Case", name)
			key1, err := keygen.Generate(c.secret1, c.context1, c.bytesCount1)
			assert.NoError(t, err)
			assert.Len(t, key1, c.bytesCount1)
			t.Logf("\tKey 1: %x", key1)
			key2, err := keygen.Generate(c.secret2, c.context2, c.bytesCount2)
			assert.NoError(t, err)
			assert.Len(t, key2, c.bytesCount2)
			t.Logf("\tKey 2: %x", key2)
			if c.mustBeEqual {
				assert.Equal(t, key1, key2)
			} else {
				assert.NotEqual(t, key1, key2)
			}
		})
	}
}

func TestKeyGen_Read(t *testing.T) {
	const secret = "mySecret"
	const context = "JWT signature key"
	const keySize = 16
	const keyCount = 3

	generateKeys := func(kg *keygen.KeyGen, keys [][]byte) {
		for i := 0; i < keyCount; i++ {
			key := make([]byte, keySize)
			n, err := kg.Read(key)
			assert.Equal(t, keySize, n)
			assert.NoError(t, err)
			keys[i] = key
			t.Logf("\t%x", key)
		}
	}

	kg1, err := keygen.New(secret, context)
	assert.NoError(t, err)
	keys1 := make([][]byte, keyCount)
	t.Log("Keys1:")
	generateKeys(kg1, keys1)

	kg2, err := keygen.New(secret, context)
	assert.NoError(t, err)
	keys2 := make([][]byte, keyCount)
	t.Log("Keys2:")
	generateKeys(kg2, keys2)

	for i := 0; i < keyCount; i++ {
		assert.Equal(t, keys1[i], keys2[i])
	}
}
