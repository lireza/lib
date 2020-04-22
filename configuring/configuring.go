// Package configuring provides configuration loading mechanism from different configuration sources. Including
// environment variables, command line arguments and JSON configuration file.
//
// The configuration should be seen as a tree like structure. For example, keys logger.level, logger.enable
// should be seen as a logger node containing two nested nodes level and enable.
// Each node itself, is a value, so the logger node is an object value (Think JSON object), because it contains
// two keys nested in. The value of level can be a string and the value of enable can be a boolean value.
//
// The configuring instance is used to load configuration from different sources mentioned. Based on our example
// the configuring instance does the steps bellow:
// 1) If the asEnv(key) is defined as environment variable, returns the value.
// 2) If the asArg(key) is defined as command line argument, returns the value.
// 3) If the configuring instance is used to load a JSON configuration file, tries to load a node from JSON.
//
// Accessor methods can be used to convert loaded value or node to different types.
package configuring

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

// ErrNotFoundOrNullValue determines a provided key not found, or the value is null.
var ErrNotFoundOrNullValue = errors.New("configuring: key not found or null value")

// Config encapsulates the configuration loading mechanism.
type Config struct {
	content map[string]interface{}
	node    interface{}
}

// New creates a new configuring instance ready to load configuration values from.
// The created instance can be used to load environment variables and/or command line arguments
// defined based on standard flag package.
func New() *Config {
	flag.Parse()
	return &Config{content: make(map[string]interface{})}
}

// LoadJSON loads JSON configuration file to the current configuring instance and returns configuring itself.
// The returned instance can be used to load environment variables, command line arguments
// defined based on standard flag package and loaded JSON configuration file.
func (c *Config) LoadJSON(filename string) (*Config, error) {
	file, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}

	if e := json.Unmarshal(file, &c.content); e != nil {
		return nil, e
	}

	return c, nil
}

// Get returns back a config instance that may be filled with an appropriate node instance.
// The accessor methods can be used to convert the node to a specific type.
func (c *Config) Get(key string) *Config {
	if v, exists := os.LookupEnv(asEnv(key)); exists {
		return &Config{content: c.content, node: v}
	}

	if f := flag.Lookup(asArg(key)); f != nil {
		for _, element := range os.Args {
			if element == "-"+asArg(key) || element == "--"+asArg(key) {
				return &Config{content: c.content, node: f.Value.String()}
			}
		}
	}

	temp := c
	for _, element := range split(key) {
		if v, exists := temp.content[element]; exists {
			if m, ok := v.(map[string]interface{}); ok {
				temp = &Config{content: m, node: v}
			} else {
				temp = &Config{content: make(map[string]interface{}), node: v}
			}
		} else {
			return c
		}
	}

	return temp
}

// String returns the string representation of a node if convertible.
func (c *Config) String() (string, error) {
	if c.node == nil {
		return "", ErrNotFoundOrNullValue
	}

	if v, ok := c.node.(string); ok {
		return v, nil
	}

	return "", errors.New(fmt.Sprintf("configuring: %T to string not supported", c.node))
}

// StringOrElse returns the string representation of a node if convertible otherwise the default value provided.
func (c *Config) StringOrElse(value string) string {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(string); ok {
		return v
	}

	return value
}

// Bool returns the boolean representation of a node if convertible.
func (c *Config) Bool() (bool, error) {
	if c.node == nil {
		return false, ErrNotFoundOrNullValue
	}

	if v, ok := c.node.(bool); ok {
		return v, nil
	}

	if v, e := strconv.ParseBool(c.StringOrElse("")); e == nil {
		return v, nil
	}

	return false, errors.New(fmt.Sprintf("configuring: %T to bool not supported", c.node))
}

// BoolOrElse returns the boolean representation of a node if convertible otherwise the default value provided.
func (c *Config) BoolOrElse(value bool) bool {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(bool); ok {
		return v
	}

	if v, e := strconv.ParseBool(c.StringOrElse("")); e == nil {
		return v
	}

	return value
}

// Int returns the integer representation of a node if convertible.
func (c *Config) Int() (int, error) {
	if c.node == nil {
		return 0, ErrNotFoundOrNullValue
	}

	if v, ok := c.node.(int); ok {
		return v, nil
	}

	if v, e := strconv.Atoi(c.StringOrElse("")); e == nil {
		return v, nil
	}

	return 0, errors.New(fmt.Sprintf("configuring: %T to int not supported", c.node))
}

// IntOrElse returns the integer representation of a node if convertible otherwise the default value provided.
func (c *Config) IntOrElse(value int) int {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(int); ok {
		return v
	}

	if v, e := strconv.Atoi(c.StringOrElse("")); e == nil {
		return v
	}

	return value
}

// Uint returns the unsigned integer representation of a node if convertible.
func (c *Config) Uint() (uint, error) {
	if c.node == nil {
		return 0, ErrNotFoundOrNullValue
	}

	if v, ok := c.node.(uint); ok {
		return v, nil
	}

	if v, e := strconv.ParseInt(c.StringOrElse(""), 10, 0); e == nil {
		return uint(v), nil
	}

	return 0, errors.New(fmt.Sprintf("configuring: %T to uint not supported", c.node))
}

// UintOrElse returns the unsigned integer representation of a node if convertible otherwise the default value provided.
func (c *Config) UintOrElse(value uint) uint {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(uint); ok {
		return v
	}

	if v, e := strconv.ParseInt(c.StringOrElse(""), 10, 0); e == nil {
		return uint(v)
	}

	return value
}

// Float32 returns the floating point representation of a node if convertible.
func (c *Config) Float32() (float32, error) {
	if c.node == nil {
		return 0, ErrNotFoundOrNullValue
	}

	if v, ok := c.node.(float32); ok {
		return v, nil
	}

	if v, e := strconv.ParseFloat(c.StringOrElse(""), 32); e == nil {
		return float32(v), nil
	}

	return 0, errors.New(fmt.Sprintf("configuring: %T to float32 not supported", c.node))
}

// Float32OrElse returns the floating point representation of a node if convertible otherwise the default value provided.
func (c *Config) Float32OrElse(value float32) float32 {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(float32); ok {
		return v
	}

	if v, e := strconv.ParseFloat(c.StringOrElse(""), 32); e == nil {
		return float32(v)
	}

	return value
}

// Float64 returns the floating point representation of a node if convertible.
func (c *Config) Float64() (float64, error) {
	if c.node == nil {
		return 0, ErrNotFoundOrNullValue
	}

	if v, ok := c.node.(float64); ok {
		return v, nil
	}

	if v, e := strconv.ParseFloat(c.StringOrElse(""), 64); e == nil {
		return v, nil
	}

	return 0, errors.New(fmt.Sprintf("configuring: %T to float64 not supported", c.node))
}

// Float64OrElse returns the floating point representation of a node if convertible otherwise the default value provided.
func (c *Config) Float64OrElse(value float64) float64 {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(float64); ok {
		return v
	}

	if v, e := strconv.ParseFloat(c.StringOrElse(""), 64); e == nil {
		return v
	}

	return value
}

// Duration returns the duration representation of a node if convertible.
func (c *Config) Duration() (time.Duration, error) {
	d, e := time.ParseDuration(c.StringOrElse(""))
	if e != nil {
		return 0, errors.New(fmt.Sprintf("configuring: %T to duration not supported", c.node))
	}

	return d, nil
}

// DurationOrElse returns the duration representation of a node if convertible otherwise the default value provided.
func (c *Config) DurationOrElse(value time.Duration) time.Duration {
	d, e := time.ParseDuration(c.StringOrElse(""))
	if e != nil {
		return value
	}

	return d
}

// SliceOfString returns the slice of string representation of a node if convertible.
func (c *Config) SliceOfString() ([]string, error) {
	if c.node == nil {
		return nil, ErrNotFoundOrNullValue
	}

	if vs, ok := c.node.([]interface{}); ok {
		ss := make([]string, 0)
		for _, v := range vs {
			if s, ok := v.(string); ok {
				ss = append(ss, s)
			} else {
				return nil, errors.New(fmt.Sprintf("configuring: %T to string not supported", v))
			}
		}

		return ss, nil
	}

	return nil, errors.New(fmt.Sprintf("configuring: %T to []string not supported", c.node))
}

// SliceOfStringOrElse returns the slice of string representation of a node if convertible, otherwise the default value provided.
func (c *Config) SliceOfStringOrElse(value []string) []string {
	if c.node == nil {
		return value
	}

	ss := make([]string, 0)
	if vs, ok := c.node.([]interface{}); ok {
		for _, v := range vs {
			if s, ok := v.(string); ok {
				ss = append(ss, s)
			} else {
				return value
			}
		}
	} else {
		return value
	}

	return ss
}

// SliceOfBool returns the slice of boolean representation of a node if convertible.
func (c *Config) SliceOfBool() ([]bool, error) {
	if c.node == nil {
		return nil, ErrNotFoundOrNullValue
	}

	if vs, ok := c.node.([]interface{}); ok {
		bs := make([]bool, 0)
		for _, v := range vs {
			if b, ok := v.(bool); ok {
				bs = append(bs, b)
			} else {
				return nil, errors.New(fmt.Sprintf("configuring: %T to bool not supported", v))
			}
		}

		return bs, nil
	}

	return nil, errors.New(fmt.Sprintf("configuring: %T to []bool not supported", c.node))
}

// SliceOfBoolOrElse returns the slice of boolean representation of a node if convertible, otherwise the default value provided.
func (c *Config) SliceOfBoolOrElse(value []bool) []bool {
	if c.node == nil {
		return value
	}

	bs := make([]bool, 0)
	if vs, ok := c.node.([]interface{}); ok {
		for _, v := range vs {
			if b, ok := v.(bool); ok {
				bs = append(bs, b)
			} else {
				return value
			}
		}
	} else {
		return value
	}

	return bs
}

// SliceOfInt returns the slice of integer representation of a node if convertible.
func (c *Config) SliceOfInt() ([]int, error) {
	if c.node == nil {
		return nil, ErrNotFoundOrNullValue
	}

	if vs, ok := c.node.([]interface{}); ok {
		is := make([]int, 0)
		for _, v := range vs {
			if i, ok := v.(int); ok {
				is = append(is, i)
			} else {
				return nil, errors.New(fmt.Sprintf("configuring: %T to int not supported", v))
			}
		}

		return is, nil
	}

	return nil, errors.New(fmt.Sprintf("configuring: %T to []int not supported", c.node))
}

// SliceOfIntOrElse returns the slice of integer representation of a node if convertible, otherwise the default value provided.
func (c *Config) SliceOfIntOrElse(value []int) []int {
	if c.node == nil {
		return value
	}

	is := make([]int, 0)
	if vs, ok := c.node.([]interface{}); ok {
		for _, v := range vs {
			if i, ok := v.(int); ok {
				is = append(is, i)
			} else {
				return value
			}
		}
	} else {
		return value
	}

	return is
}

// SliceOfUint returns the slice of unsigned integer representation of a node if convertible.
func (c *Config) SliceOfUint() ([]uint, error) {
	if c.node == nil {
		return nil, ErrNotFoundOrNullValue
	}

	if vs, ok := c.node.([]interface{}); ok {
		uis := make([]uint, 0)
		for _, v := range vs {
			if ui, ok := v.(uint); ok {
				uis = append(uis, ui)
			} else {
				return nil, errors.New(fmt.Sprintf("configuring: %T to uint not supported", v))
			}
		}

		return uis, nil
	}

	return nil, errors.New(fmt.Sprintf("configuring: %T to []uint not supported", c.node))
}

// SliceOfUintOrElse returns the slice of unsigned integer representation of a node if convertible, otherwise the default value provided.
func (c *Config) SliceOfUintOrElse(value []uint) []uint {
	if c.node == nil {
		return value
	}

	uis := make([]uint, 0)
	if vs, ok := c.node.([]interface{}); ok {
		for _, v := range vs {
			if ui, ok := v.(uint); ok {
				uis = append(uis, ui)
			} else {
				return value
			}
		}
	} else {
		return value
	}

	return uis
}

// SliceOfFloat32 returns the slice of floating point representation of a node if convertible.
func (c *Config) SliceOfFloat32() ([]float32, error) {
	if c.node == nil {
		return nil, ErrNotFoundOrNullValue
	}

	if vs, ok := c.node.([]interface{}); ok {
		fs := make([]float32, 0)
		for _, v := range vs {
			if f, ok := v.(float32); ok {
				fs = append(fs, f)
			} else {
				return nil, errors.New(fmt.Sprintf("configuring: %T to float32 not supported", v))
			}
		}

		return fs, nil
	}

	return nil, errors.New(fmt.Sprintf("configuring: %T to []float32 not supported", c.node))
}

// SliceOfFloat32OrElse returns the slice of floating point representation of a node if convertible, otherwise the default value provided.
func (c *Config) SliceOfFloat32OrElse(value []float32) []float32 {
	if c.node == nil {
		return value
	}

	fs := make([]float32, 0)
	if vs, ok := c.node.([]interface{}); ok {
		for _, v := range vs {
			if f, ok := v.(float32); ok {
				fs = append(fs, f)
			} else {
				return value
			}
		}
	} else {
		return value
	}

	return fs
}

// SliceOfFloat64 returns the slice of floating point representation of a node if convertible.
func (c *Config) SliceOfFloat64() ([]float64, error) {
	if c.node == nil {
		return nil, ErrNotFoundOrNullValue
	}

	if vs, ok := c.node.([]interface{}); ok {
		fs := make([]float64, 0)
		for _, v := range vs {
			if f, ok := v.(float64); ok {
				fs = append(fs, f)
			} else {
				return nil, errors.New(fmt.Sprintf("configuring: %T to float64 not supported", v))
			}
		}

		return fs, nil
	}

	return nil, errors.New(fmt.Sprintf("configuring: %T to []float64 not supported", c.node))
}

// SliceOfFloat64OrElse returns the slice of floating point representation of a node if convertible, otherwise the default value provided.
func (c *Config) SliceOfFloat64OrElse(value []float64) []float64 {
	if c.node == nil {
		return value
	}

	fs := make([]float64, 0)
	if vs, ok := c.node.([]interface{}); ok {
		for _, v := range vs {
			if f, ok := v.(float64); ok {
				fs = append(fs, f)
			} else {
				return value
			}
		}
	} else {
		return value
	}

	return fs
}

// asEnv converts a key to an appropriate environment variable format.
// For example it converts a to A, a.b to A_B, a_b to A_B, a.b_c to A_B_C and a_b.c to A_B_C.
func asEnv(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}

// asArg converts a key to an appropriate command line argument format.
// For example it converts a.b to a_b, a.b_c to a_b_c and a_b.c to a_b_c.
func asArg(key string) string {
	return strings.ReplaceAll(key, ".", "_")
}

// split splits a key to its separate parts.
// For example a to [a] and a.b to [a, b].
func split(key string) []string {
	return strings.Split(key, ".")
}
