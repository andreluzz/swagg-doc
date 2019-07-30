package controller

// swagg-doc:controller
// tags:
//   - name: User and Resource
//     description: Enpoints to interact with resources and users
//
// parameters:
//   user_code:
//     name: user_code
//     in: path
//     description: User Code
//     required: true
//     schema:
//       type: string
//       example: username

import (
	"encoding/json"
	"net/http"

	"github.com/andreluzz/swagg-doc/mock/response"
)

// GetUsers swagg-doc:endpoint GET /users
// summary: Returns a list of users.
// tags:
//   - User and Resource
// description: Optional extended description in CommonMark or HTML.
// parameters:
//   - $ref: '#/components/parameters/language-code'
// responses:
//   '200':
//     description: A JSON array of user objects
//     content:
//       application/json:
//         schema:
//           $ref: '#/components/schemas/swagg-doc:interface:Response:User:Array'
//   '400':
//     description: A JSON array of user objects
//     content:
//       application/json:
//         schema:
//           $ref: '#/components/schemas/swagg-doc:interface:Response:User:Array'
func GetUsers(w http.ResponseWriter, r *http.Request) {
	response := response.Response{}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

// PostUser swagg-doc:endpoint POST /users
// summary: Creates a new user.
// tags:
//   - User and Resource
// description: Optional extended description in CommonMark or HTML.
// parameters:
//   - $ref: '#/components/parameters/language-code'
// security:
//   - ApiKeyAuth: []
// requestBody:
//   description: Optional description in *Markdown*
//   required: true
//   content:
//     application/json:
//       schema:
//         $ref: '#/components/schemas/UserCreate'
// responses:
//   '200':
//     description: A JSON with the users attributes
//     content:
//       application/json:
//         schema:
//           $ref: '#/components/schemas/swagg-doc:interface:Response:User'
func PostUser(w http.ResponseWriter, r *http.Request) {
	response := response.Response{}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

// GetUser swagg-doc:endpoint GET /users/{user_code}
// summary: Returns the details of the user.
// tags:
//   - User and Resource
// description: Optional extended description in CommonMark or HTML.
// parameters:
//   - $ref: '#/components/parameters/language-code'
//   - $ref: '#/components/parameters/user_code'
// responses:
//   '200':
//     description: A JSON with the users attributes
//     content:
//       application/json:
//         schema:
//           $ref: '#/components/schemas/swagg-doc:interface:Response:User'
func GetUser(w http.ResponseWriter, r *http.Request) {
	response := response.Response{}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

// UpdateUser swagg-doc:endpoint PATCH /users/{user_code}
// summary: Returns a list of users.
// tags:
//   - User and Resource
// description: Optional extended description in CommonMark or HTML.
// parameters:
//   - $ref: '#/components/parameters/language-code'
//   - $ref: '#/components/parameters/user_code'
// security:
//   - ApiKeyAuth: []
// requestBody:
//   description: Optional description in *Markdown*
//   required: true
//   content:
//     application/json:
//       schema:
//         $ref: '#/components/schemas/UserUpdate'
//       examples:
//         Jessica:
//           value:
//             id: 10
//             name: Jessica Smith
//         Ron:
//           value:
//             id: 11
//             name: Ron Stewart
// responses:
//   '200':
//     description: A JSON with the users attributes
//     content:
//       application/json:
//         schema:
//           $ref: '#/components/schemas/swagg-doc:interface:Response:User'
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	response := response.Response{}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}

// DeleteUser swagg-doc:endpoint DELETE /users/{user_code}
// summary: Returns a list of users.
// tags:
//   - User and Resource
// description: Optional extended description in CommonMark or HTML.
// parameters:
//   - $ref: '#/components/parameters/user_code'
// security:
//   - ApiKeyAuth: []
// responses:
//   '200':
//     description: A JSON array of user objects
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	response := response.Response{}
	respBytes, _ := json.Marshal(response)
	w.Write(respBytes)
}
