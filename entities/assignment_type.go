package entities

type AssignmentType int

const (
    MultipleChoice AssignmentType = iota + 1
    FillInTheBlank
    ShortAnswer
    Essay
)

var AssignmentTypeNames = map[AssignmentType]string{
    MultipleChoice: "Multiple Choice",
    FillInTheBlank: "Fill in the Blank",
    ShortAnswer:    "Short Answer",
    Essay:          "Essay",
}

func (a AssignmentType) String() string {
    return AssignmentTypeNames[a]
}