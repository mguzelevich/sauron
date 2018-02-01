# sauron [![GoDoc](https://godoc.org/github.com/mguzelevich/sauron?status.svg)](http://godoc.org/github.com/mguzelevich/sauron) [![Build Status](https://travis-ci.org/mguzelevich/sauron.svg?branch=master)](https://travis-ci.org/mguzelevich/sauron)

Telemetry Tracking Server

## Instalation

```
$ go get -u github.com/mguzelevich/sauron/...
```

## features

telemetry sources:

- custom http url GET, POST
- udp

### custom http url fields

- `lat` - Latitude
- `lon` - Longitude
- `desc` - Annotation
- `sat` - Satellites
- `alt` - Altitude
- `spd` - Speed
- `acc` - Accuracy
- `dir` - Direction
- `prov` - Provider
- `time` - Time UTC (2011-12-25T15:27:33Z)
- `batt` - Battery
- `aid` - Android ID
- `ser` - Serial
- `act` - Activity

example:

```
$ curl localhost:8080/log?lat=53.9279421&lon=27.6437863&time=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=4de1a4a0e296ef63&acc=21.795000076293945

$ curl -X POST --data "lat=53.9279421&lon=27.6437863&time=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=4de1a4a0e296ef63&acc=21.795000076293945" localhost:8080/log
```

### UDP

example:

```
$ bash -c 'echo -e "mgu/mi5s/$GPRMC,083543,A,5355.67728,N,2738.62654,E,0.000000,0.000000,010203,,*2E" > /dev/udp/127.0.0.1/8822'
```

## Examples

run 

```
$ sauron --http :8080 --udp :8822 --ui :8081
```

web dashboard `http://localhost:8081/`

