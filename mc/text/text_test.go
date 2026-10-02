package text

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/admin-else/strom/mc/nbt"
	"github.com/google/go-cmp/cmp"
)

func TestComponent_JSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    RawComponent
		wantErr bool
	}{
		{
			name: "simple text",
			json: `"hello"`,
			want: RawComponent{Text: "hello"},
		},
		{
			name: "object text",
			json: `{"text":"hello","bold":true}`,
			want: RawComponent{Text: "hello", Bold: ptr(true)},
		},
		{
			name: "array of components",
			json: `["hello", {"text":" world","color":"red"}]`,
			want: RawComponent{
				Text: "hello",
				Extra: []RawComponent{
					{Text: " world", Color: "red"},
				},
			},
		},
		{
			name: "complex translation",
			json: `{"translate":"chat.type.text","with":["Player",{"text":"hello","italic":true}]}`,
			want: RawComponent{
				Translate: "chat.type.text",
				With: []RawComponent{
					{Text: "Player"},
					{Text: "hello", Italic: ptr(true)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got RawComponent
			if err := json.Unmarshal([]byte(tt.json), &got); (err != nil) != tt.wantErr {
				t.Fatalf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UnmarshalJSON() got = %+v, want %+v", got, tt.want)
			}

			marshaled, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("MarshalJSON() error = %v", err)
			}

			var got2 RawComponent
			if err := json.Unmarshal(marshaled, &got2); err != nil {
				t.Fatalf("UnmarshalJSON() back error = %v", err)
			}
			if !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("Roundtrip failed: got = %+v, want %+v", got2, tt.want)
			}
		})
	}
}

func TestComponent_NBT(t *testing.T) {
	tests := []struct {
		name string
		nbt  any
		want RawComponent
	}{
		{
			name: "simple text",
			nbt:  "hello",
			want: RawComponent{Text: "hello"},
		},
		{
			name: "object text",
			nbt: map[string]any{
				"text": "hello",
				"bold": int8(1),
			},
			want: RawComponent{Text: "hello", Bold: ptr(true)},
		},
		{
			name: "extra components",
			nbt: map[string]any{
				"text": "hello",
				"extra": []any{
					map[string]any{"text": " world", "color": "red"},
				},
			},
			want: RawComponent{
				Text: "hello",
				Extra: []RawComponent{
					{Text: " world", Color: "red"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got RawComponent
			if err := got.FromNBT(tt.nbt); err != nil {
				t.Fatalf("FromNBT() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromNBT() got = %+v, want %+v", got, tt.want)
			}

			nbtVal := got.ToNBT()
			var got2 RawComponent
			if err := got2.FromNBT(nbtVal); err != nil {
				t.Fatalf("FromNBT() back error = %v", err)
			}
			if !reflect.DeepEqual(got2, tt.want) {
				t.Errorf("Roundtrip failed: got = %+v, want %+v", got2, tt.want)
			}
		})
	}
}

func TestPretty(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want RawComponent
	}{
		{
			name: "simple text",
			s:    "Hello",
			want: RawComponent{Text: "Hello"},
		},
		{
			name: "color and bold",
			s:    "§c§lHello",
			want: RawComponent{
				Extra: []RawComponent{
					{Text: "Hello", Color: "red", Bold: ptr(true)},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Pretty(tt.s)
			if got.Text != tt.want.Text {
				t.Errorf("Pretty().Text = %v, want %v", got.Text, tt.want.Text)
			}
			if !reflect.DeepEqual(got.Extra, tt.want.Extra) {
				t.Errorf("Pretty().Extra = %+v, want %+v", got.Extra, tt.want.Extra)
			}
		})
	}
}

func TestPrettyF(t *testing.T) {
	got := PrettyF("Hello §b%s", "World")
	want := &RawComponent{
		Text: "Hello ",
		Extra: []RawComponent{
			{Text: "World", Color: "aqua"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("PrettyF() = \n%+v\n, want \n%+v", got, want)
	}
}

func TestComponent_String(t *testing.T) {
	tests := []struct {
		name string
		comp RawComponent
		want string
	}{
		{
			name: "simple text",
			comp: RawComponent{Text: "Hello"},
			want: "Hello",
		},
		{
			name: "color and bold",
			comp: RawComponent{
				Extra: []RawComponent{
					{Text: "Hello", Color: "red", Bold: ptr(true)},
				},
			},
			want: "§c§lHello",
		},
		{
			name: "hex color",
			comp: RawComponent{
				Color: "#ff0000",
				Text:  "Red",
			},
			want: "§x§f§f§0§0§0§0Red",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.comp.String(); got != tt.want {
				t.Errorf("RawComponent.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrettyToNbt(t *testing.T) {
	tests := []struct {
		name string
		f    string
		want nbt.Anon
	}{
		{
			name: "real data",
			f:    `§dno gooning`,
			want: nbt.Anon{Value: map[string]interface{}{"extra": []interface{}{map[string]interface{}{"color": "light_purple", "text": "no gooning"}}, "text": ""}},
		},
		{
			name: "reset thing",
			f:    "§3Any§bProtect§3 » §rffff",
			want: nbt.Anon{Value: map[string]interface{}{"extra": []interface{}{map[string]interface{}{"color": "dark_aqua", "text": "Any"}, map[string]interface{}{"color": "aqua", "text": "Protect"}, map[string]interface{}{"color": "dark_aqua", "text": " » "}, map[string]interface{}{"text": "ffff"}}, "text": ""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := Pretty(tt.f).AsNBT()
			b := tt.want
			if !cmp.Equal(a, b) {
				t.Errorf("PrettyToNbt() = %#v, want %#v", a, b)
			}
		})
	}
}
