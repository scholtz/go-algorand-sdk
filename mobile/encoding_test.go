package mobile

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

// checks json strings for equality
// inspired by https://gist.github.com/turtlemonvh/e4f7404e28387fadb8ad275a99596f67
func jsonEqual(s1, s2 string) (bool, error) {
	var o1 interface{}
	var o2 interface{}

	var err error
	err = json.Unmarshal([]byte(s1), &o1)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 1 :: %s", err.Error())
	}
	err = json.Unmarshal([]byte(s2), &o2)
	if err != nil {
		return false, fmt.Errorf("Error mashalling string 2 :: %s", err.Error())
	}

	return reflect.DeepEqual(o1, o2), nil
}

type encodingTest struct {
	msgpack string
	json    string
}

func TestTransaction(t *testing.T) {
	tests := []encodingTest{
		{
			msgpack: "iaRhcGFyiaJhbcQgZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5RmiiYW6sVGVzdCBBc3NldCAyomF1s2h0dHBzOi8vZXhhbXBsZS5jb22hY8QgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhZsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhbcQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhcsQgtJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+yhdM///////////6J1bqRUU1Qyo2ZlZc0D6KJmds4A3/ljo2dlbqx0ZXN0bmV0LXYxLjCiZ2jEIEhjtRiks8hOyBDyLU8QgcsPcfBZp6wg3sYvf3DlCToiomx2zgDf/Uukbm90ZcQOVGhpcyBpcyBhIG5vdGWjc25kxCC0kna9PsCXfquGoyHESerYAslsC9l8KVYTFRHS8R7r7KR0eXBlpGFjZmc=",
			json: `{"apar": {
				  "am": "ZkFDUE80blJnTzU1ajFuZEFLM1c2U2djNEFQa2N5Rmg=",
				  "an": "Test Asset 2",
				  "au": "https://example.com",
				  "c": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				  "f": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				  "m": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				  "r": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				  "t": 18446744073709551615,
				  "un": "TST2"
				},
				"fee": 1000,
				"fv": 14678371,
				"gen": "testnet-v1.0",
				"gh": "SGO1GKSzyE7IEPItTxCByw9x8FmnrCDexi9/cOUJOiI=",
				"lv": 14679371,
				"note": "VGhpcyBpcyBhIG5vdGU=",
				"snd": "tJJ2vT7Al36rhqMhxEnq2ALJbAvZfClWExUR0vEe6+w=",
				"type": "acfg"
			  }`,
		},
	}

	for _, test := range tests {
		expectedJson := test.json
		expectedMsgpack, err := base64.StdEncoding.DecodeString(test.msgpack)
		if err != nil {
			t.Fatal(err)
		}

		actualJson, err := TransactionMsgpackToJson(expectedMsgpack)
		if err != nil {
			t.Errorf("Could not convert transaction from msgpack to JSON: %v", err)
		}

		areJsonEqual, err := jsonEqual(expectedJson, actualJson)
		if err != nil {
			t.Error(err)
		} else if !areJsonEqual {
			t.Errorf("Expected JSON does not match actual JSON.\nExpected:\n%s\n\nActual:\n%s", expectedJson, actualJson)
		}

		actualMsgpack, err := TransactionJsonToMsgpack(expectedJson)
		if err != nil {
			t.Errorf("Could not convert transaction from JSON to msgpack: %v", err)
		}

		if !bytes.Equal(expectedMsgpack, actualMsgpack) {
			b64Expected := base64.StdEncoding.EncodeToString(expectedMsgpack)
			b64Actual := base64.StdEncoding.EncodeToString(actualMsgpack)
			t.Errorf("Expected msgpack does not match actual msgpack.\nExpected:\n%s\n\nActual:\n%s", b64Expected, b64Actual)
		}
	}
}
