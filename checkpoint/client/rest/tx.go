package rest

import (
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	restClient "github.com/maticnetwork/heimdall/client/rest"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/checkpoint/new",
		newCheckpointHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc("/checkpoint/ack", newCheckpointACKHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/checkpoint/no-ack", newCheckpointNoACKHandler(cliCtx)).Methods("POST")
}

type (
	// HeaderBlockReq struct for incoming checkpoint
	HeaderBlockReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer        hmTypes.HeimdallAddress `json:"proposer"`
		RootHash        hmTypes.HeimdallHash    `json:"rootHash"`
		AccountRootHash hmTypes.HeimdallHash    `json:"accountRootHash"`
		StartBlock      uint64                  `json:"startBlock"`
		EndBlock        uint64                  `json:"endBlock"`
	}

	// HeaderACKReq struct for sending ACK for a new headers
	// by providing the header index assigned my mainchain contract
	HeaderACKReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer    hmTypes.HeimdallAddress `json:"proposer"`
		HeaderBlock uint64                  `json:"headerBlock"`
		TxHash      hmTypes.HeimdallHash    `json:"tx_hash"`
		LogIndex    uint64                  `json:"log_index"`
	}

	// HeaderNoACKReq struct for sending no-ack for a new headers
	HeaderNoACKReq struct {
		BaseReq rest.BaseReq `json:"base_req"`

		Proposer hmTypes.HeimdallAddress `json:"proposer"`
	}
)

func newCheckpointHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req HeaderBlockReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// draft a message and send response
		msg := types.NewMsgCheckpointBlock(
			req.Proposer,
			req.StartBlock,
			req.EndBlock,
			req.RootHash,
			req.AccountRootHash,
			uint64(time.Now().UTC().Unix()),
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newCheckpointACKHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req HeaderACKReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// draft a message and send response
		msg := types.NewMsgCheckpointAck(req.Proposer, req.HeaderBlock, req.TxHash, req.LogIndex)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

func newCheckpointNoACKHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req HeaderNoACKReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// draft a message and send response
		msg := types.NewMsgCheckpointNoAck(
			req.Proposer,
			uint64(time.Now().UTC().Unix()),
		)

		// send response
		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
