```
 _______  _______  _        _______  _        _______  _______  _______ 
(  ____ \(  ___  )( \      (  ___  )( (    /|(  ____ \(       )(  ___  )
| (    \/| (   ) || (      | (   ) ||  \  ( || (    \/| () () || (   ) |
| |      | |   | || |      | |   | ||   \ | || |      | || || || |   | |
| | ____ | |   | || |      | |   | || (\ \) || | ____ | |(_)| || |   | |
| | \_  )| |   | || |      | |   | || | \   || | \_  )| |   | || | /\| |
| (___) || (___) || (____/\| (___) || )  \  || (___) || )   ( || (_\ \ |
(_______)(_______)(_______/(_______)|/    )_)(_______)|/     \|(____\/_)
                                                      
```

GoLong is a message broker queue system built in Go.

## About

GoLong now has a working broker/subscriber structure. golong/consumer contains the logic for the basic subscriber package. Currently the package only prints the message to the terminal window, but this is where a true subscriber would take the message and perform some business logic with it.

GoLong uses a TCP connection to facilitate the message brokerage. Within your Producer or Consumer project, you will need to add the code needed to make such a connection and send the appropriately formatted message to the broker.

## How to use GoLong in its current state

1. clone the repo
2. navigate into the golong/server dir and run `go run .` to start a server
3. in a new terminal window navigate into the golong/consumer dir and run `go run .` to start a consumer that is subscribed to the "test" queue (one of the two default queues offered by GoLong)
3. open a new terminal window and make a connection to the server by running `telnet 127.0.0.1 2222` or the same command and replace `127.0.0.1` with the IP of the computer running the server
4. in the control terminal just opened, send commands and they will be produced on the consumber if sent to the subscribed queue

## Message Structure

A Message sent to the server is formatted as `[queue]:[message]`
Optionally, commands can be run to either add new queues or retrieve the history of a queue. This is done with this syntax `[queueName]:[cmd]:[command_to_run]`

## Current Features

A user can run any of the following commands as well:
- nq       -- this will create a new queue using the [queueName] field as the new queue's name. (ex: `newqueue:cmd:nq`)
- hist     -- this displays the previous messages sent on a particular queue (ex: `queueName:cmd:hist`)

## Consumers

Currently, the default configuration subscribes to the "test" queue of the server. To change this, modify the `if _, err := conn.Write([]byte("sub:test\n"));` line. 
NOTE: you will need to create the new queue first if it does not already exist for this to work

Any tcp connection can also be made to subscribe to a queue if a message that follows the `sub:[queueName]` formatting is sent as a message from that client instance

## Support

Feel free to add an issue to this repo for any feature requests, bugs, or ideas! I am not currently accepting PRs right now as this project is still in its infancy.
