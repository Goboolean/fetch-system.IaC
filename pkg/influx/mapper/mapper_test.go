package mapper_test

import (
	"testing"

	"github.com/Goboolean/fetch-system.IaC/pkg/influx/mapper"
	"github.com/stretchr/testify/assert"
)

type Custom int

const (
	A Custom = iota + 1
	B
	C
	D
)

func TestMapStructToPoint(t *testing.T) {
	t.Run("depth가 1이고 기본 타입으로만 이루어진 구조체 테스트", func(t *testing.T) {
		//arrange
		type TestStruct struct {
			TestNumber int     `name:"testNumber"`
			TestString string  `name:"testString"`
			TestFloat  float32 `name:"testFloat"`
			TestBool   bool    `name:"testBool"`
		}

		in := TestStruct{
			TestNumber: 100,
			TestString: "test",
			TestFloat:  3.14,
			TestBool:   true,
		}
		//act
		out, err := mapper.StructToPoint(in)
		//assert
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, 100, out["testNumber"])
		assert.Equal(t, "test", out["testString"])
		assert.Equal(t, float32(3.14), out["testFloat"])
		assert.Equal(t, true, out["testBool"])
	})

	t.Run("depth가 2이상이고 기본 타입으로만 이루어진 구조체 테스트", func(t *testing.T) {
		//arrange
		type TestStruct struct {
			TestNumber int     `name:"testNumber"`
			TestString string  `name:"testString"`
			TestFloat  float32 `name:"testFloat"`
			TestBool   bool    `name:"testBool"`
		}

		type TestNestedStruct struct {
			TestNumber       int        `name:"testNumber"`
			TestString       string     `name:"testString"`
			TestFloat        float32    `name:"testFloat"`
			TestBool         bool       `name:"testBool"`
			TestNestedStruct TestStruct `name:"testStruct"`
		}

		in := TestNestedStruct{
			TestNumber: 100,
			TestString: "test",
			TestFloat:  3.14,
			TestBool:   true,
			TestNestedStruct: TestStruct{
				TestNumber: 100,
				TestString: "test",
				TestFloat:  3.14,
				TestBool:   true,
			},
		}
		//act
		out, err := mapper.StructToPoint(in)
		//assert
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, 100, out["testNumber"])
		assert.Equal(t, "test", out["testString"])
		assert.Equal(t, float32(3.14), out["testFloat"])
		assert.Equal(t, true, out["testBool"])

		assert.Equal(t, 100, out["testStruct.testNumber"])
		assert.Equal(t, "test", out["testStruct.testString"])
		assert.Equal(t, float32(3.14), out["testStruct.testFloat"])
		assert.Equal(t, true, out["testStruct.testBool"])
	})
	t.Run("depth가 1이고 맵과 슬라이스가 포함된 구조체 테스트", func(t *testing.T) {
		//arrange

		type TestStruct struct {
			TestNumber          int            `name:"testNumber"`
			TestString          string         `name:"testString"`
			TestFloat           float32        `name:"testFloat"`
			TestBool            bool           `name:"testBool"`
			TestMapStringNumber map[string]int `name:"testMapStringNumber"`
			TestMapEnumToNumber map[Custom]int `name:"testMapEnumToNumber"`
			TestSlice           []int          `name:"testSlice"`
		}

		in := TestStruct{
			TestNumber: 100,
			TestString: "test",
			TestFloat:  3.14,
			TestBool:   true,
			TestMapStringNumber: map[string]int{
				"test1": 1,
				"test2": 2,
				"test3": 3,
			},
			TestMapEnumToNumber: map[Custom]int{
				A: 1,
				B: 2,
				C: 3,
				D: 4,
			},
			TestSlice: []int{1, 2, 3},
		}
		//act
		out, err := mapper.StructToPoint(in)
		//assert
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, 100, out["testNumber"])
		assert.Equal(t, "test", out["testString"])
		assert.Equal(t, float32(3.14), out["testFloat"])
		assert.Equal(t, true, out["testBool"])

		assert.Equal(t, 1, out["testMapStringNumber..test1"])
		assert.Equal(t, 2, out["testMapStringNumber..test2"])
		assert.Equal(t, 3, out["testMapStringNumber..test3"])

		assert.Equal(t, 1, out["testMapEnumToNumber..1"])
		assert.Equal(t, 2, out["testMapEnumToNumber..2"])
		assert.Equal(t, 3, out["testMapEnumToNumber..3"])
		assert.Equal(t, 4, out["testMapEnumToNumber..4"])

		assert.Equal(t, 1, out["testSlice..0"])
		assert.Equal(t, 2, out["testSlice..1"])
		assert.Equal(t, 3, out["testSlice..2"])
	})

	t.Run("depth가 2이상이고 맵과 슬라이스가 포함된 구조체 테스트", func(t *testing.T) {
		// arrange
		type TestStruct struct {
			TestMapStringNumber map[string]int `name:"testMapStringNumber"`
			TestMapEnumToNumber map[Custom]int `name:"testMapEnumToNumber"`
			TestSlice           []int          `name:"testSlice"`
		}

		type TestNestedStruct struct {
			TestNumber       int        `name:"testNumber"`
			TestString       string     `name:"testString"`
			TestFloat        float32    `name:"testFloat"`
			TestBool         bool       `name:"testBool"`
			TestNestedStruct TestStruct `name:"testStruct"`
		}

		in := TestNestedStruct{
			TestNumber: 100,
			TestString: "test",
			TestFloat:  3.14,
			TestBool:   true,
			TestNestedStruct: TestStruct{
				TestMapStringNumber: map[string]int{
					"test1": 1,
					"test2": 2,
					"test3": 3,
				},
				TestMapEnumToNumber: map[Custom]int{
					A: 1,
					B: 2,
					C: 3,
				},
				TestSlice: []int{1, 2, 3},
			},
		}
		//act
		out, err := mapper.StructToPoint(in)
		//assert
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, 100, out["testNumber"])
		assert.Equal(t, "test", out["testString"])
		assert.Equal(t, float32(3.14), out["testFloat"])
		assert.Equal(t, true, out["testBool"])

		assert.Equal(t, 1, out["testStruct.testMapStringNumber..test1"])
		assert.Equal(t, 2, out["testStruct.testMapStringNumber..test2"])
		assert.Equal(t, 3, out["testStruct.testMapStringNumber..test3"])

		assert.Equal(t, 1, out["testStruct.testMapEnumToNumber..1"])
		assert.Equal(t, 2, out["testStruct.testMapEnumToNumber..2"])
		assert.Equal(t, 3, out["testStruct.testMapEnumToNumber..3"])

		assert.Equal(t, 1, out["testStruct.testSlice..0"])
		assert.Equal(t, 2, out["testStruct.testSlice..1"])
		assert.Equal(t, 3, out["testStruct.testSlice..2"])
	})

	t.Run("depth가 2이상이고 슬라이스의 원소 타입이 구조체인 경우 ", func(t *testing.T) {
		// arrange
		type TestStruct struct {
			TestNumber int    `name:"testNumber"`
			TestString string `name:"testString"`
		}

		type TestNestedStruct struct {
			TestSlice []TestStruct `name:"testSlice"`
		}

		in := TestNestedStruct{
			TestSlice: []TestStruct{
				{
					TestNumber: 1,
					TestString: "1",
				},
				{
					TestNumber: 2,
					TestString: "2",
				},
				{
					TestNumber: 3,
					TestString: "3",
				},
			},
		}

		//act
		out, err := mapper.StructToPoint(in)
		//assert
		assert.NoError(t, err)

		assert.Equal(t, 1, out["testSlice..0.testNumber"])
		assert.Equal(t, "1", out["testSlice..0.testString"])

		assert.Equal(t, 2, out["testSlice..1.testNumber"])
		assert.Equal(t, "2", out["testSlice..1.testString"])

		assert.Equal(t, 3, out["testSlice..2.testNumber"])
		assert.Equal(t, "3", out["testSlice..2.testString"])
	})

	t.Run("depth가 2이상이고 포인터가 있는 구조체 테스트", func(t *testing.T) {
		//arrange
		type TestStruct struct {
			TestNumber int     `name:"testNumber"`
			TestString string  `name:"testString"`
			TestFloat  float32 `name:"testFloat"`
			TestBool   bool    `name:"testBool"`
		}

		type TestNestedStruct struct {
			TestNumber       int         `name:"testNumber"`
			TestString       string      `name:"testString"`
			TestFloat        float32     `name:"testFloat"`
			TestBool         bool        `name:"testBool"`
			TestNestedStruct *TestStruct `name:"testStruct"`
		}

		in := TestNestedStruct{
			TestNumber: 100,
			TestString: "test",
			TestFloat:  3.14,
			TestBool:   true,
			TestNestedStruct: &TestStruct{
				TestNumber: 100,
				TestString: "test",
				TestFloat:  3.14,
				TestBool:   true,
			},
		}
		//act
		out, err := mapper.StructToPoint(in)
		//assert
		assert.NoError(t, err)
		assert.NotNil(t, out)
		assert.Equal(t, 100, out["testNumber"])
		assert.Equal(t, "test", out["testString"])
		assert.Equal(t, float32(3.14), out["testFloat"])
		assert.Equal(t, true, out["testBool"])

		assert.Equal(t, 100, out["testStruct.testNumber"])
		assert.Equal(t, "test", out["testStruct.testString"])
		assert.Equal(t, float32(3.14), out["testStruct.testFloat"])
		assert.Equal(t, true, out["testStruct.testBool"])
	})
}
