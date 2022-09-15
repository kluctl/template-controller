package controllers

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/kluctl/go-jinja2"
)

func NewJinja2(opts ...jinja2.Jinja2Opt) (*jinja2.Jinja2, error) {
	var opts2 []jinja2.Jinja2Opt
	opts2 = append(opts2, opts...)
	opts2 = append(opts2,
		jinja2.WithStrict(false),
		jinja2.WithExtension("jinja2.ext.loopcontrols"),
		jinja2.WithExtension("go_jinja2.ext.kluctl"),
	)
	return jinja2.NewJinja2("template-controller", 1, opts2...)
}

func MergeMap(a, b map[string]interface{}) {
	MergeMap2(a, b, false)
}

func MergeMap2(a, b map[string]interface{}, skipNil bool) {
	for key := range b {
		if _, ok := a[key]; ok {
			adict, adictOk := a[key].(map[string]interface{})
			bdict, bdictOk := b[key].(map[string]interface{})
			if adictOk && bdictOk {
				MergeMap2(adict, bdict, skipNil)
			} else {
				if !skipNil || b[key] != nil {
					a[key] = b[key]
				}
			}
		} else {
			if !skipNil || b[key] != nil {
				a[key] = b[key]
			}
		}
	}
}

func Sha256String(data string) string {
	return Sha256Bytes([]byte(data))
}

func Sha256Bytes(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}
