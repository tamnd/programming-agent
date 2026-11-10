// SPDX-License-Identifier: Apache-2.0

package compat

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// ResultsFile is the file of the content repository where prog compat probe records its results.
const ResultsFile = "manifests/probes.jsonl"

// Result is the record of one probe for one implementation and one version_added (spec 09.4).
type Result struct {
	// Feature is the slug of the feature, and Probe is the path of the probe file in the node directory, such as probes/match_basic.py.
	Feature string `json:"feature"`
	Impl    string `json:"impl"`
	Probe   string `json:"probe"`
	Added   string `json:"version_added"`
	// Hash is the SHA-256 of the probe file and of its entry in probes.yaml. A result with an old hash does not count.
	Hash string `json:"hash"`
	Pass Run    `json:"pass"`
	// Fail is the run on the release before version_added. It is nil when the register has no release before.
	Fail *Run   `json:"fail,omitempty"`
	Date string `json:"date"`
}

// Run is one run of a probe.
type Run struct {
	Version string `json:"version"`
	Image   string `json:"image"`
	// OK tells if the run did what the check needs: exit 0 and the output for the pass run, a failure that matches the signature for the fail run.
	OK     bool   `json:"ok"`
	Exit   int    `json:"exit"`
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

// Key identifies a result.
func (r Result) Key() string { return Key(r.Feature, r.Impl, r.Probe, r.Added) }

// Key identifies the result of a probe of a feature for an implementation and a version_added.
func Key(feature, impl, probe, added string) string {
	return feature + " " + impl + " " + probe + " " + added
}

// Verified tells if the result proves the statement: the hash is current, the pass run is OK, and the fail run is OK when there is one.
func (r Result) Verified(hash string) bool {
	return r.Hash == hash && r.Pass.OK && (r.Fail == nil || r.Fail.OK)
}

// Results maps Key to the result.
type Results map[string]Result

// ReadResults reads manifests/probes.jsonl under root. A missing file gives no results.
func ReadResults(root string) (Results, error) {
	out := Results{}
	p := filepath.Join(root, filepath.FromSlash(ResultsFile))
	f, err := os.Open(p)
	if errors.Is(err, os.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 1<<20), 1<<24)
	line := 0
	for s.Scan() {
		line++
		if len(s.Bytes()) == 0 {
			continue
		}
		var r Result
		if err := json.Unmarshal(s.Bytes(), &r); err != nil {
			return nil, fmt.Errorf("%s:%d: %w", p, line, err)
		}
		out[r.Key()] = r
	}
	return out, s.Err()
}

// Write writes the results to manifests/probes.jsonl under root, sorted by key, one result on each line.
func (rs Results) Write(root string) error {
	keys := make([]string, 0, len(rs))
	for k := range rs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	p := filepath.Join(root, filepath.FromSlash(ResultsFile))
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	var buf bytes.Buffer
	for _, k := range keys {
		b, err := json.Marshal(rs[k])
		if err != nil {
			return err
		}
		buf.Write(b)
		buf.WriteByte('\n')
	}
	return os.WriteFile(p, buf.Bytes(), 0o644)
}

// ProbeHash is the hash of a probe: the bytes of the file and its expectation.
func ProbeHash(dir, probe string, e Expect) (string, error) {
	b, err := os.ReadFile(filepath.Join(dir, filepath.FromSlash(probe)))
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(b)
	fmt.Fprintf(h, "\x00%q\x00%q", e.Output, e.FailsWith)
	return hex.EncodeToString(h.Sum(nil)), nil
}
