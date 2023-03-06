package main

import (
	"fmt"
	"google.golang.org/appengine/log"
)

func HasPermission(ctx context.Context, uin string) error {
	var err error
	defer func() { // 添加上下文信息
		if err != nil {
			err = fmt.Errorf("HasPermission(%q): %w", uin, err)
		}
	}()
	role, err := getRole(ctx, uin)
	if err != nil {
		return err
	}
	if role != admin {
		return apierr.NewUnauthorizedOperationNoPermission()
	}
	return nil
}

func (s *Service) GetData(ctx context.Context, req *Request) (*Response, error) {
	var err error
	handler := func(e error) (*Response, error) {
		log.Errorf("GetData(%q): %q", req, e.Error()) // 打印错误信息
		r := &Response{}
		var apiError *APIError
		if errors.As(e, &apiError) { // 解包错误并得到“可返回”的错误
			r.Error = apiError.ToError()
		} else { // 无法解包，使用默认的“可返回”的错误
			r.Error = apierr.NewFailedOperationError(e)
		}
	}
	if err := HasPermission(ctx, req.Uin); err != nil {
		return handler(err)
	}
	data, err := retriveData(ctx, req.Key)
	if err != nil {
		return handler(err)
	}
	// return normally
}
