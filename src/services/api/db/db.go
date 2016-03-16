package db

const (
	DATABASE_NAME = "neutrino"

	USERS_TABLE      = "users"
	DATA_TABLE       = "data_0"
	APPS_TABLE       = "apps"
	APPS_USERS_TABLE = "apps_users"

	USERS_TYPE = "users"

	APP_ID_FIELD        = "_appId"
	TYPE_FIELD          = "_type"
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

	EMAIL_INDEX            = "email"
	ITEMS_FOR_APP_INDEX    = "items_for_app_by_type"
	EMAIL_APPID_USER_INDEX = "email_appId_user"
)

func NewDbService() DbService {
	return &dbService{}
}

//
//func NewUserDbService(u, appId string) UserDbService {
//	d := NewDbService(DATABASE_NAME, USERS_TABLE)
//	return &userDbService{d, u, appId}
//}
//
//func NewDataDbService(appId, t string) DataDbService {
//	d := NewDbService(DATABASE_NAME, DATA_TABLE)
//	return &dataDbService{d, t, appId}
//}
