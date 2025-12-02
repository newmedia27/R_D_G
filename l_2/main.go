package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func FibonacciIterative(n int) int {
	// Підказка: тримайте два останні значення й оновлюйте їх у циклі.
	// Вхід вважаємо: n >= 0.
	// При отриманні негативного n повертаємо його без змін

	if n <= 1 {
		return n
	}

	prev := 0
	curr := 1

	for i := 2; i <= n; i++ {
		prev, curr = curr, prev+curr
	}
	return curr
}

func FibonacciRecursive(n int) int {
	// База: n==0 -> 0; n==1 -> 1.
	// Рекурсія: F(n-1)+F(n-2)
	// При отриманні негативного n повертаємо його без змін

	if n <= 1 {
		return n
	}
	return FibonacciRecursive(n-2) + FibonacciRecursive(n-1)
}

func IsPrime(n int) bool {
	// Підказка: n<=1 -> false; 2 -> true; парні >2 -> false;
	// Далі перевіряйте дільники до sqrt(n).

	if n <= 1 || (n != 2 && n%2 == 0) {
		return false
	}

	for i := 3; i*i <= n; i += 2 {
		if n%i == 0 {
			return false
		}
	}

	return true
}

func IsBinaryPalindrome(n int) bool {
	// Підказка: перетворіть n у строку (strconv ефективніший за fmt),
	// потім перевірте паліндромність.
	if n < 0 {
		return false
	}

	base := []rune(strconv.FormatInt(int64(n), 2))
	if len(base) == 1 {
		return true
	}
	for i, j := 0, len(base)-1; i < j; i, j = i+1, j-1 {
		if base[i] != base[j] {
			return false
		}
	}
	return true
}

func ValidParentheses(s string) bool {
	// Правила:
	// 1. Допустимі дужки (, [, {, ), ], }
	// 2. У кожної відкритої дужки є відповідна закриваюча дужка того ж типу
	// 3. Закриваючі дужки стоять у правильному порядку
	// "[{}]" - правильно
	// "[{]}" - не правильно
	// 4. Кожна закриваюча дужка має відповідну відкриваючу дужку
	// Підказка: використовуйте стек (можна зробити через масив рун []rune)

	if len(s) <= 0 {
		return false
	}

	re := regexp.MustCompile(`[^(\[{)\]}]`)
	str := re.ReplaceAllString(s, "")

	if len(str)%2 != 0 {
		return false
	}
	openStr := "([{"
	closeStr := ")]}"

	var stack []rune

	for _, c := range str {
		if strings.Contains(openStr, string(c)) {
			stack = append(stack, c)
		}
		if strings.Contains(closeStr, string(c)) {
			if len(stack) <= 0 {
				return false
			}
			index := strings.Index(closeStr, string(c))
			if stack[len(stack)-1] != rune(openStr[index]) {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}

func Increment(num string) int {
	// Тобто строка містить певне число у бінарному вигляді
	// Потрібно повернути число на один більше
	// Додайте валідацію вхідної строки, якщо вона містить недопустимі символи, повертайте 0

	if len(num) <= 0 {
		return 0
	}

	//r, err := regexp.MatchString("^[01]+$", num)
	//if err != nil || !r {
	//	return 0
	//}

	for _, s := range num {
		if s != '0' && s != '1' {
			return 0
		}
	}

	n, err := strconv.ParseInt(num, 2, 64)
	if err != nil {
		return 0
	}

	return int(n + 1)
}

func main() {
	// Невеликі демонстраційні виклики (для наочного запуску `go run .`)
	fmt.Println("FibonacciIterative(10):", FibonacciIterative(10)) // очікуємо 55
	fmt.Println("FibonacciRecursive(10):", FibonacciRecursive(10)) // очікуємо 55

	fmt.Println("IsPrime(2):", IsPrime(2))   // true
	fmt.Println("IsPrime(15):", IsPrime(15)) // false
	fmt.Println("IsPrime(29):", IsPrime(29)) // true

	fmt.Println("IsBinaryPalindrome(7):", IsBinaryPalindrome(7)) // true (111)
	fmt.Println("IsBinaryPalindrome(6):", IsBinaryPalindrome(6)) // false (110)

	fmt.Println(`ValidParentheses("[]{}()"):`, ValidParentheses("[]{}()")) // true
	fmt.Println(`ValidParentheses("[{]}"):`, ValidParentheses("[{]}"))     // false

	fmt.Println(`Increment("101") ->`, Increment("101")) // 6
}
