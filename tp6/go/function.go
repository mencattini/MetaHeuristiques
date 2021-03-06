package main

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Individu struct {
	x, y      string
	pm, pc, f float64
}

func (i *Individu) Fitness() {
	tmp_x, _ := strconv.ParseInt(i.x, 2, 32)
	tmp_y, _ := strconv.ParseInt(i.y, 2, 32)
	x := float64(tmp_x) + 0.01
	y := float64(tmp_y) + 0.01

	i.f = math.Pow(x/1000, 512) + math.Pow(10/x, 4) + math.Pow(y/1000, 512) + math.Pow(10/y, 4) - math.Abs(0.5*x*math.Sin(math.Sqrt(math.Abs(x)))) - math.Abs(y*math.Sin(30*math.Sqrt(math.Abs(x/y))))
}

func (i *Individu) Mutation() {

	pm := i.pm
	x := i.x
	for index, ele := range i.x {
		if rand.Float64() < pm {
			tmp, _ := strconv.Atoi(string(ele))
			x = x[:index] + string(strconv.FormatFloat(math.Mod(float64(tmp+1), 2), 'f', 1, 64)[0]) + x[index+1:]
		}
	}
	i.x = x

	y := i.y
	for index, ele := range i.y {
		if rand.Float64() < pm {
			tmp, _ := strconv.Atoi(string(ele))
			y = y[:index] + string(strconv.FormatFloat(math.Mod(float64(tmp+1), 2), 'f', 1, 64)[0]) + y[index+1:]
		}
	}
	i.y = y
}

func initIndividu(pm float64, pc float64, size_population int) []Individu {
	rand.Seed(int64(time.Now().Nanosecond()))
	res := make([]Individu, size_population)
	iter := size_population - 1
	for iter >= 0 {
		x := int64(10 + rand.Intn(991))
		y := int64(10 + rand.Intn(991))
		res[iter] = Individu{strconv.FormatInt(x, 2), strconv.FormatInt(y, 2), pm, pc, 0}
		iter -= 1
	}
	return res
}

func findMin(population []Individu) Individu {
	min := population[0].f
	index := 0
	iter := 1
	for iter < len(population) {
		if min > population[iter].f {
			min = population[iter].f
			index = iter
		}
		iter += 1
	}
	return population[index]
}

func tournament(number int, population []Individu) {
	res := make([]Individu, len(population))

	i := 0
	n := len(population)
	iter := 0

	// partie de séléction
	for iter < n {
		tmp := make([]Individu, number)
		// partie de tournois
		for i < number {
			tmp[i] = population[rand.Intn(n)]
			i += 1
		}
		res[iter] = findMin(tmp)
		i = 0
		iter += 1
	}
	for index, _ := range population {
		population[index] = res[index]
	}
}

func mutation(population []Individu) {
	for index, _ := range population {
		population[index].Mutation()
	}
}

func parallelFitness(population []Individu, n int) {
	pas := len(population) / n
	var vg sync.WaitGroup
	vg.Add(n)
	for i := 0; i < n; i++ {
		go fitness(population[i*pas:(i+1)*pas], &vg)
	}
	vg.Wait()
}

func fitness(population []Individu, wg *sync.WaitGroup) {
	defer wg.Done()
	for index, _ := range population {
		population[index].Fitness()
	}
}

func parallelMidPoint(population []Individu, n int) {
	pas := len(population) / n
	var vg sync.WaitGroup
	vg.Add(n)
	for i := 0; i < n; i++ {
		go crossoverMidPoint(population[i*pas:(i+1)*pas], &vg)
	}
	vg.Wait()
}

func midBreak(i1 *Individu, i2 *Individu, pc float64) {
	if pc > rand.Float64() {
		x := string(i1.x)
		y := string(i1.y)

		new_x := string(i2.x)
		new_y := string(i2.y)

		if len(new_x) >= len(x) {
			x = strings.Repeat("0", len(new_x)-len(x)) + x
		} else {
			new_x = strings.Repeat("0", len(x)-len(new_x)) + new_x
		}

		if len(new_y) >= len(y) {
			y = strings.Repeat("0", len(new_y)-len(y)) + y
		} else {
			new_y = strings.Repeat("0", len(y)-len(new_y)) + new_y
		}

		pointX := len(x) / 2
		pointY := len(y) / 2

		tmpX := x[pointX:]
		tmpY := y[pointY:]

		i1.x = x[:pointX] + new_x[pointX:]
		i1.y = y[:pointY] + new_y[pointY:]

		i2.x = new_x[:pointX] + tmpX
		i2.y = new_y[:pointY] + tmpY
	}
}

func crossoverMidPoint(population []Individu, vg *sync.WaitGroup) {
	defer vg.Done()
	pc := population[0].pc
	size := len(population) - 1

	for i := 0; i < size-1; i += 2 {
		midBreak(&population[i], &population[i+1], pc)
	}
}
