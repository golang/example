package main 
import "fmt" 
//creating structures
type student struct {
	name string
	age int
	percentage float64
}
type teacher struct {
	name string 
	age int
}
//creating same methods but different types of receivers
func (s student) show() {
	fmt.Println("name of the student: ",s.name)
	fmt.Println("age of the student: ",s.age)
	fmt.Println("percentage of the student", s.percentage)
}
func (t teacher) show() {
	fmt.Println("name of the teacher: ", t.name)
	fmt.Println("marks", t.age)
} 
// main method
func main() {
	//initialise value of structures
	val1 := student{"satya",24,78.5}
	val2 := teacher{"subbu",30}
	//calling methods
	val1.show()
	val2.show()
}
