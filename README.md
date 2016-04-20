<p align="center">
[![Build Status](https://travis-ci.org/neutrinoapp/neutrino.svg?branch=master)](https://travis-ci.org/neutrinoapp/neutrino)
[![License](https://img.shields.io/badge/license-GNU%20AGPLv3-blue.svg)](https://github.com/neutrinoapp/neutrino/blob/master/LICENSE.md)
[![Build with <3](https://img.shields.io/badge/built%20with-%E2%9D%A4-red.svg)]() [![Alpha](https://img.shields.io/badge/State-Alpha-orange.svg)]()  [![](https://img.shields.io/badge/gitter-join%20chat-green.svg)](https://gitter.im/go-neutrino/neutrino-core?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
<p/>

<p align="center">![Neutrino](http://v41.imgup.net/NEUTRINO-l3931.png)<p/>

#### Note: The project is still in its alpha stage, if you still want to see it running, check out the demo in the [website](http://neutrinoapp.com).

## Introduction

Neutrino is a service that allows developers to link their mobile and web applications with a real time cloud data storage. 
Every time you create a new application, you need to write code to authenticate users, store data, manage infrastructure, handle scalability.
Neutrino gives you out-of-the-box solution to easily manage your application’s users and data. You can be sure that your data will always be available and scale along with your application.

## Architecture

* [API Service](https://github.com/neutrinoapp/neutrino/tree/master/src/services/api) - This server handles HTTP requests. It can register and modify users as well as perform CRUD on data. The data is organized as follows:
  * One has to register an account in the service
  * Once he has an account he can create an application
  * The application has users as well, which will allow you to set permissions on different groups of users for specific data
  * The application has collections - Collections are implicit, meaning that if you create a data item in the **cars** collection, it will be automatically created. See below for the realtime specifics.
  * Each collection contains data items - See below for the realtime specifics.
*  [Realtime Service](https://github.com/neutrinoapp/neutrino/tree/master/src/services/realtime) - This server handles WebSocket connections. It can perform CRUD operations on the data. It can notify clients for changes in the following scenarios:
  * When an item in a collection is added, changed, removed etc. This allows one to bind to a whole collection as an array on the client. E.g. if you want to list all **cars** you can be sure that they will always be up to date.
  * When a single item is updated - allows listening for changes on per-item-basis.
  * Check the client (api)[https://github.com/neutrinoapp/neutrino#api] below to get a better idea of the API
* [Redis cache](http://redis.io/) - Not mandatory
* [RethinkDB](http://rethinkdb.com) - For storing the data and listening for changes.

**Note:** The API and Realtime services can be scaled horizontally. E.g. they are currently running in Google Container Engine in Kubernetes pods.

## Other components

* [Javascript client library](https://github.com/neutrinoapp/neutrino-javascript) - Handles HTTP and WebSocket communication
* CLI(TBD) - To manage everything from the command line

## API

```javascript
var app = Neutrino.app('demo');
var collection = app.collection('demo-items');
app.auth.login('demo', 'demo').then(function () {
		//get all objects in the collection
		collection.objects({realtime: true})
    	.then(function (objects) {
      	//get notified when an object is added
    		objects.on(Neutrino.ArrayEvents.added, function (item) {
    			console.log(item); //{id: 'GUID', itemText: 'some-text'}
    		});  
      });

		//lets create an object
		collection.object({
    	itemText: 'some-text'
    }).then(function () {
    	console.log('item created!');
    });
    
    //we can also get an item by id
    
    //we can get both realtime and non-realtime objects
    collection.object('some-id', {realtime: true})
    	.then(function (object) {
      		//any change that happens on the server or on the client will affect this object
          //if we change its text all clients will receive the update
          object.on(Neutrino.ObjectEvents.propertyChanged, function () {
            console.log(object);
          });
          
          object.itemText = 'some-changed-text';
      });
      
     collection.object('some-id', {realtime:false})
     	.then(function (object) {
      		//changes will not automatically propagate to this object
          //we can still update it manually and all other clients will receive the update again
          object.itemText = 'manual-text';
          return object.update();
      });
});
```

## Documentation

*Comming soon*

## Local development

**Tested only on linux**

Start rethinkdb by downloading it for your distribution or if you have docker installed, run `make dev`

Start the **api service** by running `make api`

Start the **realtime service** by running `make realtime`

You can use the **postman configurations** from the *postman* folder until docs come out

## Related projects

* [neutrino-javascript](https://github.com/neutrinoapp/neutrino-javascript) - Javascript Client
* [neutrino-todo-app](https://github.com/neutrinoapp/todo-app) - Sample todo application

