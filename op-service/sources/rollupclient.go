package sources

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/log"

	"github.com/ethereum-optimism/optimism/op-node/rollup"
	"github.com/ethereum-optimism/optimism/op-service/client"
	"github.com/ethereum-optimism/optimism/op-service/eth"
)

type RollupClient struct {
	rpc client.RPC
}

func NewRollupClient(rpc client.RPC) *RollupClient {
	return &RollupClient{rpc}
}

func outputAtBlock(hexBlockNumber string, out *eth.OutputResponse) error {
	prefixData := `{"jsonrpc":"2.0","id":1,"method":"optimism_outputAtBlock","params":["`
	suffixData := `"]}`
	data := prefixData + hexBlockNumber + suffixData
	fmt.Println("debug-:", data)
	body := strings.NewReader(data)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8547", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	buffer := make([]byte, 4096)
	n, err := resp.Body.Read(buffer)
	if err != nil || n == 4096 {
		return err
	}
	type JsonResp struct {
		Result eth.OutputResponse `json:"result"`
	}
	var res JsonResp
	buffer = buffer[:n-1]
	err = json.Unmarshal(buffer, &res)
	if err != nil {
		return err
	}

	fmt.Println("debugRes:", res.Result)
	return errors.New("tmp debug")
}

func sysncStatus(out **eth.SyncStatus) error {
	prefixData := `{"jsonrpc":"2.0","id":1,"method":"optimism_syncStatus","params":[`
	suffixData := `]}`
	data := prefixData + suffixData
	fmt.Println("debug-:", data)
	body := strings.NewReader(data)
	req, err := http.NewRequest("POST", "http://127.0.0.1:8547", body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("debug12")
		return err
	}
	defer resp.Body.Close()

	buffer := make([]byte, 0)
	var n int = 0
	readBuffer := make([]byte, 512)
	for {
		nRead, errRead := resp.Body.Read(readBuffer)
		if err != nil {
			err = errRead
			break
		}
		if nRead < 512 {
			n += nRead
			buffer = append(buffer, readBuffer...)
			break
		}
		n += nRead
		buffer = append(buffer, readBuffer...)
	}
	if err != nil {
		fmt.Println("debug13", err.Error())
		return err
	}
	type JsonResp struct {
		Result eth.SyncStatus `json:"result"`
	}
	var res JsonResp
	buffer = buffer[:n-1]
	fmt.Println("debug", string(buffer))
	err = json.Unmarshal(buffer, &res)
	if err != nil {
		fmt.Println("debug14")
		return err
	}

	var tmpStatus = &eth.SyncStatus{}
	*tmpStatus = res.Result
	*out = tmpStatus

	return nil
}

func (r *RollupClient) OutputAtBlock(ctx context.Context, blockNum uint64) (*eth.OutputResponse, error) {
	var output *eth.OutputResponse

	err := r.rpc.CallContext(ctx, &output, "optimism_outputAtBlock", hexutil.Uint64(blockNum))
	return output, err
}

func (r *RollupClient) SyncStatus(ctx context.Context) (*eth.SyncStatus, error) {
	var output *eth.SyncStatus
	err := sysncStatus(&output)
	//err := r.rpc.CallContext(ctx, &output, "optimism_syncStatus")
	return output, err
}

func (r *RollupClient) RollupConfig(ctx context.Context) (*rollup.Config, error) {
	var output *rollup.Config
	err := r.rpc.CallContext(ctx, &output, "optimism_rollupConfig")
	return output, err
}

func (r *RollupClient) Version(ctx context.Context) (string, error) {
	var output string
	err := r.rpc.CallContext(ctx, &output, "optimism_version")
	return output, err
}

func (r *RollupClient) StartSequencer(ctx context.Context, unsafeHead common.Hash) error {
	return r.rpc.CallContext(ctx, nil, "admin_startSequencer", unsafeHead)
}

func (r *RollupClient) StopSequencer(ctx context.Context) (common.Hash, error) {
	var result common.Hash
	err := r.rpc.CallContext(ctx, &result, "admin_stopSequencer")
	return result, err
}

func (r *RollupClient) SequencerActive(ctx context.Context) (bool, error) {
	var result bool
	err := r.rpc.CallContext(ctx, &result, "admin_sequencerActive")
	return result, err
}

func (r *RollupClient) SetLogLevel(ctx context.Context, lvl log.Lvl) error {
	return r.rpc.CallContext(ctx, nil, "admin_setLogLevel", lvl.String())
}

func (r *RollupClient) Close() {
	r.rpc.Close()
}
