package common

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIntgetPrameter(t *testing.T) {
	t.Run("success IntgetPrameter()", func(t *testing.T) {
		// テストケース1: 正常な整数のクエリーパラメータ
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?param1=42&param2=100", nil)

		var common CommonFetcher = NewCommonFetcher()
		paramMap, err := common.IntgetPrameter(c, "param1", "param2")

		if err != nil {
			t.Errorf("Expected an nil, but got no nil: %v", err)
		}

		if paramMap["param1"] != 42 || paramMap["param2"] != 100 {
			t.Errorf("Expected paramMap to be {param1: 42, param2: 100}, but got: %v", paramMap)
		}

		assert.Equal(t, 42, paramMap["param1"])
		assert.Equal(t, 100, paramMap["param2"])
		assert.Equal(t, nil, err)

		t.Logf("paramMap['param1']: %d", paramMap["param1"])
		t.Logf("paramMap['param2']: %d", paramMap["param2"])
		t.Logf("err: %v", err)
	})
	t.Run("success IntgetPrameter zero start 02()", func(t *testing.T) {
		// テストケース2: 0から始まる場合は削除して整数値に変換する
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?param1=042&param2=00100", nil)

		var common CommonFetcher = NewCommonFetcher()
		paramMap, err := common.IntgetPrameter(c, "param1", "param2")

		if err != nil {
			t.Errorf("Expected an nil, but got no nil: %v", err)
		}

		if paramMap["param1"] != 42 || paramMap["param2"] != 100 {
			t.Errorf("Expected paramMap to be {param1: 42, param2: 100}, but got: %v", paramMap)
		}

		assert.Equal(t, 42, paramMap["param1"])
		assert.Equal(t, 100, paramMap["param2"])
		assert.Equal(t, nil, err)

		t.Logf("paramMap['param1']: %d", paramMap["param1"])
		t.Logf("paramMap['param2']: %d", paramMap["param2"])
		t.Logf("err: %v", err)
	})
	t.Run("success IntgetPrameter zero start 02()", func(t *testing.T) {
		// テストケース3: 異なるパラメータを呼び出した際は0を取得する
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?param1=42&param2=100", nil)

		var common CommonFetcher = NewCommonFetcher()
		paramMap, err := common.IntgetPrameter(c, "param1", "param2", "param3")

		if err != nil {
			t.Errorf("Expected an nil, but got no nil: %v", err)
		}

		if paramMap["param1"] != 42 || paramMap["param2"] != 100 || paramMap["param3"] != 0 {
			t.Errorf("Expected paramMap to be {param1: 42, param2: 100, param3: 0}, but got: %v", paramMap)
		}

		assert.Equal(t, 42, paramMap["param1"])
		assert.Equal(t, 100, paramMap["param2"])
		assert.Equal(t, 0, paramMap["param3"])
		assert.Equal(t, nil, err)

		t.Logf("paramMap['param1']: %d", paramMap["param1"])
		t.Logf("paramMap['param2']: %d", paramMap["param2"])
		t.Logf("paramMap['param3']: %d", paramMap["param2"])
		t.Logf("err: %v", err)
	})

	t.Run("error case IntgetPrameter string notanumber error()", func(t *testing.T) {
		// テストケース4: 整数に変換できないクエリーパラメータ
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?param1=42&param2=notanumber", nil)

		var common CommonFetcher = NewCommonFetcher()
		paramMap, err := common.IntgetPrameter(c, "param1", "param2")

		if err == nil {
			t.Error("Expected an error, but got nil")
		}

		assert.Empty(t, paramMap)
		expectedErrorMessage := "strconv.Atoi: parsing \"notanumber\": invalid syntax"
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("paramMap[]: %v", paramMap)
		t.Logf("err: %s", err)
	})
	t.Run("error case IntgetPrameter string hoge error()", func(t *testing.T) {
		// テストケース5: 整数に変換できないクエリーパラメータ
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?param1=hoge&param2=43", nil)

		var common CommonFetcher = NewCommonFetcher()
		paramMap, err := common.IntgetPrameter(c, "param1", "param2")

		if err == nil {
			t.Error("Expected an error, but got nil")
		}

		assert.Empty(t, paramMap)
		expectedErrorMessage := "strconv.Atoi: parsing \"hoge\": invalid syntax"
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("paramMap[]: %v", paramMap)
		t.Logf("err: %s", err)
	})
	t.Run("error case IntgetPrameter float 2.1 error()", func(t *testing.T) {
		// テストケース6: 整数に変換できないクエリーパラメータ
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?param1=2.1&param2=43", nil)

		var common CommonFetcher = NewCommonFetcher()
		paramMap, err := common.IntgetPrameter(c, "param1", "param2")

		if err == nil {
			t.Error("Expected an error, but got nil")
		}

		assert.Empty(t, paramMap)
		expectedErrorMessage := "strconv.Atoi: parsing \"2.1\": invalid syntax"
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("paramMap[]: %v", paramMap)
		t.Logf("err: %s", err)
	})
	t.Run("error case IntgetPrameter float 4.3 error()", func(t *testing.T) {
		// テストケース7: 整数に変換できないクエリーパラメータ
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?param1=21&param2=4.3", nil)

		var common CommonFetcher = NewCommonFetcher()
		paramMap, err := common.IntgetPrameter(c, "param1", "param2")

		if err == nil {
			t.Error("Expected an error, but got nil")
		}

		assert.Empty(t, paramMap)
		expectedErrorMessage := "strconv.Atoi: parsing \"4.3\": invalid syntax"
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("paramMap[]: %v", paramMap)
		t.Logf("err: %s", err)
	})
}

func TestTimeToStr(t *testing.T) {
	t.Run("success TimeToStr()", func(t *testing.T) {
		// テストケース1: 日付を文字列に変換
		var common CommonFetcher = NewCommonFetcher()
		dateTime := time.Date(2022, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0))
		result := common.TimeToStr(dateTime)

		assert.Equal(t, "2022-12-23", result)

		t.Logf("dataTime repalce : '%s'", result)
	})
	t.Run("success TimeToStr()", func(t *testing.T) {
		// テストケース2: 存在しない年月日の場合、適切な日付の文字列にして変換できること
		var common CommonFetcher = NewCommonFetcher()
		dateTime := time.Date(2022, time.February, 30, 0, 0, 0, 0, time.FixedZone("", 0))
		result := common.TimeToStr(dateTime)

		assert.NotEqual(t, "2022-02-30", result)

		t.Logf("non existing date → 2022-02-30 : '%s'", result)
	})
	t.Run("error case TimeToStr()", func(t *testing.T) {
		// テストケース3: 空の場合は0001-01-01になること
		var common CommonFetcher = NewCommonFetcher()
		var emptyTime time.Time
		result := common.TimeToStr(emptyTime)

		assert.Equal(t, "0001-01-01", result)

		t.Logf("0001-01-01: %s", result)
	})
}

func TestStrToTime(t *testing.T) {
	t.Run("success StrToTime()", func(t *testing.T) {
		// テストケース1: 文字列を日付に変換する
		var common CommonFetcher = NewCommonFetcher()
		strDate := "2023-10-14"
		result, err := common.StrToTime(strDate)

		assert.NoError(t, err)

		t.Log(result)

		assert.Equal(t, time.Time(time.Date(2023, time.October, 14, 0, 0, 0, 0, time.UTC)), result)

		t.Logf("str to time repalce : '%s'", result)
	})
	t.Run("error case1 StrToTime()", func(t *testing.T) {
		// テストケース2: 文字列"hoge"は日付変換出来ないこと
		var common CommonFetcher = NewCommonFetcher()
		strDate := "hoge"
		_, err := common.StrToTime(strDate)

		assert.Error(t, err)

		expectedErrorMessage := `parsing time "hoge" as "2006-01-02": cannot parse "hoge" as "2006"`
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("no replace time 1 : '%s'", err)
	})
	t.Run("error case2 StrToTime()", func(t *testing.T) {
		// テストケース3: 文字列"3456"は日付変換出来ないこと
		var common CommonFetcher = NewCommonFetcher()
		strDate := "3456"
		_, err := common.StrToTime(strDate)

		assert.Error(t, err)

		expectedErrorMessage := `parsing time "3456" as "2006-01-02": cannot parse "" as "-"`
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("no replace time 2 : '%s'", err)
	})
}

func TestStrToInt(t *testing.T) {
	t.Run("success StrToInt()", func(t *testing.T) {
		// テストケース1: 文字列"6543"は整数値に変換する
		var common CommonFetcher = NewCommonFetcher()
		strDate := "6543"
		result, err := common.StrToInt(strDate)

		assert.NoError(t, err)

		assert.Equal(t, 6543, result)

		t.Logf("string to int replace : '%d'", result)
	})
	t.Run("error case1 StrToInt()", func(t *testing.T) {
		// テストケース2: 文字列"hoge"は整数値に変換出来ないこと
		var common CommonFetcher = NewCommonFetcher()
		strDate := "hoge"
		_, err := common.StrToInt(strDate)

		assert.Error(t, err)

		expectedErrorMessage := `strconv.Atoi: parsing "hoge": invalid syntax`
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("no replace int 1 : '%s'", err)
	})
	t.Run("error case2 StrToInt()", func(t *testing.T) {
		// テストケース3: 文字列"34hoge56"は整数値に変換出来ないこと
		var common CommonFetcher = NewCommonFetcher()
		strDate := "34hoge56"
		_, err := common.StrToInt(strDate)

		assert.Error(t, err)

		expectedErrorMessage := `strconv.Atoi: parsing "34hoge56": invalid syntax`
		assert.EqualError(t, err, expectedErrorMessage)

		t.Logf("no replace int 2 : '%s'", err)
	})
}

func TestIntToStr(t *testing.T) {
	t.Run("success IntToStr()", func(t *testing.T) {
		// テストケース1: 整数値6543は文字列に変換する
		var common CommonFetcher = NewCommonFetcher()
		intDate := 6543
		result := common.IntToStr(intDate)

		assert.Equal(t, "6543", result)

		t.Logf("int to string replace : '%s'", result)
	})
}

// func TestIntgetPrameter(t *testing.T) {
// 	t.Run("success IntgetPrameter()", func(t *testing.T) {
// 		// テストケース1: 正常な整数のクエリーパラメータ
// 		req := httptest.NewRequest("GET", "/?param1=42&param2=100", nil)

// 		paramMap, err := common.IntgetPrameter(req, "param1", "param2")

// 		if err != nil {
// 			t.Errorf("Expected an nil, but got no nil: %v", err)
// 		}

// 		if paramMap["param1"] != 42 || paramMap["param2"] != 100 {
// 			t.Errorf("Expected paramMap to be {param1: 42, param2: 100}, but got: %v", paramMap)
// 		}

// 		assert.Equal(t, 42, paramMap["param1"])
// 		assert.Equal(t, 100, paramMap["param2"])
// 		assert.Equal(t, nil, err)

// 		t.Logf("paramMap['param1']: %d", paramMap["param1"])
// 		t.Logf("paramMap['param2']: %d", paramMap["param2"])
// 		t.Logf("err: %v", err)
// 	})

// 	t.Run("error case IntgetPrameter()", func(t *testing.T) {
// 		// テストケース2: 整数に変換できないクエリーパラメータ
// 		req := httptest.NewRequest("GET", "/?param1=42&param2=notanumber", nil)

// 		paramMap, err := common.IntgetPrameter(req, "param1", "param2")

// 		if err == nil {
// 			t.Error("Expected an error, but got nil")
// 		}

// 		assert.Empty(t, paramMap)
// 		assert.EqualError(t, errors.New("notanumber"), "notanumber")

// 		t.Logf("paramMap[]: %v", paramMap)
// 		t.Logf("err: %s", err)
// 	})
// }
