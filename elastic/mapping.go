package elastic

const CommentsIndexMapping = `
{
    "mappings" : {
      "properties" : {
        "comment" : {
          "type" : "text"
        },
        "post_id" : {
          "type" : "keyword"
        }
      }
    }
  }
`

const FacesIndexMapping = `
{
	"mappings" : {
			"properties" : {
				"encoding_vector": {
					"type": "binary",
					"doc_values": true
				},
				"post_id": {
					"type": "integer"
				},
				"x": {
					"type": "integer"
				},
				"y": {
					"type": "integer"
				},
				"width": {
					"type": "integer"
				},
				"height":{
					"type": "integer"
				}
			}
	}
}
`

const PostsIndexMapping = `
	{
    "mappings" : {
      "properties" : {
        "caption" : {
          "type" : "text"
        },
        "user_id" : {
          "type" : "keyword"
        }
      }
    }
}
`

const UsersIndexMapping = `
{
    "mappings" : {
		"properties" : {
			"id": {
				"type": "integer"
			}
			"user_name": {
				"type": "text"
			}
			"real_name": {
				"type": "text"
			}
			"bio": {
				"type": "text"
			}
		}
	}
}
`
