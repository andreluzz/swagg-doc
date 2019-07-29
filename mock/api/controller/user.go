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
	"net/http"
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
//           type: array
//           items:
//             $ref: '#/components/schemas/User'
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dummy"))
}

// PostUser swagg-doc:endpoint POST /users
// summary: Creates a new user.
// tags:
//   - User and Resource
// description: Optional extended description in CommonMark or HTML.
// parameters:
//   - $ref: '#/components/parameters/language-code'
// requestBody:
//   description: Optional description in *Markdown*
//   required: true
//   content:
//     application/json:
//       schema:
//         $ref: '#/components/schemas/User'
// responses:
//   '200':
//     description: A JSON with the users attributes
//     content:
//       application/json:
//         schema:
//           $ref: '#/components/schemas/User'
func PostUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dummy"))
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
//           $ref: '#/components/schemas/User'
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dummy"))
}

// UpdateUser swagg-doc:endpoint PATCH /users/{user_code}
// summary: Returns a list of users.
// tags:
//   - User and Resource
// description: Optional extended description in CommonMark or HTML.
// parameters:
//   - $ref: '#/components/parameters/language-code'
//   - $ref: '#/components/parameters/user_code'
// requestBody:
//   description: Optional description in *Markdown*
//   required: true
//   content:
//     application/json:
//       schema:
//         $ref: '#/components/schemas/User'
// responses:
//   '200':
//     description: A JSON with the users attributes
//     content:
//       application/json:
//         schema:
//           $ref: '#/components/schemas/User'
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dummy"))
}

// DeleteUser swagg-doc:endpoint DELETE /users/{user_code}
// summary: Returns a list of users.
// tags:
//   - User and Resource
// description: Optional extended description in CommonMark or HTML.
// parameters:
//   - $ref: '#/components/parameters/user_code'
// responses:
//   '200':
//     description: A JSON array of user objects
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("dummy"))
}
