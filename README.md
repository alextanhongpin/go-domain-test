# go-domain-test

This repo demonstrates on how to structure a microservice domain layer.


- domain layer has no dependencies, and contains business logic
- usecase layer calls repositories that returns the domain objects

## Testing

Simplify tests using:

Before:

```go
func TestProductIsMine(t *testing.T) {
	type args struct {
		userID uuid.UUID
	}

	p := factories.NewProduct()
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is mine",
			args: args{userID: p.UserID},
			want: true,
		},
		{
			name: "is not mine",
			args: args{userID: uuid.New()},
			want: false,
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, p.IsMine(tc.args.userID))
		})
	}
}
```

After:
```go
func TestProductIsMine(t *testing.T) {
	p := factories.NewProduct()

	tests := make(map[string]bool)
	tests["is mine"] = p.IsMine(p.UserID) == true
	tests["is not mine"] = p.IsMine(uuid.New()) == false

	for name, ok := range tests {
		t.Run(name, func(t *testing.T) {
			assert.True(t, ok, name)
		})
	}
}
```
