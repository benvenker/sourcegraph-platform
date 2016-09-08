package jscontext

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	log15 "gopkg.in/inconshreveable/log15.v2"

	"github.com/gorilla/csrf"
	opentracing "github.com/opentracing/opentracing-go"

	"sourcegraph.com/sourcegraph/sourcegraph/api/sourcegraph"
	"sourcegraph.com/sourcegraph/sourcegraph/app/assets"
	"sourcegraph.com/sourcegraph/sourcegraph/cli/buildvar"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/conf"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/conf/feature"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/eventsutil"
	"sourcegraph.com/sourcegraph/sourcegraph/pkg/handlerutil"
)

// JSContext is made available to JavaScript code via the
// "sourcegraph/app/context" module.
type JSContext struct {
	AppURL         string                     `json:"appURL"`
	AccessToken    string                     `json:"accessToken"`
	XHRHeaders     map[string]string          `json:"xhrHeaders"`
	UserAgentIsBot bool                       `json:"userAgentIsBot"`
	AssetsRoot     string                     `json:"assetsRoot"`
	BuildVars      buildvar.Vars              `json:"buildVars"`
	Features       interface{}                `json:"features"`
	User           *sourcegraph.User          `json:"user"`
	GitHubToken    *sourcegraph.ExternalToken `json:"gitHubToken"`
	IntercomHash   string                     `json:"intercomHash"`
}

// NewJSContextFromRequest populates a JSContext struct from the HTTP
// request.
func NewJSContextFromRequest(req *http.Request, uid int, user *sourcegraph.User) (JSContext, error) {
	ctx := req.Context()
	cl := handlerutil.Client(req)

	headers := make(map[string]string)

	if span := opentracing.SpanFromContext(ctx); span != nil {
		if err := opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.TextMapCarrier(headers)); err != nil {
			return JSContext{}, err
		}
	}

	// Propagate Cache-Control no-cache and max-age=0 directives
	// to the requests made by our client-side JavaScript. This is
	// not a perfect parser, but it catches the important cases.
	if cc := req.Header.Get("cache-control"); strings.Contains(cc, "no-cache") || strings.Contains(cc, "max-age=0") {
		headers["Cache-Control"] = "no-cache"
	}

	headers["X-Csrf-Token"] = csrf.Token(req)

	var gitHubToken *sourcegraph.ExternalToken
	if user != nil {
		tok, err := cl.Auth.GetExternalToken(req.Context(), &sourcegraph.ExternalTokenSpec{
			UID:      user.UID,
			Host:     "github.com",
			ClientID: "", // defaults to GitHub client ID in environment
		})
		if err == nil {
			// No need to include the actual access token.
			tok.Token = ""
			gitHubToken = tok
		} else if grpc.Code(err) != codes.NotFound {
			log15.Warn("Error getting GitHub token in NewJSContextFromRequest", "uid", user.UID, "err", err)
		}
	}

	jsctx := JSContext{
		AppURL:         conf.AppURL(ctx).String(),
		XHRHeaders:     headers,
		UserAgentIsBot: isBot(eventsutil.UserAgentFromContext(ctx)),
		AssetsRoot:     assets.URL("/").String(),
		BuildVars:      buildvar.Public,
		Features:       feature.Features,
		User:           user,
		GitHubToken:    gitHubToken,
		IntercomHash:   intercomHMAC(uid),
	}
	cred := sourcegraph.CredentialsFromContext(ctx)
	if cred != nil {
		tok, err := cred.Token()
		if err != nil {
			return JSContext{}, err
		}
		jsctx.AccessToken = tok.AccessToken
	}

	return jsctx, nil
}

var isBotPat = regexp.MustCompile(`(?i:googlecloudmonitoring|pingdom.com|go .* package http|sourcegraph e2etest|bot|crawl|slurp|spider|feed|rss|camo asset proxy|http-client|sourcegraph-client)`)

func isBot(userAgent string) bool {
	return isBotPat.MatchString(userAgent)
}

var intercomSecretKey = os.Getenv("SG_INTERCOM_SECRET_KEY")

func intercomHMAC(uid int) string {
	if uid == 0 || intercomSecretKey == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(intercomSecretKey))
	mac.Write([]byte(strconv.Itoa(uid)))
	return hex.EncodeToString(mac.Sum(nil))
}
