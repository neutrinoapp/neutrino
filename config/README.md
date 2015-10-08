# neutrino-config

Can load configs from the following paths:

```go
"/etc/neutrino/neutrino-config"
"/var/neutrino/neutrino-config"
```

Has the following defaults:

```go
KEY_MONGO_HOST = "mongo-host"
KEY_QUEUE_HOST = "queue-host"
KEY_CORE_PORT = getVar(corePrefix, "port")
KEY_REALTIME_PORT = getVar(realtimePrefix, "port")
```

Usage:

```go
nconfig.LoadWithPath(path)
nconfig.Load()
```