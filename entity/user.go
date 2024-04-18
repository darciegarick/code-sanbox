package entity

import "strconv"

type User struct {
	ID
	Name      string `json:"name" gorm:"not null;comment:用户名称"`
	Mobile    string `json:"mobile" gorm:"not null;index;comment:用户手机号"`
	Password  string `json:"password" gorm:"not null;default:'';comment:用户密码"`
	UserRole  string `json:"user_role" gorm:"not null;default:'user';comment:用户角色"`
	AccessKey string `json:"access_key" gorm:"not null;default:'';comment:访问密钥"`
	SecretKey string `json:"secret_key" gorm:"not null;default:'';comment:密钥"`
	Timestamps
	SoftDeletes
}

func (user User) GetUid() string {
	return strconv.Itoa(int(user.ID.ID))
}
