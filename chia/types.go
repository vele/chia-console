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
			ChallengeBlockInfoHash string `json:"challenge_block_info_hash"`
			ChallengeVdfOutput     struct {
				Data string `json:"data"`
			} `json:"challenge_vdf_output"`
			Deficit                            int         `json:"deficit"`
			FarmerPuzzleHash                   string      `json:"farmer_puzzle_hash"`
			Fees                               int         `json:"fees"`
			FinishedChallengeSlotHashes        interface{} `json:"finished_challenge_slot_hashes"`
			FinishedInfusedChallengeSlotHashes interface{} `json:"finished_infused_challenge_slot_hashes"`
			FinishedRewardSlotHashes           interface{} `json:"finished_reward_slot_hashes"`
			HeaderHash                         string      `json:"header_hash"`
			Height                             int         `json:"height"`
			InfusedChallengeVdfOutput          struct {
				Data string `json:"data"`
			} `json:"infused_challenge_vdf_output"`
			Overflow                   bool   `json:"overflow"`
			PoolPuzzleHash             string `json:"pool_puzzle_hash"`
			PrevHash                   string `json:"prev_hash"`
			PrevTransactionBlockHash   string `json:"prev_transaction_block_hash"`
			PrevTransactionBlockHeight int    `json:"prev_transaction_block_height"`
			RequiredIters              int    `json:"required_iters"`
			RewardClaimsIncorporated   []struct {
				Amount         int64  `json:"amount"`
				ParentCoinInfo string `json:"parent_coin_info"`
				PuzzleHash     string `json:"puzzle_hash"`
			} `json:"reward_claims_incorporated"`
			RewardInfusionNewChallenge string      `json:"reward_infusion_new_challenge"`
			SignagePointIndex          int         `json:"signage_point_index"`
			SubEpochSummaryIncluded    interface{} `json:"sub_epoch_summary_included"`
			SubSlotIters               int         `json:"sub_slot_iters"`
			Timestamp                  int         `json:"timestamp"`
			TotalIters                 float64     `json:"total_iters"`
			Weight                     int         `json:"weight"`
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
