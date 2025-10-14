package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"e-klinik/infra/pg"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-viper/mapstructure/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func DerefString(ptr *string) string {
	if ptr != nil {
		return *ptr
	}
	return ""
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

func Int32ToStr(v int32) string {
	return strconv.FormatInt(int64(v), 10)
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
	if page <= 1 {
		return 0
	}
	return (page - 1) * limit

}

func ByteToAny[T any](b []byte) (T, error) {
	var res T
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		var zero T
		return zero, fmt.Errorf("error gob encoder")
	}
	return res, nil
}

func TextNull(s string) pgtype.Text {
	return pgtype.Text{String: s, Valid: s != ""}
}

func DumpTest(v any) {
	// Print type and kind
	t := reflect.TypeOf(v)
	fmt.Printf("Type: %T\n", v)
	fmt.Printf("Kind: %s\n", t.Kind())

	// Pretty-print as JSON (if possible)
	if data, err := json.MarshalIndent(v, "", "  "); err == nil {
		fmt.Println("JSON:")
		fmt.Println(string(data))
	} else {
		// Fallback: print raw value
		fmt.Println("Value:")
		fmt.Printf("%+v\n", v)
	}
}

func WithTransactionResult[T any](
	ctx context.Context,
	pool *pgxpool.Pool,
	fn func(q *pg.Queries, tx pgx.Tx) (T, error),
) (T, error) {
	var zero T // default zero value for T

	tx, err := pool.Begin(ctx)
	if err != nil {
		return zero, err
	}
	defer tx.Rollback(ctx)

	qtx := pg.New(tx)
	result, err := fn(qtx, tx)
	if err != nil {
		return zero, err
	}

	if err := tx.Commit(ctx); err != nil {
		return zero, err
	}

	return result, nil
}

func BoolPtr(b bool) *bool {
	return &b
}

func IntPtr(i int) *int {
	return &i
}

func StringPtr(s string) *string {

	return &s
}

func ToPtr[T any](v T) *T {
	return &v
}
func ParseBool(s string) bool {
	val, err := strconv.ParseBool(s)
	if err != nil {
		return false // or default to true if you prefer
	}
	return val
}

func JSONPtr(v any) *string {
	b, err := json.Marshal(v)
	if err != nil {
		panic("JSONPtr failed: " + err.Error())
	}
	s := string(b)
	return &s
}

func GenerateSecureKey(keyLengthBytes int) (string, error) {
	if keyLengthBytes <= 0 {
		return "", fmt.Errorf("keyLengthBytes must be a positive integer")
	}

	keyBytes := make([]byte, keyLengthBytes)
	_, err := rand.Read(keyBytes) // rand.Read reads random bytes into the slice
	if err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	// URLEncoding is recommended for keys that might be used in URLs or environment variables,
	// as it uses '-' and '_' instead of '+' and '/' and omits padding '=' characters.
	keyString := base64.URLEncoding.EncodeToString(keyBytes)

	return keyString, nil
}

func BindQueryTo[T any](c *gin.Context) (T, error) {
	var out T
	query := flattenQuery(c.Request.URL.Query())

	err := mapstructure.Decode(query, &out)
	return out, err
}

func flattenQuery(values url.Values) map[string]string {
	flat := make(map[string]string)
	for k, v := range values {
		if len(v) > 0 {
			flat[k] = v[0]
		}
	}
	return flat
}

func MustUUID(s string) pgtype.UUID {
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		panic(err) // or handle error properly
	}
	return u
}

func ParseTimestamptz(str string, layout ...string) (pgtype.Timestamptz, error) {
	tzLayout := "2006-01-02T15:04"
	if len(layout) > 0 && layout[0] != "" {
		tzLayout = layout[0]
	}

	t, err := time.Parse(tzLayout, str)
	if err != nil {
		return pgtype.Timestamptz{}, err
	}

	return pgtype.Timestamptz{
		Time:  t.UTC(),
		Valid: true,
	}, nil
}

func StringToTimestamptz(s *string) pgtype.Timestamptz {
	if s == nil || *s == "" {
		return pgtype.Timestamptz{Valid: false}
	}
	// try parse ISO8601 format (adjust if your frontend sends another format)
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		fmt.Println("⚠️ failed to parse time:", err)
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: t, Valid: true}
}

func GetJakartaDateObject() (time.Time, error) {
	// 1. Muat zona waktu Jakarta
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Timezone 'Asia/Jakarta' seharusnya selalu tersedia
		return time.Time{}, fmt.Errorf("failed to load timezone: %w", err)
	}

	// 2. Konversi waktu server saat ini ke zona waktu Jakarta
	nowJakarta := time.Now().In(loc)

	// 3. Ambil waktu awal hari (midnight) dari tanggal Jakarta tersebut.
	// Ini memastikan kita hanya menyimpan komponen DATE, tanpa komponen waktu,
	// yang paling cocok untuk kolom DATE.
	tanggalJakarta := time.Date(
		nowJakarta.Year(),
		nowJakarta.Month(),
		nowJakarta.Day(),
		0, 0, 0, 0, // Set waktu ke 00:00:00.000
		nowJakarta.Location(), // Pertahankan lokasi (Jakarta)
	)

	return tanggalJakarta, nil
}
