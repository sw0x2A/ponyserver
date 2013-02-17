ponyserver
==========

A simple TCP server to send ponies to everyone. Idea is based on [ponysay](https://github.com/erkin/ponysay) project which I also took the pony files and their quotes from.

Installation and usage on GNU/Linux (or other Unix implementations)
-------------------------------------------------------------------

To build the Go sources, [a working Go environment](http://golang.org/doc/install) is required on your computer. Clone this project and run:

	go run main.go --datadir .

The ponyserver listens on port 2000/tcp by default. You can change it by using the `--listen` parameter.

To connect to the server and receive a pony, use telnet, netcat or similar.

	echo | nc localhost 2000

Since the current version is limited and ignores whatever your request contains, the ponyserver will just send a random pony and quote.

License
-------

ponyserver is licensed under the terms of the [WTFPL](http://www.wtfpl.net/). You just DO WHAT THE FUCK YOU WANT TO.