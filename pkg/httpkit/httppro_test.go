package httpkit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

func ExampleNewProblemDetail() {
	// Example with extension by embedding ProblemDetail.
	type BalanceProblemDetail struct {
		*ProblemDetail
		Balance  int32    `json:"balance" xml:"balance"`
		Accounts []string `json:"accounts" xml:"accounts"`
	}

	prob := BalanceProblemDetail{
		ProblemDetail: NewProblemDetail(
			"https://example.com/probs/out-of-credit",
			ProbOpts.Detail("Your current balance is 30, but that costs 50."),
			ProbOpts.Instance("/account/12345/abc"),
			ProbOpts.Title("You do not have enough credit."),
		),
		Balance:  30,
		Accounts: []string{"/account/12345", "/account/67890"},
	}

	rec := httptest.NewRecorder()
	err := WriteJSONProblemDetail(rec, &prob, 403)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	if err := json.Indent(&out, rec.Body.Bytes(), "", "\t"); err != nil {
		panic(err)
	}

	fmt.Println(out.String())
	// Output:
	//{
	//	"type": "https://example.com/probs/out-of-credit",
	//	"title": "You do not have enough credit.",
	//	"status": 403,
	//	"detail": "Your current balance is 30, but that costs 50.",
	//	"instance": "/account/12345/abc",
	//	"balance": 30,
	//	"accounts": [
	//		"/account/12345",
	//		"/account/67890"
	//	]
	//}
}

type BalanceProblemDetail struct {
	*ProblemDetail
	Balance  int32    `json:"balance" xml:"balance"`
	Accounts []string `json:"accounts" xml:"accounts"`
}

func TestWriteJSONProblemDetail_WithExtension(t *testing.T) {
	data := BalanceProblemDetail{
		ProblemDetail: NewProblemDetail(
			"https://example.com/probs/out-of-credit",
			ProbOpts.Detail("Your current balance is 30, but that costs 50."),
			ProbOpts.Instance("/account/12345/abc"),
			ProbOpts.Title("You do not have enough credit."),
		),
		Balance:  30,
		Accounts: []string{"/account/12345", "/account/67890"},
	}

	rec := httptest.NewRecorder()
	err := WriteJSONProblemDetail(rec, &data, 403)
	expectTrue(t, err == nil)

	expRaw := `{"type":"https://example.com/probs/out-of-credit","title":"You do not have enough credit.","status":403,"detail":"Your current balance is 30, but that costs 50.","instance":"/account/12345/abc","balance":30,"accounts":["/account/12345","/account/67890"]}`
	gotRaw := strings.TrimSpace(rec.Body.String())

	expectTrue(t, gotRaw == expRaw)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == contentTypeJSONProblemDetail)
}

func TestWriteJSONProblemDetail_Untyped(t *testing.T) {
	data := NewProblemDetail(ProblemDetailUntyped, ProbOpts.ValidateLevel(ProbVLRequiredOnly))

	rec := httptest.NewRecorder()
	err := WriteJSONProblemDetail(rec, data, 403)
	expectTrue(t, err == nil)

	expRaw := `{"type":"about:blank","title":"Forbidden","status":403}`
	gotRaw := strings.TrimSpace(rec.Body.String())

	expectTrue(t, gotRaw == expRaw)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == contentTypeJSONProblemDetail)
}

func TestWriteJSONProblemDetail_Typed(t *testing.T) {
	data := NewProblemDetail(
		"https://example.com/probs/out-of-credit",
		ProbOpts.Detail("Your current balance is 30, but that costs 50."),
		ProbOpts.Instance("/account/12345/abc"),
		ProbOpts.Title("You do not have enough credit."),
	)

	rec := httptest.NewRecorder()
	err := WriteJSONProblemDetail(rec, data, 403)
	expectTrue(t, err == nil)

	expRaw := `{"type":"https://example.com/probs/out-of-credit","title":"You do not have enough credit.","status":403,"detail":"Your current balance is 30, but that costs 50.","instance":"/account/12345/abc"}`
	gotRaw := strings.TrimSpace(rec.Body.String())

	expectTrue(t, gotRaw == expRaw)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == contentTypeJSONProblemDetail)
}

func TestWriteJSONProblemDetail_TypedAllEmpty(t *testing.T) {
	data := NewProblemDetail("", ProbOpts.ValidateLevel(ProbVLStrict))

	rec := httptest.NewRecorder()
	err := WriteJSONProblemDetail(rec, data, 0)
	expectTrue(t, err != nil)
}

func TestWriteJSONProblemDetail_TypedAndStrictViolated(t *testing.T) {
	data := NewProblemDetail("--not-\n/a/valid/uri--",
		ProbOpts.ValidateLevel(ProbVLStrict),
		ProbOpts.Instance("\n-not/a/valid/path\n"),
		ProbOpts.Instance("\n-not/a/valid/path\n"),
	)
	rec := httptest.NewRecorder()
	err := WriteJSONProblemDetail(rec, data, 0)
	expectTrue(t, err != nil)
}

func TestWriteXMLProblemDetail_WithExtension(t *testing.T) {
	data := BalanceProblemDetail{
		ProblemDetail: NewProblemDetail(
			"https://example.com/probs/out-of-credit",
			ProbOpts.Detail("Your current balance is 30, but that costs 50."),
			ProbOpts.Instance("/account/12345/abc"),
			ProbOpts.Title("You do not have enough credit."),
		),
		Balance:  30,
		Accounts: []string{"/account/12345", "/account/67890"},
	}

	rec := httptest.NewRecorder()
	err := WriteXMLProblemDetail(rec, &data, 403)
	expectTrue(t, err == nil)

	rawExp := `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.com/probs/out-of-credit</type><title>You do not have enough credit.</title><status>403</status><detail>Your current balance is 30, but that costs 50.</detail><instance>/account/12345/abc</instance><balance>30</balance><accounts>/account/12345</accounts><accounts>/account/67890</accounts></problem>`
	rawGot := strings.TrimSpace(rec.Body.String())

	expectTrue(t, rawGot == rawExp)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == contentTypeXMLProblemDetail)
}

func TestWriteXMLProblemDetail_Untyped(t *testing.T) {
	data := NewProblemDetail(ProblemDetailUntyped, ProbOpts.ValidateLevel(ProbVLRequiredOnly))

	rec := httptest.NewRecorder()
	err := WriteXMLProblemDetail(rec, data, 403)
	expectTrue(t, err == nil)

	rawExp := `<problem xmlns="urn:ietf:rfc:7807"><type>about:blank</type><title>Forbidden</title><status>403</status></problem>`
	rawGot := strings.TrimSpace(rec.Body.String())

	expectTrue(t, rawGot == rawExp)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == contentTypeXMLProblemDetail)
}

func TestWriteXMLProblemDetail_Typed(t *testing.T) {
	data := NewProblemDetail(
		"https://example.com/probs/out-of-credit",
		ProbOpts.Detail("Your current balance is 30, but that costs 50."),
		ProbOpts.Instance("/account/12345/abc"),
		ProbOpts.Title("You do not have enough credit."),
	)

	rec := httptest.NewRecorder()
	err := WriteXMLProblemDetail(rec, data, 403)
	expectTrue(t, err == nil)

	rawExp := `<problem xmlns="urn:ietf:rfc:7807"><type>https://example.com/probs/out-of-credit</type><title>You do not have enough credit.</title><status>403</status><detail>Your current balance is 30, but that costs 50.</detail><instance>/account/12345/abc</instance></problem>`
	rawGot := strings.TrimSpace(rec.Body.String())

	expectTrue(t, rawGot == rawExp)
	expectTrue(t, rec.Code == 403)
	expectTrue(t, rec.Header().Get("Content-Type") == contentTypeXMLProblemDetail)
}

func TestWriteXMLProblemDetail_TypedAllEmpty(t *testing.T) {
	data := NewProblemDetail("", ProbOpts.ValidateLevel(ProbVLStrict))

	rec := httptest.NewRecorder()
	err := WriteXMLProblemDetail(rec, data, 0)
	expectTrue(t, err != nil)
}

func TestWriteXMLProblemDetail_TypedAndStrictViolated(t *testing.T) {
	data := NewProblemDetail("--not-\n/a/valid/uri--",
		ProbOpts.ValidateLevel(ProbVLStrict),
		ProbOpts.Instance("\n-not/a/valid/path\n"),
		ProbOpts.Instance("\n-not/a/valid/path\n"),
	)
	rec := httptest.NewRecorder()
	err := WriteXMLProblemDetail(rec, data, 0)
	expectTrue(t, err != nil)
}
