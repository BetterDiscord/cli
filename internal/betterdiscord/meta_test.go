package betterdiscord

import (
	"reflect"
	"testing"
)

func TestParseFullMeta(t *testing.T) {
    input := `
/**
 * @name ExampleAddon
 * @author YourName
 * @description Describe the basic information. Maybe a support server link.
 * @version 0.0.1
 * @invite inviteCode
 * @authorId 51512151151651
 * @authorLink https://twitter.com/Whoever
 * @donate https://paypal.me/
 * @patreon https://patreon.com/
 * @website https://github.com/BetterDiscord/BetterDiscord
 * @source https://gist.github.com/zerebos/e5f4d02fc3085a53872b0236cd6f8225
 */
`

    got := parseJSDoc(input)
    want := Meta{
        Name:        "ExampleAddon",
        Author:      "YourName",
        Description: "Describe the basic information. Maybe a support server link.",
        Version:     "0.0.1",
        Invite:      "inviteCode",
        AuthorID:    "51512151151651",
        AuthorLink:  "https://twitter.com/Whoever",
        Donate:      "https://paypal.me/",
        Patreon:     "https://patreon.com/",
        Website:     "https://github.com/BetterDiscord/BetterDiscord",
        Source:      "https://gist.github.com/zerebos/e5f4d02fc3085a53872b0236cd6f8225",
    }

    if !reflect.DeepEqual(got, want) {
        t.Errorf("parseJSDoc() = %#v\nwant %#v", got, want)
    }
}


func TestParseMissingOptionalFields(t *testing.T) {
    input := `
/**
 * @name Foo
 * @author Bar
 * @description Baz
 * @version 1.2.3
 */
`

    got := parseJSDoc(input)

    if got.Name != "Foo" || got.Author != "Bar" || got.Description != "Baz" || got.Version != "1.2.3" {
        t.Errorf("required fields not parsed correctly: %#v", got)
    }

    // optional fields should be empty
    if got.Invite != "" || got.AuthorID != "" || got.Website != "" {
        t.Errorf("optional fields should be empty: %#v", got)
    }
}


func TestParseEscapedCharacters(t *testing.T) {
    input := `
/**
 * @description Line1
 * Line2 and literal
 * \@ symbol
 */
`

    got := parseJSDoc(input)

    if got.Description != "Line1 Line2 and literal @ symbol" {
        t.Errorf("escaped characters not handled: %#v", got.Description)
    }
}


func TestParseUnknownFields(t *testing.T) {
    input := `
/**
 * @name Foo
 * @unknownField shouldBeIgnored
 * @version 1.0.0
 */
`

    got := parseJSDoc(input)

    if got.Name != "Foo" || got.Version != "1.0.0" {
        t.Errorf("known fields incorrect: %#v", got)
    }
}


func TestParseNoJSDoc(t *testing.T) {
    input := `console.log("no jsdoc here");`

    got := parseJSDoc(input)

    if got != (Meta{}) {
        t.Errorf("expected empty Meta, got %#v", got)
    }
}


func TestParseEmptyJSDoc(t *testing.T) {
    input := `
/**
 */
`
	println("Testing empty JSDoc")
    got := parseJSDoc(input)

    if got != (Meta{}) {
        t.Errorf("expected empty Meta, got %#v", got)
    }
}


func TestParseJSDoc(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected Meta
    }{
        {
            name: "Full meta block",
            input: `
/**
 * @name ExampleAddon
 * @author YourName
 * @description Describe the basic information. Maybe a support server link.
 * @version 0.0.1
 * @invite inviteCode
 * @authorId 51512151151651
 * @authorLink https://twitter.com/Whoever
 * @donate https://paypal.me/
 * @patreon https://patreon.com/
 * @website https://github.com/BetterDiscord/BetterDiscord
 * @source https://gist.github.com/zerebos/e5f4d02fc3085a53872b0236cd6f8225
 */
`,
            expected: Meta{
                Name:        "ExampleAddon",
                Author:      "YourName",
                Description: "Describe the basic information. Maybe a support server link.",
                Version:     "0.0.1",
                Invite:      "inviteCode",
                AuthorID:    "51512151151651",
                AuthorLink:  "https://twitter.com/Whoever",
                Donate:      "https://paypal.me/",
                Patreon:     "https://patreon.com/",
                Website:     "https://github.com/BetterDiscord/BetterDiscord",
                Source:      "https://gist.github.com/zerebos/e5f4d02fc3085a53872b0236cd6f8225",
            },
        },
        {
            name: "Required fields only",
            input: `
/**
 * @name Foo
 * @author Bar
 * @description Baz
 * @version 1.2.3
 */
`,
            expected: Meta{
                Name:        "Foo",
                Author:      "Bar",
                Description: "Baz",
                Version:     "1.2.3",
            },
        },
        {
            name: "Escaped characters",
            input: `
/**
 * @description Line1
 * Line2 and literal
 * \@ symbol
 */
`,
            expected: Meta{
                Description: "Line1 Line2 and literal @ symbol",
            },
        },
        {
            name: "Unknown fields ignored",
            input: `
/**
 * @name Foo
 * @unknownField shouldBeIgnored
 * @version 1.0.0
 */
`,
            expected: Meta{
                Name:    "Foo",
                Version: "1.0.0",
            },
        },
        {
            name:  "No JSDoc block",
            input: `console.log("no jsdoc here");`,
            expected: Meta{},
        },
        {
            name: "Empty JSDoc block",
            input: `
/**
 */
`,
            expected: Meta{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := parseJSDoc(tt.input)
            if !reflect.DeepEqual(got, tt.expected) {
                t.Errorf("parseJSDoc() = %#v\nwant %#v", got, tt.expected)
            }
        })
    }
}

func FuzzParseJSDoc(f *testing.F) {
    // Seed with valid examples
    f.Add("/** @name Foo */")
    f.Add("/** @description Hello\\nWorld */")
    f.Add("/** @author Zerebos */")
    f.Add("no jsdoc here")

    f.Fuzz(func(t *testing.T, input string) {
        // The only requirement: the parser must never panic
        defer func() {
            if r := recover(); r != nil {
                t.Fatalf("parseJSDoc panicked with input: %q\npanic: %v", input, r)
            }
        }()

        _ = parseJSDoc(input)
    })
}

func BenchmarkParseJSDoc(b *testing.B) {
    input := `
/**
 * @name ExampleAddon
 * @author Zerebos
 * @description Something
 * @version 1.0.0
 */
`
    for i := 0; i < b.N; i++ {
        parseJSDoc(input)
    }
}
