# Introduction 
ECS Service is a Golang application that hosts user sessions. Each session maps back to a (relatively) isolated instance of EngoEngine/ecs [https://github.com/EngoEngine/ecs] (hereafter referred to as a World)

The communication protocol used is Websockets with the messages serialized in Protobuf encoding. All inputs are received and applied as-is. Further implementations may include deduplication, message validations and other modes of latency compensation mechanics.


# Getting Started
The server listens to both HTTP/WS connections on port 9999 by default
- `cd core`
- `bazel build //protocol/src:go`
- `bazel run //server/src `

The service mainly interfaces via an HTTP api (`v1/session`). After submitting a `POST` request, a 'World' object is created and given an ID. This `World` object contains the entire state of a given session and manages all communication between connected clients.

