package rsax

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"

	"github.com/livexy/pkg/bytex"
	"github.com/livexy/pkg/strx"
)

// RSA生成公私密钥
func RSAGenKey(bits int) (pub string, pri string, ok bool) {
	if bits%1024 != 0 {
		return
	}
	privatekey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return
	}
	privatestream := x509.MarshalPKCS1PrivateKey(privatekey)
	block1 := pem.Block{Type: "private key", Bytes: privatestream}
	pri = bytex.ToStr(pem.EncodeToMemory(&block1))
	publickey := privatekey.PublicKey
	publicstream, err := x509.MarshalPKIXPublicKey(&publickey)
	if err != nil {
		return
	}
	block2 := pem.Block{Type: "public key", Bytes: publicstream}
	pub = bytex.ToStr(pem.EncodeToMemory(&block2))
	ok = true
	return
}

// RSA加密
func Encrypt(pubkey, data string) string {
	block, _ := pem.Decode(strx.ToBytes(pubkey))
	if block == nil {
		return ""
	}
	pubinterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return ""
	}
	pub := pubinterface.(*rsa.PublicKey)
	res, err := rsa.EncryptPKCS1v15(rand.Reader, pub, strx.ToBytes(data))
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(res)
}

// RSA解密
func Decrypt(prikey, data string) string {
	if len(data) < 4 {
		return ""
	}
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	block, _ := pem.Decode(strx.ToBytes(prikey))
	if block == nil {
		return ""
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return ""
	}
	text, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		return ""
	}
	return bytex.ToStr(text)
}
