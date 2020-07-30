package crypto

import (
	"fmt"
)

var errInvalidSignatureReturned = fmt.Errorf("ed25519 library returned an invalid signature")
var errFailedToCopyPK = fmt.Errorf("failed to copy the public key")
