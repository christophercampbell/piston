package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

type EngineClient struct {
	client *http.Client
	url    string
	reqID  uint64
}

func NewEngineClient(url string, jwtFile string) (*EngineClient, error) {
	secret, err := parseJWTSecretFromFile(jwtFile)
	if err != nil {
		return nil, err
	}
	authTransport := &jwtTransport{
		underlyingTransport: http.DefaultTransport,
		jwtSecret:           secret,
	}
	client := &http.Client{
		Timeout:   DefaultRPCTimeout,
		Transport: authTransport,
	}
	return &EngineClient{
		client: client,
		url:    url,
	}, nil
}

func (ec *EngineClient) Close() {
	ec.client.CloseIdleConnections()
}

func (ec *EngineClient) ForkchoiceUpdated(state *ForkChoiceState, attrs *PayloadAttributes) (*ForkchoiceUpdatedResponse, error) {
	msg, err := ec.call("engine_forkchoiceUpdatedV3", state, attrs)
	if err != nil {
		return nil, err
	}
	data, err := msg.MarshalJSON()
	if err != nil {
		return nil, err
	}
	var response ForkchoiceUpdatedResponse
	err = json.Unmarshal(data, &response)
	if err != nil {
		if e, ok := err.(*json.SyntaxError); ok {
			err = errors.Wrap(e, "invalid json response from server")
		}
		return nil, err
	}
	return &response, nil
}

func (ec *EngineClient) GetPayload(payloadId string) (*Payload, error) {
	msg, err := ec.call("engine_getPayloadV3", payloadId)
	if err != nil {
		return nil, err
	}
	var response Payload
	err = json.Unmarshal(msg, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (ec *EngineClient) NewPayload(payload ExecutionPayload, blobs []string, beacon common.Hash) (*NewPayloadResponse, error) {
	msg, err := ec.call("engine_newPayloadV3", payload, blobs, beacon)
	if err != nil {
		return nil, err
	}
	var response NewPayloadResponse
	err = json.Unmarshal(msg, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (ec *EngineClient) CheckCapabilities(requiredMethods []string) error {
	data, err := ec.call("engine_exchangeCapabilities", requiredMethods)
	if err != nil {
		return err
	}
	var response []string
	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	for _, method := range requiredMethods {
		if !contains(response, method) {
			return errors.New(fmt.Sprintf("engine API does not support method '%v'", method))
		}
	}
	return nil
}

func contains(arr []string, val string) bool {
	for _, s := range arr {
		if s == val {
			return true
		}
	}
	return false
}

// Call returns raw response of method call
func (ec *EngineClient) call(method string, params ...interface{}) (json.RawMessage, error) {
	var args []interface{}
	for _, p := range params {
		if p != nil {
			args = append(args, p)
		}
	}
	request := JsonrpcRequest{
		ID:      ec.reqID,
		JSONRPC: "2.0",
		Method:  method,
		Params:  args,
	}
	payload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	ec.reqID++

	response, err := ec.client.Post(ec.url, "application/json", bytes.NewBuffer(payload))
	if response != nil {
		defer func(Body io.ReadCloser) {
			_ = Body.Close()
		}(response.Body)
	}
	if err != nil {
		return nil, err
	}

	resp := new(JsonrpcResponse)
	if err = json.NewDecoder(response.Body).Decode(resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil
}
