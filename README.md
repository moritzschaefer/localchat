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


Info
====

For the server i used a lot from here http://gary.burd.info/go-websocket-chat. Thanks a lot!

Status
======

There is nothing working right now. Chat functionality is not tested. (I will build/copy a test client and some decent tests later). The location based features are not yet implemented

Developer Notes
==============

Please use git-flow (or the conventions of it)
