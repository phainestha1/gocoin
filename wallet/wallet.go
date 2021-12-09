package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"gocoin/utils"
	"math/big"
	"os"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address string
}

var w *wallet
const fileName string = "coin.wallet"

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleErr(err)
}

func restoreKey() *ecdsa.PrivateKey {
	keyAsByte, err := os.ReadFile(fileName)
	utils.HandleErr(err)
	key, err := x509.ParseECPrivateKey(keyAsByte)
	utils.HandleErr(err)
	return key
}

func encodeBigInts(a, b []byte) string{
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func addressFromKey(key *ecdsa.PrivateKey) string {
	return encodeBigInts(key.X.Bytes(), key.Y.Bytes())
} 

func Sign(payload string, w wallet) string {
	payloadAsByte, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsByte)
	utils.HandleErr(err)
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondHalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondHalfBytes)
	return &bigA, &bigB, nil
}

func Verify(signature, payload, address string) bool {
	r, s, err := restoreBigInts(signature)
	utils.HandleErr(err)
	
	x, y, err := restoreBigInts(address)
	utils.HandleErr(err) 
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X: x,
		Y: y,
	}

	payloadAsBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)

	ok := ecdsa.Verify(&publicKey, payloadAsBytes, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = addressFromKey(w.privateKey)
	}
	return w
}