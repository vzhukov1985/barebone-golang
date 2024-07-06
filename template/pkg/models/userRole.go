// Schema generated model - DO NOT EDIT!!!

package models

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
)

var UserRoleSlice = []UserRole{
	UserRoleAdmin,
}

func (e UserRole) Validate() bool {
	for _, v := range UserRoleSlice {
		if e == v {
			return true
		}
	}
	return false
}
