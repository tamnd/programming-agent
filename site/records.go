// SPDX-License-Identifier: Apache-2.0

package site

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
)

// The JSON records of the cards are in record packs (spec 16.13). A pack is a file of JSON lines: each line is the JSON array of the records of one collection. The index /api/records.json gives, for the URL of each collection, the number of its pack, the byte offset of its line and the length of the line without the newline. The name of a pack has the hash of its bytes, so a browser can keep a pack for a year.
const (
	recordIndex = "/api/records.json"
	recordDir   = "/api/records/"
	recordPage  = "/api/record/"
	// packTarget is the size at which a pack closes. A collection that is larger than this has a pack of its own. The host takes 25 MiB for one file, and the browser gets a whole pack when the host does not answer a range request, so a pack is small.
	packTarget = 2 << 20
)

// record is the JSON array of the records of the collection at URL.
type record struct {
	URL  string
	Data []byte
}

// RecordIndex is the file /api/records.json.
type RecordIndex struct {
	Packs   []string          `json:"packs"`
	Records map[string][3]int `json:"records"` // collection URL: pack, offset, length
}

// packRecords writes the records of each collection, in the order of recs, into the packs, and writes the index and the page that shows one record.
func (b *builder) packRecords(recs []record) error {
	idx := RecordIndex{Packs: []string{}, Records: map[string][3]int{}}
	var buf bytes.Buffer
	flush := func() error {
		if buf.Len() == 0 {
			return nil
		}
		sum := sha256.Sum256(buf.Bytes())
		name := fmt.Sprintf("%s%03d.%s.ndjson", recordDir, len(idx.Packs), hex.EncodeToString(sum[:])[:12])
		if err := b.write(name, bytes.Clone(buf.Bytes())); err != nil {
			return err
		}
		idx.Packs = append(idx.Packs, name)
		buf.Reset()
		return nil
	}
	for _, r := range recs {
		if buf.Len() > 0 && buf.Len()+len(r.Data)+1 > packTarget {
			if err := flush(); err != nil {
				return err
			}
		}
		idx.Records[r.URL] = [3]int{len(idx.Packs), buf.Len(), len(r.Data)}
		buf.Write(r.Data)
		buf.WriteByte('\n')
	}
	if err := flush(); err != nil {
		return err
	}
	data, err := json.Marshal(idx)
	if err != nil {
		return err
	}
	if err := b.write(recordIndex, data); err != nil {
		return err
	}
	p := DocPage{Base: b.base("JSON record", "The JSON records of the cards of one collection.", recordPage), Heading: "JSON record", Lead: "The JSON records of the cards of one collection, as the site keeps them in its record packs."}
	p.NoIndex, p.BodyKind = true, "record"
	p.Body = template.HTML(`<p data-record-status>This page needs script. The records of all the collections are in the packs that <a href="` + recordIndex + `">` + recordIndex + `</a> lists. Each line of a pack is the JSON array of one collection, and the index gives the pack, the byte offset and the length of each line.</p>
<p data-record-links hidden><a data-record-raw href="">Pack file</a> · <a data-record-back href="">Back to the collection</a></p>
<pre class="record" data-record hidden></pre>`)
	return b.page("doc", recordPage, p)
}
