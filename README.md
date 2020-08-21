# Lazy

[![Build Status](https://travis-ci.org/jiaxinwang/lazy.svg?branch=develop)](https://travis-ci.org/jiaxinwang/lazy)[![codecov](https://codecov.io/gh/jiaxinwang/lazy/branch/develop/graph/badge.svg)](https://codecov.io/gh/jiaxinwang/lazy)

Lazy is a package that aims to generate rest api to operate database from a small amount of configuration or zero code.

### Test Case

```golang
    // before action
	router := route()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/dogs", nil)
	q := req.URL.Query()
	q.Add("before_dog_id", `1`)
	q.Add("before_dog_id", `2`)
	req.URL.RawQuery = q.Encode()

	router.ServeHTTP(w, req)
	ret := Response{}
	err := json.Unmarshal(w.Body.Bytes(), &ret)
	assert.Equal(t, 200, w.Code)
    assert.NoError(t, err)
```

### configuration

```golang
    config := Configuration{
		DB:              gormDB,
		BeforeColumm:    "dog_id",
		BeforeStruct:    &Profile{},
		BeforeTables:    "profiles",
		BeforeResultMap: map[string]string{"dog_id": "id"},
		BeforeAction:    DefaultBeforeAction,
		Table:   "dogs",
		Columm:  "*",
		Struct:  &Dog{},
		Targets: ret,
	}
```

### Delete Method

delete records matching the primary key passed by url params.
**IgnoreAssociations** == true, and when the record has associated data (has-one, has-many, many2many), it fails.

```golang
    config := Configuration{
		DB:              	gormDB,
		Struct:  		 	&Dog{},
		IgnoreAssociations: true,
	}
```

```
```

### TODO

- [x] automated database operation 🚀🚀🚀🚀🚀
- [x] automated rest api handles 🚀🚀🚀🚀🚀
- [x] before action
- [ ] after action
- [ ] action lists
- [ ] automated code generation tools
- [ ] injector
- [ ] more practical validators
- [ ] less configuration
- [ ] simpler, more efficient middleware
