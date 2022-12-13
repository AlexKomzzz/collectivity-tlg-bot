package storage

type Bucket string

const (
	AccessTokens Bucket = "access_tokens"
	Debt         Bucket = "debt"
)

type TokenStorage interface {
	Save(chatID int64, debt string, bucket Bucket) error
	Get(chatID int64, bucket Bucket) (string, error)
	Delete(chatID int64, bucket Bucket) error
}
