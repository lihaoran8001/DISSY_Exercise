### Exercise 6.1

The first solution fails to satisfy the security policy. Assume the attacker is User B. 

When B intercept the traffic among A and D, he gets *E<sub>pk<sub>D</sub></sub>(R), S<sub>sk<sub>A</sub></sub>(E<sub>pk<sub>D</sub></sub>(R)), A*. Then he can use his own key sk<sub>B</sub> to sign the first part of this message E<sub>pk<sub>D</sub></sub>(R), and use it to replace the original signature. Also he need to replace the sender to B. So what he send to D is  *E<sub>pk<sub>D</sub></sub>(R), S<sub>sk<sub>B</sub></sub>(E<sub>pk<sub>D</sub></sub>(R)), B*. When D receive this message, he will use B's public key to verify the signature, and then send back result to B. 

This breaks the second security policy. In this situation, B can get information about A asked for.





