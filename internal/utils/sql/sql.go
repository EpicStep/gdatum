package sql

import "github.com/huandu/go-sqlbuilder"

type Builder interface {
	BuildWithFlavor(flavor sqlbuilder.Flavor, initialArg ...interface{}) (string, []interface{})
}

func Build(b Builder) (string, []interface{}) {
	return b.BuildWithFlavor(sqlbuilder.ClickHouse)
}
