package jwt

import (
	"gf-chat/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/golang-jwt/jwt/v4"
)

func init() {
	service.RegisterJwt(newJwt())
}

func newJwt() *sJwt {
	return &sJwt{}
}

type sJwt struct {
}

func getSecret() []byte {
	s, _ := g.Cfg().Get(gctx.New(), "app.jwtSecret")
	return s.Bytes()
}

func (Jwt *sJwt) CreateToken(uid string) (token string, err error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
	})
	token, err = at.SignedString(getSecret())
	if err != nil {
		return
	}
	return token, nil
}

func (Jwt *sJwt) ParseToken(token string) (uid string, err error) {
	claim, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return getSecret(), nil
	})
	if err != nil {
		return
	}
	return claim.Claims.(jwt.MapClaims)["uid"].(string),
		nil
}
