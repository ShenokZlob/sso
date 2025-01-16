package models

type App struct {
	ID   int
	Name string

	// Используем для подписи токена
	Secret string
}
