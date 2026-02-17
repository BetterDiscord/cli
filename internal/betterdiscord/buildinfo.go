package betterdiscord

import (
	"bufio"
	"os"
	"regexp"

	"github.com/betterdiscord/cli/internal/output"
	"github.com/betterdiscord/cli/internal/utils"
)

type Buildinfo struct {
	Version string
	Commit  string
	Branch  string
	Mode    string
}

func NewBuildinfo() Buildinfo {
	return Buildinfo{
		Version: "unknown",
		Commit:  "unknown",
		Branch:  "unknown",
		Mode:    "unknown",
	}
}

func (i *BDInstall) ReadBuildinfo() (bi Buildinfo, err error) {
	if !utils.Exists(i.asar) {
		return NewBuildinfo(), os.ErrNotExist
	}

	f, err := os.Open(i.asar)
	if err != nil {
		return NewBuildinfo(), err
	}

	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Compile your regexes
	versionRe := regexp.MustCompile(`version:\s?"([0-9]+\.[0-9]+\.[0-9]+)"`)
	commitRe := regexp.MustCompile(`commit:\s?"(\b[0-9a-f]{5,40}\b)"`)
	branchRe := regexp.MustCompile(`branch:\s?"([a-zA-Z0-9_\-]+)"`)
	modeRe := regexp.MustCompile(`build:\s?"([a-zA-Z]+)"`)

	regexes := map[string]*regexp.Regexp{
		"version": versionRe,
		"commit":  commitRe,
		"branch":  branchRe,
		"mode":    modeRe,
	}

	buildinfo := NewBuildinfo()
	reader := bufio.NewReader(f)

	// 64 KB chunks are a nice balance
	const chunkSize = 64 * 1024
	buf := make([]byte, chunkSize)

	// Rolling window to catch matches across chunk boundaries
	var tail []byte

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			// Combine tail + new chunk
			window := append(tail, buf[:n]...)

			// Run all regexes
			for name, re := range regexes {
				matches := re.FindAllSubmatch(window, -1)
				for _, m := range matches {
					if len(m) > 1 {
						switch name {
						case "version":
							buildinfo.Version = string(m[1])
						case "commit":
							buildinfo.Commit = string(m[1])
						case "branch":
							buildinfo.Branch = string(m[1])
						case "mode":
							buildinfo.Mode = string(m[1])
						}
					}
				}
			}

			// Keep last 1 KB as tail (enough for your patterns)
			if len(window) > 1024 {
				tail = window[len(window)-1024:]
			} else {
				tail = window
			}
		}

		if err != nil {
			break
		}
	}

	i.Buildinfo = buildinfo
	return buildinfo, nil
}

func (bdinstall *BDInstall) LogBuildinfo() {
	output.Printf("📦 BetterDiscord Information:\n")

	buildinfo, err := bdinstall.ReadBuildinfo()
	if err == nil {
		output.Printf("   Build Information:\n")
		output.Printf("     🔹 Version: %s\n", buildinfo.Version)
		output.Printf("     🔹 Commit:  %s\n", buildinfo.Commit)
		output.Printf("     🔹 Branch:  %s\n", buildinfo.Branch)
		output.Printf("     🔹 Mode:    %s\n", buildinfo.Mode)
	}

	output.Printf("   Installation Paths:\n")
	output.Printf("     📁 Base:    %s\n", bdinstall.Root())
	output.Printf("     ⚙️  Data:    %s\n", bdinstall.Data())
	output.Printf("     🔌 Plugins: %s\n", bdinstall.Plugins())
	output.Printf("     🎨 Themes:  %s\n", bdinstall.Themes())
}
