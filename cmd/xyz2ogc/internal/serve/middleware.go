/**
 * Copyright 2023 Planet Labs PBC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package serve

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func Logger(logger *logrus.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			req := ctx.Request()

			start := time.Now()
			err := next(ctx)
			if err != nil {
				ctx.Error(err)
			}

			res := ctx.Response()
			duration := int64(time.Since(start).Nanoseconds() / 1e6)

			fields := logrus.Fields{
				"http_request_method":       req.Method,
				"http_request_proto":        req.Proto,
				"http_request_uri":          req.RequestURI,
				"http_remote_addr":          ctx.RealIP(),
				"http_response_status_code": res.Status,
				"http_request_referer":      req.Referer(),
				"http_request_user_agent":   req.UserAgent(),
				"http_response_size_bytes":  res.Size,
				"http_response_duration_ms": duration,
			}

			log := logger.WithFields(fields)
			if res.Status >= http.StatusInternalServerError {
				log.WithError(err).Error("error handling request")
			} else {
				log.Debug("handled request")
			}

			return nil
		}
	}
}

func ErrorCommitter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			err := next(ctx)
			if err != nil {
				ctx.Error(err)
			}
			return err
		}
	}
}
