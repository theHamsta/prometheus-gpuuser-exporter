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

func topusers(response http.ResponseWriter, request *http.Request) {
	out, err := exec.Command(
		"sreport", "user", "top", "--parsable", "--tres=gres/gpu", "topcount=50", "--noheader",
	).Output()

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

	result := make([]string, len(records))

	for rowIdx, row := range records {
		fmt.Printf("%s\n", row[1])
		result[rowIdx] = fmt.Sprintf("%s%s{user=\"%s\", account=\"%s\"} %s",
			result,
			metricList[gpuMinutesIdx],
			row[nameIdx],
			row[accountIdx],
			row[gpuMinutesIdx])
	}

	fmt.Fprintf(response, strings.Replace(strings.Join(result[:], "\n"), ".", "_", -1))
}

func main() {
	addr := ":9104"
	if len(os.Args) > 1 {
		addr = ":" + os.Args[1]
	}

	http.HandleFunc("/topusers/", topusers)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
