package module

const (
	ErrorMessageDatabaseError              = "データベースエラーが発生しました。"
	ErrorMessageRakutanNotFound            = "「%s」は見つかりませんでした。\n【検索のヒント】%%を頭につけて検索すると部分一致検索ができます。ex.)「%%地理学」"
	ErrorMessageTooManyRakutan             = "「%s」は%d件見つかりました。検索ワードを絞ってください。"
	ErrorMessageGetRakutanInfoByIDError    = "らくたん情報の取得に失敗しました。"
	ErrorMessageGetRakutanInfoByNameError  = "らくたん情報の取得に失敗しました。"
	ErrorMessageGetOmikujiError            = "おみくじの取得に失敗しました。"
	ErrorMessageGetFavError                = "お気に入りの取得に失敗しました。"
	ErrorMessageInsertFavError             = "お気に入りの登録に失敗しました。"
	ErrorMessageAlreadyFavError            = "すでにお気に入りに登録済みです。"
	ErrorMessageFavLimitError              = "お気に入りの数が上限の50件に達しました。"
	ErrorMessageDeleteFavError             = "お気に入りの削除に失敗しました。"
	ErrorMessageCheckVerificateError       = "認証の確認に失敗しました。"
	ErrorMessageInsertVerificateError      = "認証トークンの作成に失敗しました。"
	ErrorMessageVerificationTokenSendError = "認証トークンの送信に失敗しました。"
	ErrorMessageVerificationFailed         = "すでに認証済みか、認証コードが間違っています。"

	SuccessMessageInsertFav    = "「%s」をお気に入りに登録しました。"
	SuccessMessageDeleteFav    = "「%s」をお気に入りから削除しました。"
	SuccessVericationTokenSent = "認証コードを送信しました。送られたメール内の認証リンクをクリックすると有効化されます。\n届いていない場合は、アドレスが間違っているか迷惑メールに入っている可能性があります。"
	SuccessVerified            = "認証に成功しました。"
)
