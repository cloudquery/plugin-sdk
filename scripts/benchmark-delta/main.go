package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type benchdata []struct {
	Suites []struct {
		Benchmarks []struct {
			Name    string
			Runs    int
			NsPerOp float64
			Mem     struct {
				BytesPerOp  int64
				AllocsPerOp int64
			}
			Custom map[string]float64
		}
	}
}

type deltaResult struct {
	Name   string  // file name that will be used
	Metric string  // human-friendly name for metric
	Value  float64 // value
}

func prettyName(name string) string {
	return strings.Trim(strings.TrimPrefix(name, "Benchmark"), "_")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path to benchdata JSON file>", os.Args[0])
	}
	fileName := os.Args[1]
	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	var d benchdata
	err = json.Unmarshal(b, &d)
	if err != nil {
		log.Fatalf("failed to unmarshal benchdata JSON file: %v", err)
	}
	var deltaResults []deltaResult
	for _, run := range d {
		for _, suite := range run.Suites {
			for _, bm := range suite.Benchmarks {
				if bm.NsPerOp > 0 {
					fmt.Println(bm.Name, "ns/op", bm.NsPerOp)
					deltaResults = append(deltaResults, deltaResult{
						Name:   bm.Name + "_ns_per_op",
						Metric: prettyName(bm.Name) + " " + "ns/op",
						Value:  bm.NsPerOp,
					})
				}
				for k, v := range bm.Custom {
					if k == "targetResources/s" {
						// skip, this is a calculated target that isn't expected to change
						// between runs
						continue
					}
					fmt.Println(bm.Name, k, v)
					deltaResults = append(deltaResults, deltaResult{
						Name:   bm.Name + "_" + strings.ReplaceAll(k, "/", "_per_"),
						Metric: prettyName(bm.Name) + " " + k,
						Value:  v,
					})
				}
			}
		}
	}
	for _, dr := range deltaResults {
		name := fmt.Sprintf(".delta.%s", dr.Name)
		fmt.Printf("Writing to %s\n", name)
		data := []byte(fmt.Sprintf("%v (%s)", dr.Value, dr.Metric))
		err := os.WriteFile(name, data, 0664)
		if err != nil {
			log.Fatalf("failed to write %v: %v", name, err)
		}
	}
}
