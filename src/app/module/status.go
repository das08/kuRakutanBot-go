package module

type KRBStatus int

const (
	KRBSuccess            KRBStatus = 2000
	KRBDatabaseError      KRBStatus = 4000
	KRBOmikujiError       KRBStatus = 4000
	KRBGetFavError        KRBStatus = 4003
	KRBInsertFavError     KRBStatus = 4004
	KRBDeleteFavError     KRBStatus = 4005
	KRBGetLecIDError      KRBStatus = 4006
	KRBGetLecNameError    KRBStatus = 4007
	KRBGetUidError        KRBStatus = 4008
	KRBVerifyCodeGenError KRBStatus = 4009
	KRBVerifyCodeDelError KRBStatus = 4010
)

type QueryStatus struct {
	Success bool
	Message string
	Status  KRBStatus
}
