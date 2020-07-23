# Lazy

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
