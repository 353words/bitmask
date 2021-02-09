# Using Bitmasks to Reduce Memory

You're writing an online document system that will rival Google Docs and Microsoft Office.
Business is going very well and you get a lot of customer, and you start facing performance issues.
You're servers start to consume a lot of memory, and you start to dig in and see how you can reduce memory consumption.

After some investigation, you find out that the document cache is one of the main memory hogs. 
Each document has a set of permissions defined as:

**Listing 1: DocumentPermissions struct**

```
8 // DocumentPermissions are permissions on a document.
9 type DocumentPermissions struct {
10     Locked        bool
11     GroupReadable bool
12     GroupWritable bool
13     AnyReadable   bool
14     AnyWritable   bool
15 }
```

Listing 1 shows the permission set each document has.

Let's see how much memory `DocumentPermissions` takes:

**Listing 2: Size of `DocumentPermissions`**

```
17 func main() {
18     perms := DocumentPermissions{}
19     fmt.Println(unsafe.Sizeof(perms))
20 }
```

Listing 2 shows how to measure the size of `DocumentPermissions` in memory. This program will output `5`, which is 5 bytes. Each byte is 8 bits so we use `5*8=40` bits to store this information.

Ideally, every permission flag can use a single bit, it can be either `true` (1) or `false` (0). We have 5 permissions, which mean we can use a single byte, instead of 5, to encode the same information.


## Interlude: Numbers as Bits

Compute store numbers as a sequence of bits. Every bit represents a number:

- First bit is 2⁰ = 1
- Second bit is 2¹ = 2
- Third bit is 2² = 4
- Fourth bit is 2³ = 8
- ...

To construct the number 6 we'll represent is at 4(2²) + 2(2¹).
For an 8 bit number, we're using the bit representation of `00000110`. There the smallest bit (2⁰ also called LSB - least significant bit) is on the right. We say a bit is "set" is it's value is 1, in 6 the second and third bits are set.

_Note: A fun exercise with kids is to teach them to count up to 31 on one hand. Each finger is a bit, where the pinky is 1 and the thumb is 16. Let them giggle at 4 :)_

We also define logic operations on bits:

AND (`&`) is true only if two bits are 1: 
- `0 & 0 -> 0`
- `0 & 1 -> 0`
- `1 & 0 -> 0`
- `1 & 1 -> 1`

OR (`|`) is true if one of the bits is 1:
- `0 | 0 -> 0`
- `0 | 1 -> 1`
- `1 | 0 -> 1`
- `1 | 1 -> 0`

_Note: The bitwise or is different than the "or" we use in English which means either the first bit is 1 or the second. The operator for the English "or" is known as XOR (short for eXclusive OR)._

NOT (`^`) negates the bit
- `^1 -> 0`
- `^0 -> 1`

These operators work on more than one bit at a time. To calculate 5&3 do the following:
- Represent them as bits: 5 -> 00000101, 3 -> 00000011
- Calculate AND between bits in the same position

```
00000101
00000011
--------
00000001
```

Which means that `5&3 = 0`.

SHIFT LEFT (`<<`) moves all bits one to the left.

```
00001010 << 1 -> 00010100
```

Shifting left one place is like multiplying by 2.
    

Using these operators we can do pretty complex logic. In our case we'll use them to set/unset bits and check is a bit is set.

## Back to Permissions

We're going to re-define `DocumentPermissions` as a single byte.

**Listing 3: `DocumentPermissions` as byte**

```
3 // DocumentPermissions are permissions set on a document
4 type DocumentPermissions uint8
```

Listing 3 shows the new type of `DocumentPermissions`.

Next, we're going to define the permissions. Each will be a number with one distinct bit that is 1 and all the rest 0.

**Listing 4: Permission values**

```
6  // Available permissions
7  const (
8      Locked DocumentPermissions = 1 << iota
9      GroupReadable
10     GroupWritable
11     AllReadable
12     AllWritable
13 )
```

Listing 4 show the permissions. One line 8, we use `iota` to define the first permission which is `1<<0 = 1`. On lines 8 to 12, Go will carry the same operation and type for the rest of the values. This means `GroupReadable` will be `1<<1 = 2`, `GroupWritable` will be `1<<2 = 4` ...

Next we'll define methods on `DocumentPermissions` to set/clear permissions and also to check if a permission is set.

**Listing 5: Set**

```
15 func (p *DocumentPermissions) Set(perm DocumentPermissions) {
16     *p = *p | perm
17 }
```

Listing 5 shows how to set a permission. On line 16, we use bitwise-or (`|`) to set a permission bit.

For example, if the current value of `p` is `00000010` and we'd like to set `Locked` which is `00000001` then

```
00000010
00000001
--------
000000011
```

**Listing 5: Clear**

```
19 func (p *DocumentPermissions) Clear(perm DocumentPermissions) {
20     *p = *p & (^perm)
21 }
```

Listing 5 shows how to clear a permission. On line 20, we first negate the bit (`^`) and then use bitwise and (`&`) with the current permission.

For example, if currently `Locked` and `GroupReadable` are set, then the value of `p` is `00000011`. We'd like to clear `Locked` so first we negate it:

```
00000001
--------
11111110
```

And then we do bitwise and with the current value of `p`:

```
00000011
11111110
--------
00000010
```

**Listing 6: IsSet**

```
23 func (p DocumentPermissions) IsSet(perm DocumentPermissions) bool {
24     return p&perm != 0
25 }
```

Listing 6 shows how to check if a permission is set. On line 24, we use a bitwise or (`|`) to check if a bit is set.

For example, if currently `Locked` and `GroupReadable` are set, then the value of `p` is `00000011`. Let's check against `Locked`:

```
00000011
00000001
--------
00000001
```

And if we check against `AnyReadable` which is `00001000` then

```
00000011
00001000
--------
00000000
```

## Conclusion

We managed to reduce the memory consumption of `DocumentPermissions` from 5 bytes to a single byte. It might not seem much, but as Dave Cheney [said](https://dave.cheney.net/2021/01/05/a-few-bytes-here-a-few-there-pretty-soon-youre-talking-real-memory) "A few bytes here, a few there, pretty soon you’re talking real memory."

What's nice about Go's type system is that it allows you to combine low level code such as bitmasks with high level code such as methods to give you both performance and user friendly code.

What tricks did you use to reduce memory? I'd love to hear your stories, ping me at miki@353solutions.com

The code for this blog post can be found [on GitHub](https://github.com/353words/bitmask).
