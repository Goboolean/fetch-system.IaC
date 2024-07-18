package mapper_test

import (
	"testing"

	"github.com/Goboolean/fetch-system.IaC/pkg/influx/mapper"
	"github.com/stretchr/testify/suite"
)

type Custom int

const (
	A Custom = iota + 1
	B
	C
	D
)

type StructToPointTestSuite struct {
	suite.Suite
}

func (suite *StructToPointTestSuite) TestStructToPoint_WhenDepthIsOneAndContainsOnlyPrimitiveTypes() {
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
	suite.NoError(err)
	suite.NotNil(out)
	suite.Equal(100, out["testNumber"])
	suite.Equal("test", out["testString"])
	suite.Equal(float32(3.14), out["testFloat"])
	suite.Equal(true, out["testBool"])
}

func (suite *StructToPointTestSuite) TestStructToPoint_whenDepthIsTwoOrMoreAndContainsOnlyPrimitiveTypes() {
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
	suite.NoError(err)
	suite.NotNil(out)
	suite.Equal(100, out["testNumber"])
	suite.Equal("test", out["testString"])
	suite.Equal(float32(3.14), out["testFloat"])
	suite.Equal(true, out["testBool"])

	suite.Equal(100, out["testStruct.testNumber"])
	suite.Equal("test", out["testStruct.testString"])
	suite.Equal(float32(3.14), out["testStruct.testFloat"])
	suite.Equal(true, out["testStruct.testBool"])
}

func (suite *StructToPointTestSuite) TestStructToPoint_WhenDepthIsOneContainsMapAndSlice() {
	// arrange
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
	// act
	out, err := mapper.StructToPoint(in)
	// assert
	suite.NoError(err)
	suite.NotNil(out)
	suite.Equal(100, out["testNumber"])
	suite.Equal("test", out["testString"])
	suite.Equal(float32(3.14), out["testFloat"])
	suite.Equal(true, out["testBool"])

	suite.Equal(1, out["testMapStringNumber..test1"])
	suite.Equal(2, out["testMapStringNumber..test2"])
	suite.Equal(3, out["testMapStringNumber..test3"])

	suite.Equal(1, out["testMapEnumToNumber..1"])
	suite.Equal(2, out["testMapEnumToNumber..2"])
	suite.Equal(3, out["testMapEnumToNumber..3"])
	suite.Equal(4, out["testMapEnumToNumber..4"])

	suite.Equal(1, out["testSlice..0"])
	suite.Equal(2, out["testSlice..1"])
	suite.Equal(3, out["testSlice..2"])
}
func (suite *StructToPointTestSuite) TestStructToPoint_WhenDepthIsTwoOrMoreAndContainsMapAndSlice() {

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
	suite.NoError(err)
	suite.NotNil(out)
	suite.Equal(100, out["testNumber"])
	suite.Equal("test", out["testString"])
	suite.Equal(float32(3.14), out["testFloat"])
	suite.Equal(true, out["testBool"])

	suite.Equal(1, out["testStruct.testMapStringNumber..test1"])
	suite.Equal(2, out["testStruct.testMapStringNumber..test2"])
	suite.Equal(3, out["testStruct.testMapStringNumber..test3"])

	suite.Equal(1, out["testStruct.testMapEnumToNumber..1"])
	suite.Equal(2, out["testStruct.testMapEnumToNumber..2"])
	suite.Equal(3, out["testStruct.testMapEnumToNumber..3"])

	suite.Equal(1, out["testStruct.testSlice..0"])
	suite.Equal(2, out["testStruct.testSlice..1"])
	suite.Equal(3, out["testStruct.testSlice..2"])
}

func (suite *StructToPointTestSuite) TestStructToPoint_WhenDepthIsTwoOrMoreAndSliceElementsAreStructs() {

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
	suite.NoError(err)

	suite.Equal(1, out["testSlice..0.testNumber"])
	suite.Equal("1", out["testSlice..0.testString"])

	suite.Equal(2, out["testSlice..1.testNumber"])
	suite.Equal("2", out["testSlice..1.testString"])

	suite.Equal(3, out["testSlice..2.testNumber"])
	suite.Equal("3", out["testSlice..2.testString"])
}

func (suite *StructToPointTestSuite) TestStructToPoint_WhenDepthIsTwoOrMoreAndContainsPointers() {
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
	suite.NoError(err)
	suite.NotNil(out)
	suite.Equal(100, out["testNumber"])
	suite.Equal("test", out["testString"])
	suite.Equal(float32(3.14), out["testFloat"])
	suite.Equal(true, out["testBool"])

	suite.Equal(100, out["testStruct.testNumber"])
	suite.Equal("test", out["testStruct.testString"])
	suite.Equal(float32(3.14), out["testStruct.testFloat"])
	suite.Equal(true, out["testStruct.testBool"])
}

func (suite *StructToPointTestSuite) TestStructToPoint_WhenPassedAsPointerType() {
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

	in := &TestNestedStruct{
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
	suite.NoError(err)
	suite.NotNil(out)
	suite.Equal(100, out["testNumber"])
	suite.Equal("test", out["testString"])
	suite.Equal(float32(3.14), out["testFloat"])
	suite.Equal(true, out["testBool"])

	suite.Equal(100, out["testStruct.testNumber"])
	suite.Equal("test", out["testStruct.testString"])
	suite.Equal(float32(3.14), out["testStruct.testFloat"])
	suite.Equal(true, out["testStruct.testBool"])
}

func TestStructToPoint(t *testing.T) {
	suite.Run(t, new(StructToPointTestSuite))
}
