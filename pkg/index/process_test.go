package index

import (
	"archive/tar"
	"bytes"
	"os"
	"reflect"

	// "reflect"
	"sync"
	"testing"

	"github.com/DEliasVCruz/db-indexer/pkg/data"
	"github.com/DEliasVCruz/db-indexer/pkg/search"
)

func TestExtractFS(t *testing.T) {

	ch := make(chan *search.Data)
	wg := sync.WaitGroup{}

	index := Indexer{
		FileType:   "fs",
		dataFolder: os.DirFS("../../test/fixtures"),
	}

	t.Cleanup(func() {
		close(ch)
	})

	var tests = []struct {
		name string
		file string
		want *search.Data
	}{

		{
			"extract normal file",
			"normal_extract_data",
			&search.Data{
				ID:        "5860470.1075855667730",
				MessageID: "<5860470.1075855667730.JavaMail.evans@thyme>",
				Date:      "Thu, 5 Oct 2000 06:26:00 -0700 (PDT)",
				From:      "phillip.allen@enron.com",
				To:        "david.delainey@enron.com",
				Subject:   "Hello World",
				XFileName: "don baughman 6-25-02.PST",
				Contents:  "\nSome content\n\nWith some new lines\n\n",
				FilePath:  "normal_extract_data",
			},
		},

		{
			"extract with completely empty field",
			"empty_field_data",
			&search.Data{
				ID:        "5860470.1075855667730",
				MessageID: "<5860470.1075855667730.JavaMail.evans@thyme>",
				Date:      "Thu, 5 Oct 2000 06:26:00 -0700 (PDT)",
				From:      "phillip.allen@enron.com",
				To:        "david.delainey@enron.com",
				Subject:   "",
				XFileName: "don baughman 6-25-02.PST",
				Contents:  "\nSome content\n",
				FilePath:  "empty_field_data",
			},
		},

		{
			"extract multi new line field",
			"multi_new_line_field",
			&search.Data{
				ID:        "15722007.1075840335489",
				MessageID: "<15722007.1075840335489.JavaMail.evans@thyme>",
				Date:      "Thu, 13 Dec 2001 06:39:18 -0800 (PST)",
				From:      "don.baughman@enron.com",
				Subject:   "Call Laddie for house party: Mom &dad & Mary   Janice Nieghbour",
				XFileName: "don baughman 6-25-02.PST",
				Contents:  "\nContent\n",
				FilePath:  "multi_new_line_field",
			},
		},

		{
			"extract multi line field",
			"multi_line_field",
			&search.Data{
				ID:        "33534862.1075863219076",
				MessageID: "<33534862.1075863219076.JavaMail.evans@thyme>",
				Date:      "Mon, 26 Nov 2001 12:27:12 -0800 (PST)",
				From:      "craig.breslau@enron.com",
				To:        "susan.bailey@enron.com, credit <.williams@enron.com>, legal <.taylor@enron.com>",
				Subject:   "FW: assignment",
				XFileName: "SBAILE2 (Non-Privileged).pst",
				Contents:  "\n\nContent\n\n Some more content\n",
				FilePath:  "multi_line_field",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := index.dataFolder.Open(tt.file)
			if err != nil {
				t.Fatal(err.Error())
			}

			info, err := file.Stat()
			if err != nil {
				t.Fatal(err.Error())
			}

			wg.Add(1)
			go index.extract(&data.DataInfo{RelPath: tt.file, Size: int(info.Size())}, ch, &wg)
			got := <-ch

			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("got:\n%v,\nwanted:\n%v", got, tt.want)
			}
		})
	}

}

func TestExtractTar(t *testing.T) {

	ch := make(chan *search.Data)
	wg := sync.WaitGroup{}

	index := Indexer{
		FileType: "tar",
	}

	want := &search.Data{
		ID:        "5860470.1075855667730",
		MessageID: "<5860470.1075855667730.JavaMail.evans@thyme>",
		Date:      "Thu, 5 Oct 2000 06:26:00 -0700 (PDT)",
		From:      "phillip.allen@enron.com",
		To:        "david.delainey@enron.com",
		Subject:   "Hello World",
		XFileName: "don baughman 6-25-02.PST",
		Contents:  "\nSome content\n\nWith some new lines\n\n",
		FilePath:  "normal_extract_data",
	}

	t.Cleanup(func() {
		close(ch)
	})

	file, err := os.Open("../../test/fixtures/normal_extract_data.tar")
	if err != nil {
		t.Fatal(err.Error())
	}

	tr := tar.NewReader(file)

	header, err := tr.Next()
	if err != nil {
		t.Fatal(err.Error())
	}

	buf := bytes.NewBuffer(make([]byte, 0, header.Size))
	_, err = buf.ReadFrom(tr)

	wg.Add(1)
	go index.extract(
		&data.DataInfo{
			TarBuf: &data.TarBuf{Buffer: buf, Header: header},
			Err:    err,
			Size:   int(header.Size),
		},
		ch,
		&wg,
	)
	got := <-ch

	if !reflect.DeepEqual(want, got) {
		t.Errorf("got:\n%v,\nwanted:\n%v", got, want)
	}
}

func TestExtractMissingMetadataError(t *testing.T) {

	ch := make(chan *search.Data)
	var wg sync.WaitGroup
	var got *search.Data

	index := Indexer{
		FileType:   "fs",
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

// func TestProcess(t *testing.T) {

// 	processFields := map[string]string{
// 		"message_id":   "<15722007.1075840335489.JavaMail.evans@thyme>",
// 		"content_type": "text/plain; charset=us-ascii",
// 		"x_folder":     `\ExMerge - Baughman Jr., Don\Deleted Items`,
// 	}

// 	want := map[string]string{
// 		"message_id":   "<15722007.1075840335489.JavaMail.evans@thyme>",
// 		"_id":          "15722007.1075840335489",
// 		"content_type": "text/plain",
// 		"charset":      "us-ascii",
// 		"x_folder":     `/ExMerge - Baughman Jr., Don/Deleted Items`,
// 	}

// 	got := process(processFields)

// 	if !reflect.DeepEqual(got, want) {
// 		t.Error("got ", got, " wanted ", want)
// 	}
// }

func BenchmarkExtract(b *testing.B) {
	ch := make(chan *search.Data)
	var wg sync.WaitGroup

	index := Indexer{
		FileType:   "fs",
		dataFolder: os.DirFS("../../test/fixtures"),
	}
	file, err := index.dataFolder.Open("bench_doc")
	if err != nil {
		b.Fatal(err.Error())
	}

	info, err := file.Stat()
	if err != nil {
		b.Fatal(err.Error())
	}

	fs := &data.DataInfo{RelPath: "bench_doc", Size: int(info.Size())}

	b.Cleanup(func() {
		close(ch)
	})

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go index.extract(fs, ch, &wg)
		<-ch
	}

}
