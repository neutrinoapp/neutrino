# neutrino-realtime-service
A neutrino microservice responsible for the realtime websocket updates 

# Realtime protocol:

```
{
    op: 'update|create|delete', //the operation to perform
    origin: 'api|client', //the origin of the operation
    options: {}, //additional options
    pld: { //the payload to apply
        id: '{{id}}'
        prop1: 'v',
        prop2: 'c'
    },
    type: 'type',
    app: 'appId'
}
``` 