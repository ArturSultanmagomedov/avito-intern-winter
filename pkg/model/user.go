package model

type User struct {
	Id      int     `db:"id"`
	UserId  int     `json:"id" db:"user_id"` // TODO: сделать UserId строкой
	Balance float32 `json:"balance" db:"balance"`
}

// GetFields чтобы передавать в sql.Scan() все поля структуры User
func (r *User) GetFields() []interface{} {
	return []interface{}{&r.Id, &r.UserId, &r.Balance}
}
