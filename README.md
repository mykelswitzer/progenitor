# progenitor
Progenitor is a code generator platform.

Based on the available commands, and by answering the prompts, it will create a functioning code base and eliminate a significant amount of time writing boiler plate code.

In most cases, the code generated will be a functioning service or application, to which the engineer need only add business logic.

See our [Getting Started Guide here](https://github.com/caring/progenitor/wiki/Getting-Started-with-Progenitor).

Following are a list of commands that can by run by Progenitor to create code scaffolds:

## go-grpc

The go-grpc command will generate a go microservice with a gRPC interface

We currently are building to [go 1.16](https://golang.org/)

In order to generate from a protobuf file, you must have [protoc](https://grpc.io/docs/protoc-installation/) installed

The go-grpc service will provision your infrastructure using the generated terraform code. Please see the note below regarding terraform.


### Notes on Terraform

Terraform is a language that allows you to declare your infrastructure, and using a provider, it will provision what is defined in the provider declared system (i.e. AWS)

We are using terraform 0.12.x currently. Please see these [instructions to install](https://learn.hashicorp.com/tutorials/terraform/install-cli).
