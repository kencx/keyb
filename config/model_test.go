package config

import (
	"reflect"
	"testing"
)

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
