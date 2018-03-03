Go Nano
=======

An implementation of the [Nano](https://nano.org/) protocol written from scratch in Go (golang).

About the Project
-----------------

A crypto currency has to be resilient to survive, and the network is only as resilient as the weakest link. With only one implementation of the protocol, any bugs that are found affect the entire network. The aim of this project is to create an alternative implementation that is 100% compatible with the reference implementation to create a more robust network.

Additionally, there is no reference specification for the Nano protocol, only a high level overview. I've had to learn the protocol from reading the source-code. I'm hoping a second implementation will be useful for others to learn the protocol.

Status
------
This software is in early development phase, and therefore is not suitable for use. We'll appreciate it however if you fetch yourself a clone, and start testing it (using `--testnet` flag).

Contributing
------------

Any contribution towards the advancement of Nano is much appreciated. If you see bugs, or room for improvement, please do jump in and make a pull request. We will also appreciate it if you communicate to us any comments or criticism regarding the project.

Credits
-------

This is a fork of [frankh](https://github.com/frankh/nano) repository. Kudos to him for starting this project. Also check out his [vanity address generator](https://github.com/frankh/nano-vanity) to generate a cool address for yourself :)
