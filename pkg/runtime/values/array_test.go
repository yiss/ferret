package values_test

import (
	"encoding/json"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestArray(t *testing.T) {
	Convey("#constructor", t, func() {
		Convey("Should create an empty array", func() {
			arr := values.NewArray(10)

			So(arr.Length(), ShouldEqual, 0)
		})

		Convey("Should create an array, from passed values", func() {
			arr := values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
			)

			So(arr.Length(), ShouldEqual, 3)
		})
	})

	Convey(".MarshalJSON", t, func() {
		Convey("Should serialize empty array", func() {
			arr := values.NewArray(10)
			marshaled, err := arr.MarshalJSON()

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "[]")
		})

		Convey("Should serialize full array", func() {
			arr := values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
			)
			marshaled, err := json.Marshal(arr)

			So(err, ShouldBeNil)

			So(string(marshaled), ShouldEqual, "[1,2,3]")
		})
	})

	Convey(".Type", t, func() {
		Convey("Should return type", func() {
			arr := values.NewArray(1)

			So(arr.Type(), ShouldEqual, core.ArrayType)
		})
	})

	Convey(".Unwrap", t, func() {
		Convey("Should return a an array of unwrapped values", func() {
			arr := values.NewArrayWith(
				values.ZeroInt,
				values.ZeroInt,
			)

			for _, val := range arr.Unwrap().([]interface{}) {
				So(val, ShouldHaveSameTypeAs, 0)
			}
		})
	})

	Convey(".String", t, func() {
		Convey("Should return a string representation ", func() {
			arr := values.NewArrayWith(values.ZeroInt, values.ZeroInt)

			So(arr.String(), ShouldEqual, "[0,0]")
		})
	})

	Convey(".Compare", t, func() {
		Convey("It should return 1 for all non-array and non-object values", func() {
			arr := values.NewArrayWith(values.ZeroInt, values.ZeroInt)

			So(arr.Compare(values.None), ShouldEqual, 1)
			So(arr.Compare(values.ZeroInt), ShouldEqual, 1)
			So(arr.Compare(values.ZeroFloat), ShouldEqual, 1)
			So(arr.Compare(values.EmptyString), ShouldEqual, 1)
		})

		Convey("It should return -1 for all object values", func() {
			arr := values.NewArrayWith(values.ZeroInt, values.ZeroInt)
			obj := values.NewObject()

			So(arr.Compare(obj), ShouldEqual, -1)
		})

		Convey("It should return 0 when both arrays are empty", func() {
			arr1 := values.NewArray(1)
			arr2 := values.NewArray(1)

			So(arr1.Compare(arr2), ShouldEqual, 0)
		})

		Convey("It should return 1 when other array is empty", func() {
			arr1 := values.NewArrayWith(values.ZeroFloat)
			arr2 := values.NewArray(1)

			So(arr1.Compare(arr2), ShouldEqual, 1)
		})

		Convey("It should return 1 when values are bigger", func() {
			arr1 := values.NewArrayWith(values.NewInt(1))
			arr2 := values.NewArrayWith(values.ZeroInt)

			So(arr1.Compare(arr2), ShouldEqual, 1)
		})
	})

	Convey(".Length", t, func() {
		Convey("Should return 0 when empty", func() {
			arr := values.NewArray(1)

			So(arr.Length(), ShouldEqual, 0)
		})

		Convey("Should return greater than 0 when not empty", func() {
			arr := values.NewArrayWith(values.ZeroInt, values.ZeroInt)

			So(arr.Length(), ShouldEqual, 2)
		})
	})

	Convey(".ForEach", t, func() {
		Convey("Should iterate over elements", func() {
			arr := values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
			)
			counter := 0

			arr.ForEach(func(value core.Value, idx int) bool {
				counter += 1

				return true
			})

			So(counter, ShouldEqual, arr.Length())
		})

		Convey("Should not iterate when empty", func() {
			arr := values.NewArrayWith()
			counter := 0

			arr.ForEach(func(value core.Value, idx int) bool {
				counter += 1

				return true
			})

			So(counter, ShouldEqual, arr.Length())
		})

		Convey("Should break iteration when false returned", func() {
			arr := values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
			)
			threshold := 3
			counter := 0

			arr.ForEach(func(value core.Value, idx int) bool {
				counter += 1

				return value.Compare(values.NewInt(threshold)) == -1
			})

			So(counter, ShouldEqual, threshold)
		})
	})

	Convey(".Get", t, func() {
		Convey("Should return item by index", func() {
			arr := values.NewArrayWith(
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
			)

			el := arr.Get(1)

			So(el.Compare(values.NewInt(2)), ShouldEqual, 0)
		})

		Convey("Should return None when no value", func() {
			arr := values.NewArrayWith()

			el := arr.Get(1)

			So(el.Compare(values.None), ShouldEqual, 0)
		})
	})

	Convey(".Set", t, func() {
		Convey("Should set item by index", func() {
			arr := values.NewArrayWith(values.ZeroInt)

			err := arr.Set(0, values.NewInt(1))

			So(err, ShouldBeNil)
			So(arr.Length(), ShouldEqual, 1)
			So(arr.Get(0).Compare(values.NewInt(1)), ShouldEqual, 0)
		})

		Convey("Should return an error when index is out of bounds", func() {
			arr := values.NewArray(10)

			err := arr.Set(0, values.NewInt(1))

			So(err, ShouldNotBeNil)
			So(arr.Length(), ShouldEqual, 0)
		})
	})

	Convey(".Push", t, func() {
		Convey("Should add an item", func() {
			arr := values.NewArray(10)

			src := []core.Value{
				values.ZeroInt,
				values.ZeroInt,
				values.ZeroInt,
				values.ZeroInt,
				values.ZeroInt,
			}

			for _, val := range src {
				arr.Push(val)
			}

			So(arr.Length(), ShouldEqual, len(src))
		})
	})

	Convey(".Slice", t, func() {
		Convey("Should return a slice", func() {
			arr := values.NewArrayWith(
				values.NewInt(0),
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
			)

			s := arr.Slice(0, 1)

			So(len(s), ShouldEqual, 1)
			So(s[0].Compare(values.ZeroInt), ShouldEqual, 0)

			s2 := arr.Slice(2, arr.Length())

			So(len(s2), ShouldEqual, arr.Length()-2)
		})
	})

	Convey(".Insert", t, func() {
		Convey("Should insert an item in the middle of an array", func() {
			arr := values.NewArrayWith(
				values.NewInt(0),
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
			)

			lenBefore := arr.Length()

			arr.Insert(3, values.NewInt(100))

			lenAfter := arr.Length()

			So(lenAfter, ShouldBeGreaterThan, lenBefore)
			So(arr.Get(3), ShouldEqual, 100)
		})
	})

	Convey(".RemoveAt", t, func() {
		Convey("Should remove an item from the middle", func() {
			arr := values.NewArrayWith(
				values.NewInt(0),
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
			)

			lenBefore := arr.Length()

			arr.RemoveAt(3)

			lenAfter := arr.Length()

			So(lenAfter, ShouldBeLessThan, lenBefore)
			So(arr.Get(3), ShouldEqual, 4)
		})

		Convey("Should remove an item from the end", func() {
			arr := values.NewArrayWith(
				values.NewInt(0),
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
			)

			lenBefore := arr.Length()

			arr.RemoveAt(5)

			lenAfter := arr.Length()

			So(lenAfter, ShouldBeLessThan, lenBefore)
			So(lenAfter, ShouldEqual, 5)
			So(arr.Get(4), ShouldEqual, 4)
		})

		Convey("Should remove an item from the beginning", func() {
			arr := values.NewArrayWith(
				values.NewInt(0),
				values.NewInt(1),
				values.NewInt(2),
				values.NewInt(3),
				values.NewInt(4),
				values.NewInt(5),
			)

			lenBefore := arr.Length()

			arr.RemoveAt(0)

			lenAfter := arr.Length()

			So(lenAfter, ShouldBeLessThan, lenBefore)
			So(arr.Get(0), ShouldEqual, 1)
		})
	})
}
