package badger

import (
  "fmt"
  "time"
  "encoding/json"
	b64 "encoding/base64"
  "strings"
  "errors"
  "log/slog"
)

const (
  VoucherVersionId = 1
)

type VoucherHeader struct {
  Version int
}

type VoucherBody struct {
  Issuer string
  Issued time.Time
  Expiration time.Time
}

type VoucherInfo struct {
  Hash  []byte
  Sig   []byte
}

func NewVoucher(badge Badge, expiration time.Duration) (string, error) {

  headerJson, _ := json.Marshal(VoucherHeader{
    Version: VoucherVersionId,
  })

  header := b64.StdEncoding.EncodeToString(headerJson)

  timeIssued := time.Now()
  timeExpires := timeIssued.Add(expiration)

  if timeExpires.Equal(timeIssued) ||
     timeExpires.Before(timeIssued) {
      return "", errors.New("invalid time duration")
  }

  bodyJson, _ := json.Marshal(VoucherBody{
    Issuer: badge.Id(),
    Issued: timeIssued,
    Expiration: timeExpires,
  })

  body := b64.StdEncoding.EncodeToString(bodyJson)

  result := fmt.Sprintf("%s:%s", string(header), string(body))

  phs, err := badge.Sign(&result)

  if err != nil {
    return "", err
  }

  infoJson, _ := json.Marshal(VoucherInfo{
    Hash: phs.Hash,
    Sig: phs.Sig,
  })

  info := b64.StdEncoding.EncodeToString(infoJson)

  result = fmt.Sprintf("%s:%s", result, string(info))
    
  return result, nil
}

func ValidateVoucher(publicKey string, voucher string) bool {

  slog.Debug("badger:ValidateVoucher")

  pieces := strings.Split(voucher, ":")

  if len(pieces) != 3 {
    slog.Warn("voucher of incorrect length")
    return false
  }

  keyDecoded, err := b64.StdEncoding.DecodeString(publicKey)
  if err != nil {
    slog.Warn("failed to b64 decode key")
    return false
  }
  keyActual := UnmarshalPublicKey([]byte(keyDecoded))

  headerJson, err := b64.StdEncoding.DecodeString(pieces[0])
  if err != nil {
    slog.Warn("failed to decode voucher header")
    return false
  }
  
  var header VoucherHeader
  if err := json.Unmarshal([]byte(headerJson), &header); err != nil {
    slog.Warn("failed to unmarshal voucher header")
    return false
  }

  if header.Version < VoucherVersionId {
    slog.Warn("voucher version mismatch",
      "current_version", VoucherVersionId, "voucher_version", header.Version) 
    return false
  }

  bodyJson, err := b64.StdEncoding.DecodeString(pieces[1])
  if err != nil {
    slog.Warn("failed to decode voucher body")
    return false
  }

  var body VoucherBody
  if err := json.Unmarshal([]byte(bodyJson), &body); err != nil {
    slog.Warn("failed to unmarshal voucher body")
    return false
  }

  infoJson, err := b64.StdEncoding.DecodeString(pieces[2])
  if err != nil {
    slog.Warn("failed to decode voucher info")
    return false
  }

  var info VoucherInfo
  if err := json.Unmarshal([]byte(infoJson), &info); err != nil {
    slog.Warn("failed to unmarshal voucher info")
    return false
  }
  
  evaluationTime := time.Now()

  // Just out-right invalid
  if body.Issued.After(body.Expiration) ||
     body.Issued.Equal(body.Expiration) {
    slog.Debug("invalid voucher issued/expiration times")
    return false
  }

  // Expired
  if body.Expiration.Before(evaluationTime) {
    slog.Debug("expired voucher")
    return false
  }

  // Verify signature of voucher against given pubkey
  return Verify(PubHashSig{
    PubKey: MarshalPublicKey(keyActual.X, keyActual.Y),
    Hash: info.Hash,
    Sig: info.Sig,
  })
}

