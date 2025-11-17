package passhash

import "golang.org/x/crypto/bcrypt"

type Service interface {
	Hash(password string) (string, error)
	Verify(hashed, plain string) bool
}

type bcryptService struct{ cost int }

func NewBcrypt(cost int) Service { return &bcryptService{cost: cost} }

func (b *bcryptService) Hash(password string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	return string(h), err
}
func (b *bcryptService) Verify(hashed, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}
