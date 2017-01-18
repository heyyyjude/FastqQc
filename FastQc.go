package main

import (
	"fmt"
	"log"
	"os"
	"bufio"
	"strings"
	"compress/gzip"
	"github.com/montanaflynn/stats"
	"sort"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var totalNum float64
var gcNum float64
var readLengthSlice []float64
var qualMeanS []float64
var qualMap = make(map[int][]float64)
var qualPostStatM = make(map[int]float64)

var q30 float64
var q25 float64

func Round(v float64, decimals int) float64 {
	var pow float64 = 1
	for i := 0; i < decimals; i++ {
		pow *= 10
	}
	return float64(int((v * pow) + 0.5)) / pow
}

func getQualPosMean(m map[int][]float64) (map[int]float64, []float64) {
	var resM = make(map[int]float64)
	var resS []float64

	for k, v := range m {
		mean, _ := stats.Mean(v)
		resM[k] = Round(mean, 0)
		resS = append(resS, mean)
	}
	return resM, resS
}

// calculate q30 q25
func convertQual(qual []rune, m map[int][]float64) map[int][]float64 {

	for i, j := range qual {
		q := float64(j - 33)
		m[i + 1 ] = append(m[i + 1], q)

		if q >= 30 {
			q30++
		}
		if q >= 25 {
			q25++
		}
	}
	return m
}

var readStatMap = make(map[string]float64)

func getStat(s []float64) map[string]float64 {
	var m = make(map[string]float64)
	max, _ := stats.Max(s)
	min, _ := stats.Min(s)
	median, _ := stats.Median(s)
	quarile, _ := stats.Quartile(s)
	mean, _ := stats.Mean(s)

	m["Max"] = max
	m["Min"] = min
	m["Median"] = Round(median, 2)
	m["Mean"] = Round(mean, 2)
	m["q1"] = Round(quarile.Q1, 2)
	m["q2"] = Round(quarile.Q2, 2)
	m["q3"] = Round(quarile.Q3, 2)

	return m
}

func printQualMeanOfPos(m map[int]float64) {
	var res []int

	for k := range m {
		res = append(res, k)
	}

	sort.Ints(res)
	for _, i := range res {
		fmt.Print(m[i])
		fmt.Print(",")
	}
}

func printQualCli(m map[int]float64) {

	var sortedKey = make([] int, 0)

	for k := range m {
		sortedKey = append(sortedKey, k)
	}

	sort.Ints(sortedKey)

	var tmpM = make(map[int]float64)
	var tmp []float64

	for i := range sortedKey {
		if i == 0 {
			continue
		} else if i%5 != 0 {
			tmp = append(tmp, m[i])
		} else {
			mean, _ := stats.Mean(tmp)
			tmpM[i] = Round(mean, 0)
			tmp = make([]float64, 0)
		}
	}

	var resM = make(map[int]string)

	for k, v := range tmpM {
		sth := QualCliHelper(v)
		resM[k] = sth
	}

	var res []int
	for k := range tmpM {
		res = append(res, k)
	}

	sort.Ints(res)

	for _, i := range res {
		fmt.Print(i, ":", resM[i])
		fmt.Println()
	}
}

func QualCliHelper(q float64) string {
	var tmp = ""

	for i := 1; i <= int(q); i++ {
		if i%5 != 0 {
			tmp += "-"
		} else {
			tmp += "|"
		}
	}
	return tmp

}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage")
		fmt.Println("FastQc fastq.gz")
		os.Exit(0)
	}
	fastqGz := os.Args[1]

	fq, err := os.Open(fastqGz)
	checkError(err)
	defer fq.Close()

	fqgz, err := gzip.NewReader(fq)
	//fqgz := bufio.NewReader(fq)
	checkError(err)
	defer fqgz.Close()

	scanner := bufio.NewScanner(fqgz)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var line string = strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "@") {
			continue
		} else {
			// sequence line
			gcNum += float64(strings.Count(strings.ToUpper(line), "G"))
			gcNum += float64(strings.Count(strings.ToUpper(line), "C"))
			readLengthSlice = append(readLengthSlice, float64(len(line)))
			// strand
			scanner.Scan()
			// quality
			scanner.Scan()

			var qual []rune = []rune(strings.TrimSpace(scanner.Text()))

			convertQual(qual, qualMap)
		}
	}

	qualPostStatM, qualMeanS = getQualPosMean(qualMap)
	totalNum, _ = stats.Sum(readLengthSlice)
	var gcRatio float64 = Round(gcNum/totalNum*100, 2)
	readStatMap = getStat(readLengthSlice)

	fmt.Println("----Nucleotide stat----")
	fmt.Println("Total basepair :", totalNum)
	fmt.Println("GC ratio :", gcRatio, "%")

	fmt.Println("----Reads stat----")
	fmt.Println("Total reads:", len(readLengthSlice))
	fmt.Println("Max read length:", readStatMap["Max"])
	fmt.Println("Min read length:", readStatMap["Min"])
	fmt.Println("Median read length:", readStatMap["Median"])
	fmt.Println("Mean read length:", readStatMap["Mean"])
	fmt.Println("N25:", readStatMap["q1"])
	fmt.Println("N50:", readStatMap["q2"])
	fmt.Println("N75:", readStatMap["q3"])

	fmt.Println("----Quality stat----")
	fmt.Println("Mean of each postion quality score")
	printQualMeanOfPos(qualPostStatM)
	fmt.Println()

	// they can have float values because the input is mean of each position quality score.
	qualTotalStatM := getStat(qualMeanS)
	fmt.Println("Median Quality score:", qualTotalStatM["Median"])
	fmt.Println("Mean Quality score:", qualTotalStatM["Mean"])
	fmt.Println("N25:", qualTotalStatM["q1"])
	fmt.Println("N50:", qualTotalStatM["q2"])
	fmt.Println("N75:", qualTotalStatM["q3"])
	fmt.Println("Q30:", Round(q30/totalNum*100, 2), "%")
	fmt.Println("Q20:", Round(q25/totalNum*100, 2), "%")

	fmt.Println("---quality plot---")
	printQualCli(qualPostStatM)
}
