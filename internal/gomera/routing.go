package gomera

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func createRouter() (res *mux.Router) {
	router := mux.NewRouter()

	router.HandleFunc("/notify/{cameraId}/{secret}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cameraId = vars["cameraId"]
		var secret = vars["secret"]
		log.Println("👉", "Requesting camera with id:", cameraId, "and secret:", secret)

		camera, err := FindCamera(cameraId, secret)
		if err != nil {
			_, _ = fmt.Fprintf(w, err.Error())
			return
		}

		if _, err := camera.Send(); err != nil {
			log.Println("Error sending webhook:", err)
		}
	})

	return router
}
