package opennox

import (
	"encoding/json"
	"net/http"

	"github.com/brutusmcforce/opennox/client"
)

var _ json.Marshaler = &client.Drawable{}

func init() {
	http.HandleFunc("/debug/nox/drawable", func(w http.ResponseWriter, r *http.Request) {
		writeJSONResp(w, noxClient.Objs.AllList1())
	})
}
