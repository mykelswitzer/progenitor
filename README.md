# progenitor
Progenitor is a code generator platform.

Based on the available commands, and by answering the prompts, it will create a functioning code base and eliminate a significant amount of time writing boiler plate code.

In most cases, the code generated will be a functioning service or application, to which the engineer need only add business logic.

## Quick start
1. [Initialise a new go module](https://golang.org/doc/tutorial/create-module)

       mkdir example
       cd example
       go mod init example

2. Initialise progenitor

       go run github.com/mykelswitzer/progenitor init

