
# Code refactoring
Task description is [here](task.md)  
I will reduce generated users number from 100 to 10, and logs lines from 1000 to 10 - at the development stage I do not want to generate many big files.  
Remove generated files command: "rm -fv users/*.txt"  

## Performance improvement steps
Testing on a computer with:    
goos: darwin  
goarch: amd64  
cpu: Core i9 2.30GHz  
cores: 8  

### 0. Source
```
Time 		~ 11.09 seconds
Quantity 	1
Time 		11056139596 ns/op
Size 		44336 B/op
Allocs 		351 allocs/op
```
### 1. Concurrency added
- user generate and write them into files - in several streams  
- generate all user logs - in one time as [][]logItem
Result:
```
Time 		~ 1.11 seconds
Quantity 	1
Time 		1107507368 ns/op
Size 		50728 B/op
Allocs 		501 allocs/op
```
### 2. String concatenation minimize
- bytes.Buffer + buffer.WriteString instead of fmt.Printf()  
- file.Write([]byte) instead of file.WriteString(string)
Result:
```
Time 		~ 1.11 seconds
Quantity 	1
Time 		1107307591 ns/op
Size 		40560 B/op
Allocs 		371 allocs/op
```
