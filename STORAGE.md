# sauron storage schema

## db

```
sauron
├── <B> .meta
├── <B> accounts
│   ├── <K> user_id
│   │   ├── <F> first_name
│   │   ├── <F> last_name
│   │   ├── <F> ...
│   │   └── ...
│   └── ...
├── <B> devices
│   ├── <K> device_hash
│   │   ├── <F> user_id
│   │   ├── <F> device_id
│   │   ├── <F> android_id
│   │   └── <F> serial
│   └── ...
├── <B> telemetry
│   ├── <B> <user_id>
│   │   ├── <B> <device_hash>
│   │   │   ├── <K> <timestamp>
│   │   │   │   ├── telemetry_entity
│   │   │   │   └── ...
│   │   │   └── ...
│   │   └── ...
│   └── ...
└── ...
```

```


```

### .meta

### accounts

store user accounts

```
uuid -> {
	device_id: ...,
	android_id: ...,
	serial: ...,
}
```

### devices

map device info hash to user id
```
device_hash -> user_id
```

### telemetry

store telemetry info

```
user_id -> timestamp -> point

```