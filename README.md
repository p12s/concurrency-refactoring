
# Code refactoring
Task description is [here](task.md)  
I will reduce generated users number from 100 to 10, and logs lines from 1000 to 10 - at the development stage I do not want to generate many big files.  
And I didn't remove time.Sleep() which simulate heavy operations.  
Remove generated files command: "make clean"  

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
### 2. Micro-optimizations
- String concatenation minimize: bytes.Buffer + buffer.WriteString instead of fmt.Srintf()  
- file.Write([]byte) instead of file.WriteString(string)
- passing function parameters by reference (not by value)  
- remove variables you can do without  

Result:
```
Time 		~ 1.11 seconds (100 files - 1.13 seconds)
Quantity 	1
Time 		1107307591 ns/op
Size 		40560 B/op
Allocs 		371 allocs/op
```
### 3. Pipeline
- An attempting to wrap user processing in a pipeline (failure - worse results)

Result:
```
Time 		~ 1.17 seconds (100 files - 1.34 seconds)
Quantity 	1
Time 		1101704171 ns/op
Size 		173936 B/op
Allocs 		1374 allocs/op
```
