package engine

import (
	"fmt"
	"strings"

	"github.com/appconf/engine/memkv"
	"github.com/appconf/engine/render"
	"github.com/appconf/engine/util"
	"github.com/appconf/log"
	"github.com/appconf/storage"
)

//Core 调度器核心
type Core struct {
	backend  storage.Storage
	memKv    *memkv.MemKv
	stopCh   chan bool
	keys     []string
	logger   log.Logger
	observer *observer
}

//New 创建Core实例
func New(logger log.Logger, backend storage.Storage, templates []*render.Config) (*Core, error) {
	obs := newObserver(backend)
	for _, template := range templates {
		processor, err := render.NewProcessor(template)
		if err != nil {
			return nil, err
		}
		for _, key := range template.Keys {
			obs.Add(key, processor)
		}
	}

	keys := render.GetKeys(templates)
	return &Core{
		backend:  backend,
		memKv:    memkv.GetKv(),
		stopCh:   make(chan bool),
		keys:     keys,
		logger:   logger,
		observer: obs,
	}, nil
}

//Run 启动调度引擎
func (c *Core) Run() error {
	ch, err := c.backend.Get(c.keys)
	if err != nil {
		return err
	}
	for {
		select {
		case data, ok := <-ch:
			if !ok {
				return fmt.Errorf("core: storage channel is closed")
			}
			c.process(data)
		case <-c.stopCh:
			return nil
		}
	}
}

func (c *Core) process(data []storage.Data) {
	keys := make([]string, 0)

	for _, v := range data {
		if c.memKv.Equal(v.Key, v.Value) {
			continue
		}

		c.memKv.Set(v.Key, v.Value)

		if !util.ElementInStrSlice(v.Key, keys) {
			keys = append(keys, v.Key)
		}
	}

	if len(keys) != 0 {
		c.logger.Debugf("keys (%s) update, notify template update", strings.Join(keys, ","))
		c.observer.Notify(keys)
	}
}

//Stop 停止调度
func (c *Core) Stop() error {
	c.stopCh <- true
	close(c.stopCh)
	return c.backend.Stop()
}
