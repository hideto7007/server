// validation/signin_validation.go
package validation

import (

	// "server/models"

	"regexp"
	"server/utils"

	"github.com/asaskevich/govalidator"
)

const UserPassword = "user_password"

// type (
// 	signInValidationFetcher interface {
// 		intCheck(val interface{}) error
// 		Validate() error
// 	}

// 	SinginValidation struct {
// 		models.RequestSignInData
// 	}

// 	signInValidation struct{}
// )

// func NewSignInValidationFetcher() signInValidationFetcher {
// 	return &signInValidation{}
// }

type RequestSignInData struct {
	UserName     string `json:"user_name" valid:"required~ユーザー名は必須です。,email~正しいメールアドレス形式である必要があります。"`
	UserPassword string `json:"user_password" valid:"required~パスワードは必須です。"`
}

type TemporayRequestSignUpData struct {
	UserName     string `json:"user_name" valid:"required~ユーザー名は必須です。,email~正しいメールアドレス形式である必要があります。"`
	UserPassword string `json:"user_password" valid:"required~パスワードは必須です。"`
	NickName     string `json:"nick_name" valid:"required~ニックネームは必須です。"`
}

type RequestSignUpData struct {
	UserName     string `json:"user_name" valid:"required~ユーザー名は必須です。,email~正しいメールアドレス形式である必要があります。"`
	UserPassword string `json:"user_password" valid:"required~パスワードは必須です。"`
	NickName     string `json:"nick_name" valid:"required~ニックネームは必須です。"`
}

type RequestRetryAuthEmail struct {
	UserName string `json:"user_name" valid:"required~ユーザー名は必須です。,email~正しいメールアドレス形式である必要があります。"`
	RedisKey string `json:"redis_key" valid:"required~Redisキーは必須です。"`
	NickName string `json:"nick_name" valid:"required~ニックネームは必須です。"`
}

type RequestSignInEditData struct {
	UserId       string `json:"user_id" valid:"required~ユーザーIDは必須です。"`
	UserName     string `json:"user_name" valid:"required~ユーザー名は必須です。,email~正しいメールアドレス形式である必要があります。"`
	UserPassword string `json:"user_password"`
}

type RequestSignInDeleteData struct {
	UserId   string `json:"user_id" valid:"required~ユーザーIDは必須です。"`
	UserName string `json:"user_name" valid:"required~ユーザー名は必須です。,email~正しいメールアドレス形式である必要があります。"`
}

type RequestRefreshTokenData struct {
	UserId string `json:"user_id" valid:"required~ユーザーIDは必須です。"`
}

type RequestSignOutData struct {
	UserName string `json:"user_name" valid:"required~ユーザー名は必須です。,email~正しいメールアドレス形式である必要があります。"`
}

type RequestPriceManagementData struct {
	MoneyReceived string `json:"money_received" valid:"int~月の収入は整数値のみです。"`
	Bouns         string `json:"bouns" valid:"int~ボーナスは整数値のみです。"`
	FixedCost     string `json:"fixed_cost" valid:"int~月の収入は整数値のみです。"`
	Loan          string `json:"loan" valid:"int~ローンは整数値のみです。"`
	Private       string `json:"private" valid:"int~プライベートは整数値のみです。"`
	Insurance     string `json:"insurance" valid:"int~保険は整数値のみです。"`
}

type RequestYearIncomeAndDeductiontData struct {
	UserId    string `json:"user_id" valid:"required~ユーザーIDは必須です。"`
	StartDate string `json:"start_date" valid:"required~開始期間は必須です。"`
	EndDate   string `json:"end_date" valid:"required~終了期間は必須です。"`
}

type RequestDateRangeData struct {
	UserId string `json:"user_id" valid:"required~ユーザーIDは必須です。"`
}

type RequestYearIncomeAndDeductionData struct {
	UserId string `json:"user_id" valid:"required~ユーザーIDは必須です。"`
}

// TotalAmount, DeductionAmount, TakeHomeAmountは0の値でも許容させるために
type RequestInsertIncomeData struct {
	PaymentDate     string `json:"payment_date" valid:"required~報酬日付は必須です。"`
	Age             int    `json:"age" valid:"required~年齢は必須又は整数値のみです。"`
	Industry        string `json:"industry" valid:"required~業種は必須です。"`
	TotalAmount     string `json:"total_amount" valid:"required~総支給額は必須です。"`
	DeductionAmount string `json:"deduction_amount" valid:"required~差引額は必須です。"`
	TakeHomeAmount  string `json:"take_home_amount" valid:"required~手取額は必須です。"`
	Classification  string `json:"classification" valid:"required~分類は必須です。"`
	UserId          string `json:"user_id" valid:"required~ユーザーIDは必須です。"`
}

// TotalAmount, DeductionAmount, TakeHomeAmountは0の値でも許容させるために
type RequestUpdateIncomeData struct {
	IncomeForecastID string `json:"income_forecast_id" valid:"required~年収推移IDは必須です。"`
	PaymentDate      string `json:"payment_date" valid:"required~報酬日付は必須です。"`
	Age              int    `json:"age" valid:"required~年齢は必須又は整数値のみです。"`
	Industry         string `json:"industry" valid:"required~業種は必須です。"`
	TotalAmount      string `json:"total_amount" valid:"required~総支給額は必須です。"`
	DeductionAmount  string `json:"deduction_amount" valid:"required~差引額は必須です。"`
	TakeHomeAmount   string `json:"take_home_amount" valid:"required~手取額は必須です。"`
	UpdateUser       string `json:"update_user" valid:"required~更新者は必須です。"`
	Classification   string `json:"classification" valid:"required~分類は必須です。"`
}

type RequestDeleteIncomeData struct {
	IncomeForecastID string `json:"income_forecast_id" valid:"required~年収推移IDは必須です。"`
}

// パスワードのカスタムバリデーション関数
func validPassword(password string) bool {
	// 大文字が含まれているかをチェック
	hasUpperCase := regexp.MustCompile(`[A-Z]`).MatchString(password)
	// 特殊文字が含まれているかをチェック
	hasSpecialChar := regexp.MustCompile(`[.!?/-]`).MatchString(password)
	// パスワードの長さが8～24文字かをチェック
	isCorrectLength := len(password) >= 8 && len(password) <= 24

	// すべての条件が満たされているかどうかを返す
	return hasUpperCase && hasSpecialChar && isCorrectLength
}

func validDate(date string) bool {
	dateCase := regexp.MustCompile(`^[0-9]{4}-(0[1-9]|1[0-2])-(0[1-9]|1[0-9]|2[0-9]|3[0-1])$`).MatchString(date)

	// すべての条件が満たされているかどうかを返す
	return dateCase
}

func validInt(val string) bool {
	intCase := regexp.MustCompile(`^\d+$`).MatchString(val)

	// すべての条件が満たされているかどうかを返す
	return intCase
}

func (data RequestSignInData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	validArray := [2]bool{true, true}

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if password := validPassword(data.UserPassword); !password && data.UserPassword != "" {
		var flag bool = true
		for i := range errorMessagesList {
			if errorMessagesList[i].Field == UserPassword {
				errorMessagesList[i].Message = "パスワードの形式が間違っています。"
				flag = false
			}
		}

		if flag {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   UserPassword,
				Message: "パスワードの形式が間違っています。",
			})
		}
		validArray[1] = false
	}

	for _, validCheck := range validArray {
		if !validCheck {
			valid = false
		}
	}

	return valid, errorMessagesList
}

func (data RequestRefreshTokenData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	var valid bool = true

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if UserId := validInt(data.UserId); !UserId && data.UserId != "" {
		valid = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "user_id",
			Message: "ユーザーIDは整数値のみです。",
		})
	}

	return valid, errorMessagesList
}

func (data RequestRetryAuthEmail) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	return valid, errorMessagesList
}

func (data TemporayRequestSignUpData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if password := validPassword(data.UserPassword); !password && data.UserPassword != "" {
		var flag bool = true
		for i := range errorMessagesList {
			if errorMessagesList[i].Field == UserPassword {
				errorMessagesList[i].Message = "パスワードの形式が間違っています。"
				flag = false
			}
		}

		if flag {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   UserPassword,
				Message: "パスワードの形式が間違っています。",
			})
		}
		valid = false
	}

	return valid, errorMessagesList
}

func (data RequestSignUpData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	return valid, errorMessagesList
}

func (data RequestSignInEditData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	validArray := [2]bool{true, true}

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if UserId := validInt(data.UserId); !UserId && data.UserId != "" {
		validArray[0] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "user_id",
			Message: "ユーザーIDは整数値のみです。",
		})
	}

	if password := validPassword(data.UserPassword); !password && data.UserPassword != "" {
		var flag bool = true
		for i := range errorMessagesList {
			if errorMessagesList[i].Field == UserPassword {
				errorMessagesList[i].Message = "パスワードの形式が間違っています。"
				flag = false
			}
		}

		if flag {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   UserPassword,
				Message: "パスワードの形式が間違っています。",
			})
		}
		validArray[1] = false
	}

	for _, validCheck := range validArray {
		if !validCheck {
			valid = false
		}
	}

	return valid, errorMessagesList
}

func (data RequestSignInDeleteData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	var valid bool = true

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if UserId := validInt(data.UserId); !UserId && data.UserId != "" {
		valid = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "user_id",
			Message: "ユーザーIDは整数値のみです。",
		})
	}

	return valid, errorMessagesList
}

func (data RequestSignOutData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	return valid, errorMessagesList
}

func (data RequestPriceManagementData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	return valid, errorMessagesList
}

func (data RequestYearIncomeAndDeductiontData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	validArray := [3]bool{true, true, true}

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if UserId := validInt(data.UserId); !UserId && data.UserId != "" {
		validArray[0] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "user_id",
			Message: "ユーザーIDは整数値のみです。",
		})
	}

	if date := validDate(data.StartDate); !date && data.StartDate != "" {
		validArray[1] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "start_date",
			Message: "開始日の形式が間違っています。",
		})
	}

	if date := validDate(data.EndDate); !date && data.EndDate != "" {
		validArray[2] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "end_date",
			Message: "終了日の形式が間違っています。",
		})
	}

	for _, validCheck := range validArray {
		if !validCheck {
			valid = false
		}
	}

	return valid, errorMessagesList
}

func (data RequestDateRangeData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	var valid bool = true

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if UserId := validInt(data.UserId); !UserId && data.UserId != "" {
		valid = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "user_id",
			Message: "ユーザーIDは整数値のみです。",
		})
	}

	return valid, errorMessagesList
}

func (data RequestYearIncomeAndDeductionData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	var valid bool = true

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if UserId := validInt(data.UserId); !UserId && data.UserId != "" {
		valid = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "user_id",
			Message: "ユーザーIDは整数値のみです。",
		})
	}

	return valid, errorMessagesList
}

func (data RequestInsertIncomeData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	validArray := [5]bool{true, true, true, true, true}

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if date := validDate(data.PaymentDate); !date && data.PaymentDate != "" {
		validArray[0] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "payment_date",
			Message: "給料支給日の形式が間違っています。",
		})
	}

	if TotalAmount := validInt(data.TotalAmount); !TotalAmount && data.TotalAmount != "" {
		validArray[1] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "total_amount",
			Message: "総支給額で数値文字列以外は無効です。",
		})
	}

	if DeductionAmount := validInt(data.DeductionAmount); !DeductionAmount && data.DeductionAmount != "" {
		validArray[2] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "deduction_amount",
			Message: "差引額で数値文字列以外は無効です。",
		})
	}

	if TakeHomeAmount := validInt(data.TakeHomeAmount); !TakeHomeAmount && data.TakeHomeAmount != "" {
		validArray[3] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "take_home_amount",
			Message: "手取額で数値文字列以外は無効です。",
		})
	}

	if UserId := validInt(data.UserId); !UserId && data.UserId != "" {
		validArray[4] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "user_id",
			Message: "ユーザーIDは整数値のみです。",
		})
	}

	for _, validCheck := range validArray {
		if !validCheck {
			valid = false
		}
	}

	return valid, errorMessagesList
}

func (data RequestUpdateIncomeData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages
	validArray := [4]bool{true, true, true, true}

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	if date := validDate(data.PaymentDate); !date && data.PaymentDate != "" {
		validArray[0] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "payment_date",
			Message: "給料支給日の形式が間違っています。",
		})
	}

	if TotalAmount := validInt(data.TotalAmount); !TotalAmount && data.TotalAmount != "" {
		validArray[1] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "total_amount",
			Message: "総支給額で数値文字列以外は無効です。",
		})
	}

	if DeductionAmount := validInt(data.DeductionAmount); !DeductionAmount && data.DeductionAmount != "" {
		validArray[2] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "deduction_amount",
			Message: "差引額で数値文字列以外は無効です。",
		})
	}

	if TakeHomeAmount := validInt(data.TakeHomeAmount); !TakeHomeAmount && data.TakeHomeAmount != "" {
		validArray[3] = false
		errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
			Field:   "take_home_amount",
			Message: "手取額で数値文字列以外は無効です。",
		})
	}

	for _, validCheck := range validArray {
		if !validCheck {
			valid = false
		}
	}

	return valid, errorMessagesList
}

func (data RequestDeleteIncomeData) Validate() (bool, []utils.ErrorMessages) {
	var errorMessagesList []utils.ErrorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, utils.ErrorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	return valid, errorMessagesList
}

// func (data SignInValidation) Validate() error {
// 	//NOTE: 日本語のエラー文が不要で、デフォルトの英語のエラー文で必要十分である場合、`.Error("xxx")`は不要でOK
// 	return validation.ValidateStruct(&data,
// 		validation.Field(
// 			&data.UserId,
// 			validation.Required.Error(fmt.Sprintf("%sは必須です。", userId)),
// 			validation.By(intCheck),
// 		),
// 		// validation.Field(
// 		// 	&data.UserName,
// 		// 	validation.Required.Error("著者名は必須項目です。"),
// 		// 	// validation.RuneLength(1, 50).Error("著者名は 1文字 以上 50文字 以内です。"),
// 		// ),
// 		// validation.Field(
// 		// 	&data.UserPassword,
// 		// 	validation.Required.Error("価格は必須項目です。"),
// 		// 	// validation.Max(1000000.0).Error("価格は 1,000,000円 以下で指定してください。"),
// 		// 	// validation.Min(1.0).Error("価格は 1円 以上で指定してください。"),
// 		// ),
// 	)
// }

// バリデーションで以下の場合

// [
//     {
//         "field": "user_id",
//         "message": "ユーザーIDは必須又は整数値のみです"
//     },
//     {
//         "field": "user_password",
//         "message": "パスワードは必須です。"
//     },
// ]

// {
// 	"field": "user_password",
// 	"message": "パスワードの形式が間違っています。"
// }
