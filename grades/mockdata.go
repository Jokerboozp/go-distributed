package grades

func init() {
	students = []Student{
		{
			ID:        1,
			FirstName: "Nick",
			LastName:  "Jack",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 90,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 99,
				},
			},
		},
		{
			ID:        2,
			FirstName: "Nick2",
			LastName:  "Jack2",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 78,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 79,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 80,
				},
			},
		},
		{
			ID:        3,
			FirstName: "Nick3",
			LastName:  "Jack3",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 66,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 67,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 68,
				},
			},
		},
		{
			ID:        4,
			FirstName: "Nick4",
			LastName:  "Jack4",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 55,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 56,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 57,
				},
			},
		},
		{
			ID:        5,
			FirstName: "Nick5",
			LastName:  "Jack5",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 34,
				},
				{
					Title: "Final Exam",
					Type:  GradeExam,
					Score: 35,
				},
				{
					Title: "Test 1",
					Type:  GradeTest,
					Score: 36,
				},
			},
		},
	}
}
