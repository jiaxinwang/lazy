package lazy

import (
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestBeforeParams(t *testing.T) {
	type args struct {
		params map[string][]string
	}
	tests := []struct {
		name        string
		args        args
		wantEq      map[string][]string
		wantGt      map[string]string
		wantLt      map[string]string
		wantGte     map[string]string
		wantLte     map[string]string
		wantReduced map[string][]string
	}{
		{
			"empty",
			args{
				map[string][]string{
					"before_name":           []string{"tom"},
					"before_created_at_lte": []string{"2000-01-01"},
					"before_w_lt":           []string{"0.01"},
					"w_lt":                  []string{"0.02"},
					"before_age_gt":         []string{"18"},
					"age_gt":                []string{"19"},
					"before_p_gte":          []string{"32"},
					"gte":                   []string{"33"},
					"size":                  []string{"12"},
					"page":                  []string{"2"},
					"offset":                []string{"100"},
				},
			},
			map[string][]string{"name": []string{"tom"}},
			map[string]string{"age": "18"},
			map[string]string{"w": "0.01"},
			map[string]string{"p": "32"},
			map[string]string{"created_at": "2000-01-01"},
			map[string][]string{
				"w_lt":   []string{"0.02"},
				"age_gt": []string{"19"},
				"gte":    []string{"33"},
				"size":   []string{"12"},
				"page":   []string{"2"},
				"offset": []string{"100"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotEq, gotGt, gotLt, gotGte, gotLte, gotReduced := BeforeLazy(tt.args.params)
			if !cmp.Equal(gotEq, tt.wantEq) {
				t.Errorf("BeforeParams() gotEq = %v, want %v\ndiff=%v", gotEq, tt.wantEq, cmp.Diff(gotEq, tt.wantEq))
			}
			if !cmp.Equal(gotGt, tt.wantGt) {
				t.Errorf("BeforeParams() gotGt = %v, want %v\ndiff=%v", gotGt, tt.wantGt, cmp.Diff(gotGt, tt.wantGt))
			}
			if !cmp.Equal(gotLt, tt.wantLt) {
				t.Errorf("BeforeParams() gotLt = %v, want %v\ndiff=%v", gotLt, tt.wantLt, cmp.Diff(gotLt, tt.wantLt))
			}
			if !cmp.Equal(gotGte, tt.wantGte) {
				t.Errorf("BeforeParams() gotGte = %v, want %v\ndiff=%v", gotGte, tt.wantGte, cmp.Diff(gotGte, tt.wantGte))
			}
			if !cmp.Equal(gotLte, tt.wantLte) {
				t.Errorf("BeforeParams() gotLte = %v, want %v\ndiff=%v", gotLte, tt.wantLte, cmp.Diff(gotLte, tt.wantLte))
			}
			if !cmp.Equal(gotReduced, tt.wantReduced) {
				t.Errorf("BeforeParams() gotReduced = %v, want %v\ndiff=%v", gotReduced, tt.wantReduced, cmp.Diff(gotReduced, tt.wantReduced))
			}
		})
	}
}

// func route() *gin.Engine {
// 	r := gin.Default()
// 	r.Use(gm.Trace)
// 	r.Use(gm.LazyResponse)
// 	r.GET("/dogs", func(c *gin.Context) {
// 		var ret []interface{}
// 		config := Configuration{
// 			DB:              gormDB,
// 			BeforeColumm:    "dog_id",
// 			BeforeStruct:    &Profile{},
// 			BeforeTables:    "profiles",
// 			BeforeResultMap: map[string]string{"dog_id": "id"},
// 			BeforeAction:    DefaultBeforeAction,

// 			Table:   "dogs",
// 			Columm:  "*",
// 			Model:  &Dog{},
// 			Target: ret,
// 		}
// 		c.Set("lazy-configuration", &config)
// 		if _, err := Handle(c); err != nil {
// 			c.Set("error_msg", err.Error())
// 			return
// 		}
// 		c.Set("ret", map[string]interface{}{"data": config.Target})
// 		return
// 	})

// 	return r
// }

// func TestDefaultBeforeAction(t *testing.T) {
// 	router := route()
// 	w := httptest.NewRecorder()
// 	req, _ := http.NewRequest("GET", "/dogs", nil)
// 	q := req.URL.Query()
// 	q.Add("before_dog_id", `1`)
// 	q.Add("before_dog_id", `2`)
// 	req.URL.RawQuery = q.Encode()

// 	router.ServeHTTP(w, req)
// 	ret := Response{}
// 	err := json.Unmarshal(w.Body.Bytes(), &ret)
// 	assert.Equal(t, 200, w.Code)
// 	assert.NoError(t, err)
// 	logrus.Print(ret)
// }

func Test_mergeValues(t *testing.T) {
	type args struct {
		a map[string][]string
		b map[string][]string
	}
	tests := []struct {
		name    string
		args    args
		wantRet map[string][]string
	}{
		{
			"test 1",
			args{
				a: map[string][]string{"dog": []string{"a", "b", "c"}},
				b: map[string][]string{"dog": []string{"c", "e"}},
			},
			map[string][]string{"dog": []string{"a", "b", "c", "e"}},
		},
		{
			"test 2",
			args{
				a: map[string][]string{"dog0": []string{"a", "b", "c"}},
				b: map[string][]string{"dog1": []string{"c", "e"}},
			},
			map[string][]string{
				"dog0": []string{"a", "b", "c"},
				"dog1": []string{"c", "e"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRet := mergeValues(tt.args.a, tt.args.b)
			for k, v := range gotRet {
				if vv, ok := tt.wantRet[k]; ok {
					if !cmp.Equal(v, vv, cmpopts.SortSlices(func(i, j string) bool { return i < j })) {
						t.Errorf("mergeValues() = %v, want %v\ndiff=%v", gotRet, tt.wantRet, cmp.Diff(gotRet, tt.wantRet))
					}
				} else {
					t.Errorf("mergeValues() = %v, want %v\ndiff=%v", gotRet, tt.wantRet, cmp.Diff(gotRet, tt.wantRet))
				}
			}
		})
	}
}
