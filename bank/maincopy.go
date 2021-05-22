package main

import "fmt"

func maincopy() {
	// fmt.Println("Hello world!")
	// something.SayHello()

	// const name string = "maeng"
	// fmt.Println(name)
	// name2 := "soo"
	// fmt.Println(name2)

	// fmt.Println(multiply(2, 2))

	// totalLenght, _ := lenAndUpper("maeng")
	// fmt.Println(totalLenght)

	// repeatMe("ㅁㄴㅇ", "ㅁㅁ", "ㅁㄴㅇ")

	// totalLenght, up := lenAndUpper("maeng")
	// fmt.Println(totalLenght, up)

	result := superAdd(2, 3, 2, 4, 5)
	fmt.Println(result)
}

// func multiply(a int, b int) int {
// 	return a * b
// }

// func lenAndUpper(name string) (int, string) {
// 	return len(name), strings.ToUpper(name)
// }

// func repeatMe(words ...string) {
// 	fmt.Println(words)
// }

// func lenAndUpper(name string) (lenght int, uppercase string) {
// 	defer fmt.Println("I'm done!! wow sexy")
// 	lenght = len(name)
// 	uppercase = strings.ToUpper(name)
// 	return
// }

//loop
func superAdd(numbers ...int) int {
	total := 0
	for _, a := range numbers {
		total += a
	}
	// fmt.Println(numbers)
	// for i := 0; i < len(numbers); i++ {
	// 	fmt.Println(numbers[i])
	// }
	return total
}
