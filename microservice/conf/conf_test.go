package conf

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	os.Setenv("APPLICATION_ENV", "dev")
	os.Setenv("DATA_MONGO_URL", "mongodb://localhost")
	if err := Load("APPLICATION_"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := Load("DATA_"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := Load(""); err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log("complete")
}

func TestUnmarshal(t *testing.T) {
	type Config struct {
		Env   string
		Mongo *struct {
			Url string
		}
	}
	os.Setenv("APPLICATION_ENV", "dev")
	os.Setenv("DATA_MONGO_URL", "mongodb://localhost")
	if err := Load("APPLICATION_"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := Load("DATA_"); err != nil {
		t.Error(err)
		t.FailNow()
	}
	c, err := Unmarshal[Config]("")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Logf("conf: %+v\r\n", c)
	t.Logf("mongo conf: %+v\r\n", c.Mongo)
}
