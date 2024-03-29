package blind

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"errors"
)

type ECDSAConfig struct {
	Key *ecdsa.PrivateKey
}

func (c *ECDSAConfig) Sign(data []byte) ([]byte, error) {
	signature, err := newECDSASignature(data, c.Key)
	if err != nil {
		return nil, err
	}

	o, err := signature.Marshall()
	if err != nil {
		return nil, err
	}

	return o, nil
}

func (c *ECDSAConfig) Verify(foreignPublicKey *ECDSAConfig, signatureJson []byte) (bool, error) {
	signature := ECDSASignature{}
	err := signature.Unmarshall(signatureJson)
	if err != nil {
		return false, errors.New("unable to unmarshal signature json")
	}

	return ecdsa.Verify(&foreignPublicKey.Key.PublicKey, signature.DataHash, signature.R, signature.S), nil
}

func (c *ECDSAConfig) Marshall() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"key": c.Key,
	})
}

func (c *ECDSAConfig) Unmarshall(data []byte) error {
	var fields struct {
		Key *ecdsa.PrivateKey `json:"key"`
	}
	err := json.Unmarshal(data, &fields)
	if err != nil {
		return err
	}
	c.Key = fields.Key
	return nil
}

func (c *ECDSAConfig) ECDH(foreignPublicKey *crypto.PublicKey) ([]byte, error) {

	ecdhPrivate, err := c.Key.ECDH()
	if err != nil {
		return nil, err
	}

	ecdhPublic, err := (*foreignPublicKey).(*ecdsa.PublicKey).ECDH()
	if err != nil {
		return nil, err
	}

	sharedKey, err := ecdhPrivate.ECDH(ecdhPublic)
	if err != nil {
		return nil, err
	}

	return sharedKey, nil
}

func NewECDSAConfig() *ECDSAConfig {
	k, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	return &ECDSAConfig{
		k,
	}
}
