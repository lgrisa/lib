package mgr

import "strings"

func ModuleStructName(moduleName string) string {
	return strings.Title(moduleName) + "Module"
}

func MsgKeyName(msgName, msgType string) string {
	// c2s_visitor_login
	return msgType + "_" + msgName
}

func MsgProtoName(msgName, msgType string) string {
	// C2sVisitorLoginProto，proto生成的代码好傻逼，所以做了以下处理，把C2S/S2C变成大写
	s := HumpName(MsgKeyName(msgName, msgType)) + "Proto"

	if strings.HasPrefix(msgType, "c2s") {
		return strings.Replace(s, "C2s", "C2S", 1)
	} else if strings.HasPrefix(msgType, "s2c") {
		return strings.Replace(s, "S2c", "S2C", 1)
	}

	return s
}

func MsgProcessFuncName(moduleName, msgName, msgType string) string {
	// processC2sVisitorLoginMsg
	return "process" + HumpName(MsgKeyName(msgName, msgType)) + "Msg"
}

func MsgBuildFuncName(msgName, msgType string) string {
	// NewS2cVisitorLoginMsg
	return "New" + HumpName(MsgKeyName(msgName, msgType)) + "Msg"
}

func ProtoMsgBuildFuncName(msgName, msgType string) string {
	// NewS2cVisitorLoginProtoMsg
	return "New" + HumpName(MsgKeyName(msgName, msgType)) + "ProtoMsg"
}

func MarshalMsgBuildFuncName(msgName, msgType string) string {
	// NewS2cVisitorLoginMarshalMsg
	return "New" + HumpName(MsgKeyName(msgName, msgType)) + "MarshalMsg"
}

func MsgFailCodeCacheName(msgName, codeName string) string {
	// ERR_VISITOR_LOGIN_FAIL_INVALID_UID
	return strings.ToUpper("err_" + msgName + "_fail_" + codeName)
}

func MsgFailCodeCacheField(msgName, codeName string) string {
	// ErrVisitorLoginFailInvalidUid
	return HumpName("err_" + msgName + "_fail_" + codeName)
}

func HeaderMsgName(msgName, msgType string) string {
	// VISITOR_LOGIN_S2C_SUCCESS
	return strings.ToUpper(msgName + "_" + msgType)
}

func RpcMsgProcessFuncName(moduleName, msgName, msgType string) string {
	// VisitorLogin
	return HumpName(msgName)
}

func RpcFuncName(msgName, msgType string) string {
	// ErrVisitorLoginFailInvalidUid
	return "New" + HumpName(MsgKeyName(msgName, msgType)) + "MarshalMsg"
}

func HumpName(in string) string {
	return strings.Replace(strings.Title(strings.Replace(in, "_", " ", -1)), " ", "", -1)
}

const KeySep = "#"

//const ModuleIdPrefix = "module" + KeySep
//
//func ModuleIdKey(moduleName string) (string, string) {
//	// module@battle
//	return ModuleIdPrefix + moduleName, ModuleIdPrefix
//}
//
//func MsgIdPrefix() string {
//	// MsgId#
//	return "MsgId" + KeySep
//}
//
//func MsgIdKey(moduleName, msgName, msgType string) (string, string) {
//	// MsgId#list_battle_replay#c2s
//	prefix := MsgIdPrefix()
//	return prefix + msgName + KeySep + msgType, prefix
//}
//
//func MsgFailCodeIdPrefix(moduleName, msgName string) string {
//	// MsgFailCodeId#list_battle_replay#
//	return "MsgFailCodeId" + KeySep + msgName + KeySep
//}
//
//func MsgFailCodeIdKey(moduleName, msgName, code string) (string, string) {
//	// MsgFailCodeId#list_battle_replay#invalid_uid
//	prefix := MsgFailCodeIdPrefix(moduleName, msgName)
//	return prefix + code, prefix
//}

func MsgProtoFieldIdPrefix(moduleName, msgName string) string {
	// MsgProtoFieldId#list_battle_replay#c2s#
	return "MsgProtoFieldId" + KeySep + moduleName + KeySep + msgName + KeySep
}

func MsgProtoFieldIdKey(moduleName, msgName, fieldName, fieldType string) (string, string) {
	// MsgProtoFieldId#list_battle_replay#uid#string
	prefix := MsgProtoFieldIdPrefix(moduleName, msgName)
	return prefix + fieldName + KeySep + fieldType, prefix
}

//func CommonProtoFieldIdPrefix(protoType string) string {
//	// GoodsProto^
//	return protoType + "^"
//}
//
//func CommonProtoFieldIdKey(protoType, fieldName string) (string, string) {
//	// GoodsProto^id
//	prefix := CommonProtoFieldIdPrefix(protoType)
//	return prefix + fieldName, prefix
//}
