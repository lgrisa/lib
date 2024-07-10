package dbv

import (
	"github.com/lgrisa/lib/dynamo/db/dbdef"
)

var AccountTableDefinition = &dbdef.TableDefinition{
	TableName: "LoginAccountV2",
	HashKey:   "uid",
	NewEntity: func() interface{} {
		return &Account{}
	},
}

type Account struct {
	Uid       string `dynamo:"uid,hash"  json:"uid,omitempty"`                   //主键
	CreatedAt int64  `dynamo:"created_at,omitempty" json:"created_at,omitempty"` //建立时间

	Sns map[string]string `dynamo:"sns" json:"sns"` //sns账号

	//AppleUserId     string `dynamo:"apple_user_id,omitempty" json:"apple_user_id"`
	//TapTapUserId    string `dynamo:"tap_tap_user_id,omitempty" json:"tap_tap_user_id"`
	//WxUserId        string `dynamo:"wx_user_id,omitempty" json:"wx_user_id"`
	//QqUserId        string `dynamo:"qq_user_id,omitempty" json:"qq_user_id"`
	//ThirdUserId     string `dynamo:"third_user_id,omitempty" json:"third_user_id"`
	//YxAccountId     string `dynamo:"yx_account_id,omitempty" json:"yx_account_id"`         //英雄账号ID
	//MobileAccountId string `dynamo:"mobile_account_id,omitempty" json:"mobile_account_id"` //手机号账号ID

	Mobile    string `dynamo:"mobile,omitempty" json:"mobile,omitempty"`         //手机号
	ChannelId int64  `dynamo:"channel_id,omitempty" json:"channel_id,omitempty"` // 渠道ID

	IdNum  string `dynamo:"id_num,omitempty" json:"id_num,omitempty"`   // 身份证号
	IdName string `dynamo:"id_name,omitempty" json:"id_name,omitempty"` // 身份证姓名

	AI           string      `dynamo:"ai,omitempty" json:"ai,omitempty"` // 用户实名认证唯一标识
	WlcStatus    CheckStatus `dynamo:"wlc_status,omitempty" json:"wlc_status,omitempty"`
	PI           string      `dynamo:"pi,omitempty" json:"pi,omitempty"`                       // 已通过实名认证用户的唯一标识
	UserBirthday int64       `dynamo:"user_birthday,omitempty" json:"user_birthday,omitempty"` // 用户生日

	GameLibraryMap       map[string]*GameInfo `dynamo:"game_library_map" json:"game_library_map"`                         // 游戏库 map[Game_name]GameInfo
	GameLibraryDeleteMap map[string]*GameInfo `dynamo:"game_library_delete_map" json:"game_library_delete_map,omitempty"` // 游戏库 map[Game_name]GameInfo

	EMGenTime int64 `dynamo:"em_gen_time,omitempty" json:"em_gen_time,omitempty"` // EM生成时间

	WlcErrCount     int64 `dynamo:"wlc_err_count,omitempty" json:"wlc_err_count,omitempty"`           // 实名认证错误次数
	WlcErrStartTime int64 `dynamo:"wlc_err_start_time,omitempty" json:"wlc_err_start_time,omitempty"` // 实名认证错误开始时间

	BindNewPhoneErrCount int64 `dynamo:"bind_new_phone_err_count,omitempty" json:"bind_new_phone_err_count,omitempty"` // 绑定新手机错误次数
	BindNewPhoneErrTime  int64 `dynamo:"bind_new_phone_err_time,omitempty" json:"bind_new_phone_err_time,omitempty"`   // 绑定新手机错误开始时间
}

func (d *Account) GetOrCreateSnsMap() map[string]string {
	if d.Sns == nil {
		d.Sns = make(map[string]string)
	}
	return d.Sns
}

func (d *Account) GetOrCreateGameLibraryMap() map[string]*GameInfo {
	if d.GameLibraryMap == nil {
		d.GameLibraryMap = make(map[string]*GameInfo)
	}
	return d.GameLibraryMap
}

func (d *Account) GetOrCreateGameLibraryDeleteMap() map[string]*GameInfo {
	if d.GameLibraryDeleteMap == nil {
		d.GameLibraryDeleteMap = make(map[string]*GameInfo)
	}
	return d.GameLibraryDeleteMap
}

type GameInfo struct {
	GameUserId  string       `dynamo:"game_user_id,omitempty" json:"game_user_id,omitempty"`   // 游戏用户ID
	Status      DeleteStatus `dynamo:"status,omitempty" json:"status,omitempty"`               // 账号状态	0:正常 1:冻结 2:删除
	FreezeTime  int64        `dynamo:"freeze_time,omitempty" json:"freeze_time,omitempty"`     // 冻结时间
	LastLoginAt int64        `dynamo:"last_login_at,omitempty" json:"last_login_at,omitempty"` // 最后登录时间
}

type DeleteStatus int64

const (
	Normal DeleteStatus = 0
	Freeze DeleteStatus = 1
	Delete DeleteStatus = 2
)

type CheckStatus int64

const (
	CheckStatusDefault CheckStatus = 0 // 未实名
	CheckStatusSuccess CheckStatus = 1 // 认证成功
	CheckStatusProcess CheckStatus = 2 // 认证中
	CheckStatusFailed  CheckStatus = 3 // 认证失败
)

const IOS = int64(56)
const Android = int64(18)
