package repository

//  TODO: Tests for later

// func TestMemoryRepository(t *testing.T) {
// 	repo, err := CreateRepository("unknown_repository")
// 	assert.Error(t, err)
// 	assert.Nil(t, repo)

// 	repo, err = CreateRepository("memory")
// 	assert.NoError(t, err)
// 	assert.NotNil(t, repo)

// 	t.Run("test_memory_repository_set", func(t *testing.T) {
// 		err := repo.Set("test", []byte("test"))
// 		assert.NoError(t, err)
// 	})

// 	t.Run("test_memory_repository_set_ttl", func(t *testing.T) {
// 		err := repo.SetTTL("test", []byte("test"), time.Second*1)
// 		assert.NoError(t, err)

// 		data := repo.Get("test")
// 		assert.Equal(t, []byte("test"), data)

// 		time.Sleep(time.Second * 2)

// 		data = repo.Get("test")
// 		assert.Nil(t, data)
// 	})

// 	t.Run("test_memory_repository_get", func(t *testing.T) {
// 		err := repo.Set("test", []byte("test"))
// 		assert.NoError(t, err)

// 		data := repo.Get("test")
// 		assert.Equal(t, []byte("test"), data)
// 	})

// 	t.Run("test_memory_repository_pop", func(t *testing.T) {
// 		data := repo.Pop("test")
// 		assert.Equal(t, []byte("test"), data)

// 		data = repo.Get("test")
// 		assert.Nil(t, data)
// 	})

// 	t.Run("test_memory_repository_delete", func(t *testing.T) {
// 		err := repo.Delete("test")
// 		assert.NoError(t, err)
// 	})

// 	t.Run("test_memory_repository_all", func(t *testing.T) {
// 		err := repo.Set("test:1", []byte("test"))
// 		assert.NoError(t, err)

// 		id, err := repo.Incr("_test:id")
// 		assert.NoError(t, err)
// 		assert.Equal(t, int64(1), id)

// 		err = repo.Set("test:2", []byte("test2"))
// 		assert.NoError(t, err)

// 		id, err = repo.Incr("_test:id")
// 		assert.NoError(t, err)
// 		assert.Equal(t, int64(2), id)

// 		test := repo.Get("test:1")
// 		assert.Equal(t, []byte("test"), test)

// 		keys, err := repo.Keys("^test:*")
// 		assert.NoError(t, err)
// 		assert.Len(t, keys, 2)

// 		err = repo.Delete("test:1")
// 		assert.NoError(t, err)

// 		test = repo.Get("test:1")
// 		assert.Nil(t, test)

// 		keys, err = repo.Keys("^test:*")
// 		assert.NoError(t, err)
// 		assert.Len(t, keys, 1)

// 		err = repo.Delete("test:2")
// 		assert.NoError(t, err)

// 		keys, err = repo.Keys("^test:*")
// 		assert.NoError(t, err)
// 		assert.Len(t, keys, 0)
// 	})
// }
