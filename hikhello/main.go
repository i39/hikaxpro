package main

import (
	"fmt"

	"os"

	"sync"
	"time"

	log "github.com/go-pkgz/lgr"
	"github.com/i39/hikaxprogo"
	"github.com/umputun/go-flags"
)

var revision = "latest"
var opts struct {
	HttpListen string `short:"l" long:"listen" env:"LISTEN" description:"listen on host:port" default:"0.0.0.0:8080"`

	HIKAX struct {
		Host     string `long:"host" env:"HIK_HOST" description:"host of the Hikvision AX device" required:"true"`
		Port     string `long:"port" env:"HIK_PORT" description:"port of the device" default:"80"`
		Username string `long:"username" env:"HIK_USERNAME" description:"username to access the device" required:"true"`
		Password string `long:"password" env:"HIK_PASSWORD" description:"password to access the device" required:"true"`
	} `group:"hikax" namespace:"hikax" env-namespace:"HIKAX"`

	PollingTime uint `long:"polling-time" env:"POLLING_TIME" description:"polling time in seconds" default:"10"`

	Dbg bool `long:"dbg" env:"DEBUG" description:"debug mode"`
}

type DeviceInfo struct {
	Type        string
	ID          int
	Name        string
	Signal      int
	Temperature int
	ChargeValue int
}

type HIKAXAuth struct {
	Host  string
	Port  string
	Login string
	Pass  string
}

var deviceInfoList []DeviceInfo
var mu sync.Mutex
var dataChangedToHTTP = make(chan bool)
var dataChangedToMQTT = make(chan bool)

// var dataChangedToMQTT = make(chan bool)
var wg = sync.WaitGroup{}
var pollingTime time.Duration
var hikAXAuth HIKAXAuth

func fetchData() {
	for {

		log.Printf("[DEBUG] Fetching new data from the device...")
		hik := hikaxprogo.New(hikAXAuth.Host, hikAXAuth.Port, hikAXAuth.Login, hikAXAuth.Pass)
		zoneList, err := hik.ZoneStatus()
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}
		exDev, err := hik.ExDevData()
		if err != nil {
			log.Printf("[ERROR] %v", err)
		}
		var newDeviceInfoList []DeviceInfo
		for _, zone := range zoneList.Zones {
			deviceInfo := DeviceInfo{
				Type:        "zone",
				ID:          zone.Zone.ID,
				Name:        zone.Zone.Name,
				Signal:      zone.Zone.RealSignal,
				Temperature: zone.Zone.Temperature,
				ChargeValue: zone.Zone.ChargeValue,
			}
			newDeviceInfoList = append(newDeviceInfoList, deviceInfo)
		}

		for _, siren := range exDev.ExDevStatus.SirenList {
			deviceInfo := DeviceInfo{
				Type:        "siren",
				ID:          siren.Siren.ID,
				Name:        siren.Siren.Name,
				Signal:      siren.Siren.RealSignal,
				Temperature: siren.Siren.Temperature,
				ChargeValue: siren.Siren.ChargeValue,
			}
			newDeviceInfoList = append(newDeviceInfoList, deviceInfo)
		}
		var dChanged = false
		if len(deviceInfoList) != len(newDeviceInfoList) {
			dChanged = true
		} else {
			for i, d := range deviceInfoList {
				switch {
				case d.Name != newDeviceInfoList[i].Name:
					dChanged = true
					break
				case d.Signal != newDeviceInfoList[i].Signal:
					dChanged = true
					break
				case d.Temperature != newDeviceInfoList[i].Temperature:
					dChanged = true
					break
				case d.ChargeValue != newDeviceInfoList[i].ChargeValue:
					dChanged = true
					break
				}

			}
		}
		if dChanged {
			mu.Lock()
			deviceInfoList = newDeviceInfoList
			mu.Unlock()
			dataChangedToHTTP <- true
			dataChangedToMQTT <- true
		}
		// Sleep for a specific interval before fetching data again
		time.Sleep(pollingTime)

	}
}
func main() {
	fmt.Printf("hikhello %s\n", revision)
	p := flags.NewParser(&opts, flags.PrintErrors|flags.PassDoubleDash|flags.HelpFlag)
	p.SubcommandsOptional = true
	if _, err := p.Parse(); err != nil {
		if err.(*flags.Error).Type != flags.ErrHelp {
			log.Printf("[ERROR] cli error: %v", err)
		}
		os.Exit(2)
	}
	setupLog(opts.Dbg)
	pollingTime = setPollingTime(opts.PollingTime)
	err := error(nil)
	hikAXAuth, err = setHIKAXAuth()
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}
	log.Printf("[DEBUG] options: %+v", opts)

	err = run()
	if err != nil {
		log.Fatalf("[ERROR] hikhello failed, %v", err)
	}
	log.Printf("[INFO] hikhello stopped")
}

func run() error {
	// Start the data fetching goroutine
	go fetchData()
	err := error(nil)
	// Start the HTTP polling goroutine
	wg.Add(2)
	go func() {
		log.Printf("[INFO] Starting HTTP poller on %s", listenAddress(opts.HttpListen))
		err = httpPoller()
		if err != nil {
			log.Printf("[ERROR] %v", err)
			return
		}
	}()
	if err != nil {
		return err
	}
	go func() {
		log.Print("[INFO] Starting MQTT poller")
		err = mqttPoller()
		if err != nil {
			log.Printf("[ERROR] %v", err)
			return
		}
	}()
	if err != nil {
		return err
	}
	wg.Wait()
	// Prevent the main function from exiting immediately
	select {} // Block forever

}

func setupLog(dbg bool) {
	if dbg {
		log.Setup(log.Debug, log.CallerFile, log.CallerFunc, log.Msec, log.LevelBraces)
		return
	}
	log.Setup(log.Msec, log.LevelBraces)
}

func listenAddress(addr string) string {

	// don't set default if any opts.Listen address defined by user
	if addr != "" {
		return addr
	}

	// http, set default to 8080

	return "0.0.0.0:8080"

}
func setPollingTime(ptime uint) time.Duration {
	if ptime != 0 {
		return time.Duration(ptime) * time.Second
	}
	return 10 * time.Second
}

func setHIKAXAuth() (HIKAXAuth, error) {
	hikAuth := HIKAXAuth{}
	if opts.HIKAX.Host == "" || opts.HIKAX.Username == "" || opts.HIKAX.Password == "" {
		return HIKAXAuth{}, fmt.Errorf("[ERROR] HIKAX host, username and password are required")
	}
	if opts.HIKAX.Port == "" {
		hikAuth.Port = "80"
	}
	hikAuth.Host = opts.HIKAX.Host
	hikAuth.Login = opts.HIKAX.Username
	hikAuth.Pass = opts.HIKAX.Password
	return hikAuth, nil

}
