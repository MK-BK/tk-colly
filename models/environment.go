package models

type Environment struct {
	MoviesManager
	CategoryManager
}

var GlobalEnvironment Environment
