package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID int
	FirstName string
	LastName string
	Grades []Grade
}

func (s Student)Average() float32 {
	var total float32
	for _,v := range s.Grades{
		total += v.Score
	}
	return total/float32(len(s.Grades))
}

var (
	students Students
	studentsMutex sync.Mutex
)
type Students []Student

func (ss Students) GetByID(id int) (*Student,error){

	for i:=range ss{
		if ss[i].ID == id{
			return &ss[i], nil
		}
	}
	return nil, fmt.Errorf("student id %v not found", id)
}


const   (
	GradeQuiz = GradeType("Quiz")
	GradeTest = GradeType("Test")
	GradeExam = GradeType("Exam")

)


type GradeType string
type Grade struct{
	Title string
	Type GradeType
	Score float32
}

