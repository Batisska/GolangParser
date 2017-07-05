// В этом примере реализован _набор обработчиков
// (worker pool)_ с помощью горутин и каналов.

package main

import (
	"strconv"
	"fmt"
)

type res struct {
	jobs int
	worker string
}

func main() {

	i, err := strconv.Atoi("-42")
	if err != nil {
		fmt.Println(err.Error())
	}
	s := strconv.Itoa(-42)
	fmt.Println(i)
	fmt.Println(s)
}
