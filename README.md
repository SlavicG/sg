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
if
```
is used to execute block statements based on certain conditions. We can also use 
```
else
```
optionally. For Example:
```
if(a > b) {
    return a;
} else {
b;
}
```
is the syntax that will output the maximum numbers of $a$ and $b$.

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

```
return
```
is used to return values in functions, as could have been seen above. Additionally, the programming language has implicit returns, meaning that the last expression in a scope can be considered by the return value. So,

```
let inc(x) {
    return x + 1
}
```
works the same way as
```
let inc(x) {
    x + 1;
}
```
***
