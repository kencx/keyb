package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	testKeybFile = "testkeyb.yml"
	tempKeybFile = "temp.yml"
)

func TestReadKeybFile(t *testing.T) {
	t.Run("file present", func(t *testing.T) {
		want := Apps{{
			Name: "test",
			Keybinds: []KeyBind{{
				Name: "foo",
				Key:  "bar",
			}},
		}}

		got, err := ReadKeybFile(filepath.Join(testDirPath, testKeybFile))
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("file absent", func(t *testing.T) {
		testpath := filepath.Join(testDirPath, tempKeybFile)

		_, err := ReadKeybFile(testpath)
		if err == nil {
			t.Errorf("expected err: %s does not exist", testpath)
		}

		t.Cleanup(func() {
			os.RemoveAll(filepath.Join(testDirPath, tempKeybFile))
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
		apps.addOrUpdate("test", "addFoo", "addBar", false)

		if !reflect.DeepEqual(apps, want) {
			t.Errorf("got %v, want %v", apps, want)
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
		apps.addOrUpdate("new", "addFoo", "addBar", false)

		if !reflect.DeepEqual(apps, want) {
			t.Errorf("got %v, want %v", apps, want)
		}
	})
}
