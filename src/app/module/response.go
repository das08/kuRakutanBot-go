package module

const (
	ErrorMessageDatabaseError             = "データベースエラーが発生しました。"
	ErrorMessageGetRakutanInfoByIDError   = "らくたん情報の取得に失敗しました。"
	ErrorMessageGetRakutanInfoByNameError = "らくたん情報の取得に失敗しました。"
	ErrorMessageGetOmikujiError           = "おみくじの取得に失敗しました。"
	ErrorMessageGetFavError               = "お気に入りの取得に失敗しました。"
	ErrorMessageInsertFavError            = "お気に入りの登録に失敗しました。"
	ErrorMessageAlreadyFavError           = "すでにお気に入りに登録済みです。"
	ErrorMessageFavLimitError             = "お気に入りの数が上限の50件に達しました。"
	ErrorMessageDeleteFavError            = "お気に入りの削除に失敗しました。"

	SuccessMessageInsertFav = "「%s」をお気に入りに登録しました。"
	SuccessMessageDeleteFav = "「%s」をお気に入りから削除しました。"
)
