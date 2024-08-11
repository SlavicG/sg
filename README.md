# SG (StarGust) Programming Language

This is a quite minimalistic programming language I made for learning purposes. 

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

***

## Examples:

Here I will add some programs to share how the programming language works.

### Greatest common divisor function

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

puts(gcd(10, 5), gcd(12, 16), gcd(14, 21))
```
Running this program we will receive an output of
```
5 4 7
```

### Fibonacci Sequence

```
let fib = [1, 1]
for(let i = 2; i < 30; i = i + 1) {
    let val = fib[len(fib) - 2] + last(fib)
    push(fib, val)
}
puts(fib)
```
This program computes and outputs the first $30$ fibonacci numbers in linear-time. Here is the output produced to the console.
```
[1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, 1597, 2584, 4181, 6765, 10946, 17711, 28657, 46368, 75025, 121393, 196418, 317811, 514229, 832040]
```
***
### Edit Distance

Consider the program below
```
let s = "LOVE"
let t = "MOVIE"
let n = len(s)
let m = len(t)
let dp = [[0]]
let inf = 10000000
for(let x = 0; x < m + 2; x = x + 1) {
    push(dp[0], inf)
}
for(let ii = 1; ii < n + 2; ii = ii + 1) {
    push(dp, [inf])
    for(let jj = 1; jj < m + 2; jj = jj + 1) {
        push(dp[ii], inf)
    }
}
let min = fun(a, b) {
    if(a < b) {
        return a;
    }
    b
}
for(let i = 0; i < n + 1; i = i + 1) {
    for(let j = 0; j < m + 1; j = j + 1) {
        if(i > 0) {
            if(j > 0) {
                let add = 0
                if(get(s, i - 1) != get(t, j - 1)) {
                    add = 1
                }
                let val = min(dp[i][j], dp[i - 1][j - 1] + add)
                set(dp[i], j, val)
            }
            let vall = min(dp[i][j], dp[i - 1][j] + 1)
            set(dp[i], j, vall)
        }
        if(j > 0) {
            let valll = min(dp[i][j], dp[i][j - 1] + 1)
            set(dp[i], j, valll)
        }
    }
}
puts(dp[n][m])
```

Our language can even solve quite complex dynamic programming problems! For example, this code, can solve the classic [Edit Distance](https://en.wikipedia.org/wiki/Edit_distance#:~:text=In%20computational%20linguistics%20and%20computer,one%20string%20into%20the%20other.) problem in $O(n \cdot m)$ time complexity, where $n$ and $m$ are the lengths of the strings $s$ and $t$ respectively.

Our program outputs
```
2
```
in this case, which is correct.
***

### Maximum Subarray Sum

The following code solves the classic [Maximum Subarray Sum problem](https://en.wikipedia.org/wiki/Maximum_subarray_problem) in $O(n)$ time and memory complexity where $n$ is the length of the array.

```
let arr = [-1, 3, -2, 5, 3, -5, 2, 2]
let s = [0]
let ans = [0]
for(let i = 0; i < len(arr); i = i + 1) {
    set(s, 0, s[0] + arr[i])
    if(s[0] < 0) {
        set(s, 0, 0)
    }
    if(ans[0] < s[0]) {
        set(ans, 0, s[0])
    }
}
puts(ans[0])
```

The program will output:

```
9
```
Which is the correct answer for the given test case.
Note, that we use arrays here, because they ensure a correct scope. Unfortunately, the scopes are a bit chaotic when using for-loops currently. This will be fixed eventually.

*** 
#### Future Improvements TO DO List.

- Fix Loop Scopes
- Add character Data Types
- Add bitwise operators
- Add console input option
