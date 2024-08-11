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

```
for
```
is used to create cycles. For example we can use:
```
for(let i = 0; i < 5; i = i + 1) {
    puts(i * 2 + 1)
}
```
in order to print the first five odd numbers. 

We can also nest loops. For example, the following code will print the multiplication table for all numbers between $1$ and $9$ (with a weird ordering to show different types of orderings).

```
for(let i = 1; i < 10; i = i + 1) {
    for(let j = 9; j > 0; j = j - 1) {
        puts(i, j, i * j)
    }
}
```

Currently, the for loops don't work ideally and there are some problems with the scopes. It will be solved in a later version, where I will change the structuring of scopes and use a Scope-Stack.

### Data Types:

For now, the available data types are:

- Booleans: Simple true/false values. They are returned by conditional statements such as
```
3 == 5 // 4 < 2 // 5 > 8
```
- Integers: Simple 64-bit integer values. They can be added, subtracted, divided, multiplied, and returned in functions. To declare an integer variable we can use the following format, for example:
```
let x = 2 + 3
```
- Strings: Simple array of characters. They can only be concatenated so far. To declare a string variable we can use the following format, for example:
```
let s = "abc" + "DeX"
```
- Arrays: Arrays are actually dynamic in this programming language, and update their size dynamically based on the number of elements we append to them (using builtin functions). The push(array, value) function works in amortized constant time complexity, by multiplying the array size by two every time the size goes over the corresponding capacity. Here is how to declare an array of integers, and how to access the corresponding indices:
```
let arr = [1, 2, 4]
puts(arr)
puts("Second element is:", arr[1])
```

This will output $2$.

### Operators:

For now, the only possible operators are 
- "+" (addition/concatenation)
- "-" (subtraction)
- "*" (multiplication)
- "/" (division)
- "<" (less than)
- ">" (greater than)
- "==" (equals)
- "!=" (doesn't equal)
- "=" (assignment).

### Built In Functions

This programming language also has some built-in functions.

```
len(a)
```
Takes an array or a string as an argument. This will return the length of the corresponding array $a$ or the corresponding string $a$. It works in constant time.

***

```
puts(arg1, arg2, ..., arg_n)
```
Takes any data type as arguments. This outputs $arg_1$, $arg_2$, $...$, $arg_n$, separated by space and they are followed by a new-line in the end. It works in the time it takes to output all arguments.

***

```
first(arr)
```
Takes an array as an argument. This returns the first element of the array $arr$. For example it would return the value of $6$ if $arr = [6, 2, 8, 4]$. It works in constant time.

***

```
last(arr)
```
Takes an array as an argument. This returns the last element of the array $arr$. For example it would return the value of $4$ if $arr = [6, 2, 8, 4]$. It works in constant time.

***

```
push(arr, val)
```
Takes two arguments, an array and a value that should be the same data type as the elements in the array. It appends an element to the end of the array $arr$. It works in amortized constant time.

***

```
set(arr, index, value)
```
Takes three arguments, an array, an integer representing the index, and a value that should match the arrays data type. It sets the value val to the array element on position $index$. Works in constant time.

***

```
shuffle(arr)
```
Takes an array as argument. It shuffles the array in a random manner. Works in linear time complexity.

```
reverse(arr)
```
Takes an array as argument. It reverses the ordering of the elements of the array. Works in linear time complexity.
***

```
sort(arr)
```

Takes an array as argument. It sorts the array increasingly (currently only works for integers as strings don't have a comparator as of now). It works with a time complexity $O(n \cdot log(n))$ average-case using quick-sort behind it. 
