//polymorphism
package main 
import "fmt"
type income interface {
	calculate() int
    source() string
}
type fixedbilling struct {
	projectname string
	biddedamount int
}
type timeandmaterial struct {
	projectname string
	noofhours int
	hourlyrate int
 }
 func (fb fixedbilling) calculate() int {
	 return fb.biddedamount
 }
 func (fb fixedbilling) source() string {
	 return fb.projectname
 }	 
 func (tm timeandmaterial) calculate() int {
	 return tm.noofhours * tm.hourlyrate
 }
 func (tm timeandmaterial) source() string {
	return tm.projectname
 }
 func calculatenetincome(ic [] income) {
	 var netincome int = 0
	 for _,income := range ic {
		 fmt.Printf("%s = %d" , income.source() , income.calculate())
		 netincome = income.calculate()
	 }
	 fmt.Printf("\n netincome of organization = %d",netincome)
 }
 func main() {
	 project1 := fixedbilling{projectname:"project1",biddedamount:8000}
	 project2 :=fixedbilling{projectname:"project2",biddedamount:5000}
	 project3 := timeandmaterial{projectname:"project3",noofhours:9,hourlyrate:1000}
	 totalgenerated := [] income{project1,project2,project3}
	 calculatenetincome(totalgenerated)
	 
 }
