package main

import (
	"strom"
	"strom/modules/loginclient"

	"github.com/google/uuid"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	c, err := strom.ClientConnect("127.0.0.1:25566")
	must(err)
	err = c.Start(&loginclient.LoginClient{
		Account: loginclient.Account{
			Username: "449stan_Ytb",
			Uuid:     uuid.MustParse("caec076c80c34637a31b44874ffa009d"),
			Token:    "eyJraWQiOiIwNDkxODEiLCJhbGciOiJSUzI1NiJ9.eyJ4dWlkIjoiMjUzNTQwOTA1OTUxMzgxOSIsImFnZyI6IkFkdWx0Iiwic3ViIjoiMjY4OGQ0ZTktZmZiZi00N2FkLTlkYWYtZDU3ZWY5Yzg5YWE0IiwiYXV0aCI6IlhCT1giLCJucyI6ImRlZmF1bHQiLCJyb2xlcyI6W10sImlzcyI6ImF1dGhlbnRpY2F0aW9uIiwiZmxhZ3MiOlsib3JkZXJzXzIwMjIiLCJtdWx0aXBsYXllciIsInR3b2ZhY3RvcmF1dGgiLCJtc2FtaWdyYXRpb25fc3RhZ2U0Il0sInByb2ZpbGVzIjp7Im1jIjoiY2FlYzA3NmMtODBjMy00NjM3LWEzMWItNDQ4NzRmZmEwMDlkIn0sInBsYXRmb3JtIjoiUENfTEFVTkNIRVIiLCJwZmQiOlt7InR5cGUiOiJtYyIsImlkIjoiY2FlYzA3NmMtODBjMy00NjM3LWEzMWItNDQ4NzRmZmEwMDlkIiwibmFtZSI6IjQ0OXN0YW5fWXRiIn1dLCJuYmYiOjE3NjI4NDc4OTgsImV4cCI6MTc2MjkzNDI5OCwiaWF0IjoxNzYyODQ3ODk4LCJhaWQiOiJjMzZhOWZiNi00ZjJhLTQxZmYtOTBiZC1hZTdjYzkyMDMxZWIifQ.bcyq-FrP4XXESyKBr4TFstsjw8tgjY5cz_p7ansohPV_VfbKqBeLV-E1XYySwKj6BRuc0dAuW6nGdwE9qlv94tRtapVb1q1W__-GguVWOD16TaEi6dmzDyw8--lpI6eL657lPkzFaVeNrSHNPCPn0ltGcZzlOv1XYJwZM8D1yj57x7aXaOMmJqLY8R5zCgfWp_TnGPjJrbd_gV-N5wP9z_vL6jQKj5bvaTkYqItLq5tTFWbX1wAoeSV47D3mtZfx2N4rGgeUSQ1bYanUZ_eMuH6hSwrwxzqn_-HwN2PA7j0VFVUOu54FRqJbw8D4H6tAf3SY1xIjFC5GUu22IbVP3Q",
		},
		Connection: c,
	})
	must(err)
}
