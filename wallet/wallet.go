package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/lovelycbm/bongcoin/utils"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address string
}

var w *wallet;

const (
	fileName string = "bong.wallet"
)

func hasWalletFile() bool {
	_, err := os.Stat(fileName)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey{
	privKey , err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleError(err)
	return privKey
}

func restoreKey() (key *ecdsa.PrivateKey){
	keyAsBytes, err := os.ReadFile(fileName)
	utils.HandleError(err)
	key, err = x509.ParseECPrivateKey(keyAsBytes)
	utils.HandleError(err)
	return 
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleError(err)
	err = os.WriteFile(fileName, bytes, 0644)
	utils.HandleError(err)
}

func encodeBigInts(a,b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func addressFromKey(key *ecdsa.PrivateKey) string {
	x := key.X.Bytes()
	y := key.Y.Bytes()			
	return encodeBigInts(x,y)
}

func Sign(payload string, w *wallet) string {
	payloadAsByte, err := hex.DecodeString(payload)
	utils.HandleError(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsByte)
	utils.HandleError(err)		
	return encodeBigInts(r.Bytes(), s.Bytes())
}

func restoreBigInts(payload string)(*big.Int, *big.Int, error) {
	bytes , err := hex.DecodeString(payload)	
	if err != nil{
		return nil,nil,err
	}
	firstHalfBytes := bytes[:len(bytes)/2]
	secondhalfBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstHalfBytes)
	bigB.SetBytes(secondhalfBytes)
	return &bigA, &bigB, nil
}

func Verify(signature , payload , address string) bool {	
	r,s,err:= restoreBigInts(signature)
	utils.HandleError(err)
	x,y,err := restoreBigInts(address)
	utils.HandleError(err)

	publicKey := ecdsa.PublicKey{Curve: elliptic.P256(), X:x, Y:y}

	payloadAsByte, err := hex.DecodeString(payload)
	utils.HandleError(err)
	ok:= ecdsa.Verify(&publicKey, payloadAsByte, r, s)
	return ok
}

func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile(){
			w.privateKey = restoreKey()
			// yes -> restore form file 
		} else {
			key:= createPrivKey()			
			persistKey(key)	
			w.privateKey = key
		}	

		w.Address = addressFromKey(w.privateKey)
	}
	return w
}