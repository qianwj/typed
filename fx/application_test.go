package fx

import (
	"github.com/qianwj/typed/fx/options"
	"testing"
)

func TestNewApp(t *testing.T) {
	app := NewApp(options.Address(""))
	app.Run()
	t.Log("complete")
}
