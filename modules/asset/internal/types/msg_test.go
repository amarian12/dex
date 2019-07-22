package types

import (
	"reflect"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestMsgIssueToken_ValidateBasic(t *testing.T) {

	tests := []struct {
		name string
		msg  MsgIssueToken
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgIssueToken("ABC Token", "abc", 100000, testAddr,
				false, false, false, false, "", ""),
			nil,
		},
		{
			"case-name",
			NewMsgIssueToken(string(make([]byte, 32+1)), "abc", 100000, testAddr,
				false, false, false, false, "", ""),
			ErrInvalidTokenName(string(make([]byte, 32+1))),
		},
		{
			"case-owner",
			NewMsgIssueToken("ABC Token", "abc", 100000, sdk.AccAddress{},
				false, false, false, false, "", ""),
			ErrNilTokenOwner(),
		},
		{
			"case-symbol1",
			NewMsgIssueToken("ABC Token", "1aa", 100000, testAddr,
				false, false, false, false, "", ""),
			ErrInvalidTokenSymbol("1aa"),
		},
		{
			"case-symbol2",
			NewMsgIssueToken("ABC Token", "A999", 100000, testAddr,
				false, false, false, false, "", ""),
			ErrInvalidTokenSymbol("A999"),
		},
		{
			"case-symbol3",
			NewMsgIssueToken("ABC Token", "aa1234567", 100000, testAddr,
				false, false, false, false, "", ""),
			ErrInvalidTokenSymbol("aa1234567"),
		},
		{
			"case-symbol4",
			NewMsgIssueToken("ABC Token", "a*aa", 100000, testAddr,
				false, false, false, false, "", ""),
			ErrInvalidTokenSymbol("a*aa"),
		},
		{
			"case-totalSupply1",
			NewMsgIssueToken("ABC Token", "abc", 9E18+1, testAddr,
				false, false, false, false, "", ""),
			ErrInvalidTokenSupply(9E18 + 1),
		},
		{
			"case-totalSupply2",
			NewMsgIssueToken("ABC Token", "abc", -1, testAddr,
				false, false, false, false, "", ""),
			ErrInvalidTokenSupply(-1),
		},
		{
			"case-url",
			NewMsgIssueToken("name", "coin", 2100, testAddr,
				false, false, false, false, string(make([]byte, MaxTokenURLLength+1)), ""),
			ErrInvalidTokenURL(string(make([]byte, MaxTokenURLLength+1))),
		},
		{
			"case-description",
			NewMsgIssueToken("name", "coin", 2100, testAddr,
				false, false, false, false, "", string(make([]byte, MaxTokenDescriptionLength+1))),
			ErrInvalidTokenDescription(string(make([]byte, MaxTokenDescriptionLength+1))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgIssueToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgTransferOwnership_ValidateBasic(t *testing.T) {
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")
	tests := []struct {
		name string
		msg  MsgTransferOwnership
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgTransferOwnership("abc", testAddr, addr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgTransferOwnership("123", testAddr, addr),
			ErrInvalidTokenSymbol("123"),
		},
		{
			"case-invalid1",
			NewMsgTransferOwnership("abc", sdk.AccAddress{}, testAddr),
			ErrNilTokenOwner(),
		},
		{
			"case-invalid2",
			NewMsgTransferOwnership("abc", testAddr, sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
		{
			"case-invalid3",
			NewMsgTransferOwnership("abc", testAddr, testAddr),
			ErrTransferSelfTokenOwner(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgTransferOwnership.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgMintToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgMintToken
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgMintToken("abc", 10000, testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgMintToken("()2", 10000, testAddr),
			ErrInvalidTokenSymbol("()2"),
		},
		{
			"case-invalidOwner",
			NewMsgMintToken("abc", 10000, sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidAmt1",
			NewMsgMintToken("abc", 9E18+1, testAddr),
			ErrInvalidTokenMintAmt(9E18 + 1),
		},
		{
			"case-invalidAmt2",
			NewMsgMintToken("abc", -1, testAddr),
			ErrInvalidTokenMintAmt(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgMintToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBurnToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgBurnToken
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgBurnToken("abc", 10000, testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgBurnToken("w♞", 10000, testAddr),
			ErrInvalidTokenSymbol("w♞"),
		},
		{
			"case-invalidOwner",
			NewMsgBurnToken("abc", 10000, sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidAmt1",
			NewMsgBurnToken("abc", 9E18+1, testAddr),
			ErrInvalidTokenBurnAmt(9E18 + 1),
		},
		{
			"case-invalidAmt2",
			NewMsgBurnToken("abc", -1, testAddr),
			ErrInvalidTokenBurnAmt(-1),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgBurnToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgForbidToken
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgForbidToken("abc", testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgForbidToken("*90", testAddr),
			ErrInvalidTokenSymbol("*90"),
		},
		{
			"case-invalidOwner",
			NewMsgForbidToken("abc", sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidToken_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgUnForbidToken
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgUnForbidToken("abc", testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgUnForbidToken("a¥0", testAddr),
			ErrInvalidTokenSymbol("a¥0"),
		},
		{
			"case-invalidOwner",
			NewMsgUnForbidToken("abc", sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidToken.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgAddTokenWhitelist_ValidateBasic(t *testing.T) {
	whitelist := mockAddrList()
	tests := []struct {
		name string
		msg  MsgAddTokenWhitelist
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgAddTokenWhitelist("abc", testAddr, whitelist),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgAddTokenWhitelist("abcdefghi", testAddr, whitelist),
			ErrInvalidTokenSymbol("abcdefghi"),
		},
		{
			"case-invalidOwner",
			NewMsgAddTokenWhitelist("abc", sdk.AccAddress{}, whitelist),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidWhitelist",
			NewMsgAddTokenWhitelist("abc", testAddr, []sdk.AccAddress{}),
			ErrNilTokenWhitelist(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgAddTokenWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgRemoveTokenWhitelist_ValidateBasic(t *testing.T) {
	whitelist := mockAddrList()
	tests := []struct {
		name string
		msg  MsgRemoveTokenWhitelist
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgRemoveTokenWhitelist("abc", testAddr, whitelist),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgRemoveTokenWhitelist("a℃", testAddr, whitelist),
			ErrInvalidTokenSymbol("a℃"),
		},
		{
			"case-invalidOwner",
			NewMsgRemoveTokenWhitelist("abc", sdk.AccAddress{}, whitelist),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidWhitelist",
			NewMsgRemoveTokenWhitelist("abc", testAddr, []sdk.AccAddress{}),
			ErrNilTokenWhitelist(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgRemoveTokenWhitelist.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgForbidAddr_ValidateBasic(t *testing.T) {
	addresses := mockAddrList()
	tests := []struct {
		name string
		msg  MsgForbidAddr
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgForbidAddr("abc", testAddr, addresses),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgForbidAddr("a⎝⎠", testAddr, addresses),
			ErrInvalidTokenSymbol("a⎝⎠"),
		},
		{
			"case-invalidOwner",
			NewMsgForbidAddr("abc", sdk.AccAddress{}, addresses),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidAddr",
			NewMsgForbidAddr("abc", testAddr, []sdk.AccAddress{}),
			ErrNilForbiddenAddress(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgForbidAddr.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgUnForbidAddr_ValidateBasic(t *testing.T) {
	addr := mockAddrList()
	tests := []struct {
		name string
		msg  MsgUnForbidAddr
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgUnForbidAddr("abc", testAddr, addr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgUnForbidAddr("a⥇", testAddr, addr),
			ErrInvalidTokenSymbol("a⥇"),
		},
		{
			"case-invalidOwner",
			NewMsgUnForbidAddr("abc", sdk.AccAddress{}, addr),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidAddr",
			NewMsgUnForbidAddr("abc", testAddr, []sdk.AccAddress{}),
			ErrNilForbiddenAddress(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgUnForbidAddr.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgModifyTokenURL_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgModifyTokenURL
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgModifyTokenURL("abc", "www.abc.org", testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgModifyTokenURL("a😃", "www.abc.org", testAddr),
			ErrInvalidTokenSymbol("a😃"),
		},
		{
			"case-invalidOwner",
			NewMsgModifyTokenURL("abc", "www.abc.org", sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidURL",
			NewMsgModifyTokenURL("abc", string(make([]byte, MaxTokenURLLength+1)), testAddr),
			ErrInvalidTokenURL(string(make([]byte, MaxTokenURLLength+1))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MsgModifyTokenURL.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsgModifyTokenDescription_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgModifyTokenDescription
		want sdk.Error
	}{
		{
			"base-case",
			NewMsgModifyTokenDescription("abc", "abc example description", testAddr),
			nil,
		},
		{
			"case-invalidSymbol",
			NewMsgModifyTokenDescription("a❡", "abc example description", testAddr),
			ErrInvalidTokenSymbol("a❡"),
		},
		{
			"case-invalidOwner",
			NewMsgModifyTokenDescription("abc", "abc example description", sdk.AccAddress{}),
			ErrNilTokenOwner(),
		},
		{
			"case-invalidDescription",
			NewMsgModifyTokenDescription("abc", string(make([]byte, MaxTokenDescriptionLength+1)), testAddr),
			ErrInvalidTokenDescription(string(make([]byte, MaxTokenDescriptionLength+1))),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.ValidateBasic(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMsgModifyTokenDescription.ValidateBasic() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestMsg_Route(t *testing.T) {
	want := RouterKey
	tests := []struct {
		name string
		msg  sdk.Msg
	}{
		{
			"issue-token",
			MsgIssueToken{},
		},
		{
			"transfer-ownership",
			MsgTransferOwnership{},
		},
		{
			"burn-token",
			MsgBurnToken{},
		},
		{
			"mint-token",
			MsgMintToken{},
		},
		{
			"forbid-token",
			MsgForbidToken{},
		},
		{
			"unforbid-token",
			MsgUnForbidToken{},
		},
		{
			"add_token_whitelist",
			MsgAddTokenWhitelist{},
		},
		{
			"remove-token-whitelist",
			MsgRemoveTokenWhitelist{},
		},
		{
			"forbid-addr",
			MsgForbidAddr{},
		},
		{
			"unforbid-addr",
			MsgUnForbidAddr{},
		},
		{
			"modify-token-url",
			MsgModifyTokenURL{},
		},
		{
			"modify-token-description",
			MsgModifyTokenDescription{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.Route(); got != want {
				t.Errorf("Msg.Route() = %v, want %v", got, want)
			}
		})
	}
}

func TestMsg_Type(t *testing.T) {
	tests := []struct {
		name string
		msg  sdk.Msg
		want string
	}{
		{
			"issue-token",
			MsgIssueToken{},
			"issue_token",
		},
		{
			"transfer-ownership",
			MsgTransferOwnership{},
			"transfer_ownership",
		},
		{
			"burn-token",
			MsgBurnToken{},
			"burn_token",
		},
		{
			"mint-token",
			MsgMintToken{},
			"mint_token",
		},
		{
			"forbid-token",
			MsgForbidToken{},
			"forbid_token",
		},
		{
			"unforbid-token",
			MsgUnForbidToken{},
			"unforbid_token",
		},
		{
			"add_token_whitelist",
			MsgAddTokenWhitelist{},
			"add_token_whitelist",
		},
		{
			"remove-token-whitelist",
			MsgRemoveTokenWhitelist{},
			"remove_token_whitelist",
		},
		{
			"forbid-addr",
			MsgForbidAddr{},
			"forbid_addr",
		},
		{
			"unforbid-addr",
			MsgUnForbidAddr{},
			"unforbid_addr",
		},
		{
			"modify-token-url",
			MsgModifyTokenURL{},
			"modify_token_url",
		},
		{
			"modify-token-description",
			MsgModifyTokenDescription{},
			"modify_token_description",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.Type(); got != tt.want {
				t.Errorf("Msg.Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsg_GetSigners(t *testing.T) {
	tests := []struct {
		name string
		msg  sdk.Msg
		want []sdk.AccAddress
	}{
		{
			"issue-token",
			NewMsgIssueToken("ABC Token", "abc", 100000, testAddr,
				false, false, false, false, "", ""),
			[]sdk.AccAddress{testAddr},
		},
		{
			"transfer-ownership",
			NewMsgTransferOwnership("abc", testAddr, sdk.AccAddress{}),
			[]sdk.AccAddress{testAddr},
		},
		{
			"burn-token",
			NewMsgBurnToken("abc", 100000, testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"mint-token",
			NewMsgMintToken("abc", 100000, testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"forbid-token",
			NewMsgForbidToken("abc", testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"unforbid-token",
			NewMsgUnForbidToken("abc", testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"add_token_whitelist",
			NewMsgAddTokenWhitelist("abc", testAddr, mockAddrList()),
			[]sdk.AccAddress{testAddr},
		},
		{
			"remove-token-whitelist",
			NewMsgRemoveTokenWhitelist("abc", testAddr, mockAddrList()),
			[]sdk.AccAddress{testAddr},
		},
		{
			"forbid-addr",
			NewMsgForbidAddr("abc", testAddr, mockAddrList()),
			[]sdk.AccAddress{testAddr},
		},
		{
			"unforbid-addr",
			NewMsgUnForbidAddr("abc", testAddr, mockAddrList()),
			[]sdk.AccAddress{testAddr},
		},
		{
			"modify-token-url",
			NewMsgModifyTokenURL("abc", "www.abc.com", testAddr),
			[]sdk.AccAddress{testAddr},
		},
		{
			"modify-token-description",
			NewMsgModifyTokenDescription("abc", "abc example description", testAddr),
			[]sdk.AccAddress{testAddr},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.GetSigners(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Msg.GetSigners() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsg_GetSignBytes(t *testing.T) {
	var owner, _ = sdk.AccAddressFromBech32("coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd")
	var addr, _ = sdk.AccAddressFromBech32("coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5")

	var addr1, _ = sdk.AccAddressFromBech32("coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke")
	var addr2, _ = sdk.AccAddressFromBech32("coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h")
	var addrList = []sdk.AccAddress{addr1, addr2}

	tests := []struct {
		name string
		msg  sdk.Msg
		want string
	}{
		{
			"issue-token",
			NewMsgIssueToken("ABC Token", "abc", 100000, owner,
				false, false, false, false, "", ""),
			`{"type":"asset/MsgIssueToken","value":{"addr_forbiddable":false,"burnable":false,"description":"","mintable":false,"name":"ABC Token","owner":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","token_forbiddable":false,"total_supply":"100000","url":""}}`,
		},
		{
			"transfer-ownership",
			NewMsgTransferOwnership("abc", owner, addr),
			`{"type":"asset/MsgTransferOwnership","value":{"new_owner":"coinex1e9kx6klg6z9p9ea4ehqmypl6dvjrp96vfxecd5","original_owner":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"burn-token",
			NewMsgBurnToken("abc", 100000, owner),
			`{"type":"asset/MsgBurnToken","value":{"amount":"100000","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"mint-token",
			NewMsgMintToken("abc", 100000, owner),
			`{"type":"asset/MsgMintToken","value":{"amount":"100000","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"forbid-token",
			NewMsgForbidToken("abc", owner),
			`{"type":"asset/MsgForbidToken","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"unforbid-token",
			NewMsgUnForbidToken("abc", owner),
			`{"type":"asset/MsgUnForbidToken","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"add_token_whitelist",
			NewMsgAddTokenWhitelist("abc", owner, addrList),
			`{"type":"asset/MsgAddTokenWhitelist","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","whitelist":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"]}}`,
		},
		{
			"remove-token-whitelist",
			NewMsgRemoveTokenWhitelist("abc", owner, addrList),
			`{"type":"asset/MsgRemoveTokenWhitelist","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","whitelist":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"]}}`,
		},
		{
			"forbid-addr",
			NewMsgForbidAddr("abc", owner, addrList),
			`{"type":"asset/MsgForbidAddr","value":{"addresses":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"],"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"unforbid-addr",
			NewMsgUnForbidAddr("abc", owner, addrList),
			`{"type":"asset/MsgUnForbidAddr","value":{"addresses":["coinex1y5kdxnzn2tfwayyntf2n28q8q2s80mcul852ke","coinex133w8vwj73s4h2uynqft9gyyy52cr6rg8dskv3h"],"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
		{
			"modify-token-url",
			NewMsgModifyTokenURL("abc", "www.abc.com", owner),
			`{"type":"asset/MsgModifyTokenURL","value":{"owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc","url":"www.abc.com"}}`,
		},
		{
			"modify-token-description",
			NewMsgModifyTokenDescription("abc", "abc example description", owner),
			`{"type":"asset/MsgModifyTokenDescription","value":{"description":"abc example description","owner_address":"coinex15fvnexrvsm9ryw3nn4mcrnqyhvhazkkrd4aqvd","symbol":"abc"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.msg.GetSignBytes(); !reflect.DeepEqual(string(got), tt.want) {
				t.Errorf("Msg.GetSignBytes() = %s, want %s", got, tt.want)
			}
		})
	}
}
