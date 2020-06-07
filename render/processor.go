package render

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/appconf/engine/command"
	"github.com/appconf/engine/util"
	"github.com/appconf/log"
)

//Processor 模板实际处理者
type Processor struct {
	tpl *template.Template
	cfg *Config
}

//NewProcessor 创建Processor实例
func NewProcessor(cfg *Config) (*Processor, error) {
	var (
		tpl *template.Template
		err error
	)
	tpl = template.New(filepath.Base(cfg.Src))
	tpl = tpl.Funcs(funcMap)
	if tpl, err = tpl.ParseFiles(cfg.Src); err != nil {
		return nil, err
	}
	return &Processor{
		tpl: tpl,
		cfg: cfg,
	}, nil
}

//GetInKeys 给定一些key, 判断哪些key在process, 并返回
func (p *Processor) GetInKeys(keys []string) []string {
	result := make([]string, 0)
	for _, v := range keys {
		if util.ElementInStrSlice(v, p.cfg.Keys) {
			result = append(result, v)
		}
	}
	return result
}

//GetTemplate 获取渲染模板路径
func (p *Processor) GetTemplate() string {
	return p.cfg.Src
}

//Rendering 渲染模板
func (p *Processor) Rendering() (err error) {
	var f *os.File
	f, err = os.OpenFile(p.cfg.Dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Debugf("rendering template: %s", p.cfg.Src)

	err = p.tpl.Execute(f, nil)
	if err != nil {
		return err
	}

	return p.reload()
}

//reload 执行reload_cmd
func (p *Processor) reload() (err error) {
	if len(p.cfg.ReloadCmd) == 0 {
		return nil
	}

	log.Debugf("execute reload cmd: %s", p.cfg.ReloadCmd)

	var output []byte
	output, err = command.NewBash(p.cfg.ReloadCmd).Run()
	if err != nil {
		return fmt.Errorf("failed to execute reload cmd, error: %s", string(output))
	}
	return nil
}
