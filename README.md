# Lightshow Server

A cloud-hosted server that connects controllers for LED strings with the
apps that people can use to control them. This server is part of the
[lightshow](https://github.com/guidowb/ligthshow) suite. See that page
for a description of the entire system.

## Build and Deploy

The lightshow server is currently built using Cloud Native Buildpacks
and deployed using ArgoCD.

## Design

(Some of this is aspirational. Code will approximate description over
time)

### Controllers

Controllers connect via websocket and identify by MAC address. The server
maintains configuration for each MAC address, including:

- friendly name (for logging and display)
- timezone for controller (for timer commands)
- name of program to run
- string configuration
- default pattern (when booting disconnected)

The string configuration is an array of connected strings, with for each string:

- pin
- number of LEDs in the string
- offset and direction
- color correction (array of segment, low, mid, high)
- array of client IDs allowed to manage this controller

### Clients

Clients connect via websocket and identify by server-provided secure cookie.
New clients get their cookie assigned on first connection.
