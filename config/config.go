package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var ErrorNotFoundOrNullValue = errors.New("config: key not found or null value")

type Config struct {
	content map[string]interface{}
	node    interface{}
}

func New() *Config {
	flag.Parse()
	return &Config{content: make(map[string]interface{})}
}

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

func (c *Config) Get(key string) *Config {
	env := asEnv(key)
	if v, exists := os.LookupEnv(env); exists {
		return &Config{content: make(map[string]interface{}), node: v}
	}

	arg := asArg(key)
	if f := flag.Lookup(arg); f != nil {
		for _, element := range os.Args {
			if element == "-"+arg || element == "--"+arg {
				return &Config{content: make(map[string]interface{}), node: f.Value.String()}
			}
		}
	}

	cnf := c
	for _, element := range split(key) {
		if v, exists := cnf.content[element]; exists {
			if m, ok := v.(map[string]interface{}); ok {
				cnf = &Config{content: m, node: v}
			} else {
				cnf = &Config{content: make(map[string]interface{}), node: v}
			}
		} else {
			break
		}
	}

	if cnf.node != nil {
		return cnf
	}

	return c
}

func (c *Config) String() (string, error) {
	if c.node == nil {
		return "", ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.(string); ok {
		return v, nil
	}

	return "", errors.New(fmt.Sprintf("config: %T to string not supported", c.node))
}

func (c *Config) StringOrElse(value string) string {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(string); ok {
		return v
	}

	return value
}

func (c *Config) Bool() (bool, error) {
	if c.node == nil {
		return false, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.(bool); ok {
		return v, nil
	}

	return false, errors.New(fmt.Sprintf("config: %T to bool not supported", c.node))
}

func (c *Config) BoolOrElse(value bool) bool {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(bool); ok {
		return v
	}

	return value
}

func (c *Config) Int() (int, error) {
	if c.node == nil {
		return 0, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.(int); ok {
		return v, nil
	}

	return 0, errors.New(fmt.Sprintf("config: %T to int not supported", c.node))
}

func (c *Config) IntOrElse(value int) int {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(int); ok {
		return v
	}

	return value
}

func (c *Config) Uint() (uint, error) {
	if c.node == nil {
		return 0, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.(uint); ok {
		return v, nil
	}

	return 0, errors.New(fmt.Sprintf("config: %T to uint not supported", c.node))
}

func (c *Config) UintOrElse(value uint) uint {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(uint); ok {
		return v
	}

	return value
}

func (c *Config) Float32() (float32, error) {
	if c.node == nil {
		return 0, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.(float32); ok {
		return v, nil
	}

	return 0, errors.New(fmt.Sprintf("config: %T to float32 not supported", c.node))
}

func (c *Config) Float32OrElse(value float32) float32 {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(float32); ok {
		return v
	}

	return value
}

func (c *Config) Float64() (float64, error) {
	if c.node == nil {
		return 0, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.(float64); ok {
		return v, nil
	}

	return 0, errors.New(fmt.Sprintf("config: %T to float64 not supported", c.node))
}

func (c *Config) Float64OrElse(value float64) float64 {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.(float64); ok {
		return v
	}

	return value
}

func (c *Config) Duration() (time.Duration, error) {
	d, e := time.ParseDuration(c.StringOrElse(""))
	if e != nil {
		return 0, errors.New(fmt.Sprintf("config: %T to duration not supported", c.node))
	}

	return d, nil
}

func (c *Config) DurationOrElse(value time.Duration) time.Duration {
	d, e := time.ParseDuration(c.StringOrElse(""))
	if e != nil {
		return value
	}

	return d
}

func (c *Config) SliceOfString() ([]string, error) {
	if c.node == nil {
		return nil, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.([]string); ok {
		return v, nil
	}

	return nil, errors.New(fmt.Sprintf("config: %T to []string not supported", c.node))
}

func (c *Config) SliceOfStringOrElse(value []string) []string {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.([]string); ok {
		return v
	}

	return value
}

func (c *Config) SliceOfBool() ([]bool, error) {
	if c.node == nil {
		return nil, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.([]bool); ok {
		return v, nil
	}

	return nil, errors.New(fmt.Sprintf("config: %T to []bool not supported", c.node))
}

func (c *Config) SliceOfBoolOrElse(value []bool) []bool {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.([]bool); ok {
		return v
	}

	return value
}

func (c *Config) SliceOfInt() ([]int, error) {
	if c.node == nil {
		return nil, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.([]int); ok {
		return v, nil
	}

	return nil, errors.New(fmt.Sprintf("config: %T to []int not supported", c.node))
}

func (c *Config) SliceOfIntOrElse(value []int) []int {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.([]int); ok {
		return v
	}

	return value
}

func (c *Config) SliceOfUint() ([]uint, error) {
	if c.node == nil {
		return nil, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.([]uint); ok {
		return v, nil
	}

	return nil, errors.New(fmt.Sprintf("config: %T to []uint not supported", c.node))
}

func (c *Config) SliceOfUintOrElse(value []uint) []uint {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.([]uint); ok {
		return v
	}

	return value
}

func (c *Config) SliceOfFloat32() ([]float32, error) {
	if c.node == nil {
		return nil, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.([]float32); ok {
		return v, nil
	}

	return nil, errors.New(fmt.Sprintf("config: %T to []float32 not supported", c.node))
}

func (c *Config) SliceOfFloat32OrElse(value []float32) []float32 {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.([]float32); ok {
		return v
	}

	return value
}

func (c *Config) SliceOfFloat64() ([]float64, error) {
	if c.node == nil {
		return nil, ErrorNotFoundOrNullValue
	}

	if v, ok := c.node.([]float64); ok {
		return v, nil
	}

	return nil, errors.New(fmt.Sprintf("config: %T to []float64 not supported", c.node))
}

func (c *Config) SliceOfFloat64OrElse(value []float64) []float64 {
	if c.node == nil {
		return value
	}

	if v, ok := c.node.([]float64); ok {
		return v
	}

	return value
}

func asEnv(key string) string {
	return strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
}

func asArg(key string) string {
	return strings.ReplaceAll(key, ".", "_")
}

func split(key string) []string {
	return strings.Split(key, ".")
}
