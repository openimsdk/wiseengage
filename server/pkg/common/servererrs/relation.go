package servererrs

import "github.com/openimsdk/tools/errs"

func init() {
	_ = errs.DefaultCodeRelation.Add(errs.RecordNotFoundError, UserIDNotFoundError)
	_ = errs.DefaultCodeRelation.Add(errs.RecordNotFoundError, GroupIDNotFoundError)
	_ = errs.DefaultCodeRelation.Add(errs.DuplicateKeyError, GroupIDExisted)
}
