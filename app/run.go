package app

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygon/maera/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

func Run(cli *cli.Context) error {
	jwtPath := processPath(cli.String(JWTKey))
	engineUrl := orDefaultString(cli, EngineUrlKey, EngineUrlDefault)
	ethUrl := orDefaultString(cli, EthUrlKey, EthUrlDefault)
	fmt.Printf("ethUrl: %s\n", ethUrl)
	fmt.Printf("engineUrl: %s\n", engineUrl)
	fmt.Printf("jwtPath: %s\n", jwtPath)

	return createNextBlock(cli.Context, jwtPath, ethUrl, engineUrl)
}

func createNextBlock(ctx context.Context, jwtPath, ethUrl, engineUrl string) error {

	_, latest, err := getLatestBlockInfo(ctx, ethUrl)
	if err != nil {
		return err
	}
	ec, err := engine.NewEngineClient(engineUrl, jwtPath)
	if err != nil {
		return err
	}
	defer ec.Close()

	state := engine.ForkChoiceState{
		HeadHash:           *latest,
		SafeBlockHash:      *latest,
		FinalizedBlockHash: *latest,
	}

	attrs := engine.PayloadAttributes{
		Timestamp:             hexutil.Uint64(time.Now().UnixNano()),
		PrevRandao:            common.Hash{},
		SuggestedFeeRecipient: common.Address{},
		Withdrawals:           []*engine.Withdrawal{},
	}

	fcu, err := ec.ForkchoiceUpdated(&state, &attrs)
	if err != nil {
		return err
	}
	fmt.Println("initiated forkchoice, payloadId: ", fcu.PayloadId)

	payload, err := ec.GetPayload(fcu.PayloadId)
	if err != nil {
		return err
	}
	fmt.Println("received payload")

	blobs := []string{}         // figure out what to do here
	beaconHash := common.Hash{} // figure out what to do here

	newPayload, err := ec.NewPayload(payload.ExecutionPayload, blobs, beaconHash)
	if err != nil {
		return err
	}
	fmt.Println("created new payload", "status:", newPayload.Status, "hash", newPayload.LatestValidHash)

	// forkchoice with new hashes
	state.FinalizedBlockHash = state.HeadHash
	state.HeadHash = common.HexToHash(newPayload.LatestValidHash)
	state.SafeBlockHash = common.HexToHash(newPayload.LatestValidHash)

	// reset the timestamp on params
	attrs.Timestamp = hexutil.Uint64(time.Now().UnixNano())

	finalFC, err := ec.ForkchoiceUpdated(&state, &attrs)
	if err != nil {
		return err
	}
	fmt.Println("finalized fork choice", "status:", finalFC.PayloadStatus, "hash", newPayload.LatestValidHash)
	return nil
}

func getLatestBlockInfo(ctx context.Context, ethUrl string) (*big.Int, *common.Hash, error) {
	eth, err := ethclient.Dial(ethUrl)
	if err != nil {
		return nil, nil, err
	}
	defer eth.Close()

	n, err := eth.BlockNumber(ctx)
	if err != nil {
		return nil, nil, err
	}
	number := big.NewInt(int64(n))

	block, err := eth.BlockByNumber(ctx, number)
	if err != nil {
		return nil, nil, err
	}
	hash := block.Hash()
	return number, &hash, err
}
