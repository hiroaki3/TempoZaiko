package TempoZaiko0x0

import (
	"net/url"
	"strings"
)

// 今回のサンプルに特化した構造体。
// URLをパースし、必要な値を保持。
// ※handler系から呼ばれる事を想定
// [使い方]
// p := NewParam(r.URL)
type Param struct {
	Kind  string
	Key   string
	Value string
}

// URLのGET値をベースにParamを作成。
// 1番目はindex or api、2番目は種別、3番目はキー、4番目は値、5番目以降は破棄。
func NewParam(url *url.URL) *Param {
	path := strings.Trim(url.Path, "/")
	s := strings.Split(path, "/")
	param := new(Param)
	if len(s) >= 2 {
		param.Kind = s[1]
	}
	if len(s) >= 3 {
		param.Key = s[2]
	}
	if len(s) >= 4 {
		param.Value = s[3]
	}
	return param
}

func (p *Param) HasKey() bool {
	return len(p.Key) > 0
}

func (p *Param) HasValue() bool {
	return len(p.Value) > 0
}
