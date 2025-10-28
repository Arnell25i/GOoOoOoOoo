// main.go
package main

import (
	"fmt"
	"bashnya_hw3/stack" 
)

func main() {
	st := stack.New()

	fmt.Println("empty?", st.IsEmpty()) 

	st.Push(10)
	st.Push(20)
	st.Push(30)
	fmt.Println("size after pushes:", st.Size()) 

	v, ok := st.Pop()
	fmt.Println("pop:", v, "ok:", ok) 
	v, ok = st.Pop()
	fmt.Println("pop:", v, "ok:", ok) 
	v, ok = st.Pop()
	fmt.Println("pop:", v, "ok:", ok) 

	_, ok = st.Pop()
	fmt.Println("pop on empty ok:", ok)

	st.Push(42)
	st.Push(99)
	fmt.Println("size before clear:", st.Size()) 
	st.Clear()
	fmt.Println("size after clear:", st.Size(), "empty?", st.IsEmpty())
}
