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


## Numbers as Bits

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
    

Using these operators we can do pretty complex logic. In our case we'll use them to set/unset bits and check is a bit is set.

## Back to Permissions
...

