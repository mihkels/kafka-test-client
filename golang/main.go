package main

func main() {
	var appMode = ConfigInstance.ApplicationMode

	if appMode == "producer" {
		runProducer()
	} else {
		runConsumer()
	}
}
