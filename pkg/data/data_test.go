package data

import (
	"reflect"
	"sync"
	"testing"
)

func TestExtract(t *testing.T) {

	ch := make(chan map[string]string)
	wg := sync.WaitGroup{}

	t.Cleanup(func() {
		close(ch)
	})

	path := "../../test/fixtures/normal_extract_data"

	want := map[string]string{
		"message_id": "<5860470.1075855667730.JavaMail.evans@thyme>",
		"date":       "Thu, 5 Oct 2000 06:26:00 -0700 (PDT)",
		"from":       "phillip.allen@enron.com",
		"to":         "david.delainey@enron.com",
		"subject":    "Hello World",
		"x_filename": "don baughman 6-25-02.PST",
		"contents":   "\nSome content\n\nWith some new lines\n\n",
		"file_path":  path,
	}

	wg.Add(1)
	go Extract(path, ch, &wg)
	got := <-ch

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}

}

func TestExtractWithCompletelyEmptyField(t *testing.T) {

	ch := make(chan map[string]string)
	var wg sync.WaitGroup

	t.Cleanup(func() {
		close(ch)
	})

	path := "../../test/fixtures/empty_field_data"

	want := map[string]string{
		"message_id": "<5860470.1075855667730.JavaMail.evans@thyme>",
		"date":       "Thu, 5 Oct 2000 06:26:00 -0700 (PDT)",
		"from":       "phillip.allen@enron.com",
		"to":         "david.delainey@enron.com",
		"subject":    "",
		"x_filename": "don baughman 6-25-02.PST",
		"contents":   "\nSome content\n",
		"file_path":  path,
	}

	wg.Add(1)
	go Extract(path, ch, &wg)
	got := <-ch

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}
}

func TestExtractMissingMetadataError(t *testing.T) {

	ch := make(chan map[string]string)
	var wg sync.WaitGroup
	var got map[string]string

	t.Cleanup(func() {
		close(ch)
	})

	path := "../../test/fixtures/missing_metadata"

	wg.Add(1)
	Extract(path, ch, &wg)

	select {
	case got = <-ch:
		t.Error("Expected empty map, got ", got)
	default:
	}

}

func TestExtractMultiNewLineField(t *testing.T) {

	ch := make(chan map[string]string)
	var wg sync.WaitGroup

	t.Cleanup(func() {
		close(ch)
	})

	path := "../../test/fixtures/multi_new_line_field"

	want := map[string]string{
		"message_id": "<15722007.1075840335489.JavaMail.evans@thyme>",
		"date":       "Thu, 13 Dec 2001 06:39:18 -0800 (PST)",
		"from":       "don.baughman@enron.com",
		"subject":    "Call Laddie for house party: Mom &dad & Mary   Janice Nieghbour",
		"x_filename": "don baughman 6-25-02.PST",
		"contents":   "\nContent\n",
		"file_path":  path,
	}

	wg.Add(1)
	go Extract(path, ch, &wg)
	got := <-ch

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}

}

func TestExtractMultiLineField(t *testing.T) {

	ch := make(chan map[string]string)
	var wg sync.WaitGroup

	t.Cleanup(func() {
		close(ch)
	})

	path := "../../test/fixtures/multi_line_field"

	want := map[string]string{
		"message_id": "<33534862.1075863219076.JavaMail.evans@thyme>",
		"date":       "Mon, 26 Nov 2001 12:27:12 -0800 (PST)",
		"from":       "craig.breslau@enron.com",
		"to":         "susan.bailey@enron.com, credit <.williams@enron.com>, legal <.taylor@enron.com>",
		"subject":    "FW: assignment",
		"x_filename": "SBAILE2 (Non-Privileged).pst",
		"contents":   "\n\nContent\n\n Some more content\n",
		"file_path":  path,
	}

	wg.Add(1)
	go Extract(path, ch, &wg)
	got := <-ch

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}

}

func TestProcess(t *testing.T) {

	ch := make(chan map[string]string)
	var wg sync.WaitGroup

	t.Cleanup(func() {
		close(ch)
	})

	processFields := map[string]string{
		"message_id":   "<15722007.1075840335489.JavaMail.evans@thyme>",
		"content_type": "text/plain; charset=us-ascii",
		"x_folder":     `\ExMerge - Baughman Jr., Don\Deleted Items`,
	}

	want := map[string]string{
		"message_id":   "<15722007.1075840335489.JavaMail.evans@thyme>",
		"_id":          "15722007.1075840335489",
		"content_type": "text/plain",
		"charset":      "us-ascii",
		"x_folder":     `/ExMerge - Baughman Jr., Don/Deleted Items`,
	}

	wg.Add(1)
	go Process(processFields, ch, &wg)
	got := <-ch

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}
}

func BenchmarkExtract(b *testing.B) {
	ch := make(chan map[string]string)
	var wg sync.WaitGroup

	b.Cleanup(func() {
		close(ch)
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go Extract("../../test/fixtures/bench_doc", ch, &wg)
		<-ch
	}

}

func BenchmarkProcess(b *testing.B) {

	ch := make(chan map[string]string)
	var wg sync.WaitGroup

	b.Cleanup(func() {
		close(ch)
	})

	processFields := map[string]string{
		"message_id":   "<15722007.1075840335489.JavaMail.evans@thyme>",
		"content_type": "text/plain; charset=us-ascii",
		"x_folder":     `\ExMerge - Baughman Jr., Don\Deleted Items`,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go Process(processFields, ch, &wg)
		<-ch
	}

}
