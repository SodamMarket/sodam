package service

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/lib/pq"
)

const (
	minPageSize     = 1
	defaultPageSize = 10
	maxPageSize     = 99
)

var queriesCache = make(map[string]*template.Template)

func isUniqueViolation(err error) bool {
	pqerr, ok := err.(*pq.Error)
	return ok && pqerr.Code == "23505"
}

func isForeignKeyViolation(err error) bool {
	pqerr, ok := err.(*pq.Error)
	return ok && pqerr.Code == "23503"
}

//쿼리문을 보충해주는 기능
func buildQuery(text string, data map[string]interface{}) (string, []interface{}, error) {
	//캐싱하여 템플릿 파싱
	t, ok := queriesCache[text]
	if !ok {
		var err error
		t, err = template.New("query").Parse(text)
		if err != nil {
			return "", nil, fmt.Errorf("could not parse sql query template: %v", err)
		}

		queriesCache[text] = t
	}

	//템플릿에 데이터 넣기
	var wr bytes.Buffer
	if err := t.Execute(&wr, data); err != nil {
		return "", nil, fmt.Errorf("could not apply sql query data: %v", err)
	}

	//선언된 변수 교체
	query := wr.String()
	args := []interface{}{}
	for key, val := range data {
		if !strings.Contains(query, "@"+key) {
			continue
		}
		//교체된 변수 다시 대입
		args = append(args, val)
		query = strings.Replace(query, "@"+key, fmt.Sprintf("$%d", len(args)), -1)
	}
	return query, args, nil

}

//limit the page sizes
func normailizePageSize(i int) int {
	if i == 0 {
		return defaultPageSize
	}
	if i < minPageSize {
		return minPageSize
	}
	if i > maxPageSize {
		return maxPageSize
	}
	return i
}
