package mobile

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
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

func TestAssignGroupID(t *testing.T) {
	type assignGroupIDTest struct {
		b64Txns            []string
		b64ExpectedGroupID string
	}

	tests := []assignGroupIDTest{
		{
			b64Txns: []string{
				"iqNhbXTOAA9CQKNmZWXNA+iiZnbOAOHF36NnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4A4cnHpG5vdGXEEVRlc3RpbmcgZ3JvdXAgSURzo3JjdsQgKwg17XWyS7m6iUEK87rTYF6NxV6isLU7A/xwYwuCcaOjc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlo3BheQ==",
			},
			b64ExpectedGroupID: "w2waFq6tc/5VA0ysOCk3NWBCx3ZUPkhc2T1PpMkne6g=",
		},
		{
			b64Txns: []string{
				"iaRhcGFyiaJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6sVGVzdCBBc3NldCAyomF1s2h0dHBzOi8vZXhhbXBsZS5jb22hY8QgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhZsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhbcQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhcsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhdM///////////6J1bqRUU1Qyo2ZlZc0D6KJmds4A4ciwo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToiomx2zgDhzJikbm90ZcQOVGhpcyBpcyBhIG5vdGWjc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlpGFjZmc=",
				"iKRjYWlkAqNmZWXNA+iiZnbOAOHIsKNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4A4cyYo3NuZMQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+ykdHlwZaRhY2Zn",
			},
			b64ExpectedGroupID: "TBEqLZ3z3LsE3jyt5t5Z3b/R1/XMl9Gy8Epjsoj6Pdk=",
		},
	}

	for testIndex, test := range tests {
		t.Run(fmt.Sprintf("index=%d", testIndex), func(t *testing.T) {
			encodedTxns := make([][]byte, len(test.b64Txns))
			for i := range test.b64Txns {
				txn, err := base64.StdEncoding.DecodeString(test.b64Txns[i])
				if err != nil {
					t.Fatal(err)
				}
				encodedTxns[i] = txn
			}

			expectedGroupID, err := base64.StdEncoding.DecodeString(test.b64ExpectedGroupID)
			if err != nil {
				t.Fatal(err)
			}

			txns := BytesArray{
				values: encodedTxns,
			}

			assignedTxns, err := AssignGroupID(&txns)
			if err != nil {
				t.Fatal(err)
			}

			if assignedTxns.Length() != len(encodedTxns) {
				t.Fatalf("Length of returned transactions does not match. Got %d, expected %d", assignedTxns.Length(), len(encodedTxns))
			}

			for i, atxn := range assignedTxns.Extract() {
				var assignedTxn types.Transaction
				err = msgpack.Decode(atxn, &assignedTxn)
				if err != nil {
					t.Fatal(err)
				}

				if !bytes.Equal(assignedTxn.Group[:], expectedGroupID) {
					t.Errorf("Actual group ID does not match expected for transaction at index %d. Got %s, expected %s", i, base64.StdEncoding.EncodeToString(assignedTxn.Group[:]), base64.StdEncoding.EncodeToString(expectedGroupID))
				}

				assignedTxn.Group = types.Digest{}
				encodedActualTxn := msgpack.Encode(&assignedTxn)

				if !bytes.Equal(encodedActualTxn, encodedTxns[i]) {
					t.Errorf("Returned transaction at index %d is unexpectedly modified", i)
				}
			}
		})
	}
}

func TestVerifyGroupID(t *testing.T) {
	type verifyGroupIDTest struct {
		name    string
		b64Txns []string
		valid   bool
	}

	tests := []verifyGroupIDTest{
		{
			name: "Single txn, no group",
			b64Txns: []string{
				"iqNhbXTOAA9CQKNmZWXNA+iiZnbOAOHF36NnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4A4cnHpG5vdGXEEVRlc3RpbmcgZ3JvdXAgSURzo3JjdsQgKwg17XWyS7m6iUEK87rTYF6NxV6isLU7A/xwYwuCcaOjc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlo3BheQ==",
			},
			valid: true,
		},
		{
			name: "Single txn, correct group",
			b64Txns: []string{
				"i6NhbXTOAA9CQKNmZWXNA+iiZnbOAOHF36NnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIMNsGhaurXP+VQNMrDgpNzVgQsd2VD5IXNk9T6TJJ3uoomx2zgDhycekbm90ZcQRVGVzdGluZyBncm91cCBJRHOjcmN2xCArCDXtdbJLubqJQQrzutNgXo3FXqKwtTsD/HBjC4Jxo6NzbmTEILSSdr0+wJd+q4ajIcRJ6tgCyWwL2XwpVhMVEdLxHuvspHR5cGWjcGF5",
			},
			valid: true,
		},
		{
			name: "Single txn, wrong group",
			b64Txns: []string{
				"i6NhbXTOAA9CQKNmZWXNA+iiZnbOAOHF36NnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIARlvwikIb0YIGkkP3wiLhq+D+sLipBbd2KlH4/CEHgXomx2zgDhycekbm90ZcQRVGVzdGluZyBncm91cCBJRHOjcmN2xCArCDXtdbJLubqJQQrzutNgXo3FXqKwtTsD/HBjC4Jxo6NzbmTEILSSdr0+wJd+q4ajIcRJ6tgCyWwL2XwpVhMVEdLxHuvspHR5cGWjcGF5",
			},
			valid: false,
		},
		{
			name: "Multi txn, correct group",
			b64Txns: []string{
				"iqRhcGFyiaJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6sVGVzdCBBc3NldCAyomF1s2h0dHBzOi8vZXhhbXBsZS5jb22hY8QgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhZsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhbcQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhcsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhdM///////////6J1bqRUU1Qyo2ZlZc0D6KJmds4A4ciwo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgTBEqLZ3z3LsE3jyt5t5Z3b/R1/XMl9Gy8Epjsoj6PdmibHbOAOHMmKRub3RlxA5UaGlzIGlzIGEgbm90ZaNzbmTEILSSdr0+wJd+q4ajIcRJ6tgCyWwL2XwpVhMVEdLxHuvspHR5cGWkYWNmZw==",
				"iaRjYWlkAqNmZWXNA+iiZnbOAOHIsKNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIEwRKi2d89y7BN48rebeWd2/0df1zJfRsvBKY7KI+j3Zomx2zgDhzJijc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlpGFjZmc=",
			},
			valid: true,
		},
		{
			name: "Multi txn, 1 wrong group",
			b64Txns: []string{
				"iqRhcGFyiaJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6sVGVzdCBBc3NldCAyomF1s2h0dHBzOi8vZXhhbXBsZS5jb22hY8QgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhZsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhbcQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhcsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhdM///////////6J1bqRUU1Qyo2ZlZc0D6KJmds4A4ciwo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgTBEqLZ3z3LsE3jyt5t5Z3b/R1/XMl9Gy8Epjsoj6PdmibHbOAOHMmKRub3RlxA5UaGlzIGlzIGEgbm90ZaNzbmTEILSSdr0+wJd+q4ajIcRJ6tgCyWwL2XwpVhMVEdLxHuvspHR5cGWkYWNmZw==",
				"iaRjYWlkAqNmZWXNA+iiZnbOAOHIsKNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIARlvwikIb0YIGkkP3wiLhq+D+sLipBbd2KlH4/CEHgXomx2zgDhzJijc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlpGFjZmc=",
			},
			valid: false,
		},
		{
			name: "Multi txn, all wrong group",
			b64Txns: []string{
				"iqRhcGFyiaJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6sVGVzdCBBc3NldCAyomF1s2h0dHBzOi8vZXhhbXBsZS5jb22hY8QgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhZsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhbcQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhcsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhdM///////////6J1bqRUU1Qyo2ZlZc0D6KJmds4A4ciwo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToio2dycMQgBGW/CKQhvRggaSQ/fCIuGr4P6wuKkFt3YqUfj8IQeBeibHbOAOHMmKRub3RlxA5UaGlzIGlzIGEgbm90ZaNzbmTEILSSdr0+wJd+q4ajIcRJ6tgCyWwL2XwpVhMVEdLxHuvspHR5cGWkYWNmZw==",
				"iaRjYWlkAqNmZWXNA+iiZnbOAOHIsKNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqNncnDEIARlvwikIb0YIGkkP3wiLhq+D+sLipBbd2KlH4/CEHgXomx2zgDhzJijc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlpGFjZmc=",
			},
			valid: false,
		},
		{
			name: "Multi txn, no group",
			b64Txns: []string{
				"iaRhcGFyiaJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6sVGVzdCBBc3NldCAyomF1s2h0dHBzOi8vZXhhbXBsZS5jb22hY8QgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhZsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhbcQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhcsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhdM///////////6J1bqRUU1Qyo2ZlZc0D6KJmds4A4ciwo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToiomx2zgDhzJikbm90ZcQOVGhpcyBpcyBhIG5vdGWjc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlpGFjZmc=",
				"iKRjYWlkAqNmZWXNA+iiZnbOAOHIsKNnZW6sdGVzdG5ldC12MS4womdoxCBIY7UYpLPITsgQ8i1PEIHLD3HwWaesIN7GL39w5Qk6IqJsds4A4cyYo3NuZMQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+ykdHlwZaRhY2Zn",
			},
			valid: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			encodedTxns := make([][]byte, len(test.b64Txns))
			for i := range test.b64Txns {
				txn, err := base64.StdEncoding.DecodeString(test.b64Txns[i])
				if err != nil {
					t.Fatal(err)
				}
				encodedTxns[i] = txn
			}
			txns := BytesArray{
				values: encodedTxns,
			}

			result, err := VerifyGroupID(&txns)
			if err != nil {
				t.Fatal(err)
			}

			if result != test.valid {
				t.Errorf("Unexpected result: got %v, expected %v", result, test.valid)
			}
		})
	}
}
