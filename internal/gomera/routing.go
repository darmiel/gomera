package gomera

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var (
	totalPings uint64 = 0
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
			_, _ = fmt.Fprintf(w, err.Error())
			log.Println("Error sending webhook:", err)
		} else {
			_, _ = fmt.Fprintf(w, "Success!")
		}
	})

	router.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		totalPings++
		log.Println("Client requested ping:", r.RemoteAddr, "total:", totalPings)

		if opt.DevEnvironment {
			_, _ = fmt.Fprintf(w, "Pong! "+strconv.Itoa(len(Cameras))+" camera/s loaded")
		} else {
			_, _ = fmt.Fprintf(w, "I am still here 🤗")
		}
		return
	})

	return router
}
