package handsfree

import "testing"

// Must never fail, i mean literally
func TestGenerateCodeLong(t *testing.T) {
	hash, err := GenerateCodeLong("com.iam.the.app")
	if err != nil {
		t.Fatalf("Literally said must never fail, but: %v", err)
	}

	t.Logf("Enjoy your hash: %s (%d)", hash, len(hash))
}

func TestShortenToFit(t *testing.T) {
	hash, _ := GenerateCodeLong("com.iam.the.app")
	db := make(map[string]bool)

	// First attempt
	r, err := ShortenToFit(hash, db)
	if err != nil {
		t.Fatalf("Failed to shorten on first attempt: %v", err)
	}
	t.Logf("First attempt code: %s", r)

	db[r] = true

	if r != hash[0:3]+hash[29:32] {
		t.Fatalf("Unexpected result: %s, expected: %s", r, hash[0:3]+hash[29:32])
	}

	// Second attempt
	r, err = ShortenToFit(hash, db)
	if err != nil {
		t.Fatalf("Failed to shorten on second attempt: %v", err)
	}
	t.Logf("Second attempt code: %s", r)

	db[r] = true

	if len(db) != 2 {
		t.Log("Keys map expected to be of length two by now, but it's not")
	}

	if r != hash[1:4]+hash[28:31] {
		t.Fatalf("Unexpected result: %s, expected: %s", r, hash[0:3]+hash[29:32])
	}

	// Testing "DB" overflow
	for i := 0; i < 27; i++ {
		r, err := ShortenToFit(hash, db)
		if err != nil {
			t.Fatalf("Unexpected error during keys table overflow test on iteration #%d: %v", i, err)
		}
		db[r] = true
	}

	// Final attempt (DB must be overflowed)
	r, err = ShortenToFit(hash, db)
	if err == nil {
		t.Fatalf("Function was expected to fail now, but it worked: %s", r)
	}
}
