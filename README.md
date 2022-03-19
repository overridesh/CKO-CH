## Checkout Challenger

### Docs
---
- [Bank Simulator](bank_simulator/README.md) 
- [Gateway](gateway/README.md) 
- [Insecure](gateway/certs/README.md) 
  
### How to run?
---
Requirements:
 - docker
 - docker-compose

```bash
./run-dev.sh
```

### Assumptions
--- 
- It was assumed that the service is already PCI-certified for storing credit cards. Cards are stored per transaction, we cannot use them, just as an example. They could be moved later.
- In the project phase, I decided to store cards but not tokenise them. It only works with credit cards in the body for the time being.
- Amount: integer number of cents, so if you would like to purchase $40.23, you need to pass an amount of 4023. This makes it easy to handle different currencies and always be taken as an integer in the backend.
- Bank Simulator is a static API. It does not store or register anything. It only validates the information received and returns as configured.
- It is not a real credit card or payment management project. Everything works in a simulation. However, the transactions are stored in the gateway database.

### Areas for improvement.
---
- A Code Review, I may have overlooked many things.
- Handling of the authorisation.
- A secret manager, such as Vaul or AWS Secret Manager.
- Request validations.
- Integration test and more unit and detailed testing. For example, currently, I only tested if specific errors were returned. But I didn't make any comparison of objects.
- Currently, storing idempotent keys in memory will give us a problem if we have more than one instance of the project up, as the memory between instances will not be shared. We need to centralise this.
- Event queue handling for cases where transactions were not reported as successful. We must make a refund/void.
- Card tokenisation.
- Idempotence: It is not a correct way as we store it in memory. I only simulate recovery points per step of a transaction.

### What cloud technologies you'd use and why.
---
  - We could use queues technology like SQS for creating refund/void o push payments.
  - For idempotency keys, ideally in a persistent and centralised database like SQL or NoSQL. Redis can be an excellent option to add a TTL, and we can centralise the information there.
  - Kubernetes like orchestrator, I have experience working with Kubernetes. You could implement GRPC pods and GRPC-Gateway pods separately, leave the gRPC service as private (no incoming Internet connections but outgoing to integrations), and gRPC-Gateway can only connect through the private Kubernetes network to the service.
  - Vault for secrets injection.
  - Terraform for the infrastructure ecause we need to standardise the architecture process. So that we understand all our resources and can replicate environments or apply configurations as we usually do in code.
  
### Extra
---
- Idempotency key
- There is no merchant controller. It is currently inserted into the database through a migration.
- Implemented in gRPC, and grpc-gateway was used as a proxy for HTTP output.
- The TLS certificates of this project are insecure [more info](gateway/certs/README.md).
- OpenAPI documentation auto-generated from protobuf.

### Architecture 
---
<p align="center" width="100%">
    <img width="50%" src="service.png?raw=true"> 
</p>

### Resources
- Inspired by https://github.com/grpc-ecosystem/grpc-gateway