```
 _______  _______  _        _______  _        _______ 
(  ____ \(  ___  )( \      (  ___  )( (    /|(  ____ \
| (    \/| (   ) || (      | (   ) ||  \  ( || (    \/
| |      | |   | || |      | |   | ||   \ | || |      
| | ____ | |   | || |      | |   | || (\ \) || | ____ 
| | \_  )| |   | || |      | |   | || | \   || | \_  )
| (___) || (___) || (____/\| (___) || )  \  || (___) |
(_______)(_______)(_______/(_______)|/    )_)(_______)
                                                      
```

GoLong is a message broker queue system built in Go.

## About

Currently, this project is in its early stages of creation. The code, as it stands right now, functions more like a chat between a client and a server but with messages only flowing one way to the server from the client.
As this project matures, more features will be added and eventually, clients will be able to be instantiated to listen on a certain queue and then process the messages as they come in.

GoLong uses a TCP connection to facilitate the message brokerage. Within your Producer or Consumer project, you will need to add the code needed to make such a connection and send the appropriately formatted message to the broker.

## How to use GoLong in its current state

1. clone the repo
2. navigate into the golong/server dir and run `go run .` to start a server
3. open a new terminal window and make a connection to the server by running `telnet 127.0.0.1 2222` or the same command and replace `127.0.0.1` with the IP of the computer running the server
4. in the control terminal just opened, send commands 

## Message Structure

A Message sent to the server is formatted as `[queue]:[message]`
Optionally, commands can be run to either add new queues or retrieve the history of a queue. This is done with this syntax `[queueName]:[cmd]:[command_to_run]`

## Current Features

A user can run any of the following commands as well:
- nq       -- this will create a new queue using the [queueName] field as the new queue's name. (ex: `newqueue:cmd:nq`)
- hist     -- this displays the previous messages sent on a particular queue (ex: `queueName:cmd:hist`)


## Support

Feel free to add an issue to this repo for any feature requests, bugs, or ideas! I am not currently accepting PRs right now as this project is still in its infancy.
