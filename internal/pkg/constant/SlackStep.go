package constant

const (
	SLACK_STEP_1 = 1
	SLACK_STEP_2 = 2
	SLACK_STEP_3 = 3
	SLACK_STEP_4 = 4
	SLACK_STEP_5 = 5
)

var SlackStepQuestions = map[int]string{
	SLACK_STEP_1: "What did you complete yesterday?",
	SLACK_STEP_2: "What will you do today?",
	SLACK_STEP_3: "When will you be finished with that?",
	SLACK_STEP_4: "Anything blocking your progress?",
	SLACK_STEP_5: "How do you feel today?",
}
