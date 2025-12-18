package helper

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Concat(splitcode string, args ...string) string {

	var b bytes.Buffer

	for _, arg := range args {
		b.WriteString(arg + splitcode)
	}

	f := strings.TrimRight(b.String(), splitcode)
	return f
}

func GetFormatTime(loc *time.Location, layout string) string {

	// Standard GO Constant Format :

	// ANSIC       = "Mon Jan _2 15:04:05 2006"
	// UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	// RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	// RFC822      = "02 Jan 06 15:04 MST"
	// RFC822Z     = "02 Jan 06 15:04 -0700"
	// RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	// RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	// RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"
	// RFC3339     = "2006-01-02T15:04:05Z07:00"
	// RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	// Kitchen     = "3:04PM"
	// // Handy time stamps.
	// Stamp      = "Jan _2 15:04:05"
	// StampMilli = "Jan _2 15:04:05.000"
	// StampMicro = "Jan _2 15:04:05.000000"
	// StampNano  = "Jan _2 15:04:05.000000000"

	// Using Manual Format :
	// 1. date yyyy-mm-dd = 2006-01-02
	// 2. time hhhh:ii:ss = 15:04:05

	t := time.Now()
	f := t.In(loc).Format(layout)

	return f
}

func GetCurrentTime(loc *time.Location, val string) time.Time {

	now := time.Now()
	t, _ := time.Parse(val, now.In(loc).Format(time.RFC3339))
	return t
}

func GetYesterday(loc *time.Location, day time.Duration, layout string) time.Time {

	now := time.Now()
	t, _ := time.Parse(layout, now.In(loc).Format(layout))
	return t.In(loc).Add((-24 * day) * time.Hour)
}

func GetTomorrow(loc *time.Location, day time.Duration, layout string) time.Time {

	now := time.Now()
	t, _ := time.Parse(layout, now.In(loc).Format(layout))
	return t.In(loc).Add((24 * day) * time.Hour)
}

func CounterZeroNumber(length int) string {

	var wordNumbers string

	for w := 0; w < length; w++ {
		wordNumbers += "0"
	}

	return wordNumbers
}

func ReduceWords(words string, start int, length int) string {

	runes := []rune(words)
	inputFmt := string(runes[start:length])

	return inputFmt
}

func GetUniqId(loc *time.Location) string {

	t := time.Now()
	var formatId = t.In(loc).Format("20060102150405.000000")
	uniqId := strings.Replace(formatId, ".", "", -1)

	return uniqId
}

func GetTrxID(loc *time.Location) string {

	t := time.Now()
	var formatId = t.In(loc).Format("20060102150405")

	return ReduceWords(strings.Replace(formatId, ".", "", -1)+strconv.Itoa(int(t.In(loc).UnixNano())), 0, 32)
}

func GetIpAddress(c *fiber.Ctx) string {

	ip := ""

	for k, v := range c.GetReqHeaders() {
		//fmt.Printf("(k1) %#v", k)

		if k == "Cf-Connecting-Ip" {
			for _, v2 := range v {
				//fmt.Printf("(2) %#v", v2)
				if v2 != "" {
					ip = v2
					break
				}
			}
			break
		}
	}

	if ip == "" {
		for k, v := range c.GetReqHeaders() {
			//fmt.Printf("(k1) %#v", k)

			if k == "X-Forwarded-For" {
				for _, v2 := range v {
					//fmt.Printf("(2) %#v", v2)
					if v2 != "" {
						ip = v2
						break
					}
				}
				break
			}
		}
	}

	if ip == "" {
		for k, v := range c.GetReqHeaders() {
			//fmt.Printf("(k1) %#v", k)

			if k == "X-Real-Ip" {
				for _, v2 := range v {
					//fmt.Printf("(2) %#v", v2)
					if v2 != "" {
						ip = v2
						break
					}
				}
				break
			}
		}
	}

	if ip == "" {
		ip = c.IP()
	}

	return ip
}

var (
	// We're using a 32 byte long secret key.
	// This is probably something you generate first
	// then put into and environment variable.
	secretKey string = "N1PCdw3M2B1TfJho" // 16 byte secret
)

func Encrypt(plaintext string) string {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		panic(err)
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return string(ciphertext)
}

func Decrypt(ciphertext string) string {
	aes, err := aes.NewCipher([]byte(secretKey))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		panic(err)
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext)
}

func CompressGzip(f string) string {

	// Open the original file
	originalFile, err := os.Open(f)
	if err != nil {
		panic(err)
	}
	defer originalFile.Close()

	gzf := f + ".gz"
	// Create a new gzipped file
	gzippedFile, err := os.Create(gzf)
	if err != nil {
		panic(err)
	}
	defer gzippedFile.Close()

	// Create a new gzip writer
	gzipWriter := gzip.NewWriter(gzippedFile)
	defer gzipWriter.Close()

	// Copy the contents of the original file to the gzip writer
	_, err = io.Copy(gzipWriter, originalFile)
	if err != nil {
		panic(err)
	}

	// Flush the gzip writer to ensure all data is written
	gzipWriter.Flush()

	return gzf
}

func GetDateFromInt(loc *time.Location, dateInt int64) time.Time {
	t := time.Unix(dateInt, 0)
	return t.In(loc)
}

func ParseVendorDateSend(loc *time.Location, raw string) time.Time {

	fallback := time.Date(1970, 1, 1, 12, 0, 0, 0, loc)

	if len(raw) < 8 {
		return fallback
	}

	datePart := raw[:8]

	t, err := time.ParseInLocation("20060102", datePart, loc)
	if err != nil {
		return fallback
	}

	return time.Date(
		t.Year(),
		t.Month(),
		t.Day(),
		12, 0, 0, 0,
		loc,
	)
}