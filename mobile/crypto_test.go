package mobile

import (
	"testing"

	"github.com/algorand/go-algorand-sdk/types"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/ed25519"
)

func TestKeyGeneration(t *testing.T) {
	sk := GenerateSK()

	pk := ed25519.PrivateKey(sk).Public().(ed25519.PublicKey)

	// Private key should not be empty
	require.NotEqual(t, sk, []byte{})

	// Public key should not be empty
	require.NotEqual(t, pk, []byte{})

	addr, err := GenerateAddressFromSK(sk)
	require.NoError(t, err)
	require.Len(t, addr, 58)

	addrFromPk, err := GenerateAddressFromPublicKey(pk)
	require.NoError(t, err)
	require.Equal(t, addr, addrFromPk)

	// Address should be identical to public key
	decoded, err := types.DecodeAddress(addr)
	require.NoError(t, err)
	require.Equal(t, pk, ed25519.PublicKey(decoded[:]))
}
