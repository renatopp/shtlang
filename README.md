# **oh SHT!** A new scripting language.

> **Notice!**
>
> This language is under development and this document is an inconsistent WIP.

# Hello Wolrd

```bash
print('Hello, Wolrd')
```

# The Language

## Variables

```tsx
let variable = 'Hello'
const constant = 3.14

{
  let scoped = true
  variable = 3
}
```

## Comments

```rust
# This is a comment
```

## Type System

All primitive types:

- Number
    
    ```rust
    let pi = 3.1415
    let py = 3
    let large = 1e1000
    let frac = .3
    ```
    
- String
    
    ```rust
    let pi = 3.14
    let s = 'hello $pi'
    let ms = `
    Hello,
    $pi
    World
    `
    ```
    
- Boolean
    
    ```rust
    let chance = true or false
    ```
    
- Function
    
    ```rust
    fn say(str) {
    	print(str)
    }
    
    let func = => 1
    let func1 = (x) => 1
    let func2 = (x) => { 1 }
    ```
    
- List
    
    ```rust
    let list1 = List()
    let list3 = List{ 1, 2, 3, 4, 5 }
    let list4 = List(){ 1, 2, 3, 4, 5 }
    let list5 = List(=> 2) # Default value is 2
    ```
    
- Dict
    
    ```rust
    let dict1 = Dict(=> 1) # Devault value for when the key does
    let dict2 = Dict{ 1:1, 2:2, 3:3, 4:4, 5:5 }
    ```
    
- Set
    
    ```rust
    let set = Set()
    ```
    
- Tuple
    
    ```rust
    let tuple = (1, 2)
    let tuple = (1,)
    return (1, 2)
    let tuple = Tuple{1,2}
    ```
    
- Error
    
    ```rust
    Error{ msg='Invalid!' }
    raise 'Invalid!'
    ```
    
- Maybe
    
    ```rust
    fn func()? {
    	return 1
    }
    fn func()? {
    	raise 'exception'
    }
    
    Maybe(value)
    Maybe(error)
    Maybe() // error as default
    ```
    
- Data
    
    ```rust
    data User {
    	name = ''
      email = ''
    }
    ```
    
- Regex
    
    ```rust
    /[a-zA-Z]/g
    ```
    
- Iterator
    
    ```rust
    a = Iterator(nextFn)
    a.next()
    a.finished
    ```
    
- Type
    
    ```rust
    
    ```
    

### Meta Programming

- doc
- name
- module
- file

- on call(this) …
- on set(this, property, value)
- on get(this, property) value
- on has(this, property) bool
- on new(this) this
- on index(this, …) value
- on iter(this) Iterator
- on bool(this) bool
- on string(this) string
- on repr(this) string
- on bang(this) …
- on eq
- on neq
- on gt
- …

```python

```

## Operators

****Arithmentic****

| OPERATOR | DESCRIPTION |
| --- | --- |
| + |  |
| - |  |
| * | Multiplications can be performed implicitly as 3km, where km is a variable. |
| / | Division by 0 results in runtime error |
| % | Mod |
| ** | Pow |
| // | Equivalent to floor(x/y) |
| ++ | Used as suffix only, increments 1 |
| — | Used as suffix only, decrements 1  |

**Relational**

All relational operators return 0 or 1.

| OPERATOR |
| --- |
| == |
| ≠ |
| > |
| < |
| ≥ |
| ≤ |

**********Logical**********

All logical operators return 0 or 1.

| OPERATOR |
| --- |
| ! |
| and |
| or |
| xor |
| nor |
| nand |

**************Special**************

| OPERATOR |
| --- |
| ++ |

## Functions

Functions as defined in a couple of ways:

```python
fn log(value) {
	echo value
}

log = fn(value) { ... }
log = fn { ... } # no parameter
log = () => { ... }
log = => ... # no parameter
```

Functions with arguments can be called without parenthesis if they aren’t being called inside another parenthesis-less functions:

```python
fn add(a, b) { ... }
add(a, b)
```

Function return the last statement by default.

```python
fn fib(n) {
	if n < 2 return 2
	fib(n-1) + fib(b+2)
}
```

Arguments can have a default value, and they also can be spread:

```python
fn (a, ...b, c) { ... }
fn (a, ...b, c=1) { ... }
fn (a=1, b=2, ...c) { ... }
```

Adding more parameters to a function won’t have any effect. Default value can only be if the next argument also has a default. Spread doesn’t have a default.

Functions have a special notation to return Maybe objects:

```python
fn uncertain? {
	return 2
}

uncertain() # is a maybe(number)
```

## Iterators

Every type that implements an iterator pattern can be piped:

```python
'hello'
| map char: char.toUpperCase()
| foreach char: echo char
```

notice that:

- Strings implement by default an interator over its characters
- Right after a pipe, you must provide an piped function
- These functions receives the value as first parameter, fixed
- These function may receive a function as last argument, in this case, the function have a special declaration format
- Iterators items are evaluated one by one, until finished or interrupted.
- Iterators are lazy, they are obly processed as necessary
- The pipe may return the aggregated value by its last statement.

### Iterator Pattern

Every data type can implement an meta operator called `iter`, which must return an Iterator object:

```tsx
data OneTwoThree {
  on iter(this) {
    let i = 0
    return Iterator(fn {
      i++
      if i>3 return Iteration(false)
      return Iteration(this, i)
    })
  }
}
```

Notice that an Iterator object is created with a next function, that will be triggered everytime the iterator is queried for the next element. Also notice that the iterator must return a tuple with the owner and the value. If it returns anything else, the iterator will be signaled to stop.

You can use a special expression to represent an iterator creation:

```tsx
data OneTwoThree {
  on iter(this) {
    yield 1
    yield 2
    yield 3
  }
}
```

Both examples are equivalent. 

Any function that contains the `yield` statement is called a *generator*, which returns an iterator implicitly. Notice that, the generator function return `Iteration(false)` as default after the end of the function.

### Piped Functions

Piped functions are functions that can be used in pipes.

```rust
OneTwoThree()
| map x: x + 1
```

A sample `map` function could be implemented as:

```rust
fn map(next, func) {
	for next() as iteration {
		yield func(iteration.values...)
	}
}
```

Notice that the piped function function receives a next function (from the previous iterator), and returns an iterator itself. 

A piped function may receive additional arguments:

```rust
fn reduce(next, func, initial) {
	let acc = initial
  for next() as iteration {
		value = func(acc, iteration.values...)
  }

	yield value
}

List { 1, 2, 3 } | reduce(0) acc, x: acc + x
```

### Pipe Response

Pipes respond lists by default. However you can specify its return type by doing using the expression:

```rust
word = 'word' | to String                  # 'word'
word = 'word' | to List                    # List{'w', 'o', 'r', 'd'} -- which is redudant
size = 'word' | map x: 1 | sum | to Number # 4
data = 'word' | to Custom                  # Custom object
iter = 'word' | to Iterator                # Iterator, i.e., does not execute

string = text.split(' ')
| map x: /[^a-zA-Z]+/.replaceAll(x, '')
| filter x: x
| join x: ' '
| to String
```

In order to accept this convertion, the data type must implement the meta function `iterresp`:

```rust
data Custom {
  on iterresp(next) { # no this
    s = ''
    for next() as iteration {
      s ..= iteration.values[0]
    }
    return Custom(s)
  }
}
```

## Control Flow

### Pattern Matching

```python
match var {
	1: 'hello'
  2: 'world'
  _: ''
}

match (i%3, i%5) {
	(0, 0): echo 'FizzBuzz'
  (0, _): echo 'Fizz'
  (_, 0): echo 'Buzz'
  (_, _): echo i
}

match type(x) {
	Number: echo '$x is a number'
	String: echo '$x is a string'
		   _: echo '$x is an unsupported type'
}

match x {
	valid(x): echo '$x is valid'
  _: echo '$x is not valid'
}
```

### Conditional

As a general recomentation, ifs are not recommended, but they are used as follows:

```python
if <expr> {
	...
} else if {
  ...
} else {
  ...
}
```

However, ifs are important for early returns:

```python
if <condition> return <returned value>
if <condition> raise <returned value>
```

With’s also help us addressing conditionals:

```python
if client.conn() as conn {
	...
} else {
	...
}
```

### Loops

Loops are also disencouraged, but may be useful in some situations:

```python
for i in <iterator> {
  break
  continue
}

for <condition> {}
for <expr> as <cond> {} # this is similar to while x = <expr> in other languages
for {}
```

## Error handling

We don’t have try catch since all errors are captured by default, instead, most builtin operations that can result in an error, will return a Maybe object that contains the value or the error:

```python
file = fs.open('invalid')
type(file) == Maybe # failing or not, the return is a Maybe
```

Maybes should be unwrapped in order to retrieve its value:

```python
# file.read() # Fails because Maybe doesn't have read()
file.unwrap() # Now it becomes the value OR the error
file.error # Access the error
file.value # Access the value, false if error
```

Maybes have a shortcut that help us treating it:

```python
if file! as err {
	echo 'Error'
  return 0
}

os.open()! # Ignoring the effect of the function
```

## Async

Using go routines, the usage is similar:

```python
let a = 1
const b = 2 
fn iwillrunasync(y, z) {
	# echo a # would result in an error, async functions cannot access variables outside its scope
	echo b # but constants are ok!
	return y + z
}

promise = async iwillrunasync(a, b) # the values are always copied
value = await promise # blocks the thread waiting the return

# Check how go does it
```

## Packages

Since the goal os to write a language to be used in bash, modules should be defined by file level.

```python
module math

const pi = 3.1415
```

 > sht -r ./math.sht main.sht

will register module, but wont execute it. only non modules are executed.
