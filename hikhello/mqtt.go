package main

import (
	"fmt"
	log "github.com/go-pkgz/lgr"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func publish(client mqtt.Client, topic string, payload interface{}) {
	token := client.Publish(topic, 0, true, fmt.Sprintf("%v", payload))
	token.Wait()
	if token.Error() != nil {
		log.Printf("[ERROR] Error publishing topick %s: %v", topic, token.Error())
	}
}

func mqttPoller() error {
	defer wg.Done()
	// Configure MQTT client options
	opts := mqtt.NewClientOptions().AddBroker("tcp://192.168.97.82:1883").SetClientID("hikax_mqtt_client")
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.Username = "vadim"
	opts.Password = "vad6Udkh"

	// Create and start an MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	for {
		select {
		case <-dataChangedToMQTT:
			log.Printf("[DEBUG] polling to mqtt")

			for _, d := range deviceInfoList {

				publish(client, fmt.Sprintf("hikax_devices/%s/%d/name", d.Type, d.ID), d.Name)
				publish(client, fmt.Sprintf("hikax_devices/%s/%d/signal", d.Type, d.ID), d.Signal)
				publish(client, fmt.Sprintf("hikax_devices/%s/%d/temperature", d.Type, d.ID), d.Temperature)
				publish(client, fmt.Sprintf("hikax_devices/%s/%d/charge", d.Type, d.ID), d.ChargeValue)
			}
		
		}
		// Sleep for a specific interval before fetching data again
		time.Sleep(pollingTime)

	}

}
