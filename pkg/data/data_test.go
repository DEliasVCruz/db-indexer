package data

import (
	"reflect"
	"testing"
)

func TestExtrac(t *testing.T) {

	path := "../../test/fixtures/normal_extract_data"

	want := map[string][]byte{
		"message_id": []byte("<5860470.1075855667730.JavaMail.evans@thyme>"),
		"date":       []byte("Thu, 5 Oct 2000 06:26:00 -0700 (PDT)"),
		"from":       []byte("phillip.allen@enron.com"),
		"to":         []byte("david.delainey@enron.com"),
		"subject":    []byte("Hello World"),
		"x_filename": []byte("don baughman 6-25-02.PST"),
		"contents":   []byte("Some content\n\nWith some new lines\n\n"),
	}
	got, err := Extract(path)

	if err != nil {
		t.Fatalf("unexpected error opening test store: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}

}

func TestExtractMissingMetadataError(t *testing.T) {

	path := "../../test/fixtures/missing_metadata"

	_, err := Extract(path)

	if err == nil {
		t.Errorf("expected an error for file with missing metadata")
	}

}

func TestExtractMultiNewLineField(t *testing.T) {

	path := "../../test/fixtures/multi_new_line_field"

	want := map[string][]byte{
		"message_id": []byte("<15722007.1075840335489.JavaMail.evans@thyme>"),
		"date":       []byte("Thu, 13 Dec 2001 06:39:18 -0800 (PST)"),
		"from":       []byte("don.baughman@enron.com"),
		"subject":    []byte("Call Laddie for house party: Mom &dad & Mary   Janice Nieghbour"),
		"x_filename": []byte("don baughman 6-25-02.PST"),
		"contents":   []byte("Content\n"),
	}
	got, err := Extract(path)

	if err != nil {
		t.Fatalf("unexpected error opening test store: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}

}

func TestExtractMultiLineField(t *testing.T) {

	path := "../../test/fixtures/multi_line_field"

	want := map[string][]byte{
		"message_id": []byte("<33534862.1075863219076.JavaMail.evans@thyme>"),
		"date":       []byte("Mon, 26 Nov 2001 12:27:12 -0800 (PST)"),
		"from":       []byte("craig.breslau@enron.com"),
		"to":         []byte("susan.bailey@enron.com, credit <.williams@enron.com>, legal <.taylor@enron.com>"),
		"subject":    []byte("FW: assignment"),
		"x_filename": []byte("SBAILE2 (Non-Privileged).pst"),
		"contents":   []byte("\nContent\n\n Some more content\n"),
	}
	got, err := Extract(path)

	if err != nil {
		t.Fatalf("unexpected error opening test store: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}

}

func TestProcess(t *testing.T) {

	processFields := map[string][]byte{
		"message_id":   []byte("<15722007.1075840335489.JavaMail.evans@thyme>"),
		"content_type": []byte("text/plain; charset=us-ascii"),
		"x_folder":     []byte(`\ExMerge - Baughman Jr., Don\Deleted Items`),
	}

	want := map[string][]byte{
		"message_id":   []byte("15722007.1075840335489"),
		"content_type": []byte("text/plain"),
		"charset":      []byte("us-ascii"),
		"x_folder":     []byte(`/ExMerge - Baughman Jr., Don/Deleted Items`),
	}

	got := Process(processFields)

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}

}

func BenchmarkExtract(b *testing.B) {

	for i := 0; i < b.N; i++ {
		Extract("../../test/fixtures/bench_doc")
	}

}

func BenchmarkProcess(b *testing.B) {

	processFields := map[string][]byte{
		"message_id":   []byte("<15722007.1075840335489.JavaMail.evans@thyme>"),
		"content_type": []byte("text/plain; charset=us-ascii"),
		"x_folder":     []byte(`\ExMerge - Baughman Jr., Don\Deleted Items`),
	}

	for i := 0; i < b.N; i++ {
		Process(processFields)
	}

}
