package badger

import (
	"testing"
)

func TestPasswordHashingPass(t *testing.T) {

	password := "anteater"

	hash, e := Hash([]byte(password))
	if e != nil {
		t.Fatalf("failed to hash password: %v", e)
	}

	compPass := []byte(password)
	if err := RawIsHashMatch(compPass, hash); err != nil {
		t.Fatalf("failed to check hash %v", err)
	}

	for i, _ := range compPass {
		if compPass[i] != 0 {
			t.Fatalf("failed to clear password after match")
		}
	}
}

func TestPasswordHashingZeroCheck(t *testing.T) {

	password := []byte("anteater")

	hash, e := Hash(password)
	if e != nil {
		t.Fatalf("failed to hash password: %v", e)
	}

	if err := RawIsHashMatch(password, hash); err == nil {
		t.Fatal("failed to zero-out password - expected failiure")
	}
}

func TestBadger(t *testing.T) {

	for i := 0; i < 4; i++ {
		badge, err := New(Config{"honey_badger.dgaf"})

		if err != nil {
			t.Fatalf("error:%v", err)
		}

		badgeActual := badge.(*identity)

		for _, tsd := range testSignData {

			raw, err := badge.Sign(&tsd)
			if err != nil {
				t.Fatalf("error:%v", err)
			}

			if !Verify(raw) {
				t.Fatalf("unable to verify signature")
			}

			encodedId := badge.EncodeIdentityString()

			actualPublicKey, err := publicKeyFromB64(badge.PublicKey())
			if err != nil {
				t.Fatalf("err:%v", err)
			}

			if !actualPublicKey.Equal(&badgeActual.key.PublicKey) {
				t.Fatal("public key encode/decode failure")
			}

			decodedId, err := DecodeIdentityString(encodedId)
			if err != nil {
				t.Fatalf("err:%v", err)
			}

			if decodedId.Id() != badge.Id() {
				t.Fatal("decode failure - 0")
			}

			if decodedId.Nickname() != badge.Nickname() {
				t.Fatal("decode failure - 1")
			}

			if decodedId.PublicKey() != badge.PublicKey() {
				t.Fatal("decode failure - 2")
			}

			decodedSignature, err := badge.Sign(&tsd)

			if err != nil {
				t.Fatalf("error:%v", err)
			}

			if !Verify(decodedSignature) {
				t.Fatal("unable to verify signature")
			}

			encodedDecodedId := decodedId.EncodeIdentityString()

			if encodedDecodedId != encodedId {
				t.Fatal("decode failure -3")
			}
		}
	}
}

var testSignData = []string{
	"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam sed dui dui.",
	"Pellentesque vitae mattis elit, in dapibus nunc. Sed molestie vehicula dignissim.",
	" Cras volutpat risus metus, in volutpat tellus consequat sit amet.",
	"Suspendisse mattis velit ornare vehicula pellentesque. Sed vel lectus nec velit cursus",
	" pretium. Sed fermentum rhoncus arcu non commodo. Maecenas placerat elementum accumsan. ",
	"Ut molestie felis eget dolor imperdiet rutrum. Aliquam nec maximus purus, a ullamcorper neque.",
	"Vestibulum molestie maximus elit, vel condimentum eros auctor non.",
	" Nullam cursus est et urna malesuada, a placerat nisl elementum. Praesent metus ligula, cursus",
	"ut suscipit at, malesuada a purus. Pellentesque habitant morbi tristique senectus et netus et",
	" malesuada fames ac turpis egestas. Cras vel metus dolor. Cras mattis ultricies quam nec sodales.",
	"Maecenas rhoncus, elit vel tincidunt maximus, leo urna mollis ipsum, non lacinia erat turpis eu metus.",
	" Vivamus pretium ante tristique augue molestie, eget aliquam neque fringilla. Vestibulum pharetra lobortis cursus.",
}
