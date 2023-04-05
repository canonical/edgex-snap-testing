package utils

var ServicePorts = map[string]string{
	// platform services
	"core-data":         "59880",
	"core-metadata":     "59881",
	"core-command":      "59882",
	"vault":             "8200",
	"consul":            "8500",
	"redis":             "6379",
	"support-scheduler": "59861",
	// app services
	"app-rfid-llrp-inventory":  "59711",
	"app-service-configurable": "59701",
	// device services
	"device-gpio":            "59910",
	"device-modbus":          "59901",
	"device-mqtt":            "59982",
	"device-onvif-camera":    "59984",
	"device-rest":            "59986",
	"device-rfid-llrp":       "59989",
	"device-snmp":            "59993",
	"device-usb-camera":      "59983",
	"device-usb-camera/rtsp": "8554",
	"device-virtual":         "59900",
	// others
	"ekuiper":          "20498",
	"ekuiper/rest-api": "59720",
	"ui":               "4000",
}
