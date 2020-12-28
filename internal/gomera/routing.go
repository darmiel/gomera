package gomera

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func createRouter() (res *mux.Router) {
	router := mux.NewRouter()

	router.HandleFunc("/notify/{cameraId}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		var cameraId = vars["cameraId"]
		log.Println("Requesting camera with id:", cameraId)

		// get camera
		camera, ok := Cameras[cameraId]
		if !ok {
			// camera not found
			_, _ = fmt.Fprintf(w, "Camera not found.")
			return
		}

		log.Println("Sending webhook ...")
		if resp, err := camera.sendWebhook(); err != nil {
			log.Println("Error sending webhook:", err)
		} else {
			log.Println("Probably no errors while sending", resp)
			fmt.Println()
			fmt.Println(camera.createJsonPayload())
		}
	})

	return router
}
