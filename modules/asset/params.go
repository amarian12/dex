package asset

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
)

// DefaultParamspace defines the default asset module parameter subspace
const (
	DefaultParamspace = ModuleName
	MaxTokenAmount    = 9E18 // 90 billion * 10 ^ 8

	IssueTokenFee                 = 1E12 // 10000 * 10 ^8
	FreezeAddrFee                 = 1E9  // 10 * 10 ^ 8
	UnFreezeAddrFee               = 1E9  // 10 * 10 ^ 8
	FreezeTokenFee                = 1E11 //1000 * 10 ^ 8
	UnFreezeTokenFee              = 1E11 //1000 * 10 ^ 8
	TokenFreezeWhitelistAddFee    = 1E10 // 100 * 10 ^ 8
	TokenFreezeWhitelistRemoveFee = 1E10 // 100 * 10 ^ 8
	BurnFee                       = 1E11 //1000 * 10 ^ 8
	MintFee                       = 1E11 //1000 * 10 ^ 8
)

// Parameter keys
var (
	KeyIssueTokenFee                 = []byte("IssueTokenFee")
	KeyFreezeAddrFee                 = []byte("FreezeAddrFee")
	KeyUnFreezeAddrFee               = []byte("UnFreezeAddrFee")
	KeyFreezeTokenFee                = []byte("FreezeTokenFee")
	KeyUnFreezeTokenFee              = []byte("UnFreezeTokenFee")
	KeyTokenFreezeWhitelistAddFee    = []byte("TokenFreezeWhitelistAddFee")
	KeyTokenFreezeWhitelistRemoveFee = []byte("TokenFreezeWhitelistRemoveFee")
	KeyBurnFee                       = []byte("BurnFee")
	KeyMintFee                       = []byte("MintFee")
)

var _ params.ParamSet = &Params{}

// Params defines the parameters for the asset module.
type Params struct {
	// FeeParams define the rules according to which fee are charged.
	IssueTokenFee                 sdk.Coins `json:"issue_token_fee"`
	FreezeAddrFee                 sdk.Coins `json:"freeze_address_fee"`
	UnFreezeAddrFee               sdk.Coins `json:"unfreeze_address_fee"`
	FreezeTokenFee                sdk.Coins `json:"freeze_token_fee"`
	UnFreezeTokenFee              sdk.Coins `json:"unfreeze_token_fee"`
	TokenFreezeWhitelistAddFee    sdk.Coins `json:"token_freeze_whitelist_add_fee"`
	TokenFreezeWhitelistRemoveFee sdk.Coins `json:"token_freeze_whitelist_remove_fee"`
	BurnFee                       sdk.Coins `json:"burn_fee"`
	MintFee                       sdk.Coins `json:"mint_fee"`
}

// ParamKeyTable for asset module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of asset module's parameters.
// nolint
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{

		{KeyIssueTokenFee, &p.IssueTokenFee},
		{KeyFreezeAddrFee, &p.FreezeAddrFee},
		{KeyUnFreezeAddrFee, &p.UnFreezeAddrFee},
		{KeyFreezeTokenFee, &p.FreezeTokenFee},
		{KeyUnFreezeTokenFee, &p.UnFreezeTokenFee},
		{KeyTokenFreezeWhitelistAddFee, &p.TokenFreezeWhitelistAddFee},
		{KeyTokenFreezeWhitelistRemoveFee, &p.TokenFreezeWhitelistRemoveFee},
		{KeyBurnFee, &p.BurnFee},
		{KeyMintFee, &p.MintFee},
	}
}

// Equal returns a boolean determining if two Params types are identical.
func (p Params) Equal(p2 Params) bool {
	bz1 := msgCdc.MustMarshalBinaryLengthPrefixed(&p)
	bz2 := msgCdc.MustMarshalBinaryLengthPrefixed(&p2)
	return bytes.Equal(bz1, bz2)
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {

	return Params{
		IssueTokenFee:                 CetCoin(IssueTokenFee),
		FreezeAddrFee:                 CetCoin(FreezeAddrFee),
		UnFreezeAddrFee:               CetCoin(UnFreezeAddrFee),
		FreezeTokenFee:                CetCoin(FreezeTokenFee),
		UnFreezeTokenFee:              CetCoin(UnFreezeTokenFee),
		TokenFreezeWhitelistAddFee:    CetCoin(TokenFreezeWhitelistAddFee),
		TokenFreezeWhitelistRemoveFee: CetCoin(TokenFreezeWhitelistRemoveFee),
		BurnFee:                       CetCoin(BurnFee),
		MintFee:                       CetCoin(MintFee),
	}
}

func CetCoin(amt int64) sdk.Coins {
	return sdk.Coins{
		sdk.NewCoin("cet", sdk.NewInt(amt)),
	}
}