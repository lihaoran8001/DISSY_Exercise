Assume we have one Byzantine-corrupted sender and two honest parties. The sender can send many different messages with the same broadcast ID within one round.

In our protocol, the node only keeps the first message with the same bid and signature in one round. If the sender sends 2 messages (1,m1,σ), (1,m2,σ) to P1, P1 will simply drop the second message(1,m2,σ). The rest part of the protocol is the same as the Dolev-Strong protocol.

![dolvestrongimpove](/Users/lihaoran/Downloads/dolvestrongimpove.png)