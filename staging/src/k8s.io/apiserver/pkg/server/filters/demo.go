package filters

import (
	"fmt"
	"io"
	"net/http"

	"k8s.io/klog"

	apirequest "k8s.io/apiserver/pkg/endpoints/request"
)

func WrapDemoHandler(handler http.Handler, str string) http.Handler {
	return &filterDemoHandler{handler: handler, str: str}
}

type filterDemoHandler struct {
	str     string
	handler http.Handler
}

func (h *filterDemoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestInfo, ok := apirequest.RequestInfoFrom(ctx)
	if !ok {
		klog.V(6).Infof("failed to get requestInfo, err: %+v", r)
		http.Error(w, fmt.Sprintf("failed to get requestInfo, err: %v", r), http.StatusBadRequest)
		return
	}

	klog.V(7).Infof("get requestInfo: %+v", requestInfo)
	io.WriteString(w, h.str)
	h.handler.ServeHTTP(w, r)
}
