# SG (StarGust) Programming Language

This is a very simple programming language I made for learning purposes. 

#### Credits
I learned a lot and was inspired from the Books by Thorsten Ball - "Writing an Interpreter in GO" and "Writing a Compiler in GO".

## Language Features:

### Syntax:

```
let
```

is used for declaring variables. For example we can use:

```
let x = 3 + (5 * 4 + 12) * 2
```
or 
```
let y = "abc" + "def"
```

to declare $x$ and $x$.

***

```
fun
```
is used to represent functions. For example:

```
let gcd = fun(x, y) {
    if(x == 0) {
        return y
    }
    if(y == 0) {
        return x
    }
    if(x > y) {
        return gcd(x - y, y)
    } else {
        return gcd(x, y - x)
    }
}
```

***

