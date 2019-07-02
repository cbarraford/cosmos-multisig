package types

import "strings"

type QueryWallets []MultiSigWallet

// implement fmt.Stringer
func (n QueryWallets) String() string {
	wallets := make([]string, len(n))
	for i, wallet := range n {
		wallets[i] = wallet.String()
	}
	return strings.Join(wallets[:], "\n")
}

type QueryTransactions []Transaction

// implement fmt.Stringer
func (n QueryTransactions) String() string {
	transactions := make([]string, len(n))
	for i, transaction := range n {
		transactions[i] = transaction.String()
	}
	return strings.Join(transactions[:], "\n")
}
