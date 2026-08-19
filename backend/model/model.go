package model

// request:
// {
//     "id": "1", => backend generated
//     "user": "vincent", => get from authentication in backend
//     "message": "This is a post from Vincent", => client
// 	image / video => backend saves it -> url
// 	"url": "https://example.com", => backend
// 	"type": "image" =>  obtained from the suffix of image / video file in backend
// }

type Post struct {
	Id      string `json:"id"`
	User    string `json:"user"`
	Message string `json:"message"`
	Url     string `json:"url"`
	Type    string `json:"type"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int64  `json:"age"`
	Gender   string `json:"gender"`
}
