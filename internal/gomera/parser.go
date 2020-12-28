package gomera

import "flag"

type Options struct {
	Port       uint   `json:"Port"`
	CameraFile string `json:"camera_file"`
}

func Parse() *Options {
	opt = &Options{}

	// flags
	flag.UintVar(&opt.Port, "p", 9999, "Port of api")
	flag.StringVar(&opt.CameraFile, "f", "cameras.json", "Cameras json file")

	// parse
	flag.Parse()

	return opt
}
