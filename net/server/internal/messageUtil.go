package databag

import (
  "errors"
  "crypto"
  "crypto/rand"
  "crypto/sha256"
  "crypto/rsa"
  "encoding/json"
  "encoding/hex"
  "encoding/base64"
  "time"
)

func ReadDataMessage(msg *DataMessage, obj interface{}) (string, string, int64, error) {

  var data []byte
  var hash [32]byte
  var err error
  var publicKey *rsa.PublicKey

  if msg.KeyType != APP_RSA4096 && msg.KeyType != APP_RSA2048 {
    return "", "", 0, errors.New("unsupported key type")
  }

  // extract public key
  data, err = base64.StdEncoding.DecodeString(msg.PublicKey)
  if err != nil {
    return "", "", 0, err
  }
  publicKey, err = ParseRsaPublicKeyFromPemStr(string(data))
  if err != nil {
    return "", "", 0, err
  }

  // compute guid
  hash = sha256.Sum256(data)
  guid := hex.EncodeToString(hash[:])

  // extract signature
  data, err = base64.StdEncoding.DecodeString(msg.Signature)
  if err != nil {
    return "", "", 0, err
  }
  signature := data

  // verify hash
  data, err = base64.StdEncoding.DecodeString(msg.Message)
  if err != nil {
    return "", "", 0, err
  }
  hash = sha256.Sum256(data)
  if msg.SignatureType == APP_SIGNPKCS1V15 {
    err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
    if err != nil {
      return "", "", 0, err
    }
  } else if msg.SignatureType == APP_SIGNPSS {
    err = rsa.VerifyPSS(publicKey, crypto.SHA256, hash[:], signature, nil)
    if err != nil {
      return "", "", 0, err
    }
  } else {
    return "", "", 0, errors.New("unsupported signature type")
  }

  // extract data
  var signedData SignedData
  err = json.Unmarshal(data, &signedData);
  if err != nil {
    return "", "", 0, err
  }

  // validate signer
  if signedData.Guid != guid {
    return "", "", 0, errors.New("invalid message source")
  }

  // extract data
  err = json.Unmarshal([]byte(signedData.Value), obj)
  if err != nil {
    return "", "", 0, err
  }

  return guid, signedData.MessageType, signedData.Timestamp, nil
}

func WriteDataMessage(privateKey string, publicKey string, keyType string,
    signType string, guid string, messageType string, obj interface{}) (*DataMessage, error) {

  var data []byte
  var hash [32]byte
  var err error
  var private *rsa.PrivateKey

  // create message to sign
  data, err = json.Marshal(obj)
  if err != nil {
    return nil, err
  }
  var signedData SignedData
  signedData.Guid = guid
  signedData.Timestamp = time.Now().Unix()
  signedData.MessageType = messageType
  signedData.Value = string(data)
  data, err = json.Marshal(&signedData)
  message := base64.StdEncoding.EncodeToString(data)

  if keyType != APP_RSA2048 && keyType != APP_RSA4096 {
    return nil, errors.New("unsupported key type")
  }

  // get private key
  private, err = ParseRsaPrivateKeyFromPemStr(privateKey)
  if err != nil {
    return nil, err
  }
  key := base64.StdEncoding.EncodeToString([]byte(publicKey))

  // compute signature
  hash = sha256.Sum256(data)
  if signType == APP_SIGNPKCS1V15 {
    data, err = rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA256, hash[:])
  } else if signType == APP_SIGNPSS {
    data, err = rsa.SignPSS(rand.Reader, private, crypto.SHA256, hash[:], nil)
  } else {
    return nil, errors.New("unsupported signature type")
  }
  signature := base64.StdEncoding.EncodeToString(data)

  dataMessage := DataMessage{
    PublicKey: key,
    KeyType: keyType,
    Signature: signature,
    SignatureType: signType,
    Message: message,
  }
  return &dataMessage, nil
}



