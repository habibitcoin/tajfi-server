package tapd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"tajfi-server/config"
	"tajfi-server/interfaces"
)

// TapdClientInterface defines the methods for interacting with the Tapd API.
type TapdClientInterface interface {
	GetClientConfig() *config.Config
	CallNewAddress(tapdHost, macaroon string, payload NewAddressPayload) (map[string]interface{}, error)
	DecodeAddr(tapdHost, macaroon, address string) (*DecodeAddrResponse, error)
	FundVirtualPSBT(tapdHost, macaroon, invoice string, inputs PrevIds, useTemplate bool) (*FundVirtualPSBTResponse, error)
	SignVirtualPSBT(tapdHost, macaroon, psbt string) (*SignVirtualPSBTResponse, error)
	AnchorVirtualPSBT(params AnchorVirtualPSBTParams) (*AssetTransferResponse, error)
	LogAndTransferPsbt(tapdHost, macaroon string, req CommitVirtualPsbtsRequest) (*AssetTransferResponse, error)
	GetBalances(tapdHost, macaroon string) (*WalletBalancesResponse, error)
	GetTransfers(tapdHost, macaroon string) (AssetTransfersResponse, error)
	GetUtxos(tapdHost, macaroon string) (*GetUtxosResponse, error)
	// Only used for demo mode.
	SendAssets(tapdHost, macaroon, invoice string) (*FundVirtualPSBTResponse, error)

	// Used in performing trustless swaps
	CreateInteractiveSendTemplate(tapdHost, macaroon string, request CreateInteractiveSendTemplateRequest) (*CreateInteractiveSendTemplateResponse, error)
	PrepareAnchoringTemplate(tapdHost, macaroon string, request PrepareAnchoringTemplateRequest) (*PrepareAnchoringTemplateResponse, error)
	UpdateVirtualPSBT(tapdHost, macaroon string, request UpdateVirtualPsbtRequest) (*UpdateVirtualPsbtResponse, error)
	CommitVirtualPsbts(tapdHost, macaroon string, req CommitVirtualPsbtsRequest) (*CommitVirtualPsbtsResponse, error)
}

type tapdClient struct {
	httpClient interfaces.HttpClient
	cfg        *config.Config
}

// NewTapdClient initializes a new TapdClient.
func NewTapdClient(client interfaces.HttpClient, cfg *config.Config) TapdClientInterface {
	return &tapdClient{httpClient: client, cfg: cfg}
}

// Ensure tapdClient implements TapdClientInterface.
var _ TapdClientInterface = (*tapdClient)(nil)

func (c *tapdClient) GetClientConfig() *config.Config {
	return c.cfg
}

// makeGetRequest is a helper function for GET requests.
func (c *tapdClient) makeGetRequest(url, macaroon string, response interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create GET request: %w", err)
	}

	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("GET request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return json.NewDecoder(resp.Body).Decode(response)
}

// makePostRequest simplifies making POST requests.
func (c *tapdClient) makePostRequest(url, macaroon string, payload interface{}, result interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Grpc-Metadata-macaroon", macaroon)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("tapd RPC error: %s", resp.Status)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	return nil
}
