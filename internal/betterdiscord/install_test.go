package betterdiscord

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/betterdiscord/cli/internal/models"
)

func TestNew(t *testing.T) {
	rootPath := "/test/root/BetterDiscord"
	install := New(rootPath)

	if install.Root() != rootPath {
		t.Errorf("Root() = %s, expected %s", install.Root(), rootPath)
	}

	expectedData := filepath.Join(rootPath, "data")
	if install.Data() != expectedData {
		t.Errorf("Data() = %s, expected %s", install.Data(), expectedData)
	}

	expectedAsar := filepath.Join(rootPath, "data", "betterdiscord.asar")
	if install.Asar() != expectedAsar {
		t.Errorf("Asar() = %s, expected %s", install.Asar(), expectedAsar)
	}

	expectedPlugins := filepath.Join(rootPath, "plugins")
	if install.Plugins() != expectedPlugins {
		t.Errorf("Plugins() = %s, expected %s", install.Plugins(), expectedPlugins)
	}

	expectedThemes := filepath.Join(rootPath, "themes")
	if install.Themes() != expectedThemes {
		t.Errorf("Themes() = %s, expected %s", install.Themes(), expectedThemes)
	}

	if install.HasDownloaded() {
		t.Error("HasDownloaded() should be false for new install")
	}
}

func TestBDInstall_GettersSetters(t *testing.T) {
	rootPath := "/test/path"
	install := New(rootPath)

	// Test all getters return expected paths
	tests := []struct {
		name     string
		getter   func() string
		expected string
	}{
		{"Root", install.Root, rootPath},
		{"Data", install.Data, filepath.Join(rootPath, "data")},
		{"Asar", install.Asar, filepath.Join(rootPath, "data", "betterdiscord.asar")},
		{"Plugins", install.Plugins, filepath.Join(rootPath, "plugins")},
		{"Themes", install.Themes, filepath.Join(rootPath, "themes")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.getter()
			if result != tt.expected {
				t.Errorf("%s() = %s, expected %s", tt.name, result, tt.expected)
			}
		})
	}
}

func TestBDInstall_Prepare(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	bdRoot := filepath.Join(tmpDir, "BetterDiscord")

	install := New(bdRoot)

	// Prepare should create all necessary directories
	err := install.Prepare()
	if err != nil {
		t.Fatalf("Prepare() failed: %v", err)
	}

	// Verify directories were created
	dirs := []string{
		install.Data(),
		install.Plugins(),
		install.Themes(),
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory not created: %s", dir)
		}
	}
}

func TestBDInstall_Prepare_AlreadyExists(t *testing.T) {
	// Create temporary directory with existing structure
	tmpDir := t.TempDir()
	bdRoot := filepath.Join(tmpDir, "BetterDiscord")

	install := New(bdRoot)

	// Create directories manually first
	os.MkdirAll(install.Data(), 0755)    //nolint:errcheck
	os.MkdirAll(install.Plugins(), 0755) //nolint:errcheck
	os.MkdirAll(install.Themes(), 0755)  //nolint:errcheck

	// Prepare should succeed even if directories already exist
	err := install.Prepare()
	if err != nil {
		t.Fatalf("Prepare() failed when directories already exist: %v", err)
	}

	// Verify directories still exist
	dirs := []string{
		install.Data(),
		install.Plugins(),
		install.Themes(),
	}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("Directory should exist: %s", dir)
		}
	}
}

func TestBDInstall_Repair(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	bdRoot := filepath.Join(tmpDir, "BetterDiscord")

	install := New(bdRoot)

	// Create the data directory and a test plugins.json file
	channelFolder := filepath.Join(install.Data(), models.Stable.String())
	os.MkdirAll(channelFolder, 0755) //nolint:errcheck

	pluginsJson := filepath.Join(channelFolder, "plugins.json")
	err := os.WriteFile(pluginsJson, []byte(`{"test": "data"}`), 0644)
	if err != nil {
		t.Fatalf("Failed to create test plugins.json: %v", err)
	}

	// Verify file exists before repair
	if _, err := os.Stat(pluginsJson); os.IsNotExist(err) {
		t.Fatal("plugins.json should exist before repair")
	}

	// Run repair
	err = install.Repair(models.Stable)
	if err != nil {
		t.Fatalf("Repair() failed: %v", err)
	}

	// Verify file was removed
	if _, err := os.Stat(pluginsJson); !os.IsNotExist(err) {
		t.Error("plugins.json should be removed after repair")
	}
}

func TestBDInstall_Repair_NoPluginsFile(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	bdRoot := filepath.Join(tmpDir, "BetterDiscord")

	install := New(bdRoot)

	// Don't create any files - repair should succeed without error
	err := install.Repair(models.Stable)
	if err != nil {
		t.Fatalf("Repair() should succeed when plugins.json doesn't exist: %v", err)
	}
}

func TestBDInstall_Repair_MultipleChannels(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	bdRoot := filepath.Join(tmpDir, "BetterDiscord")

	install := New(bdRoot)

	// Create plugins.json for multiple channels
	channels := []models.DiscordChannel{models.Stable, models.Canary, models.PTB}
	pluginsFiles := make(map[models.DiscordChannel]string)

	for _, channel := range channels {
		channelFolder := filepath.Join(install.Data(), channel.String())
		os.MkdirAll(channelFolder, 0755) //nolint:errcheck

		pluginsJson := filepath.Join(channelFolder, "plugins.json")
		os.WriteFile(pluginsJson, []byte(`{}`), 0644) //nolint:errcheck
		pluginsFiles[channel] = pluginsJson
	}

	// Repair only Stable channel
	err := install.Repair(models.Stable)
	if err != nil {
		t.Fatalf("Repair(Stable) failed: %v", err)
	}

	// Verify only Stable's plugins.json was removed
	if _, err := os.Stat(pluginsFiles[models.Stable]); !os.IsNotExist(err) {
		t.Error("Stable plugins.json should be removed")
	}

	// Verify other channels' files still exist
	if _, err := os.Stat(pluginsFiles[models.Canary]); os.IsNotExist(err) {
		t.Error("Canary plugins.json should still exist")
	}
	if _, err := os.Stat(pluginsFiles[models.PTB]); os.IsNotExist(err) {
		t.Error("PTB plugins.json should still exist")
	}
}

func TestGetInstallation_WithBase(t *testing.T) {
	basePath := "/test/config"
	install := GetInstallation(basePath)

	expectedRoot := filepath.Join(basePath, "BetterDiscord")
	if install.Root() != expectedRoot {
		t.Errorf("GetInstallation(%s).Root() = %s, expected %s", basePath, install.Root(), expectedRoot)
	}
}

func TestGetInstallation_Singleton(t *testing.T) {
	// Reset the global instance to ensure clean test
	// Note: In a real test, you might want to add a reset function
	// For now, we'll just test that multiple calls work

	install1 := GetInstallation()
	install2 := GetInstallation()

	// Both should return the same instance (singleton pattern)
	if install1 != install2 {
		t.Error("GetInstallation() should return the same instance (singleton)")
	}

	// Both should have the same root path
	if install1.Root() != install2.Root() {
		t.Errorf("Singleton instances have different roots: %s vs %s", install1.Root(), install2.Root())
	}
}

func TestMakeDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	testDir := filepath.Join(tmpDir, "test", "nested", "directory")

	// Make the directory
	err := makeDirectory(testDir)
	if err != nil {
		t.Fatalf("makeDirectory() failed: %v", err)
	}

	// Verify it exists
	if _, err := os.Stat(testDir); os.IsNotExist(err) {
		t.Error("Directory was not created")
	}

	// Test making the same directory again (should not error)
	err = makeDirectory(testDir)
	if err != nil {
		t.Errorf("makeDirectory() should succeed when directory already exists: %v", err)
	}
}

func TestBDInstall_HasDownloaded(t *testing.T) {
	install := New("/test/path")

	// Initially should be false
	if install.HasDownloaded() {
		t.Error("HasDownloaded() should initially be false")
	}

	// After setting hasDownloaded (internal state)
	install.hasDownloaded = true
	if !install.HasDownloaded() {
		t.Error("HasDownloaded() should be true after download")
	}
}

func TestBDInstall_PathStructure(t *testing.T) {
	// Test with different root paths
	tests := []struct {
		name     string
		rootPath string
	}{
		{"Unix absolute path", "/home/user/.config/BetterDiscord"},
		{"Windows absolute path", "C:\\Users\\User\\AppData\\Roaming\\BetterDiscord"},
		{"Relative path", "BetterDiscord"},
		{"Path with spaces", "/path with spaces/BetterDiscord"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			install := New(tt.rootPath)

			// Verify all paths are constructed correctly relative to root
			if install.Root() != tt.rootPath {
				t.Errorf("Root() incorrect")
			}

			if install.Data() != filepath.Join(tt.rootPath, "data") {
				t.Errorf("Data() path incorrect: %s", install.Data())
			}

			if install.Asar() != filepath.Join(tt.rootPath, "data", "betterdiscord.asar") {
				t.Errorf("Asar() path incorrect: %s", install.Asar())
			}

			if install.Plugins() != filepath.Join(tt.rootPath, "plugins") {
				t.Errorf("Plugins() path incorrect: %s", install.Plugins())
			}

			if install.Themes() != filepath.Join(tt.rootPath, "themes") {
				t.Errorf("Themes() path incorrect: %s", install.Themes())
			}
		})
	}
}
