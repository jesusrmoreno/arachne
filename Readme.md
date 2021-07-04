# What is Arachne?
> Arachne, (Greek: “Spider”) in [[greek/mythology]], the [[Arachne:daughter of:Idmon of Colophon]] in Lydia, a dyer in purple. Arachne was a weaver who acquired such skill in her art that she ventured to challenge Athena, goddess of war, handicraft, and practical reason.


Arachne is a simple application written in Go that traverses a directory and parses markdown files looking for relationship syntax, and creates a graph structure from what it finds. 

It also starts a websocket server and simple endpoint to access the structure.

**Takes this syntax**

```
// Will generate a relationship between the current page and "target"
[[target]]

// Will generate a relationship between the "source" and "target", with a label of "action"
ex: [[Sally:Is Studying:Computer Science]]
[[source:action:target]]

// Will generate a relationship between the current page and "target" with a label of "action"
[[:action:target]]

// Will generate a relationship between "source" and the current page with a label of "action"
[[source:action:]]
```

**And transforms it into this shape** 

```

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
	Content map[string]string `json:"content"`
	Base string `json:"base"`
}
```

**And serves it over (REST)**

```
GET http:localhost:PORT/graph
```

**And serves it over (Websockets)**

```
// javscript example
let socket = new WebSocket(`ws://${window.location.hostname}:8080/ws`);

socket.onmessage = evt => {
  try {
    const graph = JSON.parse(evt.data);
	console.log(graph) // your graph data
  } catch (e) {
    console.warn(e);
  }
};

```

## So what does this enable?
This is your notes and the relationships you define in those notes as a real time API, do whatever you want with it. Build awesome things. Don't get tied down to any particular editor.

## Why
I needed a weekend project.

Because the whole point of using markdown for your files is that it's supposed to be open and platform independent. 

As a developer I hate having to figure out how to work with someone else's plugin system for quick extensions / functionality additions,



## Example 
```
---
name: Arachne
description: Markdown Powered Knowledge Base API
created: July 2021
---

# What is Arachne?
> Arachne, (Greek: “Spider”) in [[greek/mythology]], the [[Arachne:daughter of:Idmon of Colophon]] in Lydia, a dyer in purple. Arachne was a weaver who acquired such skill in her art that she ventured to challenge Athena, goddess of war, handicraft, and practical reason.
```

```
{
  "nodes": [
    {
      "name": "main",
      "properties": {
        "created": "July 2021",
        "description": "Markdown Powered Knowledge Base API",
        "name": "Arachne"
      },
      "fileEntry": {
        "path": "notes/main.md",
        "fileName": "main.md",
        "directory": "notes",
        "updatedAt": 1625440392,
        "root": "./notes"
      }
    },
    {
      "name": "greek/mythology",
      "properties": null,
      "fileEntry": {
        "path": "",
        "fileName": "",
        "directory": "",
        "updatedAt": 0,
        "root": ""
      }
    },
    {
      "name": "Arachne",
      "properties": null,
      "fileEntry": {
        "path": "",
        "fileName": "",
        "directory": "",
        "updatedAt": 0,
        "root": ""
      }
    },
    {
      "name": "Idmon of Colophon",
      "properties": null,
      "fileEntry": {
        "path": "",
        "fileName": "",
        "directory": "",
        "updatedAt": 0,
        "root": ""
      }
    }
  ],
  "edges": [
    {
      "source": "main",
      "target": "greek/mythology",
      "context": "> Arachne, (Greek: “Spider”) in [[greek/mythology]], the [[Arachne:daughter of:Idmon of Colophon]] in Lydia, a dyer in purple. Arachne was a weaver who acquired such skill in her art that she ventured to challenge Athena, goddess of war, handicraft, and practical reason.",
      "label": "",
      "foundIn": "notes/main.md"
    },
    {
      "source": "Arachne",
      "target": "Idmon of Colophon",
      "context": "> Arachne, (Greek: “Spider”) in [[greek/mythology]], the [[Arachne:daughter of:Idmon of Colophon]] in Lydia, a dyer in purple. Arachne was a weaver who acquired such skill in her art that she ventured to challenge Athena, goddess of war, handicraft, and practical reason.",
      "label": "daughter of",
      "foundIn": "notes/main.md"
    }
  ],
  "properties": {
    "main": {
      "created": "July 2021",
      "description": "Markdown Powered Knowledge Base API",
      "name": "Arachne"
    }
  },
  "content": {
    "Arachne": "",
    "Idmon of Colophon": "",
    "greek/mythology": "",
    "main": "---\nname: Arachne\ndescription: Markdown Powered Knowledge Base API\ncreated: July 2021\n---\n\n# What is Arachne?\n> Arachne, (Greek: “Spider”) in [[greek/mythology]], the [[Arachne:daughter of:Idmon of Colophon]] in Lydia, a dyer in purple. Arachne was a weaver who acquired such skill in her art that she ventured to challenge Athena, goddess of war, handicraft, and practical reason.\n"
  },
  "base": "./notes"
}
```

You'll notice that the above JSON file has some empty fileEntry properties. These are known as "virtual" nodes, and they're nodes that are referenced in notes but have no actual markdown file backing them.
 
## Usage
There are three executable files in the /bin folder in this repo.

```
macOs = arachne-amd64-darwin
windows = arachne-386.exe
linux = arachne-amd64-linux
```

Create a config.toml file (must be a sibling to the executable) that looks like the following:

```
port = "8080"
root = "./notes"
```

Then in your terminal run (again make sure you're in the same directory as the config.toml file).

```
./arachne-amd64-darwin
```

**From any browser you can now visit**

```
// Replace port with the port defined in the config file
// default: http://localhost:8080/graph
http://localhost:<PORT>/graph
```

And you should see a graph representation of your notes.

# TODO
- [ ] Better docs
- [ ] Query API
- [ ] User friendly UI for config updates
- [ ] Build cool stuff on top of this