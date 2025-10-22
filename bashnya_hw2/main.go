package main

import "fmt"

var 


func main() {
	var num int

	num = get_num()

	num = operating_num(num)

	print_result(num)
}

func get_num() int {
	var num int
	
	fmt.Println("Введите целое число (num < 12307): ")
	_, err := fmt.Scanf("%d", &num)

	for (err != nil) {
		fmt.Println("Некорректный ввод, попробуйте снова: ")
		_, err = fmt.Scanf("%d", &num)
	}

	return num
}

func operating_num(num int) int {
	for (num < 12307) {
		if num < 0 {
			num *= -1
		} else if num % 7 == 0 {
			num *= 39
		} else if num % 9 == 0 {
			num = num * 13 + 1
			continue
		} else {
			num = (num + 2) * 3
		}

		if num % 9 == 0 && num % 13 == 0 {
			fmt.Println("service error")
			break
		} else {
			num += 1
		}
	}

	return num
}

func print_result(num int) {
	fmt.Printf("Полученное число: %d", num)
}