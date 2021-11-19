# Compile
```
go build machine_schedule.go
```
# Execute
```
./machine_schedule -i input_file_name
```
## Input File Format
```
machine_num;job_num;job0_start_time,job0_end_time;job1_start_time,job2_start_time
```
## example
### input.txt
3;6;0,3;4,5;1,3;2,4;5,6;2,4
### output
```
machine[2]
        job[3] 2  4
machine[1]
        job[2] 1  3
        job[1] 4  5
machine[0]
        job[0] 0  3
        job[4] 5  6
job[5][2][4] miss 
tatal job num[6] schedule[5] miss[1]
```

# Test1
## exec
```
./machine_schedule -i test_input_1.txt
```
## input file
```
3;6;0,3;4,5;1,3;2,4;5,6;7,8
```
## output
```
machine#0 has jobs:
(0,3)
(5,6)
machine#1 has jobs:
(1,3)
(4,5)
machine#2 has jobs:
(2,4)
(7,8)
[6] number of jobs out of [6] total jobs
```
# Test2
## exec
```
./machine_schedule -i test_input_2.txt
```
## input
```
2;4;0,2;1,4;1,7;2,8
```
## output
```
machine#0 has jobs:
(0,2)
machine#1 has jobs:
(1,4)
Jobs not processed:
(1,7)
(2,8)
[2] number of jobs out of [4] total jobs
```
# Test3
## exec
```
./machine_schedule -i test_input_3.txt 
```
## input
```
2;7;0,2;1,4;4,7;100,103;101,113;104,110;0,8
```
## output
``` 
machine#0 has jobs:
(0,2)
(4,7)
(104,110)
machine#1 has jobs:
(1,4)
(100,103)
Jobs not processed:
(0,8)
(101,113)
[5] number of jobs out of [7] total jobs
```

# Test4
## exec
```
./machine_schedule -i test_input_4.txt 
```
## input
```
3;15;0,2;0,1;1,4;3,4;5,6;0,1;7,8;5,6;1,5;0,4;5,6;1,12;0,2;11,15;11,15
```
## output
```
machine#0 has jobs:
(0,1)
(5,6)
(11,15)
machine#1 has jobs:
(0,1)
(3,4)
(5,6)
(7,8)
machine#2 has jobs:
(0,2)
(5,6)
(11,15)
Jobs not processed:
(0,2)
(1,4)
(0,4)
(1,5)
(1,12)
[10] number of jobs out of [15] total jobs
```

