package types

// Key представляет собой пользовательский тип ключа для контекста.
type Key int

const (
	// UserID - ключ для доступа к user_id в контексте
	UserID Key = iota
)
