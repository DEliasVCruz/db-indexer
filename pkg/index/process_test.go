package index

import (
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/DEliasVCruz/db-indexer/pkg/data"
)

func TestExtractFS(t *testing.T) {

	ch := make(chan map[string]string)
	wg := sync.WaitGroup{}

	index := Indexer{
		FileType:   "folder",
		dataFolder: os.DirFS("../../test/fixtures"),
	}

	t.Cleanup(func() {
		close(ch)
	})

	var tests = []struct {
		name string
		file string
		want map[string]string
	}{

		{
			"extract normal file",
			"normal_extract_data",
			map[string]string{
				"_id":        "5860470.1075855667730",
				"message_id": "<5860470.1075855667730.JavaMail.evans@thyme>",
				"date":       "Thu, 5 Oct 2000 06:26:00 -0700 (PDT)",
				"from":       "phillip.allen@enron.com",
				"to":         "david.delainey@enron.com",
				"subject":    "Hello World",
				"x_filename": "don baughman 6-25-02.PST",
				"contents":   "\nSome content\n\nWith some new lines\n\n",
				"file_path":  "normal_extract_data",
			},
		},

		{
			"extract with completely empty field",
			"empty_field_data",
			map[string]string{
				"_id":        "5860470.1075855667730",
				"message_id": "<5860470.1075855667730.JavaMail.evans@thyme>",
				"date":       "Thu, 5 Oct 2000 06:26:00 -0700 (PDT)",
				"from":       "phillip.allen@enron.com",
				"to":         "david.delainey@enron.com",
				"subject":    "",
				"x_filename": "don baughman 6-25-02.PST",
				"contents":   "\nSome content\n",
				"file_path":  "empty_field_data",
			},
		},

		{
			"extract multi new line field",
			"multi_new_line_field",
			map[string]string{
				"_id":        "15722007.1075840335489",
				"message_id": "<15722007.1075840335489.JavaMail.evans@thyme>",
				"date":       "Thu, 13 Dec 2001 06:39:18 -0800 (PST)",
				"from":       "don.baughman@enron.com",
				"subject":    "Call Laddie for house party: Mom &dad & Mary   Janice Nieghbour",
				"x_filename": "don baughman 6-25-02.PST",
				"contents":   "\nContent\n",
				"file_path":  "multi_new_line_field",
			},
		},

		{
			"extract multi line field",
			"multi_line_field",
			map[string]string{
				"_id":        "33534862.1075863219076",
				"message_id": "<33534862.1075863219076.JavaMail.evans@thyme>",
				"date":       "Mon, 26 Nov 2001 12:27:12 -0800 (PST)",
				"from":       "craig.breslau@enron.com",
				"to":         "susan.bailey@enron.com, credit <.williams@enron.com>, legal <.taylor@enron.com>",
				"subject":    "FW: assignment",
				"x_filename": "SBAILE2 (Non-Privileged).pst",
				"contents":   "\n\nContent\n\n Some more content\n",
				"file_path":  "multi_line_field",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wg.Add(1)
			go index.extract(&data.DataInfo{RelPath: tt.file}, ch, &wg)
			got := <-ch

			if !reflect.DeepEqual(got, tt.want) {
				t.Error("got ", got, " wanted ", tt.want)
			}
		})
	}

}

func TestExtractMissingMetadataError(t *testing.T) {

	ch := make(chan map[string]string)
	var wg sync.WaitGroup
	var got map[string]string

	index := Indexer{
		FileType:   "folder",
		dataFolder: os.DirFS("../../test/fixtures"),
	}
	file := &data.DataInfo{RelPath: "missing_metadata"}

	t.Cleanup(func() {
		close(ch)
	})

	wg.Add(1)
	index.extract(file, ch, &wg)

	select {
	case got = <-ch:
		t.Error("Expected empty map, got ", got)
	default:
	}

}

func TestProcess(t *testing.T) {

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

	got := process(processFields)

	if !reflect.DeepEqual(got, want) {
		t.Error("got ", got, " wanted ", want)
	}
}

func BenchmarkExtract(b *testing.B) {
	ch := make(chan map[string]string)
	var wg sync.WaitGroup

	index := Indexer{FileType: "folder"}
	file := &data.DataInfo{RelPath: "../../test/fixtures/multi_line_field"}

	b.Cleanup(func() {
		close(ch)
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go index.extract(file, ch, &wg)
		<-ch
	}

}

func BenchmarkProcess(b *testing.B) {

	processFields := map[string]string{
		"message_id":   "<15722007.1075840335489.JavaMail.evans@thyme>",
		"content_type": "text/plain; charset=us-ascii",
		"x_folder":     `\ExMerge - Baughman Jr., Don\Deleted Items`,
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		process(processFields)
	}

}
