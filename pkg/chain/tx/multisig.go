package wallettx

import (
	"errors"

	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
	walletdb "github.com/p9c/pod/pkg/wallet/db"
)

// MakeMultiSigScript creates a multi-signature script that can be redeemed with
// nRequired signatures of the passed keys and addresses.  If the address is a
// P2PKH address, the associated pubkey is looked up by the wallet if possible,
// otherwise an error is returned for a missing pubkey.
//
// This function only works with pubkeys and P2PKH addresses derived from them.
func (w *Wallet) MakeMultiSigScript(addrs []util.Address, nRequired int) ([]byte, error) {
	pubKeys := make([]*util.AddressPubKey, len(addrs))
	var dbtx walletdb.ReadTx
	var addrmgrNs walletdb.ReadBucket
	defer func() {
		if dbtx != nil {
			err := dbtx.Rollback()
			if err != nil {
				log.ERROR(err)
			}
		}
	}()
	// The address list will made up either of addreseses (pubkey hash), for
	// which we need to look up the keys in wallet, straight pubkeys, or a
	// mixture of the two.
	for i, addr := range addrs {
		switch addr := addr.(type) {
		default:
			return nil, errors.New("cannot make multisig script for " +
				"a non-secp256k1 public key or P2PKH address")
		case *util.AddressPubKey:
			pubKeys[i] = addr
		case *util.AddressPubKeyHash:
			if dbtx == nil {
				var err error
				dbtx, err = w.db.BeginReadTx()
				if err != nil {
					log.ERROR(err)
					return nil, err
				}
				addrmgrNs = dbtx.ReadBucket(waddrmgrNamespaceKey)
			}
			addrInfo, err := w.Manager.Address(addrmgrNs, addr)
			if err != nil {
				log.ERROR(err)
				return nil, err
			}
			serializedPubKey := addrInfo.(waddrmgr.ManagedPubKeyAddress).
				PubKey().SerializeCompressed()
			pubKeyAddr, err := util.NewAddressPubKey(
				serializedPubKey, w.chainParams)
			if err != nil {
				log.ERROR(err)
				return nil, err
			}
			pubKeys[i] = pubKeyAddr
		}
	}
	return txscript.MultiSigScript(pubKeys, nRequired)
}

// ImportP2SHRedeemScript adds a P2SH redeem script to the wallet.
func (w *Wallet) ImportP2SHRedeemScript(script []byte) (*util.AddressScriptHash, error) {
	var p2shAddr *util.AddressScriptHash
	err := walletdb.Update(w.db, func(tx walletdb.ReadWriteTx) error {
		addrmgrNs := tx.ReadWriteBucket(waddrmgrNamespaceKey)
		// TODO(oga) blockstamp current block?
		bs := &waddrmgr.BlockStamp{
			Hash:   *w.ChainParams().GenesisHash,
			Height: 0,
		}
		// As this is a regular P2SH script, we'll import this into the
		// BIP0044 scope.
		bip44Mgr, err := w.Manager.FetchScopedKeyManager(
			waddrmgr.KeyScopeBIP0084,
		)
		if err != nil {
			log.ERROR(err)
			return err
		}
		addrInfo, err := bip44Mgr.ImportScript(addrmgrNs, script, bs)
		if err != nil {
			log.ERROR(err)
			// Don't care if it's already there, but still have to
			// set the p2shAddr since the address manager didn't
			// return anything useful.
			if waddrmgr.IsError(err, waddrmgr.ErrDuplicateAddress) {
				// This function will never error as it always
				// hashes the script to the correct length.
				p2shAddr, _ = util.NewAddressScriptHash(script,
					w.chainParams)
				return nil
			}
			return err
		}
		p2shAddr = addrInfo.Address().(*util.AddressScriptHash)
		return nil
	})
	return p2shAddr, err
}