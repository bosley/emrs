package badger

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	b64 "encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"math/big"
	"time"
)

type PubHashSig struct {
	PubKey []byte
	Hash   []byte
	Sig    []byte
}

type Badge interface {
	Id() string
	Nickname() string

	Sign(data *string) (PubHashSig, error)

	GenerateVoucher(expiration time.Duration) (string, error)
	ValidateVoucher(voucher string) bool

	PublicKey() string
	EncodeIdentity() EncodedIdentity
	EncodeIdentityString() string
}

type EncodedIdentity struct {
	Id         string `json:id`
	Nickname   string `json:nickname`
	PublicKey  string `json:public_key`
	PrivateKey string `json:private_key`
}

func New(nickname string) (Badge, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return &identity{
		nickname: nickname,
		key:      generateKeyPair(),
		uid: fmt.Sprintf(
			"%x-%x-%x-%x-%x",
			b[0:4],
			b[4:6],
			b[6:8],
			b[8:10],
			b[10:]),
	}, nil
}

func zeroArr(raw []byte) {
	for i := 0; i < len(raw); i++ {
		raw[i] = 0
	}
}

func RawIsHashMatch(raw []byte, hashed []byte) error {
	defer zeroArr(raw)
	return bcrypt.CompareHashAndPassword(hashed, raw)
}

func Hash(raw []byte) ([]byte, error) {
	defer zeroArr(raw)
	return bcrypt.GenerateFromPassword(raw, bcrypt.DefaultCost)
}

func ToFormattedString(badge Badge) string {
	return fmt.Sprintf(
		"\n\tNICK: %s\n\t UID: %s\n\t KEY: %s\n",
		badge.Nickname(),
		badge.Id(),
		badge.PublicKey())
}

func getCurve() elliptic.Curve {
	return elliptic.P256()
	/*
	   return elliptic.P224()
	   return elliptic.P384()
	   return elliptic.P521()
	*/
}

func generateKeyPair() *ecdsa.PrivateKey {
	key, err := ecdsa.GenerateKey(getCurve(), rand.Reader)
	if err != nil {
		panic(err.Error())
	}
	return key
}

type identity struct {
	nickname string
	uid      string
	key      *ecdsa.PrivateKey
	ks       int
}

func (id *identity) MarshalPublicKeyBytes() []byte {
	return MarshalPublicKey(id.key.PublicKey.X, id.key.PublicKey.Y)
}

func MarshalPublicKey(x *big.Int, y *big.Int) []byte {
	return elliptic.MarshalCompressed(getCurve(), x, y)
}

func UnmarshalPublicKey(publicKey []byte) ecdsa.PublicKey {
	curve := getCurve()
	x, y := elliptic.UnmarshalCompressed(curve, publicKey)
	return ecdsa.PublicKey{
		curve,
		x,
		y,
	}
}

func (id *identity) Id() string {
	return id.uid
}

func (id *identity) Nickname() string {
	return id.nickname
}

func (id *identity) Sign(data *string) (PubHashSig, error) {
	if data == nil {
		return PubHashSig{}, errors.New("nil data given")
	}

	hash := sha256.New()
	hash.Write([]byte(*data))

	hashBytes := hash.Sum(nil)
	sigBytes, err := ecdsa.SignASN1(rand.Reader, id.key, hashBytes)
	if err != nil {
		return PubHashSig{}, err
	}

	return PubHashSig{
		id.MarshalPublicKeyBytes(),
		hashBytes,
		sigBytes,
	}, nil
}

func Verify(phs PubHashSig) bool {
	publicKey := UnmarshalPublicKey(phs.PubKey)
	return ecdsa.VerifyASN1(&publicKey, phs.Hash, phs.Sig)
}

func (id *identity) PublicKey() string {
	return b64.StdEncoding.EncodeToString(
		elliptic.MarshalCompressed(getCurve(), id.key.PublicKey.X, id.key.PublicKey.Y))
}

func publicKeyFromB64(key string) (*ecdsa.PublicKey, error) {
	dk, err := b64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	x, y := elliptic.UnmarshalCompressed(getCurve(), dk)
	return &ecdsa.PublicKey{
		getCurve(),
		x,
		y,
	}, nil
}

func (id *identity) PrivateKey() string {
	x509Encoded, _ := x509.MarshalECPrivateKey(id.key)
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded}))
}

func (id *identity) EncodeIdentity() EncodedIdentity {

	eid := EncodedIdentity{
		Id:         id.uid,
		Nickname:   id.nickname,
		PublicKey:  id.PublicKey(),
		PrivateKey: id.PrivateKey(),
	}
	return eid
}

func (id *identity) EncodeIdentityString() string {

	eid := EncodedIdentity{
		Id:         id.uid,
		Nickname:   id.nickname,
		PublicKey:  id.PublicKey(),
		PrivateKey: id.PrivateKey(),
	}

	b, _ := json.Marshal(eid)
	return string(b)
}

func DecodeIdentity(eid EncodedIdentity) (Badge, error) {
	block, _ := pem.Decode([]byte(eid.PrivateKey))
	if block == nil {
		return nil, errors.New("no key extracted from identity")
	}

	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, err
	}
	id := identity{
		nickname: eid.Nickname,
		uid:      eid.Id,
		key:      privateKey,
	}
	return &id, nil
}

func DecodeIdentityString(encodedId string) (Badge, error) {

	did := EncodedIdentity{}

	if err := json.Unmarshal([]byte(encodedId), &did); err != nil {
		return nil, err
	}

	block, _ := pem.Decode([]byte(did.PrivateKey))
	if block == nil {
		return nil, errors.New("no key extracted from identity")
	}
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		return nil, err
	}

	id := identity{
		nickname: did.Nickname,
		uid:      did.Id,
		key:      privateKey,
	}
	return &id, nil
}

func (id *identity) GenerateVoucher(expiration time.Duration) (string, error) {
	return NewVoucher(id, expiration)
}

func (id *identity) ValidateVoucher(voucher string) bool {
	return ValidateVoucher(id.PublicKey(), voucher)
}
