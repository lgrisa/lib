package orderid

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"github.com/tinylib/msgp/msgp"
	"math/rand"
	"strings"
)

const (
	VersionBitCount = 24

	V1Bit    = 1 << VersionBitCount
	RandMask = (1 << VersionBitCount) - 1
)

func NewCpOrderId(heroId int64, chargeId uint64, moneyFen uint64, ctime int64) string {
	// version(1) + random(3)
	randNum := rand.Uint32()&RandMask | V1Bit

	// 5  + 5 + 5 + 5 + 5
	b := make([]byte, 0, 25)
	b = msgp.AppendUint32(b, randNum)
	b = msgp.AppendInt64(b, heroId)
	b = msgp.AppendUint64(b, chargeId)
	b = msgp.AppendUint64(b, moneyFen)
	b = msgp.AppendInt64(b, ctime)

	return base64.RawURLEncoding.EncodeToString(b)
}

func NewCpOrderSign(cpOrderId, key string) string {
	sum := md5.Sum([]byte(cpOrderId + "_" + key))
	return hex.EncodeToString(sum[:])
}

func NewCpOrderExtend(sign, namespace string) string {
	return sign + "_" + namespace
}

func ParsingExtend(extend string) (string, string) {
	sign, namespace, _ := strings.Cut(extend, "_")
	return sign, namespace
}

func ParseCpOrderId(orderId string) (version, randNum uint32, heroId int64, chargeId uint64, moneyFen uint64, ctime int64, err error) {
	b, err := base64.RawURLEncoding.DecodeString(orderId)
	if err != nil {
		return
	}

	randNum, b, err = msgp.ReadUint32Bytes(b)
	if err != nil {
		return
	}

	heroId, b, err = msgp.ReadInt64Bytes(b)
	if err != nil {
		return
	}

	chargeId, b, err = msgp.ReadUint64Bytes(b)
	if err != nil {
		return
	}

	moneyFen, b, err = msgp.ReadUint64Bytes(b)
	if err != nil {
		return
	}

	ctime, b, err = msgp.ReadInt64Bytes(b)
	if err != nil {
		return
	}

	return randNum >> VersionBitCount, randNum, heroId, chargeId, moneyFen, ctime, nil
}
