package entities

type EnglishLevel int

const (
	A1Beginner EnglishLevel = iota + 1
	A2Elementary
	B1Intermediate
	B2UpperIntermediate
	C1Advanced
	C2Proficient
)

var EnglishLevelNames = map[EnglishLevel]string{
	A1Beginner:          "A1 - Beginner",
	A2Elementary:        "A2 - Elementary",
	B1Intermediate:      "B1 - Intermediate",
	B2UpperIntermediate: "B2 - Upper Intermediate",
	C1Advanced:          "C1 - Advanced",
	C2Proficient:        "C2 - Proficient",
}

func (e EnglishLevel) String() string {
	return EnglishLevelNames[e]
}
