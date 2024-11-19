package app

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/signal"
	"time"

	"github.com/0xPolygon/maera/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
)

func Run(cli *cli.Context) error {
	jwtPath := processPath(cli.String(JWTKey))
	engineUrl := orDefaultString(cli, EngineUrlKey, EngineUrlDefault)
	ethUrl := orDefaultString(cli, EthUrlKey, EthUrlDefault)
	fmt.Printf("ethUrl: %s\n", ethUrl)
	fmt.Printf("engineUrl: %s\n", engineUrl)
	fmt.Printf("jwtPath: %s\n", jwtPath)

	runner := &Runner{
		ctx:       cli.Context,
		jwtPath:   jwtPath,
		ethUrl:    ethUrl,
		engineUrl: engineUrl,
		stop:      make(chan struct{}),
	}

	go runner.createBlocks()

	waitSignal()

	close(runner.stop)

	return runner.lastError
}

func waitSignal() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	for sig := range signals {
		switch sig {
		case os.Interrupt, os.Kill:
			log.Info("terminating..")
			return
		}
	}
}

type Runner struct {
	ctx                        context.Context
	jwtPath, ethUrl, engineUrl string
	stop                       chan struct{}
	lastError                  error
}

func (r *Runner) createBlocks() {
	r.createNextBlockSafe()
	for {
		select {
		case <-time.NewTimer(time.Second * 10).C:
			r.createNextBlockSafe()
		case <-r.stop:
			return
		}
	}
}

func (r *Runner) createNextBlockSafe() {
	if err := r.createNextBlock(); err != nil {
		fmt.Println("error creating block", err)
		r.lastError = err
	}
}

func (r *Runner) createNextBlock() error {
	_, latest, err := getLatestBlockInfo(r.ctx, r.ethUrl)
	if err != nil {
		return err
	}
	ec, err := engine.NewEngineClient(r.engineUrl, r.jwtPath)
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
