package config

import (
	"flag"
	"io/ioutil"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	config := New()

	if config.content == nil {
		t.Errorf("content == nil")
	}

	if config.node != nil {
		t.Errorf("node != nil")
	}
}

func TestConfig_LoadJSON(t *testing.T) {
	config := New()

	config, e := config.LoadJSON("config.json")
	if e == nil {
		t.Errorf("error == nil")
	}
}

func TestConfig_LoadJSON2(t *testing.T) {
	config := New()

	file, e := ioutil.TempFile("", "config_*.json")
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}
	defer file.Close()

	_, e = file.WriteString(`{ db : { "user": "root", "password": "root" } }`)
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	config, e = config.LoadJSON(file.Name())
	if e == nil {
		t.Errorf("error == nil")
	}
}

func TestConfig_LoadJSON3(t *testing.T) {
	config := New()

	file, e := ioutil.TempFile("", "config_*.json")
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}
	defer file.Close()

	_, e = file.WriteString(`{ "db" : { "user": "root", "password": "root" } }`)
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	config, e = config.LoadJSON(file.Name())
	if e != nil {
		t.Errorf(e.Error())
	}
}

func TestConfig_Get(t *testing.T) {
	config := New()

	if e := os.Setenv("TEST_DB_USER", "root"); e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	user, e := config.Get("db.user").String()
	if e == nil {
		t.Errorf("error == nil")
		t.FailNow()
	}

	if user != "" {
		t.Errorf("user == %v", user)
		t.FailNow()
	}

	if e != ErrNotFoundOrNullValue {
		t.Errorf("error != ErrNotFoundOrNullValue")
	}
}

func TestConfig_Get2(t *testing.T) {
	config := New()

	if e := os.Setenv("TEST_DB_USER", "root"); e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	user, e := config.Get("test.db.user").String()
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	if user != "root" {
		t.Errorf("user != root")
	}
}

func TestConfig_Get3(t *testing.T) {
	os.Args = append(os.Args, "--test_db_user", "root")
	flag.String("test_db_user", "default", "")

	config := New()
	user, e := config.Get("db.user").String()
	if e == nil {
		t.Errorf("error == nil")
		t.FailNow()
	}

	if user != "" {
		t.Errorf("user == %v", user)
		t.FailNow()
	}

	if e != ErrNotFoundOrNullValue {
		t.Errorf("error != ErrNotFoundOrNullValue")
	}
}

func TestConfig_Get4(t *testing.T) {
	os.Args = append(os.Args, "--test_db_password", "password")
	flag.String("test_db_password", "default", "")

	config := New()
	password, e := config.Get("test.db.password").String()
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	if password != "password" {
		t.Errorf("password != password")
	}
}

func TestConfig_Get5(t *testing.T) {
	os.Args = append(os.Args, "-test_db_ssl", "enable")
	flag.String("test_db_ssl", "default", "")

	config := New()
	enable, e := config.Get("db.ssl").String()
	if e == nil {
		t.Errorf("error == nil")
		t.FailNow()
	}

	if enable != "" {
		t.Errorf("enable == %v", enable)
		t.FailNow()
	}

	if e != ErrNotFoundOrNullValue {
		t.Errorf("error != ErrNotFoundOrNullValue")
	}
}

func TestConfig_Get6(t *testing.T) {
	os.Args = append(os.Args, "-test_db_connection", "localhost")
	flag.String("test_db_connection", "default", "")

	config := New()
	connection, e := config.Get("test.db.connection").String()
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	if connection != "localhost" {
		t.Errorf("connection != localhost")
	}
}

func TestConfig_Get7(t *testing.T) {
	flag.String("test_db_migration", "default", "")

	config := New()
	migration, e := config.Get("test.db.migration").String()
	if e == nil {
		t.Errorf("error == nil")
		t.FailNow()
	}

	if migration != "" {
		t.Errorf("migration == %v", migration)
		t.FailNow()
	}

	if e != ErrNotFoundOrNullValue {
		t.Errorf("error != ErrNotFoundOrNullValue")
	}
}

func TestConfig_Get8(t *testing.T) {
	os.Args = append(os.Args, "--test_db_schema", "public")
	flag.String("test_db_schema", "default", "")

	config := New()
	schema, e := config.Get("test.db.schema").Bool()
	if e == nil {
		t.Errorf("error == nil")
		t.FailNow()
	}

	if schema != false {
		t.Errorf("schema != false")
	}

	if e.Error() != "config: string to bool not supported" {
		t.Errorf("error != config: string to bool not supported")
	}
}

func TestConfig_Get9(t *testing.T) {
	config := New()

	file, e := ioutil.TempFile("", "config_*.json")
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}
	defer file.Close()

	_, e = file.WriteString(`{}`)
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	config, e = config.LoadJSON(file.Name())
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	level, e := config.Get("logger.level").String()
	if e == nil {
		t.Errorf("e == nil")
		t.FailNow()
	}

	if level != "" {
		t.Errorf("user == %v", level)
	}
}

func TestConfig_Get10(t *testing.T) {
	config := New()

	file, e := ioutil.TempFile("", "config_*.json")
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}
	defer file.Close()

	_, e = file.WriteString(`{"provider_url" : "localhost:8080"}`)
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	config, e = config.LoadJSON(file.Name())
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	url, e := config.Get("provider_url").String()
	if e != nil {
		t.Errorf(e.Error())
		t.FailNow()
	}

	if url != "localhost:8080" {
		t.Errorf("user != localhost:8080")
	}
}

func TestAsEnv(t *testing.T) {
	sample1 := asEnv("a")
	if sample1 != "A" {
		t.Errorf("A == %v", sample1)
	}

	sample2 := asEnv("ab")
	if sample2 != "AB" {
		t.Errorf("AB == %v", sample2)
	}

	sample3 := asEnv("a.b")
	if sample3 != "A_B" {
		t.Errorf("A_B == %v", sample3)
	}

	sample4 := asEnv("a_b")
	if sample4 != "A_B" {
		t.Errorf("A_B == %v", sample4)
	}

	sample5 := asEnv("a.b_c")
	if sample5 != "A_B_C" {
		t.Errorf("A_B_C == %v", sample5)
	}

	sample6 := asEnv("a_b.c")
	if sample6 != "A_B_C" {
		t.Errorf("A_B_C == %v", sample6)
	}
}

func TestAsArg(t *testing.T) {
	sample1 := asArg("a.b")
	if sample1 != "a_b" {
		t.Errorf("a_b == %v", sample1)
	}

	sample2 := asArg("a.b_c")
	if sample2 != "a_b_c" {
		t.Errorf("a_b_c == %v", sample2)
	}

	sample3 := asArg("a_b.c")
	if sample3 != "a_b_c" {
		t.Errorf("a_b_c == %v", sample3)
	}
}

func TestSplit(t *testing.T) {
	sample1 := "a"
	ks := split(sample1)
	if len(ks) != 1 {
		t.Errorf("len != 1")
	}
	if ks[0] != "a" {
		t.Errorf("ks[0] != a")
	}

	sample2 := "a.b"
	ks = split(sample2)
	if len(ks) != 2 {
		t.Errorf("len != 2")
	}
	if ks[0] != "a" {
		t.Errorf("ks[0] != a")
	}
	if ks[1] != "b" {
		t.Errorf("ks[1] != b")
	}
}
