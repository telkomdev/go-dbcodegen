# Json Schema Spec
```
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    },
    "fields": {
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            },
            "type": {
              "type": "string"
            },
            "options": {
              "type": "array",
              "items": [
                {
                  "type": "string"
                }
              ]
            }
          },
          "required": [
            "name",
            "type",
            "options"
          ]
        }
      ]
    },
    "indexes": {
      "type": "array",
      "items": [
        {
          "type": "object",
          "properties": {
            "name": {
              "type": "string"
            },
            "fields": {
              "type": "array",
              "items": [
                {
                  "type": "object",
                  "properties": {
                    "column": {
                      "type": "string"
                    },
                    "order": {
                      "type": "string"
                    }
                  },
                  "required": [
                    "column"
                  ]
                }
              ]
            },
            "unique": {
              "type": "boolean"
            }
          },
          "required": [
            "name",
            "fields",
            "unique"
          ]
        }
      ]
    }
  },
  "required": [
    "name",
    "fields",
    "indexes"
  ]
}
```