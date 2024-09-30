package devconfig

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Dialosoft/src/domain/models"
	"github.com/Dialosoft/src/pkg/utils/security"
	"gorm.io/gorm"
)

func StartWithAdmin(db *gorm.DB) error {
	var adminRole models.RoleEntity
	result := db.Where("role_type = ?", "administrator").First(&adminRole)
	if result.Error != nil {
		return result.Error
	}

	passwordHashed, err := security.HashPassword("admin")
	if err != nil {
		return err
	}

	defaultUser := models.UserEntity{
		Username: "administrator",
		Email:    "administrator@dialosoft.com",
		Password: passwordHashed,
		Name:     "Administrator",
		RoleID:   adminRole.ID,
		Banned:   false,
		Role:     adminRole,
	}

	result = db.Create(&defaultUser)
	if result.Error != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(result.Error.Error(), "duplicate key value violates unique constraint") {
			fmt.Println("default User already exists")
			return nil
		}
		return result.Error
	}

	fmt.Println("default User created by id: " + defaultUser.ID.String())
	return nil
}
