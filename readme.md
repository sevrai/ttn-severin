### TheThingsNetwork - OpenSensors.io integration

#### features

This integration enables to forward messages published into a TheThingsNetwork's topic to another topic of OpenSensors.io.

* Using the MQTT protocol, it retrieves uplink messages from a dev topic *via* TTN application's login.
* The integration wait for a given number of messages before disconnecting.
* Messages are forward and published to OpenSensors using the realtime HTTP endpoint.

#### Configuration

The integration allows any TTN app with any dev's topic to be connected to any OpenSensors topic by any device with an API key.

The number of messages to link is also configurable. (-1 for an infinite loop)

All the parameters can be set in the [config.json](config.json) file.
