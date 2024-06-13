package constant

//post and comment service error
const (
	LongContentError = "content is long"
	EmptyAuthorError = "empty author error"

	WrongIdError = "wrong id error"

	CreatingPostError = "error with creating new post"

	PostNotFoundError = "error with not found post"
	GettingPostError  = "error with getting post"

	WrongPageError     = "wrong page number error"
	WrongPageSizeError = "wrong page size error"

	CommentsNotAllowedError = "comments not allowed for this post"

	CreatingCommentError = "error with creating new comment"
	GettingCommentError  = "error with getting comments"
)

const (
	InternalErrorType = "Internal Server Error"
	BadRequestType    = "Bad Request"
	NotFoundType      = "Not Found Error"
)

const (
	// WrongLimitOffsetError   = "limit and offset must be not negative"
	ThereIsNoSubscriptionError = "there is no connection to the observer"
)
