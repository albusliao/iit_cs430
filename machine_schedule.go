package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Job struct {
	JobId     int
	StartTime int
	EndTime   int
	MachineId int
}

type MachineInfo struct {
	MachineId  int
	FinishTime int
	JobList    []Job
}

type MachineSchedule struct {
}

func (s *MachineSchedule) SortJobList(jobList []Job) error {
	sort.SliceStable(jobList, func(i, j int) bool {
		if jobList[i].EndTime < jobList[j].EndTime {
			return true
		}
		return false
	})
	return nil
}

func (s *MachineSchedule) SortMachineListByFinishTime(machineList []MachineInfo) error {
	sort.SliceStable(machineList, func(i, j int) bool {
		if machineList[i].FinishTime < machineList[j].FinishTime {
			return true
		}
		return false
	})
	return nil
}
func (s *MachineSchedule) SortMachineListByMachineId(machineList []MachineInfo) error {
	sort.SliceStable(machineList, func(i, j int) bool {
		if machineList[i].MachineId < machineList[j].MachineId {
			return true
		}
		return false
	})
	return nil
}

//read input data from file
func (s *MachineSchedule) ReadFile() error {
	return nil
}

func (s *MachineSchedule) Run(machineNum int, jobList []Job) (machineList []MachineInfo, err error) {
	//sort jobs by end time, coat logn
	_ = s.SortJobList(jobList)

	//init machine list
	machineList = make([]MachineInfo, machineNum)
	for i := 0; i < machineNum; i++ {
		machineList[i].MachineId = i
		machineList[i].FinishTime = -1
	}

	for i, job := range jobList {
		jobList[i].MachineId = -1
		for j, machine := range machineList {
			if machine.FinishTime < job.StartTime {
				machineList[j].JobList = append(machineList[j].JobList, job)
				machineList[j].FinishTime = job.EndTime
				jobList[i].MachineId = machine.MachineId
				s.SortMachineListByFinishTime(machineList)
				break
			}
		}
	}
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
	content, err := ioutil.ReadAll(file)
	input := strings.Split(strings.Trim(string(content), "\n"), ";")
	if len(input) < 3 {
		fmt.Printf("[%s] format error", string(content))
		return
	}
	machineNum, _ := strconv.Atoi(input[0])
	jobNum, _ := strconv.Atoi(input[1])
	if jobNum != len(input)-2 {
		fmt.Printf("job num[%d] != job list num[%d]", jobNum, len(input)-2)
		return
	}

	var jobList []Job
	for i := 2; i < len(input); i++ {
		job := Job{
			JobId: i - 2,
		}
		jobData := strings.Split(input[i], ",")
		if len(jobData) != 2 {
			fmt.Printf("job[%d][%s] format error", job.JobId, input[i])
			return
		}
		job.StartTime, _ = strconv.Atoi(jobData[0])
		job.EndTime, _ = strconv.Atoi(jobData[1])
		jobList = append(jobList, job)
	}

	var machineSchedule MachineSchedule
	machineList, err := machineSchedule.Run(machineNum, jobList)
	if err != nil {
		fmt.Println(err)
		return
	}

	scheduleNum := 0
	for _, machine := range machineList {
		fmt.Printf("machine#%d has jobs:\n", machine.MachineId)
		for _, job := range machine.JobList {
			scheduleNum++
			fmt.Printf("(%d,%d)\n", job.StartTime, job.EndTime)
		}
	}

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
	fmt.Printf("[%d] number of jobs out of [%d] total jobs\n", scheduleNum, jobNum)

}
