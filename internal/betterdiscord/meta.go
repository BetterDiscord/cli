package betterdiscord

import (
	"regexp"
	"strings"
)

var splitRegex = regexp.MustCompile(`(?m)[^\S\r\n]*?\r?(?:\r\n|\n)[^\S\r\n]*?\*[^\S\r\n]?`)
var escapedAtRegex = regexp.MustCompile(`^\\@`)

type Meta struct {
    Name        string
    Author      string
    Description string
    Version     string
    Invite      string
    AuthorID    string
    AuthorLink  string
    Donate      string
    Patreon     string
    Website     string
    Source      string
}

func parseJSDoc(fileContent string) Meta {
    meta := Meta{}

    parts := strings.SplitN(fileContent, "/**", 2)
    if len(parts) < 2 {
        return meta
    }
    blockParts := strings.SplitN(parts[1], "*/", 2)
    if len(blockParts) < 1 {
        return meta
    }
    block := blockParts[0]

    field := ""
    accum := ""

    lines := splitRegex.Split(block, -1)
    for _, line := range lines {
        line = strings.TrimSpace(line)
        if len(line) == 0 {
            continue
        }

        if strings.HasPrefix(line, "@") && (len(line) > 1 && line[1] != ' ') {
            // flush previous field
            if accum != "" {
                assignField(&meta, field, strings.TrimSpace(accum))
            }

            // new field
            l := strings.Index(line, " ")
            if l == -1 {
                field = strings.TrimPrefix(line, "@")
                accum = ""
            } else {
                field = line[1:l]
                accum = line[l+1:]
            }
        } else {
            // accumulate prose
            if escapedAtRegex.MatchString(line) {
                line = strings.Replace(line, `\@`, "@", 1)
            }
            line = strings.ReplaceAll(line, `\n`, "\n")
            accum += " " + line
        }
    }

    // flush last field
    if accum != "" {
        assignField(&meta, field, strings.TrimSpace(accum))
    }

    return meta
}

func assignField(meta *Meta, field, value string) {
    switch strings.ToLower(field) {
    case "name":
        meta.Name = value
    case "author":
        meta.Author = value
    case "description":
        meta.Description = value
    case "version":
        meta.Version = value
    case "invite":
        meta.Invite = value
    case "authorid":
        meta.AuthorID = value
    case "authorlink":
        meta.AuthorLink = value
    case "donate":
        meta.Donate = value
    case "patreon":
        meta.Patreon = value
    case "website":
        meta.Website = value
    case "source":
        meta.Source = value
    }
}
