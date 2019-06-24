package multisig

import (
	"testing"

	"github.com/cbarraford/cosmos-multisig/x/multisig/types"
	"github.com/cbarraford/parsec"
	ctype "github.com/cosmos/cosmos-sdk/types"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type KeeperSuite struct{}

var _ = Suite(&KeeperSuite{})

func (s *KeeperSuite) TestKeeper(c *C) {
	cdc := parsec.MakeCodec()
	storeKey := ctype.NewKVStoreKey("test")
	bank := parsec.NewBankKeeper()

	keeper := NewKeeper(bank, storeKey, cdc)

	wallet := types.NewMultiSigWallet("test-wallet", nil, 5)

	keeper.SetWallet(parsec.MockContext{}, "test-wallet", wallet)
	wallet = keeper.GetWallet(parsec.MockContext{}, "test-wallet")
	c.Assert(wallet.MinSigTx, Equals, 5)
}
