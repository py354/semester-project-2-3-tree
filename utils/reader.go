package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var versions = [...]string{"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10"}
var counts = [...]int{100, 1000, 10000, 100000, 200000, 400000, 600000, 800000, 1000000}

const (
	FindOperation   = "find.csv"
	InsertOperation = "insert.csv"
	DeleteOperation = "delete.csv"
)

// Average - функция для вычисления среднего времени
func Average(durs []time.Duration) time.Duration {
	if len(durs) == 0 {
		return 0
	}

	s := time.Duration(0)
	for _, d := range durs {
		s += d
	}
	return s / time.Duration(len(durs))
}

// Data - структура, содержащая наборы данных (например по 100 элементов)
type Data struct {
	Count int        // кол-во элементов в наборах
	Nums  chan []int // канал для передачи этих наборов
}

// ReadData - передает канал для передачи наборов данных (по 100, 1000...)
func ReadData() chan Data {
	c := make(chan Data, 5)

	// выполняется параллельно
	go func() {
		for _, count := range counts {
			data := Data{
				Count: count,
				Nums:  make(chan []int, 10),
			}

			for _, version := range versions {
				data.Nums <- ReadFile(version, count)
			}
			close(data.Nums)

			c <- data
		}
		close(c)
	}()

	return c
}

// ReadFile - читает числа из файла
func ReadFile(version string, count int) []int {
	res := make([]int, 0, count)
	path := "dataset/" + version + "/" + strconv.Itoa(count) + ".csv"
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	for _, binNum := range bytes.Split(b, []byte(",")) {
		if len(binNum) == 0 {
			continue
		}
		num, err := strconv.Atoi(string(binNum))
		if err != nil {
			panic(err)
		}

		res = append(res, num)
	}

	return res
}

type Result struct {
	Count int
	Dur   time.Duration
}

// WriteResult - записывает итоги бенчмаркинга
func WriteResult(operation string, results []Result) {
	file, err := os.Create("result/" + operation)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, r := range results {
		_, err = file.WriteString(fmt.Sprintf("%d,%d\n", r.Count, r.Dur.Nanoseconds()))
		if err != nil {
			panic(err)
		}
	}
}

// RandElem возвращает рандомный элемент из данных
func RandElem(nums []int) int {
	return nums[rand.Intn(len(nums))]
}

// Ставим рандомное зерно для рандома
func init() {
	rand.Seed(time.Now().Unix())

}
