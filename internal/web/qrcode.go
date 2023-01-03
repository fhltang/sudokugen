package web

import (
	"fmt"
	"net/http"

	"github.com/skip2/go-qrcode"
)

func QrcodeHandler(w http.ResponseWriter, r *http.Request) {
	query, ok := r.URL.Query()["q"]
	if !ok || len(query) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	png, err := qrcode.Encode(fmt.Sprintf("%s/?q=%s", r.Host, query[0]), qrcode.Medium, 256)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}


	w.Write(png)
}
