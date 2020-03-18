package http_t

import "net/textproto"

type Header map[string][]string

// textproto.MIMEHeader 里面封装了Add方法，不仅仅是将k/v添加进map
// 还需要判断是否有已存在的key
// 因为Header是map类型（引用传递），所以不需要func (h *Header) ...
func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}

func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}

func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}

func (h Header) get(key string) string {
	if v := h[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}

func (h Header) Del(key string) {
	textproto.MIMEHeader(h).Del(key)
}