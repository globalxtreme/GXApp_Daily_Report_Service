package core

type BaseParser interface {
	Get() []interface{}
	First() interface{}
}

type BaseActivityPropertyParserInterface interface {
	CreateActivity(action string) interface{}
	UpdateActivity(action string) interface{}
	DeleteActivity(action string) interface{}
	GeneralActivity(action string) interface{}
}
