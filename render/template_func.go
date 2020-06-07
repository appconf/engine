package render

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/appconf/engine/memkv"
	"github.com/appconf/engine/util"
)

//newFuncMap 模板渲染函数
func newFuncMap() map[string]interface{} {
	fnMap := memkv.GetKv().FuncMap()
	fnMap["hasprefix"] = strings.HasPrefix
	fnMap["hassuffix"] = strings.HasSuffix
	fnMap["join"] = strings.Join
	fnMap["contains"] = strings.Contains
	fnMap["trim"] = strings.Trim
	fnMap["trimprefix"] = strings.TrimPrefix
	fnMap["trimleft"] = strings.TrimLeft
	fnMap["trimright"] = strings.TrimRight
	fnMap["trimspace"] = strings.TrimSpace
	fnMap["trimsuffix"] = strings.TrimSuffix
	fnMap["title"] = strings.Title
	fnMap["index"] = strings.Index
	fnMap["replace"] = strings.Replace
	fnMap["replaceall"] = strings.ReplaceAll
	fnMap["equalfold"] = strings.EqualFold
	fnMap["count"] = strings.Count
	fnMap["repeat"] = strings.Repeat
	fnMap["splitn"] = strings.SplitN
	fnMap["split"] = strings.Split
	fnMap["tolower"] = strings.ToLower
	fnMap["toupper"] = strings.ToUpper
	fnMap["json"] = jsonUnmarshal
	fnMap["strftime"] = strftime
	fnMap["getenv"] = os.Getenv
	return fnMap
}

func jsonUnmarshal(s string) (v interface{}, err error) {
	err = json.Unmarshal([]byte(s), &v)
	return
}

func strftime(value, layout, fmt string) (string, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return "", err
	}
	return t.Format(fmt), nil
}

var (
	funcMap  = newFuncMap()
	funcLock = new(sync.RWMutex)
)

//Register 注册模板渲染函数
func Register(name string, fn interface{}) {
	funcLock.Lock()
	defer funcLock.Unlock()

	if !util.IsFn(fn) {
		panic("template:" + name + "is not a function type")
	}

	_, dup := funcMap[name]
	if dup {
		panic("template: register called twice for " + name)
	}

	funcMap[name] = fn
}
