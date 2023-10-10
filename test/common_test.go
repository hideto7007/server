// test/common_test.go
package common

import (
	"net/http/httptest"
	"server/common"
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

		var common common.CommonFetcher = common.NewCommonFetcher()
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

	t.Run("error case IntgetPrameter()", func(t *testing.T) {
		// テストケース2: 整数に変換できないクエリーパラメータ
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Request = httptest.NewRequest("GET", "/?param1=42&param2=notanumber", nil)

		var common common.CommonFetcher = common.NewCommonFetcher()
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
}

func TestTimeToStr(t *testing.T) {
	t.Run("success TimeToStr()", func(t *testing.T) {

		var common common.CommonFetcher = common.NewCommonFetcher()
		dateTime := time.Date(2022, time.December, 23, 0, 0, 0, 0, time.FixedZone("", 0))
		result := common.TimeToStr(dateTime)

		assert.Equal(t, "2022-12-23", result)

		t.Logf("dataTime repalce : '%s'", result)
	})
	t.Run("error case1 TimeToStr()", func(t *testing.T) {

		var common common.CommonFetcher = common.NewCommonFetcher()
		dateTime := time.Date(2022, time.February, 30, 0, 0, 0, 0, time.FixedZone("", 0))
		result := common.TimeToStr(dateTime)

		assert.NotEqual(t, "2022-02-30", result)

		t.Logf("non existing date → 2022-02-30 : '%s'", result)
	})
	t.Run("error case2 TimeToStr()", func(t *testing.T) {

		var common common.CommonFetcher = common.NewCommonFetcher()
		var emptyTime time.Time
		result := common.TimeToStr(emptyTime)

		assert.Equal(t, "0001-01-01", result)

		t.Logf("0001-01-01: %s", result)
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
