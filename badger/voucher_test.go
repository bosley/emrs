package badger

import (
	"crypto/rand"
	b64 "encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"
)

type testDataStruct struct {
	badge   Badge
	voucher string
}

func TestVoucherValid(t *testing.T) {
	badge, _ := New("voucher-test")
	voucher, err := NewVoucher(badge, 30*time.Minute)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	if !ValidateVoucher(badge.PublicKey(), voucher) {
		t.Fatalf("faild to validate valid voucher")
	}
}

func TestVoucherInvalidDurationInit(t *testing.T) {
	badge, _ := New("voucher-test")
	_, err := NewVoucher(badge, -30*time.Minute)
	if err == nil {
		t.Fatalf("expected error for invalid duration")
	}
}

func TestVoucherInvalidExpired(t *testing.T) {
	badge, _ := New("voucher-test")
	voucher, err := NewVoucher(badge, 1*time.Second)
	if err != nil {
		t.Fatalf("err:%v", err)
	}
	time.Sleep(2 * time.Second)
	if ValidateVoucher(badge.PublicKey(), voucher) {
		t.Fatalf("expired voucher was validated")
	}
}

func TestVoucherInvalidEmptyVoucher(t *testing.T) {
	badge, _ := New("voucher-test")
	if ValidateVoucher(badge.PublicKey(), "") {
		t.Fatalf("empty validated")
	}
}

func TestVoucherInvalidData(t *testing.T) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelError,
				})))
	badge, _ := New("voucher-test")
	for x := 0; x < 100; x++ {
		data := make([]byte, 434)
		_, err := rand.Read(data)
		if err != nil {
			t.Fatalf("err:%v", err)
		}
		voucher := string(data)
		if ValidateVoucher(badge.PublicKey(), voucher) {
			t.Fatal("either voucher is broken or you just one the lottery 1.5M times over")
		}
	}
}

func TestVoucherSpecificInvalidData(t *testing.T) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(os.Stdout,
				&slog.HandlerOptions{
					Level: slog.LevelError,
				})))
	badge, _ := New("voucher-test")

	for x := 0; x < 100; x++ {

		header := make([]byte, 20)
		body := make([]byte, 196)
		info := make([]byte, 216)

		_, err := rand.Read(header)
		if err != nil {
			t.Fatalf("err:%v", err)
		}

		_, err = rand.Read(body)
		if err != nil {
			t.Fatalf("err:%v", err)
		}

		_, err = rand.Read(info)
		if err != nil {
			t.Fatalf("err:%v", err)
		}

		voucherA := fmt.Sprintf(
			"%s:%s:%s",
			string(header),
			string(body),
			string(info))

		voucherB := fmt.Sprintf(
			"%s:%s:%s",
			b64.StdEncoding.EncodeToString(header),
			b64.StdEncoding.EncodeToString(body),
			b64.StdEncoding.EncodeToString(info))

		if ValidateVoucher(badge.PublicKey(), voucherA) {
			t.Fatal("validated invalid loosly manually-crafted voucher")
		}
		if ValidateVoucher(badge.PublicKey(), voucherB) {
			t.Fatal("validated invalid specific manually-crafted voucher")
		}
	}
}
