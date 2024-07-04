package hikaxprogo

type ZoneList struct {
	Zones []struct {
		Zone struct {
			ID               int    `json:"id"`
			Name             string `json:"name"`
			Status           string `json:"status"`
			SensorStatus     string `json:"sensorStatus"`
			TamperEvident    bool   `json:"tamperEvident"`
			Shielded         bool   `json:"shielded"`
			Bypassed         bool   `json:"bypassed"`
			Armed            bool   `json:"armed"`
			IsArming         bool   `json:"isArming"`
			Alarm            bool   `json:"alarm"`
			Charge           string `json:"charge"`
			ChargeValue      int    `json:"chargeValue"`
			Signal           int    `json:"signal"`
			RealSignal       int    `json:"realSignal"`
			Temperature      int    `json:"temperature"`
			SubSystemNo      int    `json:"subSystemNo"`
			LinkageSubSystem []int  `json:"linkageSubSystem"`
			DetectorType     string `json:"detectorType"`
			Model            string `json:"model"`
			StayAway         bool   `json:"stayAway"`
			ZoneType         string `json:"zoneType"`
			IsViaRepeater    bool   `json:"isViaRepeater"`
			ZoneAttrib       string `json:"zoneAttrib"`
			Version          string `json:"version"`
			DeviceNo         int    `json:"deviceNo"`
			AbnormalOrNot    bool   `json:"abnormalOrNot"`
		} `json:"Zone"`
	} `json:"ZoneList"`
}
