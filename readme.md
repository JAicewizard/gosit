# Gosit
Gosit is a posit library for go. There is currently one other go implementation but this one has been archived and wont receive updates.

Posits are a representation of real numbers that tries to compete with floating point numbers. See [the wikipedia page][positWiki] and [the paper][positPaper] if you want to learn more about posits.

## Testing
Gosit is fuzzed against goposit, the fuzzing finds a case where the bits returned form this library are diferent from goposit, it fails. As of the initial commit all functions have been fuzzed at least 1 hout, or around 26000 cases. This is not a lot and there might be some edge cases it hasnt found yet, but it gives a good indication that in general this library does the right thing. Only ES=2 has been fuzzed, mostly because this is the ES used in softposit for julia and rust, allowing me to easily verify and debug with a maintained version of softposit. Other values will be tested later on.

## Suport
Gosit currently only supports 32 bit posits, this should be enough for most use cases. 64 bit, and probably some smaller sizes, are on the roadmap.
Currently only es=2 is fuzzed as states in the [Testing](#Testing) section.

## How to use.

```go
p := gosit.Newposit32FromBits(0b0_001_10_11_10000011_01111110_10010111)
q := gosit.Newposit32FromBits(0b0_000001_11_0000001_00000111_10100011)
fmt.Println(gosit.Getfloat(p), "*", gosit.Getfloat(q), "=", gosit.Getfloat(gosit.MulPositSameES(p,q)))
```

## Benchmarks

Currently I only benchmark against the slow goposit, since its the only other go library for posits.
### Goposit

All tests are ran with the exect same bench cases to eliminate favouring one library over another because of coincidence.
These are rotated out every iteration. No other code is ran besides getting the case from an array, and running the corresponding function on it.
All benches ware using ES=2

```
go test --run=X --bench=. -benchtime 30s
BenchmarkAddSlow-10             1000000000              24.33 ns/op
BenchmarkAddSlowGoposit-10       9959919              3691 ns/op
BenchmarkMulSlow-10             1000000000              24.30 ns/op
BenchmarkMulSlowGoposit-10       9517118              3746 ns/op
BenchmarkDivSlow-10             1000000000              28.83 ns/op
BenchmarkDivSlowGoposit-10       9478204              3688 ns/op
```
### softposit-rs

No direct comparisons exist, but taking some averages from cargo bench on my machine:
| Operation | ns/op on my machine | 
|:---------:|:-------------------:|
|    Add    | 7.5ns/op            | 
|    -      | 6.9ns/op            |
|    *      | 6.3ns/op            |
|    /      | 9.8ns/op            |

[positWiki]:https://en.wikipedia.org/wiki/Unum_(number_format)#Unum_III 
[positPaper]:http://www.johngustafson.net/pdfs/BeatingFloatingPoint.pdf 