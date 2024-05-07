package playername

import (
	"encoding/hex"
	"github.com/lgrisa/lib/utils/math/i64"
	"math/rand"
)

var PlayerNamePrefix = "Player_"

func PlayerName(id int64) string {
	//sid := GetSid(id)
	//accountId := GetAccountId(id)
	return PlayerNamePrefix + hex.EncodeToString(i64.ToBytes(id))
}

func RandomPlayerName(id int64) string {
	//sid := GetSid(id)
	//accountId := GetAccountId(id)
	r := rand.Int63() % 65535
	return PlayerNamePrefix + hex.EncodeToString(i64.ToBytes(id)) + "_" + hex.EncodeToString(i64.ToBytes(r))
}
