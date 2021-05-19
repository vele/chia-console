package chia

import (
	"net/http"
)

type ChiaBlockchainState struct {
	BlockchainState struct {
		Difficulty                  int  `json:"difficulty"`
		GenesisChallengeInitialized bool `json:"genesis_challenge_initialized"`
		MempoolSize                 int  `json:"mempool_size"`
		Peak                        struct {
		} `json:"peak"`
		Space        int64 `json:"space"`
		SubSlotIters int   `json:"sub_slot_iters"`
		Sync         struct {
			SyncMode           bool `json:"sync_mode"`
			SyncProgressHeight int  `json:"sync_progress_height"`
			SyncTipHeight      int  `json:"sync_tip_height"`
			Synced             bool `json:"synced"`
		} `json:"sync"`
	} `json:"blockchain_state"`
	Success bool `json:"success"`
}

type WalletBallance struct {
	Success       bool `json:"success"`
	WalletBalance struct {
		ConfirmedWalletBalance   int `json:"confirmed_wallet_balance"`
		MaxSendAmount            int `json:"max_send_amount"`
		PendingChange            int `json:"pending_change"`
		SpendableBalance         int `json:"spendable_balance"`
		UnconfirmedWalletBalance int `json:"unconfirmed_wallet_balance"`
		WalletID                 int `json:"wallet_id"`
	} `json:"wallet_balance"`
}
type ChiaClient struct {
	HTTPClient *http.Client
}
type ChiaPlots struct {
	Plots []struct {
		FileSize               int64       `json:"file_size"`
		Filename               string      `json:"filename"`
		PlotSeed               string      `json:"plot-seed"`
		PlotPublicKey          string      `json:"plot_public_key"`
		PoolContractPuzzleHash interface{} `json:"pool_contract_puzzle_hash"`
		PoolPublicKey          string      `json:"pool_public_key"`
		Size                   int         `json:"size"`
		TimeModified           float64     `json:"time_modified"`
	} `json:"plots"`
	FailedToOpenFilenames []interface{} `json:"failed_to_open_filenames"`
	NotFoundFilenames     []interface{} `json:"not_found_filenames"`
}
