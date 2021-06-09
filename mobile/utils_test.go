package mobile

import (
	"math/rand"
	"testing"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/stretchr/testify/require"
)

func randomBytes(s []byte) {
	_, err := rand.Read(s)
	if err != nil {
		panic(err)
	}
}

func TestEncodeDecode(t *testing.T) {
	a := types.Address{}
	for i := 0; i < 1000; i++ {
		randomBytes(a[:])
		addr := a.String()
		b, err := types.DecodeAddress(addr)
		require.NoError(t, err)
		require.Equal(t, a, b)
		require.True(t, IsValidAddress(a.String()))

		require.False(t, IsValidAddress("SONNPE7I3TYWE7VQQA7VCGZK54WXEICYROMWQYXCOB5A5RHM46LVNGZLRU"))

	}
}
