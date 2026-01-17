package xctrl

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/michurin/cnbot/pkg/ctxlog"
	"github.com/michurin/cnbot/pkg/xbot"
	"github.com/michurin/cnbot/pkg/xjson"
	"github.com/michurin/cnbot/pkg/xlog"
	"github.com/michurin/cnbot/pkg/xproc"
)

//nolint:nestif // reason to refactor
func Handler(bot *xbot.Bot, cmd *xproc.Cmd, loggingPatch ctxlog.PatchAttrs) http.HandlerFunc { //nolint:gocognit,cyclop,funlen // reason to refactor
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := ctxlog.ApplyPatch(r.Context(), loggingPatch)
		// TODO mark ctx for logging?
		// TODO put http method to ctx
		// TODO put http content-type to ctx
		body, err := io.ReadAll(r.Body)
		if err != nil {
			xlog.L(ctx, fmt.Errorf("body reading: %w", err))
		}
		method := r.URL.Path[strings.LastIndex(r.URL.Path, "/")+1:]
		data := []byte(nil)
		switch r.Method {
		case http.MethodGet:
			fileID := r.URL.Query().Get("file_id")
			if fileID == "" {
				data, err = bot.API(ctx, &xbot.Request{Method: method})
			} else {
				req, err := xbot.RequestStruct("getFile", map[string]string{"file_id": fileID})
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
				x, err := bot.API(ctx, req)
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
				xlog.L(ctx, x) // TODO!!!!!!
				s := any(nil)
				err = json.Unmarshal(x, &s)
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
				ok, err := xjson.Bool(s, "ok")
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
				filePath, err := xjson.String(s, "result", "file_path")
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
				xlog.L(ctx, fmt.Sprintf("%s %t %s", "ok/filePath", ok, filePath)) // TODO remove
				w.WriteHeader(http.StatusOK)
				err = bot.Download(ctx, filePath, w)
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
				return
			}
		case http.MethodPost:
			ct := r.Header.Get("Content-Type")
			sct, _, err := mime.ParseMediaType(ct)
			if err != nil {
				xlog.L(ctx, err) // TODO response!
				return
			}
			if sct == "application/json" || sct == "multipart/form-data" {
				data, err = bot.API(ctx, &xbot.Request{
					Method:      method,
					ContentType: ct,
					Body:        body,
				})
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
			} else {
				var to int64          // TODO refactor
				var req *xbot.Request // TODO refactor
				to, err = strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
				// TODO add `to` to log context
				req, err = xbot.RequestFromBinary(body, to) //nolint:contextcheck
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
				data, err = bot.API(ctx, req)
				if err != nil {
					xlog.L(ctx, err) // TODO response!
					return
				}
			}
		case "RUN":
			q := r.URL.Query()
			to, err := strconv.ParseInt(r.URL.Query().Get("to"), 10, 64)
			if err != nil {
				xlog.L(ctx, err) // TODO response!
				return
			}
			ctx := xlog.User(ctx, to)
			logCtxPatch := ctxlog.Patch(ctx)
			go func() { //nolint:contextcheck // TODO: limit concurrency
				ctx := ctxlog.ApplyPatch(context.Background(), logCtxPatch)
				// TODO refactor. it is similar to processMessage
				body, err := cmd.Run(ctx, xproc.SanitizeArgs(q["a"]), []string{"tg_x_to=" + strconv.FormatInt(to, 10)})
				if err != nil {
					xlog.L(ctx, err)
					return
				}
				req, err := xbot.RequestFromBinary(body, to)
				if err != nil {
					xlog.L(ctx, err)
					return
				}
				if req == nil { // TODO hmm... it happens?
					xlog.L(ctx, "Script response skipped")
					return
				}
				_, err = bot.API(ctx, req) // TODO check body?
				if err != nil {
					xlog.L(ctx, err)
					return
				}
			}()
			return
		default:
			xlog.L(ctx, fmt.Errorf("method not allowed: %q", r.Method))
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if err != nil {
			xlog.L(ctx, err)
			return
		}
		w.WriteHeader(http.StatusOK)
		// TODO consider `silent=true` parameter and skip writing if present
		_, err = w.Write(data) // TODO consider error
		if err != nil {
			xlog.L(ctx, err)
			return
		}
	}
}
