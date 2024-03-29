package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"encoding/json"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// InSlice check whether a element is in the given slice
func InSlice[T comparable](s []T, e interface{}) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

// ReprBitsLen returns the minimum bit length to represent an uint64 number
func ReprBitsLen(num uint64) int {
	for i := 1; ; i++ {
		num >>= 1

		if num == 0 {
			return i
		}
	}
}

// GetStrSet returns ascii char set as a string
func GetStrSet(n int) string {
	strSet := []string{
		"0123456789",
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~",
	}

	if n < 0 {
		n = 0
	} else if n > len(strSet) {
		n = len(strSet)
	}

	return strings.Join(strSet[:n], "")
}

// GenRandStr generate random string in length n with characters from the
// given set represented by s
func GenRandStr[T string | int](n int, s T) string {

	// reseed should be performed
	// rand.Seed(time.Now().UnixNano())

	var ss string

	// get source set of characters
	var i interface{} = s
	if v, ok := i.(string); ok {
		ss = v
	} else if v, ok := i.(int); ok {
		ss = GetStrSet(v)
	}

	// check length
	if len(ss) == 0 || n == 0 {
		return ""
	}

	// rand select
	sb := strings.Builder{}
	sb.Grow(n)
	for i := 0; i < n; i++ {
		sb.WriteByte(ss[rand.Intn(len(ss))])
	}
	return sb.String()
}

// IsFileExist returns whether file exist
func IsFileExist(path string) bool {
	_, err := os.Stat(path)

	return !errors.Is(err, os.ErrNotExist)
}

// ReadFile reads up to n bytes from given path
func ReadFile(path string, n int) ([]byte, int) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	buf := make([]byte, n)
	nread, err := file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

	return buf, nread
}

// Dumps dump object to string for the given format: toml|yaml|json
func Dumps(o interface{}, format string) string {
	var b []byte
	var err error

	switch format {
	case "toml":
		b, err = toml.Marshal(o)
	case "yaml":
		b, err = yaml.Marshal(o)
	case "json":
		b, err = json.MarshalIndent(o, "", "\t")
	default:
		return ""
	}

	if err != nil {
		log.Fatal(err)
	}

	return string(b)
}

// ShowTable print table with preset simple style
func ShowTable(t [][]interface{}) {
	tt := table.NewWriter()

	tt.SetAutoIndex(false)
	tt.Style().Options.DrawBorder = false
	tt.Style().Options.SeparateColumns = false
	tt.Style().Options.SeparateFooter = false
	tt.Style().Options.SeparateHeader = false
	tt.Style().Options.SeparateRows = false

	for i := range t {
		tt.AppendRow(t[i])
	}

	fmt.Println(tt.Render())
}
