package hikaxprogo

// ExDevStatus Define the data structures
type ExDevStatus struct {
	OutputModList   []interface{} `json:"OutputModList"`
	OutputList      []interface{} `json:"OutputList"`
	SirenList       []SirenList   `json:"SirenList"`
	RepeaterList    []interface{} `json:"RepeaterList"`
	CardReaderList  []interface{} `json:"CardReaderList"`
	KeypadList      []interface{} `json:"KeypadList"`
	RemoteList      []interface{} `json:"RemoteList"`
	TransmitterList []interface{} `json:"TransmitterList"`
}

type SirenList struct {
	Siren Siren `json:"Siren"`
}

type Siren struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Seq           string `json:"seq"`
	Status        string `json:"status"`
	TamperEvident bool   `json:"tamperEvident"`
	Charge        string `json:"charge"`
	ChargeValue   int    `json:"chargeValue"`
	Signal        int    `json:"signal"`
	RealSignal    int    `json:"realSignal"`
	SignalType    string `json:"signalType"`
	Model         string `json:"model"`
	Temperature   int    `json:"temperature"`
	SubSystemList []int  `json:"subSystemList"`
	SirenColor    string `json:"sirenColor"`
	IsViaRepeater bool   `json:"isViaRepeater"`
	Version       string `json:"version"`
	DeviceNo      int    `json:"deviceNo"`
	AbnormalOrNot bool   `json:"abnormalOrNot"`
}

type ExDevData struct {
	ExDevStatus ExDevStatus `json:"ExDevStatus"`
}
