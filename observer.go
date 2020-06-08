package engine

import (
	"github.com/appconf/engine/memkv"
	"github.com/appconf/engine/render"
	"github.com/appconf/storage"
)

type observer struct {
	watch   map[string][]*render.Processor
	backend storage.Storage
	memkv   *memkv.MemKv
}

func newObserver(backend storage.Storage) *observer {
	return &observer{
		watch:   make(map[string][]*render.Processor),
		backend: backend,
		memkv:   memkv.GetKv(),
	}
}

//Add 关联key与模板处理器
func (obs *observer) Add(key string, processor *render.Processor) {
	processors, ok := obs.watch[key]
	if !ok {
		processors = make([]*render.Processor, 0)
	}
	obs.watch[key] = append(processors, processor)
}

func removeRepeatProcessor(processors []*render.Processor) []*render.Processor {
	result := make([]*render.Processor, 0)
	for _, processor := range processors {
		if !processorInSlice(processor, result) {
			result = append(result, processor)
		}
	}
	return result
}

func processorInSlice(processor *render.Processor, processors []*render.Processor) bool {
	for _, v := range processors {
		if v == processor {
			return true
		}
	}
	return false
}

//Notify 当有keys发生更新时，通知这些keys对应的模板处理进行重新渲染
func (obs *observer) Notify(keys []string) {
	//获取需要渲染的模板的处理器
	processors := make([]*render.Processor, 0)
	for _, key := range keys {
		value, ok := obs.watch[key]
		if !ok {
			continue
		}
		processors = append(processors, value...)
	}
	processors = removeRepeatProcessor(processors)

	for _, processor := range processors {
		//获取当前处理器渲染的当前变更的所有Key的value
		kvs := make([]map[string]interface{}, 0)
		ks := processor.GetInKeys(keys)
		for _, v := range ks {
			pair, err := obs.memkv.Get(v)
			if err != nil {
				continue
			}
			kvs = append(kvs, map[string]interface{}{v: pair.Value})
		}

		//获取当前渲染模板的路径
		template := processor.GetTemplate()

		//渲染模板，并将结果通知给Storage
		if err := processor.Rendering(); err != nil {
			obs.backend.Error(template, kvs, err)
		} else {
			obs.backend.Success(template, kvs)
		}
	}
}
