milestone 1:
============

features
--------

- Receive messages from  7 km radius
- Update own position
- Update own name
- Send messages

Protocol:
JSON based on WebSockets.

1. connect to ws://server/init (server creates new user for that  connection with uninitiaized location and name)
2. client sends update username
3. client sends update position
4. ready to go. client can now send messages and will receive messages

Clients send to server:
{"action": <>, "field": <>, "value": <>}
actions:
- update(field is username or position)
    - position value is in form "52.1242352345,13.123"
- message
    - "field" stays empty or is not present
    - message is sent in "value"

status
------

done

milestone 2:
============

- Write a mobile-ready client. (Android or HTML5)

### to discuss:
- Should message-radius be variable/user-set?
- Implement tests in this milestone?
- Security model


Info
====

For the server and the webclient i used a lot of code from http://gary.burd.info/go-websocket-chat. Thanks a lot!

Status
======

The prove of concept is running, though not really tested.

Developer Notes
==============

Please use git-flow (or the conventions of it)
