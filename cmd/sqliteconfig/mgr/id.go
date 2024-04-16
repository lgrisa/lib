package mgr

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"sync"
)

type MessageIdGen struct {
	mux sync.Mutex

	idMap     map[string]int
	dataExist bool
}

func loadGen(idGenPath string) (*MessageIdGen, error) {
	idData, err := ioutil.ReadFile(idGenPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, "MessageIdGen read yaml fail")
	}

	return newGen(idData)
}

func newGen(data []byte) (*MessageIdGen, error) {
	g := &MessageIdGen{
		idMap: map[string]int{},
	}

	if len(data) > 0 {
		err := yaml.Unmarshal(data, g.idMap)
		if err != nil {
			return nil, errors.Wrap(err, "yaml.unmarshal fail")
		}

		g.dataExist = true
	}

	return g, nil
}

func (g *MessageIdGen) encode() ([]byte, error) {
	return yaml.Marshal(g.idMap)
}

func (g *MessageIdGen) newId(key string) int {

	id, ok := g.idMap[key]
	if !ok {
		g.idMap[key] = 1
		return 1
	}

	g.idMap[key] = id + 1

	return id + 1
}

func (g *MessageIdGen) get(key string) (int, bool) {
	id, ok := g.idMap[key]
	return id, ok
}

func (g *MessageIdGen) getOrCreate(key, prefix string) int {

	id, ok := g.idMap[key]
	if !ok {
		id = g.newId(prefix)
		g.idMap[key] = id
		return id
	}

	return id
}

//func (g *MessageIdGen) ModuleId(moduleName string) int {
//	return g.getOrCreate(ModuleIdKey(moduleName))
//}
//
//func (g *MessageIdGen) MsgId(moduleName, msgName, msgType string) int {
//	return g.getOrCreate(MsgIdKey(moduleName, msgName, msgType))
//}
//
//func (g *MessageIdGen) MsgFailCodeId(moduleName, msgName, code string) int {
//	return g.getOrCreate(MsgFailCodeIdKey(moduleName, msgName, code))
//}

func (g *MessageIdGen) MsgProtoFieldId(moduleName, msgName, fieldName, fieldType string) int {
	g.mux.Lock()
	defer g.mux.Unlock()

	key, prefix := MsgProtoFieldIdKey(moduleName, msgName, fieldName, fieldType)
	return g.getOrCreate(key, prefix)
}

//var idMap map[string]int
//
//func init() {
//	idMap = map[string]int{}
//}
//
//func ModuleId(moduleName string) int {
//	return getOrCreate(ModuleIdKey(moduleName))
//}
//
//func MsgId(moduleName, msgName, msgType string) int {
//	return getOrCreate(MsgIdKey(moduleName, msgName, msgType))
//}
//
//func MsgFailCodeId(moduleName, msgName, code string) int {
//	return getOrCreate(MsgFailCodeIdKey(moduleName, msgName, code))
//}
//
//func MsgProtoFieldId(moduleName, msgName, msgType, fieldName string) int {
//	return getOrCreate(MsgProtoFieldIdKey(moduleName, msgName, msgType, fieldName))
//}
//
//func CommonProtoFieldId(protoType, fieldName string) int {
//	return getOrCreate(CommonProtoFieldIdKey(protoType, fieldName))
//}
//
//func newId(key string) int {
//	id, ok := idMap[key]
//	if !ok {
//		idMap[key] = 1
//		return 1
//	}
//
//	idMap[key] = id + 1
//
//	return id + 1
//}
//
//func getOrCreate(key, newKey string) int {
//	id, ok := idMap[key]
//	if !ok {
//		id = newId(newKey)
//		idMap[key] = id
//		return id
//	}
//
//	return id
//}
