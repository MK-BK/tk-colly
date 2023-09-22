package models

type Environment struct {
	CategoryManager
	MoviesManager
}

var GlobalEnvironment Environment
