# Kaboom
**Kaboom** is a decentralized [dead man switch](https://en.wikipedia.org/wiki/Dead_man%27s_switch#Software):
* it encrypts your secrets
* it splits the encrypted secrets  between a customizable number of sensor nodes
* it allows you to send a regular heartbeat to the nodes
* if your heartbeat stops coming, the nodes coordinate, unencrypt and release the payload to your selected audience.

## Warning
Cobsider this alpha software and a theoretical exercise. Imagine it was written on trams, trains and in waiting rooms. Be aware it's not peer reviewed. In short: **don't use this software in scenarios where your life depends on it; it might do the job, but if doesn't: I told you**.

## Roadmap
- [x] Encryption/Decryption with threshold cryptography
- [x] Authenticated messaging
- [x] IPFS payload storage
- [x] IPFs/pubsub heartbeat
- [x] dynamic heartbeat topic/channel
- [x] health chack logging
- [ ] multi node release protocol
- [ ] release to non-nodes
