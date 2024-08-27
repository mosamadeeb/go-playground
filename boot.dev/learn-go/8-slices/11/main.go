package main

type cost struct {
	day   int
	value float64
}

func getCostsByDay(costs []cost) []float64 {
	// Not sure if this is a good practice, but it works
	totalCosts := make([]float64, 0)

	for i := 0; i < len(costs); i++ {
		for costs[i].day >= len(totalCosts) {
			totalCosts = append(totalCosts, 0.0)
		}

		totalCosts[costs[i].day] += costs[i].value
	}

	return totalCosts
}
