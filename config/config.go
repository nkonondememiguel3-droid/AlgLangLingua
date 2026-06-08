package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

// Meta holds language metadata.
type Meta struct {
	Name       string `toml:"name"`
	Version    string `toml:"version"`
	DecimalSep string `toml:"decimal_sep"`
}

// Keywords holds the surface-name mapping for every logical keyword.
// Each field name matches the TOML key exactly.
type Keywords struct {
	Algorithm string `toml:"algorithm"`
	Variable  string `toml:"variable"`
	Constant  string `toml:"constant"`
	Type      string `toml:"type"`
	Begin     string `toml:"begin"`
	End       string `toml:"end"`

	Function    string `toml:"function"`
	EndFunction string `toml:"end_function"`
	Method      string `toml:"method"`
	EndMethod   string `toml:"end_method"`
	Return      string `toml:"return"`

	Structure string `toml:"structure"`
	EndStruct string `toml:"end_struct"`

	If       string `toml:"if"`
	Then     string `toml:"then"`
	Else     string `toml:"else"`
	ElseIf   string `toml:"elseif"`
	EndIf    string `toml:"endif"`
	For      string `toml:"for"`
	To       string `toml:"to"`
	Step     string `toml:"step"`
	EndFor   string `toml:"endfor"`
	While    string `toml:"while"`
	Do       string `toml:"do"`
	EndWhile string `toml:"endwhile"`
	Repeat   string `toml:"repeat"`
	Until    string `toml:"until"`

	IntegerType   string `toml:"integer_type"`
	DoubleType    string `toml:"double_type"`
	StringType    string `toml:"string_type"`
	CharacterType string `toml:"character_type"`
	BooleanType   string `toml:"boolean_type"`
	Table         string `toml:"table"`
	Of            string `toml:"of"`

	And string `toml:"and"`
	Or  string `toml:"or"`
	Not string `toml:"not"`
	Mod string `toml:"mod"`

	Write string `toml:"write"`
	Read  string `toml:"read"`

	True  string `toml:"true"`
	False string `toml:"false"`

	Nil   string `toml:"nil"`
	Class string `toml:"class"`
}

// Config is the top-level structure of a .toml language file.
type Config struct {
	Meta     Meta     `toml:"meta"`
	Keywords Keywords `toml:"keywords"`
}

// Load reads and parses a TOML language config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: cannot read %q: %w", path, err)
	}

	var cfg Config
	if _, err := toml.Decode(string(data), &cfg); err != nil {
		return nil, fmt.Errorf("config: cannot parse %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("config: invalid %q: %w", path, err)
	}

	return &cfg, nil
}

// validate checks that no keyword surface name is empty or duplicated.
func (c *Config) validate() error {
	seen := make(map[string]string) // surface -> logical name

	check := func(logical, surface string) error {
		if surface == "" {
			return fmt.Errorf("keyword %q has no surface name", logical)
		}
		if prev, clash := seen[surface]; clash {
			return fmt.Errorf("surface name %q is used by both %q and %q", surface, prev, logical)
		}
		seen[surface] = logical
		return nil
	}

	kw := c.Keywords
	pairs := [][2]string{
		{"algorithm", kw.Algorithm},
		{"variable", kw.Variable},
		{"constant", kw.Constant},
		{"type", kw.Type},
		{"begin", kw.Begin},
		{"end", kw.End},
		{"function", kw.Function},
		{"end_function", kw.EndFunction},
		{"method", kw.Method},
		{"end_method", kw.EndMethod},
		{"return", kw.Return},
		{"structure", kw.Structure},
		{"end_struct", kw.EndStruct},
		{"if", kw.If},
		{"then", kw.Then},
		{"else", kw.Else},
		{"elseif", kw.ElseIf},
		{"endif", kw.EndIf},
		{"for", kw.For},
		{"to", kw.To},
		{"step", kw.Step},
		{"endfor", kw.EndFor},
		{"while", kw.While},
		{"do", kw.Do},
		{"endwhile", kw.EndWhile},
		{"repeat", kw.Repeat}, {"until", kw.Until},
		{"integer_type", kw.IntegerType},
		{"double_type", kw.DoubleType},
		{"string_type", kw.StringType},
		{"character_type", kw.CharacterType},
		{"boolean_type", kw.BooleanType},
		{"table", kw.Table},
		{"of", kw.Of},
		{"and", kw.And},
		{"or", kw.Or},
		{"not", kw.Not},
		{"mod", kw.Mod},
		{"write", kw.Write},
		{"read", kw.Read},
		{"true", kw.True},
		{"false", kw.False},
		{"nil", kw.Nil},
		{"class", kw.Class},
	}

	for _, p := range pairs {
		if err := check(p[0], p[1]); err != nil {
			return err
		}
	}
	return nil
}
