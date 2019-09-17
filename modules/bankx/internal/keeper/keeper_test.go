package keeper_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"github.com/coinexchain/dex/modules/authx"
	"github.com/coinexchain/dex/modules/bankx"
	"github.com/coinexchain/dex/modules/bankx/internal/keeper"
	"github.com/coinexchain/dex/testapp"
	"github.com/coinexchain/dex/testutil"
	"github.com/coinexchain/dex/types"
)

var myaddr = testutil.ToAccAddress("myaddr")

func defaultContext() (keeper.Keeper, sdk.Context) {
	app := testapp.NewTestApp()
	ctx := sdk.NewContext(app.Cms, abci.Header{}, false, log.NewNopLogger())
	return app.BankxKeeper, ctx
}

func givenAccountWith(ctx sdk.Context, keeper keeper.Keeper, addr sdk.AccAddress, coinsString string) {
	coins, _ := sdk.ParseCoins(coinsString)

	acc := auth.NewBaseAccountWithAddress(addr)
	_ = acc.SetCoins(coins)
	keeper.Ak.SetAccount(ctx, &acc)

	accX := authx.AccountX{
		Address: addr,
	}
	keeper.Axk.SetAccountX(ctx, accX)
}

func coinsOf(ctx sdk.Context, keeper keeper.Keeper, addr sdk.AccAddress) string {
	return keeper.Ak.GetAccount(ctx, addr).GetCoins().String()
}

func frozenCoinsOf(ctx sdk.Context, keeper keeper.Keeper, addr sdk.AccAddress) string {
	accX, _ := keeper.Axk.GetAccountX(ctx, addr)
	return accX.FrozenCoins.String()
}

func TestFreezeMultiCoins(t *testing.T) {
	bkx, ctx := defaultContext()

	givenAccountWith(ctx, bkx, myaddr, "1000000000cet,100abc")

	freezeCoins, _ := sdk.ParseCoins("300000000cet, 20abc")
	err := bkx.FreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "80abc,700000000cet", coinsOf(ctx, bkx, myaddr))
	require.Equal(t, "20abc,300000000cet", frozenCoinsOf(ctx, bkx, myaddr))

	err = bkx.UnFreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "100abc,1000000000cet", coinsOf(ctx, bkx, myaddr))
	require.Equal(t, "", frozenCoinsOf(ctx, bkx, myaddr))
}

func TestFreezeUnFreezeOK(t *testing.T) {

	bkx, ctx := defaultContext()

	givenAccountWith(ctx, bkx, myaddr, "1000000000cet")

	freezeCoins := types.NewCetCoins(300000000)
	err := bkx.FreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "700000000cet", coinsOf(ctx, bkx, myaddr))
	require.Equal(t, "300000000cet", frozenCoinsOf(ctx, bkx, myaddr))

	err = bkx.UnFreezeCoins(ctx, myaddr, freezeCoins)

	require.Nil(t, err)
	require.Equal(t, "1000000000cet", coinsOf(ctx, bkx, myaddr))
	require.Equal(t, "", frozenCoinsOf(ctx, bkx, myaddr))
}

func TestFreezeUnFreezeInvalidAccount(t *testing.T) {

	bkx, ctx := defaultContext()

	freezeCoins := types.NewCetCoins(500000000)
	err := bkx.FreezeCoins(ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("insufficient account funds;  < 500000000cet"), err)

	err = bkx.UnFreezeCoins(ctx, myaddr, freezeCoins)
	require.Equal(t, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", myaddr)), err)
}

func TestFreezeUnFreezeInsufficientCoins(t *testing.T) {
	bkx, ctx := defaultContext()

	givenAccountWith(ctx, bkx, myaddr, "10cet")

	InvalidFreezeCoins := types.NewCetCoins(50)
	err := bkx.FreezeCoins(ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("insufficient account funds; 10cet < 50cet"), err)

	freezeCoins := types.NewCetCoins(5)
	err = bkx.FreezeCoins(ctx, myaddr, freezeCoins)
	require.Nil(t, err)

	err = bkx.UnFreezeCoins(ctx, myaddr, InvalidFreezeCoins)
	require.Equal(t, sdk.ErrInsufficientCoins("account has insufficient coins to unfreeze"), err)
}

func TestGetTotalCoins(t *testing.T) {
	bkx, ctx := defaultContext()
	givenAccountWith(ctx, bkx, myaddr, "100cet, 20bch, 30btc")

	lockedCoins := authx.LockedCoins{
		authx.NewLockedCoin("bch", sdk.NewInt(20), 1000),
		authx.NewLockedCoin("eth", sdk.NewInt(30), 2000),
	}

	frozenCoins := sdk.NewCoins(sdk.Coin{Denom: "btc", Amount: sdk.NewInt(50)},
		sdk.Coin{Denom: "eth", Amount: sdk.NewInt(10)},
	)

	accX := authx.AccountX{
		Address:     myaddr,
		LockedCoins: lockedCoins,
		FrozenCoins: frozenCoins,
	}

	bkx.Axk.SetAccountX(ctx, accX)

	expected := sdk.NewCoins(
		sdk.Coin{Denom: "bch", Amount: sdk.NewInt(40)},
		sdk.Coin{Denom: "btc", Amount: sdk.NewInt(80)},
		sdk.Coin{Denom: "cet", Amount: sdk.NewInt(100)},
		sdk.Coin{Denom: "eth", Amount: sdk.NewInt(40)},
	)
	expected = expected.Sort()
	coins := bkx.GetTotalCoins(ctx, myaddr)

	require.Equal(t, expected, coins)
}

func TestKeeper_TotalAmountOfCoin(t *testing.T) {

	bkx, ctx := defaultContext()
	amount := bkx.TotalAmountOfCoin(ctx, "cet")
	require.Equal(t, int64(0), amount.Int64())

	givenAccountWith(ctx, bkx, myaddr, "100cet")

	lockedCoins := authx.LockedCoins{
		authx.NewLockedCoin("cet", sdk.NewInt(100), 1000),
	}
	frozenCoins := sdk.NewCoins(sdk.Coin{Denom: "cet", Amount: sdk.NewInt(100)})

	accX := authx.AccountX{
		Address:     myaddr,
		LockedCoins: lockedCoins,
		FrozenCoins: frozenCoins,
	}
	bkx.Axk.SetAccountX(ctx, accX)
	amount = bkx.TotalAmountOfCoin(ctx, "cet")
	require.Equal(t, int64(300), amount.Int64())
}

func TestKeeper_AddCoins(t *testing.T) {
	bkx, ctx := defaultContext()
	coins := sdk.NewCoins(
		sdk.Coin{Denom: "aaa", Amount: sdk.NewInt(10)},
		sdk.Coin{Denom: "bbb", Amount: sdk.NewInt(20)},
	)

	coins2 := sdk.NewCoins(
		sdk.Coin{Denom: "aaa", Amount: sdk.NewInt(5)},
		sdk.Coin{Denom: "bbb", Amount: sdk.NewInt(10)},
	)

	err := bkx.AddCoins(ctx, myaddr, coins)
	require.Equal(t, nil, err)
	err = bkx.SubtractCoins(ctx, myaddr, coins2)
	require.Equal(t, nil, err)
	cs := bkx.GetTotalCoins(ctx, myaddr)
	require.Equal(t, coins2, cs)

	coins3 := sdk.NewCoins(
		sdk.Coin{Denom: "aaa", Amount: sdk.NewInt(15)},
		sdk.Coin{Denom: "bbb", Amount: sdk.NewInt(10)},
	)
	err = bkx.SubtractCoins(ctx, myaddr, coins3)
	require.Error(t, err)
}

func TestParamGetSet(t *testing.T) {
	bkx, ctx := defaultContext()

	//expect DefaultActivationFees=1
	defaultParam := bankx.DefaultParams()
	require.Equal(t, int64(100000000), defaultParam.ActivationFee)

	//expect SetParam don't panic
	require.NotPanics(t, func() { bkx.SetParams(ctx, defaultParam) }, "bankxKeeper SetParam panics")

	//expect GetParam equals defaultParam
	require.Equal(t, defaultParam, bkx.GetParams(ctx))
}

func TestKeeper_SendCoins(t *testing.T) {
	bkx, ctx := defaultContext()
	coins := sdk.NewCoins(
		sdk.Coin{Denom: "aaa", Amount: sdk.NewInt(10)},
	)
	addr2 := testutil.ToAccAddress("addr2")
	_ = bkx.AddCoins(ctx, myaddr, coins)
	exist := bkx.HasCoins(ctx, myaddr, coins)
	assert.True(t, exist)
	err := bkx.SendCoins(ctx, myaddr, addr2, coins)
	require.Equal(t, nil, err)
	cs := bkx.GetTotalCoins(ctx, addr2)
	require.Equal(t, coins, cs)
}
