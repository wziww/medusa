package stream

import (
	"encoding/json"
	"github/wziww/medusa/config"
	"github/wziww/medusa/log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

func init() {
	go func() {
		log.FMTLog(log.LOGINFO, "api server enable:", config.C.Base.API.Enable)
		if !config.C.Base.API.Enable {
			return
		}
		server := &http.Server{
			Addr: func() string {
				addr := "0.0.0.0:"
				if config.C.Base.API.Port == 0 {
					addr += "8083"
				} else {
					addr += strconv.Itoa(config.C.Base.API.Port)
				}
				return addr
			}(),
			WriteTimeout: 5 * time.Second,
			ReadTimeout:  5 * time.Second,
		}
		http.HandleFunc("/v1/date", func(w http.ResponseWriter, r *http.Request) {
			type res struct {
				FlowIn  uint64
				FlowOut uint64
			}
			d := &res{}
			d.FlowIn = atomic.LoadUint64(FlowIn)
			d.FlowOut = atomic.LoadUint64(FlowOut)
			w.WriteHeader(200)
			j, _ := json.Marshal(d)
			w.Write(j)
			return
		})
		err := server.ListenAndServe()
		if err != nil {
			log.FMTLog(log.LOGERROR, err)
			os.Exit(1)
		}
	}()
}
