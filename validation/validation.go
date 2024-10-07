// validation/singin_validation.go
package validation

import (

	// "server/models"

	"regexp"

	"github.com/asaskevich/govalidator"
)

const UserPassword = "user_password"

// type (
// 	singInValidationFetcher interface {
// 		intCheck(val interface{}) error
// 		Validate() error
// 	}

// 	SinginValidation struct {
// 		models.RequestSingInData
// 	}

// 	singInValidation struct{}
// )

// func NewSingInValidationFetcher() singInValidationFetcher {
// 	return &singInValidation{}
// }

type RequestSingInData struct {
	UserId       int    `json:"user_id" valid:"required~ユーザーIDは必須です"`
	UserName     string `json:"user_name" valid:"required~ユーザー名は必須です,email~正しいメールアドレス形式である必要があります"`
	UserPassword string `json:"user_password" valid:"required~パスワードは必須です"`
}

type RequestSingUpData struct {
	UserName     string `json:"user_name" valid:"required~ユーザー名は必須です,email~正しいメールアドレス形式である必要があります"`
	UserPassword string `json:"user_password" valid:"required~パスワードは必須です"`
	NickName     string `json:"nick_name" valid:"required~ニックネームは必須です"`
}

type RequestSingInEditData struct {
	UserId       int    `json:"user_id" valid:"required~ユーザーIDは必須です"`
	UserName     string `json:"user_name" valid:"email~正しいメールアドレス形式である必要があります"`
	UserPassword string `json:"user_password"`
}

type RequestSingInDeleteData struct {
	UserId int `json:"user_id" valid:"required~ユーザーIDは必須です"`
}

type RequestRefreshTokenData struct {
	UserId int `json:"user_id" valid:"required~ユーザーIDは必須です"`
}

type RequestPriceManagementData struct {
	MoneyReceived string `json:"money_received" valid:"int~月の収入は整数値のみです"`
	Bouns         string `json:"bouns" valid:"int~ボーナスは整数値のみです"`
	FixedCost     string `json:"fixed_cost" valid:"int~月の収入は整数値のみです"`
	Loan          string `json:"loan" valid:"int~ローンは整数値のみです"`
	Private       string `json:"private" valid:"int~プライベートは整数値のみです"`
	Insurance     string `json:"insurance" valid:"int~保険は整数値のみです"`
}

type errorMessages struct {
	Field   string `json:"field"`
	Message string `json:"message"`
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

func (data RequestSingInData) Validate() (bool, []errorMessages) {
	var errorMessagesList []errorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, errorMessages{
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
			errorMessagesList = append(errorMessagesList, errorMessages{
				Field:   UserPassword,
				Message: "パスワードの形式が間違っています。",
			})
		}
		valid = false
	}

	return valid, errorMessagesList
}

func (data RequestRefreshTokenData) Validate() (bool, []errorMessages) {
	var errorMessagesList []errorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, errorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	return valid, errorMessagesList
}

func (data RequestSingUpData) Validate() (bool, []errorMessages) {
	var errorMessagesList []errorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, errorMessages{
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
			errorMessagesList = append(errorMessagesList, errorMessages{
				Field:   UserPassword,
				Message: "パスワードの形式が間違っています。",
			})
		}
		valid = false
	}

	return valid, errorMessagesList
}

func (data RequestSingInEditData) Validate() (bool, []errorMessages) {
	var errorMessagesList []errorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, errorMessages{
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
			errorMessagesList = append(errorMessagesList, errorMessages{
				Field:   UserPassword,
				Message: "パスワードの形式が間違っています。",
			})
		}
		valid = false
	}

	return valid, errorMessagesList
}

func (data RequestSingInDeleteData) Validate() (bool, []errorMessages) {
	var errorMessagesList []errorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, errorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	return valid, errorMessagesList
}

func (data RequestPriceManagementData) Validate() (bool, []errorMessages) {
	var errorMessagesList []errorMessages

	valid, err := govalidator.ValidateStruct(data)

	if err != nil {
		errorMap := govalidator.ErrorsByField(err)
		for field, msg := range errorMap {
			errorMessagesList = append(errorMessagesList, errorMessages{
				Field:   field,
				Message: msg,
			})
		}
	}

	return valid, errorMessagesList
}

// func (data SingInValidation) Validate() error {
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
//         "message": "ユーザーIDは必須です"
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
