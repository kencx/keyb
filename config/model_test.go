package config

import (
	"os"
	"path"
	"reflect"
	"testing"
)

var (
	testKeybFile = "testkeyb.yml"
	tempKeybFile = "temp.yml"
)

func TestParseApps(t *testing.T) {

	t.Run("file present", func(t *testing.T) {
		want := Apps{{
			Name: "test",
			Keybinds: []KeyBind{{
				Name: "foo",
				Key:  "bar",
			}},
		}}

		got, err := ParseApps(path.Join(parentDir, testKeybFile))
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("file absent", func(t *testing.T) {
		testpath := path.Join(parentDir, tempKeybFile)

		_, err := ParseApps(testpath)
		if err == nil {
			t.Errorf("expected err: %s does not exist", testpath)
		}

		t.Cleanup(func() {
			os.RemoveAll(path.Join(parentDir, tempKeybFile))
		})

		_, err = os.Stat(testpath)
		if err == nil {
			t.Errorf("expected err: stat %s: no such file", testpath)
		}
	})
}

func TestAddOrUpdate(t *testing.T) {

	t.Run("update existing", func(t *testing.T) {
		apps := Apps{{
			Name: "test",
			Keybinds: []KeyBind{{
				Name: "foo",
				Key:  "bar",
			}},
		}}
		want := Apps{{
			Name: "test",
			Keybinds: []KeyBind{{
				Name: "foo",
				Key:  "bar",
			}, {
				Name: "addFoo",
				Key:  "addBar",
			}},
		}}
		got := apps.addOrUpdate("test", "addFoo", "addBar", false)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("add new", func(t *testing.T) {
		apps := Apps{{
			Name: "test",
			Keybinds: []KeyBind{{
				Name: "foo",
				Key:  "bar",
			}},
		}}
		want := Apps{
			{
				Name: "test",
				Keybinds: []KeyBind{{
					Name: "foo",
					Key:  "bar",
				}},
			}, {
				Name: "new",
				Keybinds: []KeyBind{{
					Name: "addFoo",
					Key:  "addBar",
				}},
			},
		}
		got := apps.addOrUpdate("new", "addFoo", "addBar", false)

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
