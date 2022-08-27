package gm

import (
	"github.com/CuteReimu/FengSheng/config"
	"github.com/CuteReimu/FengSheng/utils"
	"net/http"
	"net/url"
)

var logger = utils.GetLogger("gm")

func Init() {
	if config.IsGmEnable() && len(handlers) > 0 {
		go func() {
			for name, handler := range handlers {
				http.HandleFunc("/"+name, func(writer http.ResponseWriter, request *http.Request) {
					if request.Method != "GET" {
						writer.WriteHeader(http.StatusMethodNotAllowed)
						_, _ = writer.Write([]byte(`{"error": "invalid method"}`))
						return
					}
					if request.ParseForm() != nil {
						writer.WriteHeader(http.StatusBadRequest)
						_, _ = writer.Write([]byte(`{"error": "parse form failed"}`))
						return
					}
					_, _ = writer.Write(handler(request.Form))
				})
			}
			_ = http.ListenAndServe(config.GetGmListenAddress(), nil)
		}()
	}
}

var handlers = make(map[string]func(values url.Values) []byte)
