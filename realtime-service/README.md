# neutrino-realtime-service
A neutrino microservice responsible for the realtime websocket updates 

# Realtime protocol:

### From neutrino-core

```
{
    op: 'update|create|delete',
    payload: {
        _id: '{{id}}'
        prop1: 'v',
        prop2: 'c'
    }
}
``` 