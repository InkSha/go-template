package template

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Engine 模板引擎
type Engine struct {
	templateDir string
	templates   map[string]*template.Template
	mutex       sync.RWMutex
	hotReload   bool // 是否热重载（开发模式）
}

// Config 模板配置
type Config struct {
	TemplateDir string // 模板目录
	HotReload   bool   // 热重载开关
}

// New 创建模板引擎
func New(cfg Config) (*Engine, error) {
	e := &Engine{
		templateDir: cfg.TemplateDir,
		templates:   make(map[string]*template.Template),
		hotReload:   cfg.HotReload,
	}

	if err := e.loadTemplates(); err != nil {
		return nil, err
	}

	return e, nil
}

// loadTemplates 加载所有模板
func (e *Engine) loadTemplates() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	baseFile := filepath.Join(e.templateDir, "base.html")

	// 查找所有非 base.html 的模板文件
	files, err := filepath.Glob(filepath.Join(e.templateDir, "*.html"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Base(file) == "base.html" {
			continue
		}

		name := filepath.Base(file)
		name = name[:len(name)-5] // 去掉 .html

		tmpl, err := template.ParseFiles(baseFile, file)
		if err != nil {
			return err
		}

		e.templates[name] = tmpl
	}

	return nil
}

// Render 渲染模板
func (e *Engine) Render(w io.Writer, name string, data interface{}) error {
	// 开发模式下热重载
	if e.hotReload {
		if err := e.loadTemplates(); err != nil {
			return err
		}
	}

	e.mutex.RLock()
	tmpl, ok := e.templates[name]
	e.mutex.RUnlock()

	if !ok {
		return os.ErrNotExist
	}

	return tmpl.Execute(w, data)
}

// Reload 手动重新加载模板
func (e *Engine) Reload() error {
	return e.loadTemplates()
}

// SetHotReload 设置热重载
func (e *Engine) SetHotReload(enabled bool) {
	e.hotReload = enabled
}

// AddTemplate 动态添加模板
func (e *Engine) AddTemplate(name, content string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	baseFile := filepath.Join(e.templateDir, "base.html")
	baseTmpl, err := template.ParseFiles(baseFile)
	if err != nil {
		return err
	}

	tmpl, err := baseTmpl.Parse(content)
	if err != nil {
		return err
	}

	e.templates[name] = tmpl
	return nil
}

// RemoveTemplate 移除模板
func (e *Engine) RemoveTemplate(name string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	delete(e.templates, name)
}

// ListTemplates 列出所有模板
func (e *Engine) ListTemplates() []string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	names := make([]string, 0, len(e.templates))
	for name := range e.templates {
		names = append(names, name)
	}
	return names
}
