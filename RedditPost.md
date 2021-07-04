**TL:DR: Arachne traverses a directory looking for .md files and creates a structure to represent the relationships it finds inside those files, then serves it over Websockets and REST so you can build cool things.**

https://github.com/jesusrmoreno/arachne

Hey guys, I recently took Obsidian out for a spin and it _almost_ hit everything I've wanted from an editor, but I kept thinking "man what if I could just treat my notes as an API and build whatever I want on top of it.."

I had nothing to work on this weekend so Arachne was born. 

Long story short it turns

    # What is Arachne?
    > Arachne, (Greek: “Spider”) in [[greek/mythology]], the [[Arachne:daughter of:Idmon of Colophon]] in Lydia, a dyer in purple. Arachne was a weaver who acquired such skill in her art that she ventured to challenge Athena, goddess of war, handicraft, and practical reason.

into

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
            ...
          }
        },
        {
          "name": "Arachne",
          "properties": null,
          "fileEntry": {
            ...
          }
        },
        {
          "name": "Idmon of Colophon",
          "properties": null,
          "fileEntry": {
            ...
          }
        }
      ],
      "edges": [
        ...
        {
          "source": "Arachne",
          "target": "Idmon of Colophon",
          "label": "daughter of",
          
          "context": "> Arachne, (Greek: “Spider”) in [[greek/mythology]], the [[Arachne:daughter of:Idmon of Colophon]] in Lydia, a dyer in purple. Arachne was a weaver who acquired such skill in her art that she ventured to challenge Athena, goddess of war, handicraft, and practical reason.",
          
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
