package wallet

import (
	"fmt"
	"log"
	"tajfi-server/wallet/lnd"
	"tajfi-server/wallet/tapd"
)

// Receive initializes the generate invoice process.
func Receive(params ReceiveParams) (map[string]interface{}, error) {
	// Step 1: Call LND to get the internal key
	log.Println("Getting internal key with params", params)
	internalKey, err := lnd.GetInternalKey(params.LNDHost, params.LNMacaroon)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal key: %w", err)
	}

	log.Println("Got internal key", internalKey)

	// compress pubkey
	rawKeyBytes, err := CompressPubKey(params.PubKey)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("failed to compress pubkey: %w", err)
	}

	// Step 2: Prepare the Tapd payload
	payload := tapd.NewAddressPayload{
		AssetID: params.AssetID,
		Amt:     params.Amount,
		ScriptKey: map[string]interface{}{
			"pub_key": params.PubKey,
			"key_desc": map[string]interface{}{
				"raw_key_bytes": rawKeyBytes,
			},
		},
		InternalKey: *internalKey,
	}

	log.Println("Prepared Tapd payload", payload)

	// Step 3: Call the Tapd API to create the address
	response, err := tapd.CallNewAddress(params.TapdHost, params.TapdMacaroon, payload)
	if err != nil {
		log.Println("Failed to call Tapd API", err)
		return nil, fmt.Errorf("failed to call Tapd API: %w", err)
	}

	return response, nil
}
