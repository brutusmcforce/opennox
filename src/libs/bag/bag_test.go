package bag

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/brutusmcforce/opennox/libs/noxtest"
)

// The hash of all decompressed sprite data of the game's video.bag, to check the NXZ decoder.
const videoBagHash = "1f5ff02f86c909e196980ea3c627709c30a5066d7685f11c162df4be1af112cc"

func TestVideoBagDecode(t *testing.T) {
	path := filepath.Join(noxtest.DataPath(t), "video.bag")
	if _, err := os.Stat(path); err != nil {
		t.Skip(err)
	}
	f, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	imgs, err := f.Images()
	if err != nil {
		t.Fatal(err)
	}
	h := sha256.New()
	for _, img := range imgs {
		data, err := img.Raw()
		if err != nil {
			t.Fatal(img.Name, err)
		}
		h.Write(data)
	}
	sum := hex.EncodeToString(h.Sum(nil))
	t.Logf("%d images, hash %s", len(imgs), sum)
	if videoBagHash != "" && sum != videoBagHash {
		t.Fatalf("hash %s, want %s", sum, videoBagHash)
	}
}
