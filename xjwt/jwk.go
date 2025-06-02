// Copyright 2022 The imkuqin-zw Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package xjwt 包提供了一些用于生成JWT密钥的函数
package xjwt

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/x25519"
)

// GenerateRsaKey 生成2048位的RSA私钥
// 返回: *rsa.PrivateKey 生成的RSA私钥
// 返回: error 生成过程中遇到的错误
func GenerateRsaKey() (*rsa.PrivateKey, error) {
	return rsa.GenerateKey(rand.Reader, 2048)
}

// GenerateRsaJwk 生成RSA类型的JWK(JSON Web Key)
// 返回: jwk.Key 生成的JWK密钥
// 返回: error 生成过程中遇到的错误
func GenerateRsaJwk() (jwk.Key, error) {
	key, err := GenerateRsaKey()
	if err != nil {
		return nil, fmt.Errorf(`failed to generate RSA private key: %w`, err)
	}

	k, err := jwk.FromRaw(key)
	if err != nil {
		return nil, fmt.Errorf(`failed to generate jwk.RSAPrivateKey: %w`, err)
	}

	return k, nil
}

// GenerateRsaPublicJwk 从RSA私钥生成对应的公钥JWK
// 返回: jwk.Key 生成的公钥JWK
// 返回: error 生成过程中遇到的错误
func GenerateRsaPublicJwk() (jwk.Key, error) {
	key, err := GenerateRsaJwk()
	if err != nil {
		return nil, fmt.Errorf(`failed to generate jwk.RSAPrivateKey: %w`, err)
	}

	return jwk.PublicKeyOf(key)
}

// GenerateEcdsaKey 根据指定的椭圆曲线算法生成ECDSA私钥
// 参数: alg 椭圆曲线算法类型
// 返回: *ecdsa.PrivateKey 生成的ECDSA私钥
// 返回: error 生成过程中遇到的错误
func GenerateEcdsaKey(alg jwa.EllipticCurveAlgorithm) (*ecdsa.PrivateKey, error) {
	var crv elliptic.Curve
	if tmp, ok := CurveForAlgorithm(alg); ok {
		crv = tmp
	} else {
		return nil, fmt.Errorf(`invalid curve algorithm %s`, alg)
	}

	return ecdsa.GenerateKey(crv, rand.Reader)
}

// GenerateEcdsaJwk 生成P521曲线ECDSA的JWK
// 返回: jwk.Key 生成的JWK密钥
// 返回: error 生成过程中遇到的错误
func GenerateEcdsaJwk() (jwk.Key, error) {
	key, err := GenerateEcdsaKey(jwa.P521)
	if err != nil {
		return nil, fmt.Errorf(`failed to generate ECDSA private key: %w`, err)
	}

	k, err := jwk.FromRaw(key)
	if err != nil {
		return nil, fmt.Errorf(`failed to generate jwk.ECDSAPrivateKey: %w`, err)
	}

	return k, nil
}

// GenerateEcdsaPublicJwk 从ECDSA私钥生成对应的公钥JWK
// 返回: jwk.Key 生成的公钥JWK
// 返回: error 生成过程中遇到的错误
func GenerateEcdsaPublicJwk() (jwk.Key, error) {
	key, err := GenerateEcdsaJwk()
	if err != nil {
		return nil, fmt.Errorf(`failed to generate jwk.ECDSAPrivateKey: %w`, err)
	}

	return jwk.PublicKeyOf(key)
}

// GenerateSymmetricKey 生成64字节的对称密钥
// 返回: []byte 生成的对称密钥字节数组
func GenerateSymmetricKey() []byte {
	sharedKey := make([]byte, 64)
	_, _ = rand.Read(sharedKey)
	return sharedKey
}

// GenerateSymmetricJwk 生成对称密钥的JWK
// 返回: jwk.Key 生成的JWK密钥
// 返回: error 生成过程中遇到的错误
func GenerateSymmetricJwk() (jwk.Key, error) {
	key, err := jwk.FromRaw(GenerateSymmetricKey())
	if err != nil {
		return nil, fmt.Errorf(`failed to generate jwk.SymmetricKey: %w`, err)
	}

	return key, nil
}

// GenerateEd25519Key 生成Ed25519签名算法的私钥
// 返回: ed25519.PrivateKey 生成的私钥
// 返回: error 生成过程中遇到的错误
func GenerateEd25519Key() (ed25519.PrivateKey, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	return priv, err
}

// GenerateEd25519Jwk 生成Ed25519的JWK
// 返回: jwk.Key 生成的JWK密钥
// 返回: error 生成过程中遇到的错误
func GenerateEd25519Jwk() (jwk.Key, error) {
	key, err := GenerateEd25519Key()
	if err != nil {
		return nil, fmt.Errorf(`failed to generate Ed25519 private key: %w`, err)
	}

	k, err := jwk.FromRaw(key)
	if err != nil {
		return nil, fmt.Errorf(`failed to generate jwk.OKPPrivateKey: %w`, err)
	}

	return k, nil
}

// GenerateX25519Key 生成X25519密钥交换算法的私钥
// 返回: x25519.PrivateKey 生成的私钥
// 返回: error 生成过程中遇到的错误
func GenerateX25519Key() (x25519.PrivateKey, error) {
	_, priv, err := x25519.GenerateKey(rand.Reader)
	return priv, err
}

// GenerateX25519Jwk 生成X25519的JWK
// 返回: jwk.Key 生成的JWK密钥
// 返回: error 生成过程中遇到的错误
func GenerateX25519Jwk() (jwk.Key, error) {
	key, err := GenerateX25519Key()
	if err != nil {
		return nil, fmt.Errorf(`failed to generate X25519 private key: %w`, err)
	}

	k, err := jwk.FromRaw(key)
	if err != nil {
		return nil, fmt.Errorf(`failed to generate jwk.OKPPrivateKey: %w`, err)
	}

	return k, nil
}
