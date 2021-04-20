# Using Bitmasks In Go

### Introduction

You write a server for a massively multiplayer online role-playing game ([MMORPG](https://en.wikipedia.org/wiki/Massively_multiplayer_online_role-playing_game)).

In the game, players collect keys and you want to design how to store the set of keys each player has.

As an example, imagine the set of keys are `copper`, `jade` and `crystal`. You consider the following options for storing a player key set:

- `[]string`
- `map[string]bool`

Both options will work, but did you consider a third option of using a [bitmask](https://en.wikipedia.org/wiki/Mask_(computing))? Using a bitmask will make storing and processing keys more efficient. Once you learn the mechanics, it will be readable and maintainable as well.

### Interlude: Numbers as Bits

Let’s start off with the understanding that computers store numbers as a sequence of 8 bits called a byte, where every byte represents a number.
How is a base 2 number converted into a base 10 number?

Every bit position represents a [power of two](https://en.wikipedia.org/wiki/Power_of_two).
Let’s look at a byte, which is 8 bits long:

**Listing 1** 
```
| 2⁷| 2⁶| 2⁵| 2⁴| 2³| 2²| 2¹| 2⁰|  <- Bit Position
|---|---|---|---|---|---|---|---|
|128| 64| 32| 16|  8|  4|  2|  1|  <- Base 10 Value
```

The rightmost bit (known as the least significant bit, or LSB) represents 2⁰ or in base 10 the number 1, the second bit represents 2¹ or 2, and so fourth until we reach the leftmost bit (called MSB = most significant bit) which representsis 2⁷ or 128.

For example, to represent the number 13, we break it down to setting bit positions for the numbers 8, 4, and 1.

**Listing 2** 
```
| 0 | 0 | 0 | 0 | 1 | 1 | 0 | 1 |  <- Bit Position
|---|---|---|---|---|---|---|---|
|128| 64| 32| 16| 8 | 4 | 2 | 1 |  <- Base 10 Value
=================================
 0+  0+  0+  0+  8+  4+  0+  1    = 13
```

_Note: You can use the %b verb to print the binary representation of a number.
`fmt.Printf("%08b\n", 13)` will print `00001101`._

This encoding scheme means the maximal number that can be represented by a byte is 255.

**Listing 3** 
```
| 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |  <- Bit Position
|---|---|---|---|---|---|---|---|
|128| 64| 32| 16| 8 | 4 | 2 | 1 |  <- Base 10 Value
=================================
128+ 64+ 32+ 16+ 8+  4+  2+  1     = 255
```

_Note: A fun exercise with kids is to teach them to count up to 31 on one hand. Each finger is a bit, where the pinky is 1 and the thumb is 16. Let them giggle at 4 :)_

We can also perform logical operations on bits.

AND (`&`) is true only if both bits are 1:

**Listing 4** 
```
0 & 0 -> 0 (false)
0 & 1 -> 0 (false)
1 & 0 -> 0 (false)
1 & 1 -> 1 (true)
```

OR (`|`) is true if one of the bits is 1:

**Listing 5** 
```
0 | 0 -> 0 (false)
0 | 1 -> 1 (true)
1 | 0 -> 1 (true)
1 | 1 -> 0 (true)
```

NOT (`^`) reverses the bit:

**Listing 6** 
```
^1 -> 0
^0 -> 1
```

These operators work on more than one bit at a time. To calculate which bits are common between the numbers 5 AND 3 do the following.

**Listing 7** 
```
00000101 AND  (4, 1)
00000011      (2, 1)
--------------------
00000001      (1)
```

To join the bits between the numbers 5 OR 3 do the following.

**Listing 8** 
```
00000101 OR  (4, 1)
00000011     (2, 1)
-------------------
00000111     (4, 2, 1)
```

To multiply a number by 2, SHIFT all the bits one position to the left using the `<<` operator.

**Listing 9** 
```
00001010 (10) << 1
------------------
00010100 (20)
```

There’s also a right SHIFT operator (`>>`) which divides a number by 2.

**Listing 10** 
```
00010100 (20) >> 1
------------------
00001010 (10)
```

Using these operators we can perform complex logic. In our case, we'll use them to set/unset bits and check if a bit is set.

### Back To  Our Problem

To support 3 keys in the application, we only need 3 bits. This is great because that means we only need to allocate 1 byte of memory.

**Listing 11**
```
8 // KeySet is a set of keys in the game.
9 type KeySet byte
```

On line 11, we define the `KeySet` type using an underlying type of `byte`. A byte in Go is an alias for `uint8`.

**Listing 12**
```
11 const (
12     Copper  KeySet = 1 << iota // 1
13     Jade                       // 2
14     Crystal                    // 4
15     maxKey                     // 8
16 )
```

Listing 12 shows the available keys we will support. On line 11, we start to define our keys by shifting the value of `iota` to the left by 1. When the `iota` keyword is used in a constant group, the compiler will automatically apply the formula for each subsequent line and increment `iota` by one.

In order to have nice string representation for our keys, we'll have it implement the [fmt.Stringer](https://golang.org/pkg/fmt/#Stringer) interface.

**Listing 13**
```
18 // String implements the fmt.Stringer interface
19 func (k KeySet) String() string {
20     if k >= maxKey {
21         return fmt.Sprintf("<unknown key: %d>", k)
22     }
23
24     switch k {
25     case Copper:
26         return "copper"
27     case Jade:
28         return "jade"
29     case Crystal:
30         return "crystal"
31     }
32
33     // multiple keys
34     var names []string
35     for key := Copper; key < maxKey; key <<= 1 {
36         if k&key != 0 {
37             names = append(names, key.String())
38         }
39     }
40     return strings.Join(names, "|")
41 }
```

Listing 13 shows the `String` method implementation. On line 20, we check that the value is valid. On lines 24 to 31, we return string value for a single key (bit) and on lines 34 to 40, we construct a string representation for multiple keys (bits).

Now we're ready to use `KeySet` in our `Player` struct.

**Listing 14**
```
43 // Player is a player in the game
44 type Player struct {
45     Name string
46     Keys KeySet
47 }
```

Listing 14 shows the `Player` struct implementation. On line 45, we have the player name and on line 46, we have the set of keys it holds. As the game is developed, we'll add more fields.

**Listing 15**
```
49 // AddKey adds a key to the player keys
50 func (p *Player) AddKey(key KeySet) {
51     p.Keys |= key
52 }
```

Listing 15 shows how to add a key to the bitmask. On line 51, we use bitwise OR to set the `key` bit in the `Keys` field.

How does the function work? If the KeySet already has Copper and we want to add Crystal, we can pass Crystal to the AddKey method and the OR operation will do the rest.

**Listing 16**
```
p.Keys : 00000001 OR  (Copper)
key    : 00000100     (Crystal)
---------------------------------------
result : 00000101     (Copper, Crystal)
```

We can see the resulting bit pattern includes the Crystal bit.

**Listing 17**
```
54 // HasKey returns true if player has a key
55 func (p *Player) HasKey(key KeySet) bool {
56     return p.Keys & key != 0
57 }
```

Listing 17 shows how to check for a key in the KeySet. On line 56, we use bitwise AND to check if the key bit is set in the `Keys` field.

How does it work? If the KeySet already has Copper and Crystal, and we want to check if Crystal exists, we can pass Crystal to the HasKey method and the AND operation will do the rest.

**Listing 18**
```
p.Keys : 00000101 AND  (Copper, Crystal)
key    : 00000100      (Crystal)
----------------------------------------
result : 00000100      (Crystal)
```
We can see the resulting bit pattern includes the Crystal bit, so there is a match.

On the other hand, when we check for Jade we get a different result.

**Listing 19**
```
p.Keys : 00000101 AND  (Copper, Crystal)
key    : 00000010      (Jade)
----------------------------------------
result : 00000000      Nothing
```

We can see the resulting bit pattern doesn't include the Jade bit, so there is no match.

**Listing 20**
```
59 // RemoveKey removes key from player
60 func (p *Player) RemoveKey(key KeySet) {
61     p.Keys &= ^key
62 }
```

Listing 20 shows removing a key. On line 61, we first use bitwise NOT to flip the `key` bits and then use bitwise AND to unset the key bit in the `Keys` field.

How does it work? If the KeySet already has Copper and Crystal, and we want to remove Crystal, we can pass Crystal to the RemoveKey method and the use of the NOT and AND operations will do the rest.

**Listing 21**
```
p.Keys : 00000101      (Copper, Crystal)
^key   : 11111011 AND  (org: 00000100 Crystal)
----------------------------------------
result : 00000001      (Copper)
```

We can see the resulting bit pattern doesn't include the Crystal bit anymore.

## Conclusion

Go's type system allows you to combine low level code such as bitmasks with high level code such as methods to give you both performance and user friendly code.

How much did we save? I wrote [a benchmark](https://github.com/353words/bitmask/blob/master/bench_test.go) that checks the three approaches: using a `[]string`, using a `map[string]bool` and using our `byte` based implementation.

**Listing 22**
```
$ go test -bench . -benchmem
goos: linux
goarch: amd64
pkg: github.com/353words/bitmask
cpu: Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz
BenchmarkMap-4       249177445            4.800 ns/op           0 B/op          0 allocs/op
BenchmarkSlice-4     243485120            4.901 ns/op           0 B/op          0 allocs/op
BenchmarkBits-4      1000000000           0.2898 ns/op          0 B/op          0 allocs/op
BenchmarkMemory-4    21515095            52.25 ns/op       32 B/op          1 allocs/op
PASS
ok   github.com/353words/bitmask 4.881s
```

Listing 20 shows the results of the benchmark on my machine. On line 1, we run the benchmark with the `-benchmem` flag to benchmark memory allocations. On lines 06-09, we see that our `byte` based implementation is about 16 times faster than the alternatives. On line 09, we see that allocating `[]string{"copper", "jade"}` consumes 32 bytes, which is 32 times more memory than our single byte implementation.

_Note: You should optimize *only* after you have performance requirements and you have profiled your application. See the [Rules Of Optimization Club](https://wiki.c2.com/?RulesOfOptimizationClub) for more great advice._

What tricks did you use to reduce memory? I'd love to hear your stories, ping me at miki@353solutions.com
The code for this blog post can be found [on GitHub](https://github.com/353words/bitmask).


### Using Bitmasks

You write a server for a massively multiplayer online role-playing game ([MMORPG](https://en.wikipedia.org/wiki/Massively_multiplayer_online_role-playing_game)).

In the game, players collect keys and you want to design how to store the set of keys each player has.

As an example, imagine the set of keys are `copper`, `jade` and `crystal`. You consider the following options for storing a player key set:

- A `[]string`
- A `map[string]bool`

Both options will work, but did you consider a third option of using a [bitmask](https://en.wikipedia.org/wiki/Mask_(computing))? Using a bitmask will make storing and processing keys more efficient. Once you learn the mechanics, it will be readable and maintainable as well.

Let’s see how.

### Interlude: Numbers as Bits

Let’s start off with the understanding that computers store numbers as a sequence of 8 bits called a byte, where every byte represents a number.
How is a base 2 number converted into a base 10 number?

Every bit position represents a [power of two](https://en.wikipedia.org/wiki/Power_of_two). 
Let’s look at a byte, which is 8 bits long:

| 2⁷| 2⁶| 2⁵| 2⁴| 2³| 2²| 2¹| 2⁰|
|---|---|---|---|---|---|---|---|
|128| 64| 32| 16|  8|  4|  2|  1|

The rightmost bit (known as the least significant bit, or LSB) represents 2⁰ or in base 10 the number 1, the second bit represents 2¹ or 2, and so fourth until we reach the leftmost bit (called MSB = most significant bit) which representsis 2⁷ or 128.

For example, to represent the number 13, we break it down to 8 + 4 + 1 and then set the corresponding bits to 1:

| 0 | 0 | 0 | 0 | 1 | 1 | 0 | 1 |
|---|---|---|---|---|---|---|---|
|128| 64| 32| 16|  8|  4|  2|  1|

_Note: You can use the %b verb to print the binary representation of a number. 
`fmt.Printf("%08b\n", 13)` will print `00001101`._

This encoding scheme means the maximal number that can be represented by a byte is 

| 1 | 1 | 1 | 1 | 1 | 1 | 1 | 1 |
|---|---|---|---|---|---|---|---|
|128| 64| 32| 16|  8|  4|  2|  1|

Which in base 10 means: 128 + 64 + 32 + 16 + 8 + 4 + 2 + 1 = 255.

_Note: A fun exercise with kids is to teach them to count up to 31 on one hand. Each finger is a bit, where the pinky is 1 and the thumb is 16. Let them giggle at 4 :)_

We can perform logical operations on bits:

AND (`&`) is true only if both bits are 1: 
```
0 & 0 -> 0 (false)
0 & 1 -> 0 (false)
1 & 0 -> 0 (false)
1 & 1 -> 1 (true)
```

OR (`|`) is true if one of the bits is 1:
```
0 | 0 -> 0 (false)
0 | 1 -> 1 (true)
1 | 0 -> 1 (true)
1 | 1 -> 0 (false)
```

NOT (`^`) reverses the bit:
```
^1 -> 0
^0 -> 1
```

These operators work on more than one bit at a time. To calculate 5&3 do the following:

- Represent them as bits: 5 -> 00000101, 3 -> 00000011
- Calculate AND between bits in the same position

```
00000101
00000011
--------
00000001
```

Which means that `5&3 = 1`.

SHIFT LEFT (`<<`) moves all the bits one position to the left

```
00001010 << 1 -> 00010100
    10   << 1 ->    20
```

Shifting left one place is like multiplying by 2.

There’s also a SHIFT RIGHT operator (`>>`) which moves all the bits on position to the right.

Using these operators we can perform complex logic. In our case, we'll use them to set/unset bits and check if a bit is set.

### Back To  Our Problem

To support 3 keys in the application, we only need 3 bits. This is great because that means we only need to allocate 1 byte of memory.

**Listing 1: KeySet type**
```
8 // KeySet is a set of keys in the game.
9 type KeySet byte
```

On line 9, we define the `KeySet` type using an underlying type of `byte`. A byte in Go is an alias for `uint8`.

**Listing 2: The Keys**
```
11 const (
12     Copper  KeySet = 1 << iota // 1
13     Jade                    // 2
14     Crystal                 // 4
15     maxKey
16 )
```

Listing 2 shows the available keys we will support. On line 11, we start to define our keys by shifting the value of `iota` to the left by 1. When the `iota` keyword is used in a constant group, the compiler will automatically apply the formula for each subsequent line and increment `iota` by one.

Listing 2 shows the available keys we will support. On line 11, we start to define our keys using a constant block. When the `iota` keyword is used in a constant group, the compiler will automatically apply the formula for each subsequent line and increment `iota` by one, starting `iota` at zero.. On line 12, we define the value for the first key as `1 << iota` which is `1 << 0 = 1`. On line 13, `Jade` will have the value of `1 << 2 = 2`. On line 15, we define `maxKey` which is the maximal key value, this value is not exported.

In order to have nice string representation for our keys, we'll have it implement the [fmt.Stringer](https://golang.org/pkg/fmt/#Stringer) interface.

**Listing 2: String() implementation**
```
18 // String implements the fmt.Stringer interface
19 func (k KeySet) String() string {
20     if k >= maxKey {
21         return fmt.Sprintf("<unknown key: %d>", k)
22     }
23 
24     switch k {
25     case Copper:
26         return "copper"
27     case Jade:
28         return "jade"
29     case Crystal:
30         return "crystal"
31     }
32 
33     // multiple keys
34     var names []string
35     for key := Copper; key < maxKey; key <<= 1 {
36         if k&key != 0 {
37             names = append(names, key.String())
38         }
39     }
40     return strings.Join(names, "|")
41 }

```

Listing 3 shows the `String` method implementation. On line 20, we check that the value is valid. On lines 24 to 31, we return string value for a single key (bit) and on lines 34 to 40, we construct a string representation for multiple keys (bits).

Now we're ready to use `KeySet` in our `Player` struct.

**Listing 4: Player struct**
```
43 // Player is a player in the game
44 type Player struct {
45     Name string
46     Keys KeySet
47 }
```

Listing 4 shows the `Player` struct implementation. On line 45, we have the player name and on line 46, we have the set of keys it holds. As the game is developed, we'll add more fields.

**Listing 5: Adding a key**
```
49 // AddKey adds a key to the player keys
50 func (p *Player) AddKey(key KeySet) {
51     p.Keys |= key
52 }
```

Listing 5 shows adding a key. On line 51, we use bitwise OR to set the `key` bit in the `Keys` field.

How does it work? Say we have the Copper (1) key and we add the Crystal (4) key.
The `Keys` field before is: 00000001, then we do

```
00000001 OR
00000100
--------
00000101
```


**Listing 6: Checking for a Key**
```
54 // HasKey returns true if player has a key
55 func (p *Player) HasKey(key KeySet) bool {
56     return p.Keys&key != 0
57 }
```

Listing 6 shows checking for a key. On line 56, we use bitwise AND to check if the key bit is set in the `Keys` field.

How does it work? Say we have Copper (1) and Crystal (4), then the `Keys` field is: 00000101. We’d like to check if we have the Crystal key then we do

```
00000101 AND
00000100
--------
00000100
```

The result does not equal 0. On the other hand, when we check for Jade (2), then we do:
```
00000101 AND
00000010
--------
00000000
```

Which does equal to 0.

**Listing 7: Removing a Key**
```
59 // RemoveKey removes key from player
60 func (p *Player) RemoveKey(key KeySet) {
61     p.Keys &= ^key
62 }
```

Listing 7 shows removing a key. On line 61, we first use bitwise NOT to flip the `key` bits and then use bitwise AND to unset the key bit in the `Keys` field.

How does it work? Say we have Copper (1) and Crystal (4), then the `Keys` field is: 00000101. We’d like to remove the Crystal key then first apply bitwise NOT on the Crystal key

```
00000100 NOT
--------
11111011
```

And then do a bitwise AND with the current values of `Keys`:

The result does not equal 0. On the other hand, when we check for Jade (2), then we do:
```
00000101 AND
11111011
--------
00000001
```

## Conclusion

Go's type system allows you to combine low level code such as bitmasks with high level code such as methods to give you both performance and user friendly code.

How much did we save? I wrote [a benchmark](https://github.com/353words/bitmask/blob/master/bench_test.go) that checks the three approaches: using a `[]string`, using a `map[string]bool` and using our `byte` based implementation.

```
01 $ go test -bench . -benchmem
02 goos: linux
03 goarch: amd64
04 pkg: github.com/353words/bitmask
05 cpu: Intel(R) Core(TM) i7-7500U CPU @ 2.70GHz
06 BenchmarkMap-4      	249177445	         4.800 ns/op	       0 B/op	       0 allocs/op
07 BenchmarkSlice-4    	243485120	         4.901 ns/op	       0 B/op	       0 allocs/op
08 BenchmarkBits-4     	1000000000	         0.2898 ns/op	       0 B/op	       0 allocs/op
09 BenchmarkMemory-4   	21515095	        52.25 ns/op	      32 B/op	       1 allocs/op
10 PASS
11 ok  	github.com/353words/bitmask	4.881s
```

Listing 8 shows the results of the benchmark on my machine. On line 1, we run the benchmark with the `-benchmem` flag to benchmark memory allocations. On lines 06-09, we see that our `byte` based implementation is about 16 times faster than the alternatives. On line 09, we see that allocating `[]string{"copper", "jade"}` consumes 32 bytes, which is 32 times more memory than our single byte implementation.

_Note: You should optimize *only* after you have performance requirements and you have profiled your application. See the [Rules Of Optimization Club](https://wiki.c2.com/?RulesOfOptimizationClub) for more great advice._

What tricks did you use to reduce memory? I'd love to hear your stories, ping me at miki@353solutions.com
The code for this blog post can be found [on GitHub](https://github.com/353words/bitmask).
