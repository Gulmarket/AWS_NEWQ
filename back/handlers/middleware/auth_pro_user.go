package middleware

import (
    "github.com/gulmarket/model"
    "github.com/gulmarket/model/apperrors"
    "log"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/go-playground/validator/v10"
)

func AuthProUser(s model.TokenService) gin.HandlerFunc {
    return func(c *gin.Context) {
        h := authHeader{}

        if err := c.ShouldBindHeader(&h); err != nil {
            if errs, ok := err.(validator.ValidationErrors); ok {

                var invalidArgs []invalidArgument

                for _, err := range errs {
                    invalidArgs = append(invalidArgs, invalidArgument{
                        err.Field(),
                        err.Value().(string),
                        err.Tag(),
                        err.Param(),
                    })
                }

                err := apperrors.NewBadRequest("Invalid request parameters. See invalidArgs")

                c.JSON(err.Status(), gin.H{
                    "error":       err,
                    "invalidArgs": invalidArgs,
                })
                c.Abort()
                return
            }

            err := apperrors.NewInternal()
            c.JSON(err.Status(), gin.H{
                "error": err,
            })
            c.Abort()
            return
        }

        idTokenHeader := strings.Split(h.IDToken, "Bearer ")

        if len(idTokenHeader) < 2 {
            err := apperrors.NewAuthorization("Must provide Authorization header with format `Bearer {token}`")

            c.JSON(err.Status(), gin.H{
                "error": err,
            })
            c.Abort()
            return
        }

        proUser, err := s.ValidateProUserIDToken(idTokenHeader[1])

        if err != nil {
            err := apperrors.NewAuthorization("Provided token is invalid")
            c.JSON(err.Status(), gin.H{
                "error": err,
            })
            c.Abort()
            return
        }

        registrationUrls := []string{
            "/api/update-initials",
            "/api/update-pro-user-city",
            "/api/add-shop",
            "/api/update-pro-user-role",
            "/api/onboarding",
        }

        log.Printf("ahahah %v", proUser.Id)

        if len(proUser.Role) == 0 && !containsString(registrationUrls, c.Request.URL.String()) {
            err := apperrors.NewAuthorization("Role is not chosen")
            c.JSON(err.Status(), gin.H{
                "error": err,
            })
            c.Abort()
            return
        }

        c.Set("pro_user", proUser)

        c.Next()
    }
}

func containsString(slice []string, item string) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}
