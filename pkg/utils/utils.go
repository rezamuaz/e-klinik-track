package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"runtime"
	"strconv"

	"github.com/jackc/pgx/v5/pgtype"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ErrorWrapper(format string, args ...interface{}) error {
	err := fmt.Errorf(format, args...)
	_, file, line, _ := runtime.Caller(1)
	log.Printf("ERROR: [%s:%d] %v \n", file, line, err)
	return err
}

func StrToObjectID(id string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	return objectId
}

func CountTrues(boolSlice []bool) int {
	count := 0
	for _, val := range boolSlice {
		if val {
			count++
		}
	}
	return count
}

func CreateBoolList(size int, sublistSize int) [][]bool {
	list := make([][]bool, size)
	for i := range list {
		list[i] = make([]bool, sublistSize)
	}
	return list
}

// ANCHOR - Pgx Converter
func FloatToPgNum(value float64) pgtype.Numeric {
	var x pgtype.Numeric
	x.Scan(strconv.FormatFloat(value, 'f', 2, 64))
	return x
}

func PgNumToFloat(src pgtype.Numeric) float64 {
	f, err := strconv.ParseFloat(src.Int.String(), 64)
	if err != nil {
		return 0
	}
	var roundTo float64 = 1
	if src.Exp > 0 {
		for i := 0; i < int(src.Exp); i++ {
			f *= 10
			roundTo *= 10
		}
	} else if src.Exp < 0 {
		for i := 0; i > int(src.Exp); i-- {
			f /= 10
			roundTo /= 10
		}
	}
	return math.Round(f/roundTo) * roundTo
}
func StrtoFloat32(str string) float32 {
	float64Value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		fmt.Println("Error converting string to float:", err)
	}

	// Cast the float64 to float32
	return float32(float64Value)
}

func StructToBytes[T any](v T) []byte {
	data, err := json.Marshal(v)
	if err != nil {
		log.Println("error marshaling struct to bytes:", err)

	}
	return data
}

func ByteToJson[T any](data []byte) T {
	var result T
	err := json.Unmarshal(data, &result)
	if err != nil {
		log.Println("error marshaling", err)
	}
	return result
}

func StrToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.Println("Invalid id")
	}

	return num
}

func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

func StrIsEmpty(value string) bool {
	if value == "" {
		return false
	} else {
		return true
	}
}

func Int64IsEmpty(value int64) bool {
	if value == 0 {
		return false
	} else {
		return true
	}
}

func StrToInt32(str string) int32 {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Println("Error converting string to int:", err)
	}
	return int32(i)
}

func StrToInt(str string) int {
	i, err := strconv.Atoi(str)
	if err != nil {
		log.Println("Error converting string to int:", err)
	}
	return i
}

func GetOffset(page int32, limit int32) int32 {
	if page < 1 {
		page = 1 // Ensure page is at least 1
	}
	if limit < 1 {
		return 0 // Avoid negative or zero limit causing unexpected behavior
	}
	return int32((page - 1) * limit)
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
