# sauron

Telemetry Tracking Server

## Instalation

```
$ go get -u github.com/mguzelevich/sauron/...
```

## availible fields

- lat
- lon
- sat
- desc
- alt
- acc
- dir
- prov
- spd
- time
- battery
- androidId
- serial
- activity
- epoch

## Examples

run 

```
$ sauron -p 8080
```

send location

```
$ curl -X POST --data "lat=53.9279421&lon=27.6437863&time=2017-12-27T12%3A30%3A30.338Z&s=0.0&prov=network&aid=4de1a4a0e296ef63&acc=21.795000076293945" localhost:8080/log
```

web dashboard `http://localhost:8080/ui`

