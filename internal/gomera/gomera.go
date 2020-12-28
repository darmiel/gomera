package gomera

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func New(opt *Options) {
	// print banner
	fmt.Println(banner)

	// print settings
	optjs, _ := json.Marshal(opt)
	fmt.Println()
	fmt.Println("👉 with:", string(optjs))
	fmt.Println()

	// try to parse cameras
	if _, err := parseCameras(); err != nil {
		log.Fatalln("Error parsing cameras:", err)
		return
	}
	log.Println("📷 Parsed", len(Cameras), "cameras.")

	// http server
	if err := http.ListenAndServe(":"+strconv.Itoa(int(opt.Port)), createRouter()); err != nil {
		log.Println("Error listening and serving:")
	}
}
