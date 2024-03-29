package sq

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWherePartsAppendToSql(t *testing.T) {
	parts := []Sqlizer{
		newWherePart("x = ?", 1),
		newWherePart(nil),
		newWherePart(Eq{"y": 2}),
	}
	var sql strings.Builder
	args, _ := appendToSql(parts, &sql, " AND ", []any{})
	require.Equal(t, "x = ? AND y = ?", sql.String())
	require.Equal(t, []any{1, 2}, args)
}

func TestWherePartsAppendToSqlErr(t *testing.T) {
	parts := []Sqlizer{newWherePart(1)}
	_, err := appendToSql(parts, &strings.Builder{}, "", []any{})
	require.Error(t, err)
}

func TestWherePartNil(t *testing.T) {
	sql, _, _ := newWherePart(nil).ToSql()
	require.Equal(t, "", sql)
}

func TestWherePartErr(t *testing.T) {
	_, _, err := newWherePart(1).ToSql()
	require.Error(t, err)
}

func TestWherePartString(t *testing.T) {
	sql, args, _ := newWherePart("x = ?", 1).ToSql()
	require.Equal(t, "x = ?", sql)
	require.Equal(t, []any{1}, args)
}

func TestWherePartMap(t *testing.T) {
	test := func(pred any) {
		sql, _, _ := newWherePart(pred).ToSql()
		expect := []string{"x = ? AND y = ?", "y = ? AND x = ?"}
		if sql != expect[0] && sql != expect[1] {
			t.Errorf("expected one of %#v, got %#v", expect, sql)
		}
	}
	m := map[string]any{"x": 1, "y": 2}
	test(m)
	test(Eq(m))
}
