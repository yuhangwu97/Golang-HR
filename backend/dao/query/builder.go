package query

import (
	"fmt"
	"strings"
)

type SQLBuilder struct {
	table      string
	selectFields []string
	whereClause []string
	joinClause  []string
	orderClause []string
	groupClause []string
	havingClause []string
	limitValue  int
	offsetValue int
	distinct    bool
	args        []interface{}
}

func NewSQLBuilder(table string) *SQLBuilder {
	return &SQLBuilder{
		table:       table,
		selectFields: []string{"*"},
		whereClause: make([]string, 0),
		joinClause:  make([]string, 0),
		orderClause: make([]string, 0),
		groupClause: make([]string, 0),
		havingClause: make([]string, 0),
		args:        make([]interface{}, 0),
	}
}

func (b *SQLBuilder) Select(fields ...string) *SQLBuilder {
	b.selectFields = fields
	return b
}

func (b *SQLBuilder) Distinct() *SQLBuilder {
	b.distinct = true
	return b
}

func (b *SQLBuilder) Where(field string, operator string, value interface{}) *SQLBuilder {
	clause := fmt.Sprintf("%s %s ?", field, operator)
	b.whereClause = append(b.whereClause, clause)
	b.args = append(b.args, value)
	return b
}

func (b *SQLBuilder) WhereIn(field string, values []interface{}) *SQLBuilder {
	if len(values) == 0 {
		return b
	}
	
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
		b.args = append(b.args, values[i])
	}
	
	clause := fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ", "))
	b.whereClause = append(b.whereClause, clause)
	return b
}

func (b *SQLBuilder) WhereNotIn(field string, values []interface{}) *SQLBuilder {
	if len(values) == 0 {
		return b
	}
	
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
		b.args = append(b.args, values[i])
	}
	
	clause := fmt.Sprintf("%s NOT IN (%s)", field, strings.Join(placeholders, ", "))
	b.whereClause = append(b.whereClause, clause)
	return b
}

func (b *SQLBuilder) WhereBetween(field string, start, end interface{}) *SQLBuilder {
	clause := fmt.Sprintf("%s BETWEEN ? AND ?", field)
	b.whereClause = append(b.whereClause, clause)
	b.args = append(b.args, start, end)
	return b
}

func (b *SQLBuilder) WhereNull(field string) *SQLBuilder {
	clause := fmt.Sprintf("%s IS NULL", field)
	b.whereClause = append(b.whereClause, clause)
	return b
}

func (b *SQLBuilder) WhereNotNull(field string) *SQLBuilder {
	clause := fmt.Sprintf("%s IS NOT NULL", field)
	b.whereClause = append(b.whereClause, clause)
	return b
}

func (b *SQLBuilder) OrWhere(field string, operator string, value interface{}) *SQLBuilder {
	if len(b.whereClause) == 0 {
		return b.Where(field, operator, value)
	}
	
	clause := fmt.Sprintf("OR %s %s ?", field, operator)
	b.whereClause = append(b.whereClause, clause)
	b.args = append(b.args, value)
	return b
}

func (b *SQLBuilder) Join(table, condition string) *SQLBuilder {
	clause := fmt.Sprintf("JOIN %s ON %s", table, condition)
	b.joinClause = append(b.joinClause, clause)
	return b
}

func (b *SQLBuilder) LeftJoin(table, condition string) *SQLBuilder {
	clause := fmt.Sprintf("LEFT JOIN %s ON %s", table, condition)
	b.joinClause = append(b.joinClause, clause)
	return b
}

func (b *SQLBuilder) RightJoin(table, condition string) *SQLBuilder {
	clause := fmt.Sprintf("RIGHT JOIN %s ON %s", table, condition)
	b.joinClause = append(b.joinClause, clause)
	return b
}

func (b *SQLBuilder) InnerJoin(table, condition string) *SQLBuilder {
	clause := fmt.Sprintf("INNER JOIN %s ON %s", table, condition)
	b.joinClause = append(b.joinClause, clause)
	return b
}

func (b *SQLBuilder) OrderBy(field string, direction string) *SQLBuilder {
	if direction == "" {
		direction = "ASC"
	}
	clause := fmt.Sprintf("%s %s", field, strings.ToUpper(direction))
	b.orderClause = append(b.orderClause, clause)
	return b
}

func (b *SQLBuilder) GroupBy(fields ...string) *SQLBuilder {
	b.groupClause = append(b.groupClause, fields...)
	return b
}

func (b *SQLBuilder) Having(condition string, args ...interface{}) *SQLBuilder {
	b.havingClause = append(b.havingClause, condition)
	b.args = append(b.args, args...)
	return b
}

func (b *SQLBuilder) Limit(limit int) *SQLBuilder {
	b.limitValue = limit
	return b
}

func (b *SQLBuilder) Offset(offset int) *SQLBuilder {
	b.offsetValue = offset
	return b
}

func (b *SQLBuilder) Raw(sql string, args ...interface{}) *SQLBuilder {
	// For raw SQL, we'll reset everything and use the provided SQL
	b.whereClause = []string{sql}
	b.args = args
	return b
}

func (b *SQLBuilder) ToSQL() (string, []interface{}) {
	var sql strings.Builder
	
	// SELECT clause
	sql.WriteString("SELECT ")
	if b.distinct {
		sql.WriteString("DISTINCT ")
	}
	sql.WriteString(strings.Join(b.selectFields, ", "))
	
	// FROM clause
	sql.WriteString(fmt.Sprintf(" FROM %s", b.table))
	
	// JOIN clauses
	if len(b.joinClause) > 0 {
		sql.WriteString(" ")
		sql.WriteString(strings.Join(b.joinClause, " "))
	}
	
	// WHERE clause
	if len(b.whereClause) > 0 {
		sql.WriteString(" WHERE ")
		sql.WriteString(strings.Join(b.whereClause, " AND "))
	}
	
	// GROUP BY clause
	if len(b.groupClause) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(b.groupClause, ", "))
	}
	
	// HAVING clause
	if len(b.havingClause) > 0 {
		sql.WriteString(" HAVING ")
		sql.WriteString(strings.Join(b.havingClause, " AND "))
	}
	
	// ORDER BY clause
	if len(b.orderClause) > 0 {
		sql.WriteString(" ORDER BY ")
		sql.WriteString(strings.Join(b.orderClause, ", "))
	}
	
	// LIMIT clause
	if b.limitValue > 0 {
		sql.WriteString(fmt.Sprintf(" LIMIT %d", b.limitValue))
	}
	
	// OFFSET clause
	if b.offsetValue > 0 {
		sql.WriteString(fmt.Sprintf(" OFFSET %d", b.offsetValue))
	}
	
	return sql.String(), b.args
}

func (b *SQLBuilder) Reset() *SQLBuilder {
	b.selectFields = []string{"*"}
	b.whereClause = make([]string, 0)
	b.joinClause = make([]string, 0)
	b.orderClause = make([]string, 0)
	b.groupClause = make([]string, 0)
	b.havingClause = make([]string, 0)
	b.limitValue = 0
	b.offsetValue = 0
	b.distinct = false
	b.args = make([]interface{}, 0)
	return b
}

func (b *SQLBuilder) Clone() *SQLBuilder {
	clone := &SQLBuilder{
		table:       b.table,
		selectFields: make([]string, len(b.selectFields)),
		whereClause: make([]string, len(b.whereClause)),
		joinClause:  make([]string, len(b.joinClause)),
		orderClause: make([]string, len(b.orderClause)),
		groupClause: make([]string, len(b.groupClause)),
		havingClause: make([]string, len(b.havingClause)),
		limitValue:  b.limitValue,
		offsetValue: b.offsetValue,
		distinct:    b.distinct,
		args:        make([]interface{}, len(b.args)),
	}
	
	copy(clone.selectFields, b.selectFields)
	copy(clone.whereClause, b.whereClause)
	copy(clone.joinClause, b.joinClause)
	copy(clone.orderClause, b.orderClause)
	copy(clone.groupClause, b.groupClause)
	copy(clone.havingClause, b.havingClause)
	copy(clone.args, b.args)
	
	return clone
}

// Helper functions for common queries
func (b *SQLBuilder) BuildInsert(data map[string]interface{}) (string, []interface{}) {
	var fields []string
	var placeholders []string
	var args []interface{}
	
	for field, value := range data {
		fields = append(fields, field)
		placeholders = append(placeholders, "?")
		args = append(args, value)
	}
	
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		b.table,
		strings.Join(fields, ", "),
		strings.Join(placeholders, ", "))
	
	return sql, args
}

func (b *SQLBuilder) BuildUpdate(data map[string]interface{}) (string, []interface{}) {
	var setParts []string
	var args []interface{}
	
	for field, value := range data {
		setParts = append(setParts, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	
	sql := fmt.Sprintf("UPDATE %s SET %s", b.table, strings.Join(setParts, ", "))
	
	if len(b.whereClause) > 0 {
		sql += " WHERE " + strings.Join(b.whereClause, " AND ")
		args = append(args, b.args...)
	}
	
	return sql, args
}

func (b *SQLBuilder) BuildDelete() (string, []interface{}) {
	sql := fmt.Sprintf("DELETE FROM %s", b.table)
	
	if len(b.whereClause) > 0 {
		sql += " WHERE " + strings.Join(b.whereClause, " AND ")
	}
	
	return sql, b.args
}