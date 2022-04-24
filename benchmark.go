package main

import (
	"log"
	"semester-project-2-3-tree/tree"
	"semester-project-2-3-tree/utils"
	"time"
)

// CYCLES - Количество повторений
const CYCLES = 1000

func main() {
	//// Операция поиска
	findResults := make([]utils.Result, 0)
	for data := range utils.ReadData() {
		timings := make([]time.Duration, 0)

		// для каждого набора данных
		for nums := range data.Nums {
			//подготовительный этап
			t := tree.Tree{}
			for _, num := range nums {
				t.Insert(num)
			}

			// измерение поиска случайного элемента
			timer := utils.Timer{}
			for i := 0; i < CYCLES; i++ {
				elem := utils.RandElem(nums)
				timer.Start()
				_, ok := t.Find(elem)
				timer.Stop()
				if !ok {
					panic("cant find element")
				}
			}

			timings = append(timings, timer.Passed()/CYCLES)
		}

		// сохранение результатов
		findResults = append(findResults, utils.Result{
			Count: data.Count,
			Dur:   utils.Average(timings),
		})
		log.Println("find", data.Count, utils.Average(timings))
	}

	// Операция вставки и удаления
	insertResults := make([]utils.Result, 0)
	deleteResults := make([]utils.Result, 0)
	for data := range utils.ReadData() {
		insertTimings := make([]time.Duration, 0)
		deleteTimings := make([]time.Duration, 0)

		// для каждого набора данных
		for nums := range data.Nums {
			//подготовительный этап
			t := tree.Tree{}
			for _, num := range nums {
				t.Insert(num)
			}

			// измерение вставки рандомного элемена
			insertTimer := utils.Timer{}
			deleteTimer := utils.Timer{}

			for i := 0; i < CYCLES; i++ {
				elem := utils.RandElem(nums)

				deleteTimer.Start()
				ok1 := t.Delete(elem)
				deleteTimer.Stop()

				insertTimer.Start()
				ok2 := t.Insert(elem)
				insertTimer.Stop()

				if !ok1 || !ok2 {
					panic("cant delete/insert element")
				}
			}

			insertTimings = append(insertTimings, insertTimer.Passed()/CYCLES)
			deleteTimings = append(deleteTimings, deleteTimer.Passed()/CYCLES)
		}

		log.Println("insert", data.Count, utils.Average(insertTimings))
		log.Println("delete", data.Count, utils.Average(deleteTimings))

		// сохранение результатов
		insertResults = append(insertResults, utils.Result{
			Count: data.Count,
			Dur:   utils.Average(insertTimings),
		})

		deleteResults = append(deleteResults, utils.Result{
			Count: data.Count,
			Dur:   utils.Average(deleteTimings),
		})
	}

	// запись результатов на файл
	utils.WriteResult(utils.FindOperation, findResults)
	utils.WriteResult(utils.InsertOperation, insertResults)
	utils.WriteResult(utils.DeleteOperation, deleteResults)
}
