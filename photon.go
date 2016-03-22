package photon

import (
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"google.golang.org/cloud/storage"
)

func main() {
	r := mux.NewRouter()

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))
}
