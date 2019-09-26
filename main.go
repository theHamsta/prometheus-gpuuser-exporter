package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

// name, index, temperature.gpu, utilization.gpu,
// utilization.memory, memory.total, memory.free, memory.used

func metrics(response http.ResponseWriter, request *http.Request) {
	out, err := exec.Command("sreport",
		"user",
		"top",
		"--parsable",
		"--tres=gres/gpu",
		"topcount=50",
		"--noheader").Output()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	csvReader := csv.NewReader(bytes.NewReader(out))
	csvReader.TrimLeadingSpace = true
	csvReader.Comma = '|'
	records, err := csvReader.ReadAll()

	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	metricList := []string{
		"Cluster",
		"Login",
		"proper.name",
		"account",
		"tres.name",
		"gpu.minutes"}

	accountIdx := 3
	nameIdx := 2
	gpuMinutesIdx := 5

	result := ""
	for _, row := range records {
		fmt.Printf("%s\n", row[1])
		for idx, metric := range metricList {
			if idx == gpuMinutesIdx {
				result = fmt.Sprintf("%s%s{user=\"%s\", account=\"%s\"} %s\n",
					result,
					metric,
					row[nameIdx],
					row[accountIdx],
					row[idx])
			}
		}

	}

	fmt.Fprintf(response, strings.Replace(result, ".", "_", -1))
}

func main() {
	addr := ":9104"
	if len(os.Args) > 1 {
		addr = ":" + os.Args[1]
	}

	http.HandleFunc("/metrics/", metrics)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
