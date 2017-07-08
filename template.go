package TempoZaiko0x0

import (
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// template用の構造体の定義
type templateHandler struct {
	once     sync.Once
	filename string
	tpl      *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//templateをコンパイルするのは1回だけで良いのでonceを利用。
	t.once.Do(func() {
		t.tpl = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.tpl.Execute(w, nil)
}
