package repository

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/gofrs/uuid"
	jsoniter "github.com/json-iterator/go"
	traQrouterV1 "github.com/traPtitech/traQ/router/v1"
	traQrouterV3 "github.com/traPtitech/traQ/router/v3"
)

var traQjson = jsoniter.Config{
	EscapeHTML:             true,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
	TagKey:                 "traq",
}.Froze()

type UserMetaRepository interface {
	SaveUser(userID uuid.UUID, isAdmin, istraQ bool) (*UserMeta, error)
	GetUser(userID uuid.UUID) (*UserMeta, error)
	GetAllUsers() ([]*UserMeta, error)
	ReplaceToken(userID uuid.UUID, token string) error
	GetToken(userID uuid.UUID) (string, error)
	ReplaceiCalSecret(userID uuid.UUID, secret string) error
	GetiCalSecret(userID uuid.UUID) (string, error)
}

type UserBodyRepository interface {
	CreateUser(name, displayName, password string) (*UserBody, error)
	GetUser(userID uuid.UUID) (*UserBody, error)
	GetAllUsers() ([]*UserBody, error)
}

// GormRepository implements UserRepository

func (repo *GormRepository) SaveUser(userID uuid.UUID, isAdmin, istraQ bool) (*UserMeta, error) {
	user := UserMeta{
		ID:     userID,
		Admin:  isAdmin,
		IsTraq: istraQ,
	}
	if err := repo.DB.Create(&user).Error; err != nil {
		var me *mysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}
	return &user, nil
}

// GetUser ユーザー情報を取得します
func (repo *GormRepository) GetUser(userID uuid.UUID) (*UserMeta, error) {
	user := &UserMeta{}
	if err := repo.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil

}

func (repo *GormRepository) GetAllUsers() ([]*UserMeta, error) {
	users := make([]*UserMeta, 0)
	err := repo.DB.Find(&users).Error
	return users, err
}

// GCM encryption
func encryptByGCM(key []byte, plainText string) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize()) // Unique nonce is required(NonceSize 12byte)
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(plainText), nil)
	cipherText = append(nonce, cipherText...)

	return cipherText, nil
}

// Decrypt by GCM
func decryptByGCM(key []byte, cipherText []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := cipherText[:gcm.NonceSize()]
	plainByte, err := gcm.Open(nil, nonce, cipherText[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}

	return string(plainByte), nil
}

func (repo *GormRepository) ReplaceiCalSecret(userID uuid.UUID, secret string) error {
	if userID == uuid.Nil {
		return ErrNilID
	}
	if err := repo.DB.Model(&UserMeta{ID: userID}).Update("ical_secret", secret).Error; err != nil {
		return err
	}
	return nil
}
func (repo *GormRepository) GetiCalSecret(userID uuid.UUID) (string, error) {
	user := UserMeta{
		ID: userID,
	}
	err := repo.DB.First(&user).Error
	return user.IcalSecret, err
}

func (repo *GormRepository) ReplaceToken(userID uuid.UUID, token string) error {
	user := UserMeta{
		ID: userID,
	}
	var cipherText []byte
	var err error
	if token != "" {
		cipherText, err = encryptByGCM(repo.TokenKey, token)
		if err != nil {
			return err
		}
	}
	return repo.DB.Model(&user).Update("token", cipherText).Error
}

func (repo *GormRepository) GetToken(userID uuid.UUID) (string, error) {
	user := UserMeta{
		ID: userID,
	}
	err := repo.DB.First(&user).Error
	if err != nil {
		return "", err
	}
	var token string
	if user.Token != "" {
		token, err = decryptByGCM(repo.TokenKey, []byte(user.Token))
	}

	return token, err
}

// traQRepository implements UserRepository

// CreateUser 新たにユーザーを作成する
func (repo *TraQRepository) CreateUser(name, password, displayName string) (*UserBody, error) {
	if repo.Version != TraQv1 {
		repo.Version = TraQv1
		defer func() {
			repo.Version = TraQv3
		}()
	}
	reqUser := &traQrouterV1.PostUserRequest{
		Name:     name,
		Password: password,
	}
	body, _ := json.Marshal(reqUser)
	resBody, err := repo.postRequest("/users", body)
	if err != nil {
		return nil, err
	}
	traQuser := struct {
		ID uuid.UUID `json:"id"`
	}{}
	err = json.Unmarshal(resBody, &traQuser)
	if err != nil {
		return nil, err
	}
	return &UserBody{ID: traQuser.ID}, nil
}

// GetUser get from /users/{userID}
func (repo *TraQRepository) GetUser(userID uuid.UUID) (*UserBody, error) {
	data, err := repo.getRequest(fmt.Sprintf("/users/%s", userID))
	if err != nil {
		return nil, err
	}
	traQuser := new(traQrouterV3.User)
	err = json.Unmarshal(data, &traQuser)
	return formatV3User(traQuser), err
}

// GetAllUsers get from /users
func (repo *TraQRepository) GetAllUsers() ([]*UserBody, error) {
	data, err := repo.getRequest("/users")
	if err != nil {
		return nil, err
	}
	traQusers := make([]*traQrouterV3.User, 0)
	err = traQjson.Unmarshal(data, &traQusers)
	users := make([]*UserBody, len(traQusers))
	for i, u := range traQusers {
		users[i] = formatV3User(u)
	}
	return users, err
}
func (repo *TraQRepository) UpdateiCalSecretUser(userID uuid.UUID, secret string) error {
	return ErrForbidden
}

func formatV3User(u *traQrouterV3.User) *UserBody {
	return &UserBody{
		ID:          u.ID,
		Name:        u.Name,
		DisplayName: u.DisplayName,
	}
}
