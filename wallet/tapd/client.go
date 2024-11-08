package tapd

import (
	"tajfi-server/interfaces"
)

// TapdClientInterface defines the methods for interacting with the Tapd API.
type TapdClientInterface interface {
	CallNewAddress(tapdHost, macaroon string, payload NewAddressPayload) (map[string]interface{}, error)
	DecodeAddr(tapdHost, macaroon, address string) (*DecodeAddrResponse, error)
	FundVirtualPSBT(tapdHost, macaroon, invoice string, inputs PrevIds) (fundedPsbt *FundVirtualPSBTResponse, err error)
	SignVirtualPSBT(tapdHost, macaroon, psbt string) (fundedPsbt *SignVirtualPSBTResponse, err error)
	AnchorVirtualPSBT(params AnchorVirtualPSBTParams) (*AssetTransferResponse, error)
	// Only used for demo mode
	SendAssets(tapdHost, macaroon, invoice string) (fundedPsbt *FundVirtualPSBTResponse, err error)
}

type tapdClient struct {
	httpClient interfaces.HttpClient
}

func NewTapdClient(client interfaces.HttpClient) TapdClientInterface {
	return &tapdClient{
		httpClient: client,
	}
}

// Ensure tapdClient implements TapdClientInterface.
var _ TapdClientInterface = (*tapdClient)(nil)
