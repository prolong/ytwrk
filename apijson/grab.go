package apijson

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

const videoURL = "https://www.youtube.com/watch?pbj=1&v="
const channelURL = "https://www.youtube.com/browse_ajax?ctoken="
const commentURL = "https://www.youtube.com/comment_service_ajax?action_get_comments=1&pbj=1&ctoken=%[1]s&continuation=%[1]s"

func GrabVideo(videoID string) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(videoURL + videoID)
	setHeaders(&req.Header)

	return req
}

func GrabCommentPage(continuation *CommentContinuation) *fasthttp.Request {
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("POST")
	req.Header.Set("Cookie", continuation.Cookie)
	req.SetRequestURI(fmt.Sprintf(commentURL, continuation.Token))
	req.SetBodyString("session_token=" + continuation.XSRF)
	setHeaders(&req.Header)
	return req
}

func GrabChannelPage(channelID string, page uint) *fasthttp.Request {
	// Generate page URL
	token := GenChannelPageToken(channelID, uint64(page))
	url := channelURL + token

	// Prepare request
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	setHeaders(&req.Header)

	return req
}

func setHeaders(h *fasthttp.RequestHeader) {
	h.Add("Accept-Language", "en-US")
	h.Add("Host", "www.youtube.com")
	h.Add("X-YouTube-Client-Name", "1")
	h.Add("X-YouTube-Client-Version", "2.20170707")
}
