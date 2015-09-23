package client

import (
	"fmt"
)

//Actionable failure conditions, covered in detail in their relevant sections,
//are reported as part of 4xx responses, in a json response body
type Errors struct {
	Errors []Error `json:errors`
}

type Error struct {
	//The code field will be a unique identifier, all caps with underscores by
	// convention.
	code string `json:code`
	//The message field will be a human readable string
	message string `json:message`
	//The optional detail field may contain arbitrary json data providing
	//information the client can use to resolve the issue.
	detail string `json:detail`
}

func (e *Error) Error() {
	fmt.Sprintf("%s: %s", e.code, e.message)
}

//Repositories is the set of repository infomation.
type Repositories struct {
	Repositories []string `json:repositories`
}

//Manifest
//type Manifest struct {
//	Name     string `json:name`
//	Tag      string `josn:tag`
//	FsLayers string `json:fsLayers`
//	History
//	SchemaVersion int `json:schemaVersion`
//	Signatures     string
//}

//Tags is the info of repository tag
type Tags struct {
	Name string   `json:name`
	Tags []string `json:tags`
}

type BlobSum struct {
	Blobsum string
}
