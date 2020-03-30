# go_socketserver

Socket communication for custom protocols via GO

There are two versions, the master is the ordinary version, and the room is the version for clients.

The main core communication file is the socket directory, which contains the solution of socket sticky packets.

'config.json' config port and log level


##protocol

the head char is '*',next have 4 chars are 'Data' length,and 2 chars are checksum.


| head | Data Length | Check Sum |
|:----:| :----: | :----: |
| 1 byte | 4 bytes | 2 bytes |
