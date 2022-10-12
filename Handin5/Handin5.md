## Exercise 9.3

1. Our algorithm has several steps.
   1. First, since we have `h(pw)`, we can get *F<sup>1</sup>(pw)* = `g(h(pw))`. Then we do F  n times to get a list L = [F<sup>1</sup>(pw), F<sup>2</sup>(pw), ..., F<sup>n</sup>(pw)].
   2. For i = 1 to t, we will find out if *y<sub>i</sub>* is in list L. Fortunately, if *y<sub>i</sub>* == F<sup>k</sup>(pw) == L[k], we can conclude that the password is located in this row i. 
   3. We calculate M[i, n - k], which is F<sup>n-k</sup>(pw<sub>i</sub>), if there is no conflict, this value is equal to pw. (According to the question, In matrix M, i starts from 1, and index of column starts from 0).

2. In step 1.a, to get list L, we compute F n times. In step 1.c, to get M[i, n - k], we calculate F (n - k) times. So we compute F (2n-k) times in total.

3. If not all y<sub>i</sub> 's are distinct, we might be able to find several different candidates which might be correct password. They are located in different rows.

   To ensure all y<sub>i</sub> 's are distinct, we need to ensure that our function `g(x)` and function `g(h(x))` is collision resistent, which means (ideally) there's no same output given different input.