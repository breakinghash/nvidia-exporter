package watch

import (
	"nvidia-exporter/nvml"
	"time"

	log "github.com/sirupsen/logrus"
)

// Watch starts status watching
func Watch(prefix string) {
	err := nvml.Init()

	check(err)

	deviceCount, err := nvml.GetDeviceCount()
	check(err)

	metrics := Metrics{}.Init(prefix, deviceCount)

	go metrics.ListenAndServe()

	for {
		for i := uint(0); i < deviceCount; i++ {
			device, err := nvml.NewDevice(i)
			check(err)

			status, err := device.Status()
			check(err)

			metrics.temp[i].Update(float64(*status.Temperature))
			metrics.fan[i].Update(float64(*status.Fan))

			log.WithFields(
				log.Fields{"GPU": i},
			).Infof(
				"t=%dC, fan=%d%%",
				*status.Temperature,
				*status.Fan,
			)
		}

		time.Sleep(30 * time.Second)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
