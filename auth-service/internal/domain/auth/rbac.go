package domain

var rolePermissions = map[Role][]Permission{
	RoleAdmin: {
		PermRead, PermWrite, PermDelete, PermUpdate,
	},
	RoleUser: {
		PermRead,
	},
	RoleManager: {
		PermRead, PermWrite, PermUpdate,
	},
}

func HasPermission(role Role, permission Permission) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}

	for _, p := range perms {
		if p == permission {
			return true
		}
	}
	return false
}
