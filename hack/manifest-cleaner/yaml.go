/*
Copyright 2022 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"fmt"
	"io"
	"regexp"

	"github.com/dprotaso/go-yit"
	"gopkg.in/yaml.v3"
)

// yamlNodes decodes the yaml.Nodes from the given input.
func yamlNodes(data []byte) ([]*yaml.Node, error) {
	var docs []*yaml.Node

	decoder := yaml.NewDecoder(bytes.NewReader(data))

	// Read all YAML documents contained in the input data.
	for i := 0; ; i++ {
		var doc yaml.Node
		if err := decoder.Decode(&doc); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("decoding YAML document: %w", err)
		}

		// perform sanity checks, even though Decode() shouldn't read
		// anything but YAML documents
		if doc.Kind != yaml.DocumentNode {
			return nil, fmt.Errorf("decoded YAML node %d has kind %s, expected only documents",
				i, doc.Tag)
		}
		if nNodes := len(doc.Content); nNodes != 1 {
			return nil, fmt.Errorf("decoded YAML node %d contains %d nodes instead of 1",
				i, nNodes)
		}

		docs = append(docs, &doc)
	}

	return docs, nil
}

var copyrightHeader = regexp.MustCompile(`^# Copyright 20[0-9]{2}`)

// findCopyrightHeader returns the first found document-level comment that
// represents a copyright header in the given list of YAML nodes.
func findCopyrightHeader(docs []*yaml.Node) string {
	it := yit.FromNodes(docs...)

	for node, ok := it(); ok; node, ok = it() {
		// Comments located at the top of a file are parsed as head comment.
		if copyrightHeader.MatchString(node.HeadComment) {
			return node.HeadComment
		}
		// Comments located directly after a "---" marker are parsed as foot comment.
		if copyrightHeader.MatchString(node.FootComment) {
			return node.FootComment
		}
	}

	return ""
}

// filterOutCopyrightHeaders returns the given list of YAML nodes with all
// document-level copyright headers filtered out.
func filterOutCopyrightHeaders(docs []*yaml.Node) []*yaml.Node {
	it := yit.FromNodes(docs...)

	for node, ok := it(); ok; node, ok = it() {
		// Comments located at the top of a file are parsed as head comment.
		if copyrightHeader.MatchString(node.HeadComment) {
			node.HeadComment = ""
		}

		// Comments located directly after a "---" marker are parsed as foot comment.
		if copyrightHeader.MatchString(node.FootComment) {
			node.FootComment = ""
		}
	}

	return docs
}

// filterOutNullDocs returns the given list of YAML nodes with all entries
// representing a null document filtered out.
// A null document is either empty or contains a single node of type "!!null".
func filterOutNullDocs(docs []*yaml.Node) []*yaml.Node {
	filteredDocs := docs[:0]

	it := yit.FromNodes(docs...).Filter(yit.Negate(nullDocument))

	for node, ok := it(); ok; node, ok = it() {
		filteredDocs = append(filteredDocs, node)
	}

	return filteredDocs
}

// nullDocument is a yit.Predicate which returns whether the given node
// represents a null document.
func nullDocument(n *yaml.Node) bool {
	if n.Kind != yaml.DocumentNode || len(n.Content) > 1 {
		return false
	}
	return len(n.Content) == 0 || n.Content[0].Tag == "!!null"
}
