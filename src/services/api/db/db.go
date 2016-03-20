package db

const (
	DATABASE_NAME = "neutrino"

	USERS_TABLE      = "users"
	DATA_TABLE       = "data0"
	APPS_TABLE       = "apps"
	APPS_USERS_TABLE = "apps_users"

	USERS_TYPE = "users"

	APP_ID_FIELD        = "appId"
	TYPE_FIELD          = "type"
	TYPES_FIELD         = "types"
	APPS_FIELD          = "apps"
	USERS_FIELD         = "users"
	ID_FIELD            = "id"
	PASSWORD_FIELD      = "password"
	EMAIL_FIELD         = "email"
	REGISTERED_AT_FIELD = "registeredAt"
	NAME_FIELD          = "name"
	OWNER_FIELD         = "owner"
	MASTER_KEY_FIELD    = "masterKey"

	USERS_TABLE_EMAIL_INDEX           = "email"
	DATA_TABLE_APPIDTYPE_INDEX        = "items_for_app_by_type"
	APPS_USERS_TABLE_EMAILAPPID_INDEX = "email_appId_user"
)

var (
	DB_FIELDS = []string{TYPE_FIELD, APP_ID_FIELD}
)

func NewDbService() DbService {
	return &dbService{}
}
