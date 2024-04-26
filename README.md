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
There is a message limit of 10 hardcoded into each queue for the time being as well to allow for easy testing of the history feature. 
As this project matures, more features will be added and eventually, clients will be able to be instantiated to listen on a certain queue and then process the messages as they come in.

## How to use GoLong in its current state

1. clone the repo
2. navigate into the golong dir and run `go run .` to start a server
3. open a new terminal window and make a connection to the server by running `ssh -o "StrictHostKeyChecking=no"  -p 2222 your_name_here@127.0.0.1` or the same command and replace `127.0.0.1` with the IP of the computer running the server
4. in the control terminal just opened, send commands 


## Current Features

Currently, a connection can be made to a running server via SSH and messages can be sent to the server. 

A user can run any of the following commands as well:
- /nq    -- this prompts the user for a name and then creates a new queue with that name
- /sq    -- this prompts the user with a list of existing queues and allows the user to type the name of that queue to switch to it
- /h     -- this displays up to the last 10 messages sent on a particular queue
- /exit  -- this closes the SSH connection to the server


## Support

Feel free to add an issue to this repo for any feature requests, bugs, or ideas! I am not currently accepting PRs right now as this project is still in its infancy.
