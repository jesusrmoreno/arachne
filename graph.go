package main

import (
	"io/fs"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/gernest/front"
)

// Edge represents a relationship between two nodes
type Edge struct {
	id      string
	Source  string `json:"source"`
	Target  string `json:"target"`
	Context string `json:"context"`
	Label   string `json:"label"`
	FoundIn string `json:"foundIn"`
}

// Node represents a discrete "thing" derived from notes
type Node struct {
	Name       string      `json:"name"`
	Properties interface{} `json:"properties"`
	FileEntry  FileEntry   `json:"fileEntry"`
}

// Graph ...
type Graph struct {
	Nodes      []Node                 `json:"nodes"`
	Edges      []Edge                 `json:"edges"`
	Keys       map[string]bool        `json:"-"`
	Properties map[string]interface{} `json:"properties"`
	Content    map[string]string      `json:"content"`
	Base       string                 `json:"base"`
}

func NewGraph() Graph {
	return Graph{
		Keys:       map[string]bool{},
		Nodes:      []Node{},
		Properties: map[string]interface{}{},
		Edges:      []Edge{},
		Content:    map[string]string{},
	}
}

func (g *Graph) AddNode(nodes ...Node) {
	newNodes := []Node{}

	for _, node := range nodes {
		exists := g.Keys[node.Name]
		if !exists {
			g.Content[node.Name] = ""
			g.Keys[node.Name] = true
			newNodes = append(newNodes, Node{
				Name:       node.Name,
				FileEntry:  node.FileEntry,
				Properties: g.Properties[node.Name],
			})
		}
	}
	g.Nodes = append(g.Nodes, newNodes...)
}

func (g *Graph) AddEdge(e Edge) {
	id := e.Source + e.Label + e.Target + e.FoundIn + e.Context
	exists := g.Keys[id]
	g.AddNode(Node{
		Name: e.Source,
	}, Node{
		Name: e.Target,
	})
	if !exists {
		g.Keys[e.id] = true
		e.id = id
		g.Edges = append(g.Edges, e)
	}
}

// FileEntry holds the metadata that we need to create the node
type FileEntry struct {
	Path      string `json:"path"`
	FileName  string `json:"fileName"`
	Directory string `json:"directory"`
	UpdatedAt int64  `json:"updatedAt"`
	Root      string `json:"root"`
}

func hasRelationshipString(text string) bool {
	return strings.Contains(text, "[[") && strings.Contains(text, "]]")
}

func normalizeRelationshipString(def string, entry FileEntry) (string, bool) {
	characters := strings.Split(def, "")
	if len(characters) == 0 {
		return "", true
	}

	isAppend := characters[len(characters)-1] == ":"
	isPrepend := characters[0] == ":" || !isAppend

	isValid := isPrepend && !isAppend || !isPrepend && !isAppend || !isPrepend && isAppend

	if !isValid {
		return "", true
	}

	normalized := ""
	rawParts := strings.Split(def, ":")
	parts := []string{}
	for _, part := range rawParts {
		if part != "" {
			parts = append(parts, part)
		}
	}

	if len(parts) == 1 {
		return entry.normalizedName() + ":" + "" + ":" + parts[0], false
	}

	if len(parts) == 2 {
		if isPrepend {
			relationshipStringParts := append([]string{entry.normalizedName()}, parts...)
			normalized = strings.Join(relationshipStringParts, ":")
		} else {
			relationshipStringParts := append(parts, entry.normalizedName())
			normalized = strings.Join(relationshipStringParts, ":")
		}
		return normalized, false
	}

	if len(parts) == 3 {
		return strings.Join(parts, ":"), false
	}

	return "", true
}

func makeRelationships(g *Graph, chunk string, entry FileEntry, fileLookup map[string]bool) {
	// We have a valid section
	if hasRelationshipString(chunk) {
		state := "empty"
		section := ""
		rels := []string{}

		for _, r := range chunk {
			if state == "empty" {
				section = ""
			}
			char := string(r)
			if char == "[" {
				section += char
				if state == "empty" {
					state = "opened"
				}
				continue
			}
			if char == "]" {
				section += char
				if state == "opened" {
					state = "closed"
				} else if state == "closed" {
					rels = append(rels, section)
					state = "empty"
				}
				continue
			}
			if state == "opened" {
				section += char
			}
		}

		for _, token := range rels {
			if hasRelationshipString(token) {
				start := strings.Index(token, "[[")
				end := strings.Index(token, "]]")
				relationshipString := string([]rune(token)[start+2 : end])

				normalized, err := normalizeRelationshipString(relationshipString, entry)
				if err {
					continue
				}
				parts := strings.Split(normalized, ":")
				source := parts[0]
				label := parts[1]
				target := parts[2]

				g.AddEdge(Edge{
					FoundIn: entry.Path,
					Source:  source,
					Label:   label,
					Target:  target,
					Context: chunk,
				})
			}
		}
	}
}

func makeTagsRelationships(g *Graph, chunk string, entry FileEntry, fileLookup map[string]bool) {

}

func makeFrontmatter(g *Graph, chunk string, entry FileEntry) interface{} {
	m := front.NewMatter()
	m.Handle("---", front.YAMLHandler)
	f, _, err := m.Parse(strings.NewReader(chunk))
	if err != nil {
		return make(map[string]interface{})
	}
	return f
}

func initFrontMatter(g *Graph, entry FileEntry) {
	contentBytes, err := ioutil.ReadFile(entry.Path)
	if err != nil {
		log.Fatal(err)
	}
	contentString := string(contentBytes)
	frontMatter := makeFrontmatter(g, contentString, entry)

	// Set a lookup for convenience in the top level graph obj
	g.Properties[entry.normalizedName()] = frontMatter
}

func ParseFile(g *Graph, entry FileEntry, fileLookup map[string]bool) {
	contentBytes, err := ioutil.ReadFile(entry.Path)
	if err != nil {
		log.Fatal(err)
	}
	contentString := string(contentBytes)

	// set the properties on the node
	g.AddNode(Node{
		Name:      entry.normalizedName(),
		FileEntry: entry,
	})

	for _, chunk := range strings.Split(contentString, "\n") {
		makeRelationships(g, chunk, entry, fileLookup)
		makeTagsRelationships(g, chunk, entry, fileLookup)
	}
}

func (f *FileEntry) normalizedName() string {
	p := strings.Replace(f.Path, path.Clean(f.Root)+"/", "", 1)
	return f.normalizedPathRoot(p)
}

func (f *FileEntry) normalizedPathRoot(fileName string) string {
	return path.Join(strings.Replace(fileName, ".md", "", 1))
}

func DeriveGraphStructure(entries []FileEntry, output chan Graph, rootFolder string) {
	g := NewGraph()
	fileLookup := map[string]bool{}

	for _, entry := range entries {
		fileLookup[entry.normalizedName()] = true
		initFrontMatter(&g, entry)
	}

	for _, entry := range entries {
		ParseFile(&g, entry, fileLookup)
		content, err := ioutil.ReadFile(entry.Path)
		if err != nil {
			log.Fatal(err)
		}
		g.Content[entry.normalizedName()] = string(content)
		g.Base = rootFolder
	}
	output <- g
}

func StartGraphStructureService(rootFolder string, output chan Graph) {
	tracker := make(map[string]int64)
	tick := time.Tick(time.Millisecond * 150)
	for range tick {
		dirty := false
		entries := []FileEntry{}
		previousFileCount := len(tracker)
		currentCount := 0
		if err := filepath.WalkDir(rootFolder, func(path string, d fs.DirEntry, err error) error {
			// We only want to check markdown files for link content
			if !d.IsDir() && filepath.Ext(path) == ".md" {
				currentCount += 1
				info, err := d.Info()
				if err != nil {
					log.Fatal(err)
				}
				updatedAt := tracker[path]
				if updatedAt != info.ModTime().Unix() {
					tracker[path] = info.ModTime().Unix()
					dirty = true
				}

				if err != nil {
					log.Fatal(err)
				}
				entries = append(entries, FileEntry{
					FileName:  d.Name(),
					Path:      path,
					Root:      rootFolder,
					Directory: filepath.Dir(path),
					UpdatedAt: info.ModTime().Unix(),
				})
			}
			return nil
		}); err != nil {
			log.Fatal(err)
		}
		hasDeletedFiles := previousFileCount != currentCount
		if hasDeletedFiles {
			dirty = true
			tracker = make(map[string]int64)
			for _, entry := range entries {
				tracker[entry.Path] = entry.UpdatedAt
			}
		}

		if dirty {
			DeriveGraphStructure(entries, output, rootFolder)
		}
	}
}
