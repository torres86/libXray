package libXray

import (
	"os"
	"strings"
	"testing"
)

const (
	stableGoVersion       = "1.24.6"
	stableGomobileVersion = "v0.0.0-20250813145510-f12310a0cfd9"
	stableXrayCoreVersion = "v1.250803.0"
)

func TestAndroidReleaseBuildUsesStableXrayToolchain(t *testing.T) {
	goMod := readProjectFile(t, "go.mod")
	workflow := readProjectFile(t, ".github/workflows/release-android.yml")
	buildScript := readProjectFile(t, "build/app/build.py")

	assertContains(t, goMod, "github.com/xtls/xray-core "+stableXrayCoreVersion)
	assertContains(t, goMod, "golang.org/x/mobile "+stableGomobileVersion)
	assertContains(t, workflow, "GO_VERSION: '"+stableGoVersion+"'")
	assertContains(t, workflow, "GOMOBILE_VERSION: '"+stableGomobileVersion+"'")
	assertContains(t, workflow, "XRAY_CORE_VERSION: '"+stableXrayCoreVersion+"'")
	assertContains(t, workflow, "- **Xray-core**: `${{ env.XRAY_CORE_VERSION }}`")
	assertContains(t, buildScript, `("gomobile", "gobind")`)
	assertContains(t, buildScript, "golang.org/x/mobile/cmd/{tool}@{gomobile_version}")
}

func readProjectFile(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}

func assertContains(t *testing.T, content string, expected string) {
	t.Helper()
	if !strings.Contains(content, expected) {
		t.Fatalf("expected file content to contain %q", expected)
	}
}
