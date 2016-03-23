package photon

import (
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"github.com/satori/go.uuid"
	"google.golang.org/cloud/storage"
)

func main() {
	r := mux.NewRouter()

	r.Methods("GET").Path("/users/{id}").Handler(appHandler(uploadImage))

	http.Handle("/", handlers.CombinedLoggingHandler(os.Stderr, r))

	log.Fatal(http.ListenAndServe(":4040", nil))
}

func uploadImage(w http.ResponseWriter, r *http.Request) *appError {
	f, fh, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return appErrorf(err, "no image found in request: %v", err)
	}
	if err != nil {
		return appErrorf(err, "could not add image: %v", err)
	}
	name := uuid.NewV4().String() + path.Ext(fh.Filename)
	return nil
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

type appError struct {
	Error   error
	Message string
	Code    int
}

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil {
		log.Printf("Handler error: status code: %d, message: %s, underlying error: %#v",
			e.Code, e.Message, e.Error)

		http.Error(w, e.Message, e.Code)
	}
}

func appErrorf(err error, format string, v ...interface{}) *appError {
	return &appError{
		Error:   err,
		Message: fmt.Sprintf(format, v...),
		Code:    500,
	}
}
