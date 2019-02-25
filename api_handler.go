package jam

import "net/http"

type ApiHandler interface {
	HandleIndex(w http.ResponseWriter, r *http.Request)
	HandleS3Upload(w http.ResponseWriter, r *http.Request)
}
