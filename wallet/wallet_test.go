package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testKey string = "30770201010420dd7dccf9afbc67278d548526070a4fe51a669d92c13b3ab8ad42b231a813c82ca00a06082a8648ce3d030107a14403420004b34901f6caf0f963ac23f4a2bfd63691fbc0af1143d28b0b0c2a1fa386a20a8a92fc34d09afa1ac8973b8a61842783a4efbe3cb64ca59f20db56b83c27ddf806"
	testPayload string= "00e404d3e8d6a48301f553194640c951d0f91dcc8132fbc41d7c938c788bd068"
	testSig string= "bdcf6b12497868c2b76bb743b66622713cef8b48c99d4d6addc3b7ad33f599c4f0c36e1ac7b1e64457d6a65d08ade4b7b1a6a5e9288700e91e8604298c97830c"
)

// 이부분이 좀 헷갈린다.
type testLayer struct{
	fakeHasWalletFile func() bool
}

func(t testLayer) hasWalletFile() bool {
	// 원래 코드의 if files.hasWalletFile() 을 구현하기 위해 만들어짐
	return t.fakeHasWalletFile() 
}

func(testLayer) writeFile(name string, data []byte, perm fs.FileMode) error {	
	return nil
}

func(testLayer) readFile(name string) ([]byte, error){	
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

// 기존 wallet은 sideEffect가 많은 구조이므로 테스트용 지갑 만들기부터 진행
func makeTestWallet() *wallet{
	w := &wallet{}
	b , _ := hex.DecodeString(testKey)
	key , _ :=x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = addressFromKey(key)
	return w
}

func TestWallet(t *testing.T){
	t.Run("Wallet is created", func(t *testing.T){
		files = testLayer{
			fakeHasWalletFile: func() bool {
				t.Log("I have been called")
				return false
			},
		}
		tw:= Wallet()

		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}){
			t.Error("New Wallet should return new wallet instance")
		}
	})
	t.Run("Wallet is restored", func(t *testing.T){
		files = testLayer{
			fakeHasWalletFile: func() bool {
				t.Log("I have been called")
				return true				
			},
		}
		w = nil
		tw:= Wallet()

		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}){
			t.Error("New Wallet should return new wallet instance")
		}
	})
}


func TestSign(t *testing.T) {
	s := Sign(testPayload, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct{
		input string
		ok bool
	}
	tests := []test{
		{testPayload, true},
		{"40e404d3e8d6a48301f553194640c951d0f91dcc8132fbc41d7c938c788bd068", false},
	}

	for _, tc := range tests{
		w:= makeTestWallet()	
		ok := Verify(testSig,tc.input,w.Address)
		if ok != tc.ok{
			t.Error("Verify() could not verify testSignature and Payload")
		}
	}	
}

func TestRestoreBigInts(t *testing.T){
	_,_,err := restoreBigInts("xx")
	if err == nil {
		t.Error("restoreBigInts() should return an error when payload is not hex")
	}
}