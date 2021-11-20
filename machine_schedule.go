package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

//job info
type Job struct {
	JobId     int
	StartTime int
	EndTime   int
	MachineId int
}

//machine info
type MachineInfo struct {
	MachineId  int
	FinishTime int
	JobList    []Job
}

type MachineSchedule struct {
}

//sort the job list by end time
func (s *MachineSchedule) SortJobList(jobList []Job) error {
	sort.SliceStable(jobList, func(i, j int) bool {
		if jobList[i].EndTime < jobList[j].EndTime {
			return true
		}
		return false
	})
	return nil
}

//sort the machine list by finish time
func (s *MachineSchedule) SortMachineListByFinishTime(machineList []MachineInfo) error {
	sort.SliceStable(machineList, func(i, j int) bool {
		if machineList[i].FinishTime < machineList[j].FinishTime {
			return true
		}
		return false
	})
	return nil
}

//sort the machine list by machine id
func (s *MachineSchedule) SortMachineListByMachineId(machineList []MachineInfo) error {
	sort.SliceStable(machineList, func(i, j int) bool {
		if machineList[i].MachineId < machineList[j].MachineId {
			return true
		}
		return false
	})
	return nil
}

//input machine num and job list, return machine schedule list
func (s *MachineSchedule) Run(machineNum int, jobList []Job) (machineList []MachineInfo, err error) {
	//sort jobs by end time, coat logn
	_ = s.SortJobList(jobList)

	//init machine list
	machineList = make([]MachineInfo, machineNum)
	for i := 0; i < machineNum; i++ {
		machineList[i].MachineId = i
		machineList[i].FinishTime = -1
	}

	//for all the job
	for i, job := range jobList {
		jobList[i].MachineId = -1
		//get the first free machine
		for j, machine := range machineList {
			if machine.FinishTime < job.StartTime {
				machineList[j].JobList = append(machineList[j].JobList, job)
				machineList[j].FinishTime = job.EndTime
				jobList[i].MachineId = machine.MachineId
				//restore the machine
				s.SortMachineListByFinishTime(machineList)
				break
			}
		}
	}
	//finish. return the machine list bye machine id.
	s.SortMachineListByMachineId(machineList)
	return machineList, nil
}

type Cmd struct {
	Help       bool
	InputFile  string
	OutputFile string
}

func (c *Cmd) Usage() {
	fmt.Fprintf(os.Stderr, `Usage: ./machine_schedule [-h] [-i filename] 
	Options:
	`)
	flag.PrintDefaults()
}

var cmd Cmd

//init
func init() {

	flag.BoolVar(&cmd.Help, "h", false, "this help")
	flag.StringVar(&cmd.InputFile, "i", "./input.txt", "input file name")
	flag.Usage = cmd.Usage
}

func main() {
	flag.Parse()
	if cmd.Help || len(cmd.InputFile) == 0 {
		flag.Usage()
		return
	}

	//read data from file
	file, err := os.Open(cmd.InputFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err == io.EOF {
		fmt.Printf("read error:$s", err.Error())
		return
	}

	//first line is machine num
	machineNum, _ := strconv.Atoi(string(line))

	//other lines are jobs
	var jobList []Job
	jobNum := 0
	for {
		line, _, err = reader.ReadLine()
		if err == io.EOF {
			break
		}
		job := Job{
			JobId: jobNum,
		}
		data := strings.Trim(string(line), "\n")
		jobData := strings.Split(data, " ")
		if len(jobData) != 2 {
			fmt.Printf("job[%d][%s] format error", job.JobId, data)
			continue
		}
		job.StartTime, _ = strconv.Atoi(jobData[0])
		job.EndTime, _ = strconv.Atoi(jobData[1])
		jobList = append(jobList, job)
		jobNum++
	}

	//run and get the result
	var machineSchedule MachineSchedule
	machineList, err := machineSchedule.Run(machineNum, jobList)
	if err != nil {
		fmt.Println(err)
		return
	}

	//print the machine & jobs
	scheduleNum := 0
	for _, machine := range machineList {
		fmt.Printf("machine#%d has jobs:\n", machine.MachineId)
		for _, job := range machine.JobList {
			scheduleNum++
			fmt.Printf("(%d,%d)\n", job.StartTime, job.EndTime)
		}
	}

	//print the missed jobs
	hasPrintHead := false
	for _, job := range jobList {
		if job.MachineId < 0 {
			if !hasPrintHead {
				fmt.Printf("Jobs not processed:\n")
				hasPrintHead = true
			}
			fmt.Printf("(%d,%d)\n", job.StartTime, job.EndTime)
		}
	}
	//result
	fmt.Printf("[%d] number of jobs out of [%d] total jobs\n", scheduleNum, jobNum)
}
