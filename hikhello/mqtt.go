package main

import (
	"fmt"
	"time"

	log "github.com/go-pkgz/lgr"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func publish(client mqtt.Client, topic string, payload interface{}) {
	token := client.Publish(topic, 0, true, fmt.Sprintf("%v", payload))
	token.Wait()
	if token.Error() != nil {
		log.Printf("[ERROR] Error publishing topick %s: %v", topic, token.Error())
	}
}

func mqttPoller(config MQTTConfig) error {
	defer wg.Done()
	// Configure MQTT client options

	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", config.Host, config.Port))
	opts.SetKeepAlive(config.KeepAlive)
	opts.SetPingTimeout(config.PingTimeout)
	opts.Username = config.Login
	opts.Password = config.Pass

	// Create and start an MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		//panic(token.Error())
		return token.Error()
	}
	for {
		select {
		case <-dataChangedToMQTT:
			log.Printf("[DEBUG] polling to mqtt")

			for _, d := range deviceInfoList {
				publish(client, fmt.Sprintf("%s/%s/%d/name", config.Topic, d.Type, d.ID), d.Name)
				publish(client, fmt.Sprintf("%s/%s/%d/signal", config.Topic, d.Type, d.ID), d.Signal)
				publish(client, fmt.Sprintf("%s/%s/%d/temperature", config.Topic, d.Type, d.ID), d.Temperature)
				publish(client, fmt.Sprintf("%s/%s/%d/charge", config.Topic, d.Type, d.ID), d.ChargeValue)
			}

		}
		// Sleep for a specific interval before fetching data again
		time.Sleep(pollingTime)

	}

}
